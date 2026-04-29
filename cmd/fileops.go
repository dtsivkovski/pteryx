package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// checks the file signature
func runFileCheck(filePath string) error {
	// check if signature map known
	ext := strings.ToLower(filepath.Ext(filePath))
	signatures, ok := fileSignatures[ext]

	if !ok {
		return fmt.Errorf("unsupported file extension %q", ext)
	}

	// try to open file
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open %q: %w", filePath, err)
	}
	defer f.Close()

	// make 16 byte buffer to account for largest signatures
	buffer := make([]byte, 16)

	// read attempt
	n, err := f.Read(buffer)
	if err != nil {
		return fmt.Errorf("read %q: %w", filePath, err)
	}

	// check file signature against known
	matched := false
	for _, signature := range signatures {
		if n < len(signature) { // file too small
			continue // try next signature
		}

		// check if sig matches for current one
		signatureMatches := true
		for i, b := range signature {
			if buffer[i] != b {
				signatureMatches = false
				break
			}
		}

		if signatureMatches { // return early if matched
			matched = true
			break
		}
	}

	if !matched {
		fmt.Printf("%s✗%s %s %sis not a %s file%s\n", red, reset, filePath, red, ext, reset)
			actualExts := findSignatureMatches(buffer)
			if len(actualExts) == 0 {
				fmt.Print("└─ Unknown signature\n")
			} else {
				fmt.Printf("%s└─ File signature matches: %s", red, reset)
				for i := range actualExts {
					if i == len(actualExts) - 1 {
						fmt.Printf("%s\n", actualExts[i])
					} else {
						fmt.Printf("%s, ", actualExts[i])
					}
				}
			}
		return nil
	}

	fmt.Printf("%s✓%s %s %sis a %s file%s\n", cyan, reset, filePath, cyan, ext, reset)
	return nil
}

// check if signature matches anything else 
func findSignatureMatches(buffer []byte) []string {
	var matches []string // account for multiple signature possibilities

	for ext, signatures := range fileSignatures {
		for _, signature := range signatures {
			if len(buffer) < len(signature) { // if length is shorter, skip
				continue
			}

			// match signature
			matched := true
			for i, b := range signature {
				if buffer[i] != b {
					matched = false
					break
				}
			}

			if matched {
				matches = append(matches, ext)
				break
			}
		}
	}

	return matches
}
