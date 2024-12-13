package fileutil

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

// ZipFiles takes a list of files and zips them into a single zip archive.
func ZipFiles(zipName string, files map[string]string) error {
	zipFile, err := os.Create(zipName)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Iterate over each file and add it to the zip
	for filename, filePath := range files {
		err := addFileToZip(zipWriter, filename, filePath)
		if err != nil {
			return fmt.Errorf("failed to add file %s to zip: %v", filename, err)
		}
	}

	return nil
}

// addFileToZip adds a single file to the zip archive.
func addFileToZip(zipWriter *zip.Writer, filename, filePath string) error {
	zipEntry, err := zipWriter.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create zip entry: %v", err)
	}

	// Open the file to add to the zip archive
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	// Copy the content of the file to the zip archive
	_, err = io.Copy(zipEntry, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content to zip: %v", err)
	}

	return nil
}
