package zipdir

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ZipToBytes(path string) ([]byte, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir '%s': %v", path, err)
	}
	buf := bytes.NewBuffer(nil)
	zipWriter := zip.NewWriter(buf)

	err = doZip(path, "", entries, zipWriter)
	if err != nil {
		zipWriter.Close()
		return nil, err
	}

	zipWriter.Close()
	return buf.Bytes(), nil
}

func ZipToDir(inPath, outPath string) error {
	data, err := ZipToBytes(inPath)
	if err != nil {
		return fmt.Errorf("failed to zip directory: %v", err)
	}

	outFile, err := os.OpenFile(outPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output zip: %v", err)
	}
	defer outFile.Close()

	_, err = outFile.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write zip data to file: %v", err)
	}
	return nil
}

func doZip(sourcePath, targetPath string, entries []os.DirEntry, writer *zip.Writer) error {
	for _, entry := range entries {
		sourceName := filepath.Join(sourcePath, entry.Name())
		targetName := filepath.Join(targetPath, entry.Name())
		if entry.IsDir() {
			dirEntries, err := os.ReadDir(entry.Name())
			if err != nil {
				return fmt.Errorf("failed to read dir: %v", err)
			}

			err = doZip(sourceName, targetName, dirEntries, writer)
			if err != nil {
				return err
			}
			continue
		}

		data, err := ioutil.ReadFile(sourceName)
		if err != nil {
			return fmt.Errorf("failed to read file '%s': %v", sourceName, err)
		}

		f, err := writer.Create(targetName)
		if err != nil {
			return fmt.Errorf("failed to add file '%s' to zip: %v", sourceName, err)
		}
		if _, err = f.Write(data); err != nil {
			return fmt.Errorf("failed to write file data for file '%s': %v", sourceName, err)
		}
	}
	return nil
}
