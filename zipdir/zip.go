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

type File struct {
	Name string
	Data []byte
}

func UnzipToDir(destDir string, zipData []byte) error {
	files, err := UnzipToFileData(zipData)
	if err != nil {
		return err
	}

	for _, file := range files {
		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to make dirs '%s': %v", destDir, err)
		}

		path := filepath.Join(destDir, file.Name)
		outFile, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
		if err != nil {
			if outFile != nil {
				outFile.Close()
			}
			return fmt.Errorf("failed to open or create file '%s': %v", path, err)
		}

		_, err = outFile.Write(file.Data)
		if err != nil {
			return fmt.Errorf("failed to write data to file '%s': %v", path, err)
		}
	}
	return nil
}

func UnzipToFileData(zipData []byte) ([]File, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, fmt.Errorf("failed to read zip data: %v", err)
	}

	var files []File
	for _, file := range zipReader.File {
		fileReader, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to unzip file '%s': %v", file.Name, err)
		}

		data, err := ioutil.ReadAll(fileReader)
		fileReader.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read file data for file '%s': %v", file.Name, err)
		}
		files = append(files, File{
			Name: file.Name,
			Data: data,
		})
	}
	return files, nil
}

func doZip(sourcePath, targetPath string, entries []os.DirEntry, writer *zip.Writer) error {
	for _, entry := range entries {
		sourceName := filepath.Join(sourcePath, entry.Name())
		targetName := filepath.Join(targetPath, entry.Name())
		if entry.IsDir() {
			dirPath := filepath.Join(sourcePath, entry.Name())
			dirEntries, err := os.ReadDir(dirPath)
			if err != nil {
				return fmt.Errorf("failed to read dir: dirPath=%s error=%v", dirPath, err)
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
