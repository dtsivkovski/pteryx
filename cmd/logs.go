package cmd

import (
	"fmt"
	"os"
	"strings"
)

// open log file for appending
func openLogFile(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file %q: %w", path, err)
	}

	return file, nil
}

// close log file if it is open
func closeLogFile(file *os.File) error {
	if file == nil {
		return nil
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("close log file: %w", err)
	}

	return nil
}

// write mismatch log entry for a file given found extensions
func writeMismatchLog(file *os.File, filePath string, expectedExt string, actualExts []string) error {
	if file == nil {
		return nil
	}

	_, err := fmt.Fprintf(file, "MISMATCH path=%q expected=%q actual=%q\n",
		filePath,
		expectedExt,
		strings.Join(actualExts, ","),
	)
	if err != nil {
		return fmt.Errorf("write log: %w", err)
	}

	return nil
}
