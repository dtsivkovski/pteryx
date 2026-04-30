package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// hash structs
type HashBaseline struct {
	GeneratedAt string       `json:"generated_at"`
	Root        string       `json:"root"`
	Algorithm   string       `json:"algorithm"`
	Files       []HashRecord `json:"files"`
}

type HashRecord struct {
	Path     string `json:"path"`
	SHA256   string `json:"sha256"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
}

type HashCompareStats struct {
	Unchanged int
	Changed   int
	Added     int
	Deleted   int
}

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Create or compare SHA-256 file hash baselines",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var hashCreateCmd = &cobra.Command{
	Use:   "create <file-or-directory>",
	Short: "Create a baseline file for SHA-256 file hashes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		allowDirectory, err := cmd.Flags().GetBool("directory")
		if err != nil {
			return err
		}

		recursive, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		return runHashCreate(args[0], allowDirectory, recursive, output)
	},
}

var hashCompareCmd = &cobra.Command{
	Use:   "compare <file-or-directory>",
	Short: "Compare files against an existing SHA-256 baseline file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		allowDirectory, err := cmd.Flags().GetBool("directory")
		if err != nil {
			return err
		}

		recursive, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}

		input, err := cmd.Flags().GetString("input")
		if err != nil {
			return err
		}

		return runHashCompare(args[0], allowDirectory, recursive, input)
	},
}

func init() {
	hashCmd.AddCommand(hashCreateCmd)
	hashCmd.AddCommand(hashCompareCmd)

	hashCreateCmd.Flags().BoolP("directory", "d", false, "hash files in a directory")
	hashCreateCmd.Flags().BoolP("recursive", "r", false, "recursively hash directories")
	hashCreateCmd.Flags().StringP("output", "o", "pteryx.hash", "hash baseline output file") // optional change output file

	hashCompareCmd.Flags().BoolP("directory", "d", false, "compare files in a directory")
	hashCompareCmd.Flags().BoolP("recursive", "r", false, "recursively compare directories")
	hashCompareCmd.Flags().StringP("input", "i", "pteryx.hash", "hash baseline input file") // optional change input file
}

func runHashCreate(path string, allowDirectory bool, recursive bool, outputPath string) error {
	logFile, err := openLogFile("pteryx.log") // open log file
	if err != nil {
		return err
	}
	defer func() {
		if err := closeLogFile(logFile); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	// collect all hash records for this path
	records, root, err := collectHashRecords(path, allowDirectory, recursive, outputPath)
	if err != nil {
		return err
	}

	baseline := HashBaseline{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Root:        root,
		Algorithm:   "sha256",
		Files:       records,
	}

	raw, err := json.MarshalIndent(baseline, "", "  ")
	if err != nil {
		return fmt.Errorf("encode hash baseline: %w", err)
	}

	if err := os.WriteFile(outputPath, append(raw, '\n'), 0644); err != nil {
		return fmt.Errorf("write hash baseline %q: %w", outputPath, err)
	}

	if err := writeAuditLog(logFile, "HASH_CREATE path=%q output=%q files=%d", path, outputPath, len(records)); err != nil {
		return err
	}

	fmt.Printf("%s✵ Hash baseline created:%s %s\n", magenta, reset, outputPath)
	fmt.Printf("%sFiles hashed:%s %d\n", cyan, reset, len(records))

	return nil
}

func runHashCompare(path string, allowDirectory bool, recursive bool, inputPath string) error {
	logFile, err := openLogFile("pteryx.log")
	if err != nil {
		return err
	}
	defer func() {
		if err := closeLogFile(logFile); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	raw, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read hash baseline %q: %w", inputPath, err)
	}

	var baseline HashBaseline
	if err := json.Unmarshal(raw, &baseline); err != nil {
		return fmt.Errorf("parse hash baseline %q: %w", inputPath, err)
	}

	if baseline.Algorithm != "sha256" {
		return fmt.Errorf("unsupported hash algorithm %q", baseline.Algorithm)
	}

	currentRecords, _, err := collectHashRecords(path, allowDirectory, recursive, inputPath)
	if err != nil {
		return err
	}

	stats, err := compareHashRecords(baseline.Files, currentRecords, logFile)
	if err != nil {
		return err
	}

	if err := writeAuditLog(logFile, "HASH_COMPARE_COMPLETED path=%q input=%q unchanged=%d changed=%d added=%d deleted=%d",
		path, inputPath, stats.Unchanged, stats.Changed, stats.Added, stats.Deleted); err != nil {
		return err
	}

	printHashCompareSummary(stats)
	return nil
}

// runs through specified directory or file to get hash value
func collectHashRecords(path string, allowDirectory bool, recursive bool, baselinePath string) ([]HashRecord, string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, "", fmt.Errorf("stat %q: %w", path, err)
	}

	if info.IsDir() { // directory hash crawl
		if !allowDirectory {
			return nil, "", fmt.Errorf("%q is a directory; use -d or --directory to hash directories", path)
		}

		return collectDirectoryHashRecords(path, recursive, baselinePath)
	}

	// if recursive but no -d value
	if recursive {
		return nil, "", fmt.Errorf("-r or --recursive can only be used with -d or --directory")
	}

	// hash just one file
	record, err := hashFile(path, filepath.Base(path))
	if err != nil {
		return nil, "", err
	}

	root, err := filepath.Abs(path)
	if err != nil {
		root = path
	}

	return []HashRecord{record}, root, nil
}

// get all hash records for a directory
func collectDirectoryHashRecords(dirPath string, recursive bool, baselinePath string) ([]HashRecord, string, error) {
	root, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, "", fmt.Errorf("resolve directory %q: %w", dirPath, err)
	}

	// get baseline file
	baselineAbs, err := filepath.Abs(baselinePath)
	if err != nil {
		baselineAbs = baselinePath
	}

	logAbs, err := filepath.Abs("pteryx.log")
	if err != nil {
		logAbs = "pteryx.log"
	}

	// iterathe through files and hash each one
	var records []HashRecord
	addRecord := func(path string) error {
		fileAbs, err := filepath.Abs(path)
		if err == nil {
			if fileAbs == baselineAbs || fileAbs == logAbs {
				return nil
			}
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("make relative path for %q: %w", path, err)
		}

		record, err := hashFile(path, filepath.ToSlash(rel))
		if err != nil {
			return err
		}

		records = append(records, record)
		return nil
	}

	if recursive { // walk dir and all inner dirs if recursive
		err = filepath.WalkDir(dirPath, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("walk %q: %w", path, err)
			}

			if entry.IsDir() {
				return nil
			}

			// add record to path
			return addRecord(path)
		})
		if err != nil {
			return nil, "", err
		}
	} else { // read through all entries in current directory level
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return nil, "", fmt.Errorf("read directory %q: %w", dirPath, err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			if err := addRecord(filepath.Join(dirPath, entry.Name())); err != nil {
				return nil, "", err
			}
		}
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Path < records[j].Path
	})

	return records, root, nil
}

// computer the hash of the file
func hashFile(path string, displayPath string) (HashRecord, error) {
	f, err := os.Open(path)
	if err != nil {
		return HashRecord{}, fmt.Errorf("open %q: %w", path, err)
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return HashRecord{}, fmt.Errorf("hash %q: %w", path, err)
	}

	info, err := f.Stat()
	if err != nil {
		return HashRecord{}, fmt.Errorf("stat %q: %w", path, err)
	}

	return HashRecord{
		Path:     displayPath,
		SHA256:   hex.EncodeToString(hash.Sum(nil)),
		Size:     info.Size(),
		Modified: info.ModTime().Format(time.RFC3339),
	}, nil
}

// compare hash records to expected values
func compareHashRecords(expected []HashRecord, current []HashRecord, logFile *os.File) (HashCompareStats, error) {
	expectedByPath := make(map[string]HashRecord)
	currentByPath := make(map[string]HashRecord)

	// compare by path
	for _, record := range expected {
		expectedByPath[record.Path] = record
	}
	for _, record := range current {
		currentByPath[record.Path] = record
	}

	var stats HashCompareStats

	// iterate through all and compare
	for _, record := range expected {
		currentRecord, ok := currentByPath[record.Path]
		if !ok { // file must have been deleted
			stats.Deleted++
			fmt.Printf("%sDELETED%s   %s\n", red, reset, record.Path)
			continue
		}

		// if hashes do not match, they were changed
		if !strings.EqualFold(record.SHA256, currentRecord.SHA256) {
			stats.Changed++
			fmt.Printf("%sCHANGED%s   %s\n", red, reset, record.Path)
			if err := writeHashMismatchLog(logFile, record, currentRecord); err != nil {
				return stats, err
			}
			continue
		}

		// if hashes fully match, then unchanged
		stats.Unchanged++
		fmt.Printf("%sUNCHANGED%s %s\n", cyan, reset, record.Path)
	}

	for _, record := range current {
		if _, ok := expectedByPath[record.Path]; !ok {
			stats.Added++
			fmt.Printf("%sADDED%s     %s\n", magenta, reset, record.Path)
		}
	}

	return stats, nil
}

// print hash summary
func printHashCompareSummary(stats HashCompareStats) {
	fmt.Printf("\n\n")
	fmt.Print(`             █▓▓        
             ▒▒▒▒▒▒     
            ▒▓     ▓▓▓  
           ▒▒           
          ▒▒            
       ▓▓▒▒▓            
    ▓ ▓▒▓▒▓▓            
    ▓ ▓▒ ▓▓▒▒      ▓    
    ▓▓▓    ▓▓   █     ▓  
    ▓▓▓     ██    ▓ █   
      ▒█                
`)
	fmt.Printf("\n%sPteryx counted all its eggs:%s\n", magenta, reset)
	fmt.Printf("%sUnchanged:%s %d\n", cyan, reset, stats.Unchanged)
	fmt.Printf("%sChanged:%s %d\n", red, reset, stats.Changed)
	fmt.Printf("%sAdded:%s %d\n", magenta, reset, stats.Added)
	fmt.Printf("%sDeleted:%s %d\n", red, reset, stats.Deleted)
}
