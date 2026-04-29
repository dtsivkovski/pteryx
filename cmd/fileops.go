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
		return nil
	}

	fmt.Printf("%s✓%s %s %sis a %s file%s\n", cyan, reset, filePath, cyan, ext, reset)
	return nil
}