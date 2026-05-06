package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type ScanStats struct {
	Checked int
	Passed  int
	Failed  int
	Skipped int
}

// creates stats object and runs the path check
func runPathCheck(path string, allowDirectory bool, recursive bool, logResults bool, verbose bool) error {
	stats := &ScanStats{} // initialize stats to track results

	logFile := (*os.File)(nil) // init empty log file variable

	// open log file (or create) for appending
	if logResults {
		var err error
		logFile, err = openLogFile("pteryx.log")
		if err != nil {
			return err
		}
		defer func() {
			if err := closeLogFile(logFile); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}()
	}

	// run path check with stats tracking
	err := runPathCheckWithStats(path, allowDirectory, recursive, stats, logFile, verbose)
	if err != nil {
		return err
	}

	printScanSummary(stats)

	return nil
}

// check if path is file or dir and runs appropriate check, while tracking stats
func runPathCheckWithStats(path string, allowDirectory bool, recursive bool, stats *ScanStats, logFile *os.File, verbose bool) error {
	info, err := os.Stat(path)
	if err != nil { // return error if stat fails
		return fmt.Errorf("stat %q: %w", path, err)
	}

	if info.IsDir() { // if directory, check if allowed
		if !allowDirectory {
			return fmt.Errorf("%q is a directory; use -d or --directory to check directories", path)
		}

		return runDirectoryCheck(path, recursive, stats, logFile, verbose)
	}

	if recursive {
		return fmt.Errorf("-r or --recursive can only be used with -d or --directory")
	}

	return runFileCheck(path, false, stats, logFile, verbose)
}

// checks each file signature in the entire directory
func runDirectoryCheck(dirPath string, recursive bool, stats *ScanStats, logFile *os.File, verbose bool) error {
	if recursive {
		// walk entire dir recursively
		fmt.Printf("%s✵ Swooping into directory (RECURSIVE):%s %s\n", magenta, reset, dirPath)
		return filepath.WalkDir(dirPath, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("walk %q: %w", path, err)
			}

			if entry.IsDir() {
				return nil
			}

			return runFileCheck(path, true, stats, logFile, verbose)
		})
	}

	// read all entries in directory
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("read directory %q: %w", dirPath, err)
	}

	fmt.Printf("%s✵ Flying over directory (NON-RECURSIVE):%s %s\n", magenta, reset, dirPath)

	// check each file in directory
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, entry.Name())
		if err := runFileCheck(filePath, true, stats, logFile, verbose); err != nil {
			return err
		}
	}

	return nil
}

// checks the file signature
func runFileCheck(filePath string, indent bool, stats *ScanStats, logFile *os.File, verbose bool) error {
	// check if signature map known
	ext := strings.ToLower(filepath.Ext(filePath))
	signatures, ok := fileSignatures[ext]

	if !ok {
		stats.Skipped++

		if indent {
			return nil
		}

		return fmt.Errorf("unsupported file extension %q", ext)
	}

	// try to open file
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open %q: %w", filePath, err)
	}
	defer f.Close()

	// read enough bytes for expected or alternative sigs
	readLength := maxSignatureReadLength(signatures)
	if maxFileSignatureReadBytes > readLength {
		readLength = maxFileSignatureReadBytes
	}
	buffer := make([]byte, readLength)

	stats.Checked++

	// read attempt
	n, err := io.ReadFull(f, buffer)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		return fmt.Errorf("read %q: %w", filePath, err)
	}

	// check file signature against known
	matched := false
	for _, signature := range signatures {
		if signatureMatches(buffer[:n], signature) {
			matched = true
			break
		}
	}

	if !matched {
		stats.Failed++
		if indent { // indent for dir output
			fmt.Print("  ")
		}

		fmt.Printf("%s✗%s %s %sis not a %s file%s\n", red, reset, filePath, red, ext, reset)
		actualExts := findSignatureMatches(buffer[:n])
		sortExtensions(actualExts)

		if indent { // indent for dir output
			fmt.Print("  ")
		}

		if len(actualExts) == 0 {
			fmt.Print("└─ Unknown signature\n")
		} else {
			fmt.Printf("%s└─ File signature matches: %s%s\n", red, reset, strings.Join(actualExts, ", "))
		}

		if logFile != nil {
			if err := writeMismatchLog(logFile, filePath, ext, actualExts); err != nil {
				return err
			}
		}

		return nil
	}
	stats.Passed++

	if !verbose {
		return nil
	}

	if indent { // indent for dir output
		fmt.Print("  ")
	}

	fmt.Printf("%s✓%s %s %sis a %s file%s\n", cyan, reset, filePath, cyan, ext, reset)
	return nil
}

// check if signature matches anything else
func findSignatureMatches(buffer []byte) []string {
	var matches []string // account for multiple signature possibilities
	bestScore := 0       // use scoring to prioritize more specific signatures

	for ext, signatures := range fileSignatures {
		for _, signature := range signatures {
			if signatureMatches(buffer, signature) {
				// use its specificity score to see if better match than before
				score := signatureSpecificityScore(signature)
				if score < bestScore {
					break
				}

				if score > bestScore {
					bestScore = score
					matches = nil
				}

				matches = append(matches, ext)
				break
			}
		}
	}

	return matches
}

// print all data to summarize entire scan results
func printScanSummary(stats *ScanStats) {
	fmt.Printf("\n\n")
	fmt.Print(`             █▓▓        
             ▒▒▒▒▒▒     
            ▒▓     ▓▓▓  
           ▒▒           
          ▒▒         ▒   ▒
       ▓▓▒▒▓      ▒ ▒  ▓▒
    ▓ ▓▒▓▒▓▓       ▒▓▒▓ ▓ ▒
    ▓ ▓▒ ▓▓▒▒       ▒▓▓▓
    ▓▓▓    ▓▓     █▒▓▓▓▓▓█
    ▓▓▓     ██     ██████    
      ▒█                
`)
	fmt.Printf("\n%s✵ Back from the hunt (checked signatures)! Here's what Pteryx caught:%s\n", magenta, reset)
	fmt.Printf("%sChecked:%s %d\n", magenta, reset, stats.Checked)
	fmt.Printf("%sPassed:%s %d\n", cyan, reset, stats.Passed)
	fmt.Printf("%sFailed:%s %d\n", red, reset, stats.Failed)
	fmt.Printf("%sSkipped:%s %d\n", magenta, reset, stats.Skipped)
}
