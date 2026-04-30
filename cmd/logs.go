package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// open log file for appending
func openLogFile(path string) (*os.File, error) {
	// open file
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file %q: %w", path, err)
	}

	// write start of log entry
	_, err = fmt.Fprintln(file, "LOG STARTED AT", time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("write log start: %w", err)
	}

	return file, nil
}

// close log file if it is open
func closeLogFile(file *os.File) error {
	if file == nil {
		return nil
	}

	// write end of log entry
	_, err := fmt.Fprintln(file, "LOG ENDED AT", time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("write log end: %w", err)
	}

	// close file
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
