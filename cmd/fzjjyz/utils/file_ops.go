// Package utils provides file operation utilities (DRY principle: eliminate repetition).
package utils

import (
	"fmt"
	"os"

	"codeberg.org/jiangfire/fzjjyz/internal/i18n"
)

// FileExists checks if a file exists (eliminates 12 repetitions).
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ValidateInputFile validates input file (eliminates 12 repetitions).
func ValidateInputFile(path string) error {
	if !FileExists(path) {
		return fmt.Errorf(i18n.T("error.input_file_not_exists"), path)
	}
	return nil
}

// ValidateInputDir validates input directory (eliminates 2 repetitions).
func ValidateInputDir(path string) error {
	if !FileExists(path) {
		return fmt.Errorf(i18n.T("error.source_dir_not_exists"), path)
	}
	info, err := os.Stat(path)
	if err != nil {
		//nolint:wrapcheck
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf(i18n.T("error.input_not_dir"), path)
	}
	return nil
}

// CheckOutputConflict checks output file conflict (eliminates 5 repetitions).
func CheckOutputConflict(output string, force bool) error {
	if !force && FileExists(output) {
		return fmt.Errorf(i18n.T("error.output_file_exists"), output)
	}
	return nil
}

// GetFileSize gets file size (eliminates 4 repetitions).
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		//nolint:wrapcheck
		return 0, err
	}
	return info.Size(), nil
}

// MustOpen opens a file, panics on failure (safety guaranteed by upper validation).
//
//nolint:gosec
func MustOpen(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return f
}

// MustCreate creates a file, panics on failure.
//
//nolint:gosec
func MustCreate(path string) *os.File {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	return f
}
