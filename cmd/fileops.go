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
func runPathCheck(path string, allowDirectory bool, recursive bool) error {
	stats := &ScanStats{} // initialize stats to track results

	// run path check with stats tracking
	err := runPathCheckWithStats(path, allowDirectory, recursive, stats)
	if err != nil {
		return err
	}

	printScanSummary(stats)

	return nil
}

// check if path is file or dir and runs appropriate check, while tracking stats
func runPathCheckWithStats(path string, allowDirectory bool, recursive bool, stats *ScanStats) error {
	info, err := os.Stat(path)
	if err != nil { // return error if stat fails
		return fmt.Errorf("stat %q: %w", path, err)
	}

	if info.IsDir() { // if directory, check if allowed
		if !allowDirectory {
			return fmt.Errorf("%q is a directory; use -d or --directory to check directories", path)
		}

		return runDirectoryCheck(path, recursive, stats)
	}

	if recursive {
		return fmt.Errorf("-r or --recursive can only be used with -d or --directory")
	}

	return runFileCheck(path, false, stats)
}


// checks each file signature in the entire directory
func runDirectoryCheck(dirPath string, recursive bool, stats *ScanStats) error {
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

			return runFileCheck(path, true, stats)
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
		if err := runFileCheck(filePath, true, stats); err != nil {
			return err
		}
	}

	return nil
}

// checks the file signature
func runFileCheck(filePath string, indent bool, stats *ScanStats) error {
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
			fmt.Printf("%s└─ File signature matches: %s", red, reset)
			for i := range actualExts {
				if i == len(actualExts)-1 {
					fmt.Printf("%s\n", actualExts[i])
				} else {
					fmt.Printf("%s, ", actualExts[i])
				}
			}
		}
		return nil
	}
	stats.Passed++

	if indent { // indent for dir output
		fmt.Print("  ")
	}

	fmt.Printf("%s✓%s %s %sis a %s file%s\n", cyan, reset, filePath, cyan, ext, reset)
	return nil
}

// check if signature matches anything else
func findSignatureMatches(buffer []byte) []string {
	var matches []string // account for multiple signature possibilities
	bestScore := 0 // use scoring to prioritize more specific signatures

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
          ▒▒            
       ▓▓▒▒▓            
    ▓ ▓▒▓▒▓▓            
    ▓ ▓▒ ▓▓▒▒           
    ▓▓▓    ▓▓           
    ▓▓▓     ██          
      ▒█                
`)
	fmt.Printf("\n%sBack from the hunt! Here's what Pteryx caught:%s\n", magenta, reset)
	fmt.Printf("%sChecked:%s %d\n", magenta, reset, stats.Checked)
	fmt.Printf("%sPassed:%s %d\n", cyan, reset, stats.Passed)
	fmt.Printf("%sFailed:%s %d\n", red, reset, stats.Failed)
	fmt.Printf("%sSkipped:%s %d\n", magenta, reset, stats.Skipped)
}