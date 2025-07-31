package archiver

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Create(archivePath string, sourceFiles []string) error {
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("could not create archive file: %w", err)
	}
	defer archiveFile.Close()

	gzipWriter := gzip.NewWriter(archiveFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	for _, filePath := range sourceFiles {
		if err := addFileToTar(tarWriter, filePath); err != nil {
			return fmt.Errorf("failed to add file '%s' to archive: %w", filePath, err)
		}
	}

	return nil
}

// addFileToTar - это вспомогательная функция для добавления одного файла в tar-архив.
func addFileToTar(tarWriter *tar.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}
	header.Name = filepath.Base(filePath)

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(tarWriter, file); err != nil {
		return err
	}

	return nil
}

// Extract распаковывает .tar.gz архив (archivePath) в указанную директорию (destDir).
func Extract(archivePath, destDir string) error {
	archiveFile, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("could not open archive file: %w", err)
	}
	defer archiveFile.Close()

	gzipReader, err := gzip.NewReader(archiveFile)
	if err != nil {
		return fmt.Errorf("could not create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	absDestDir, err := filepath.Abs(destDir)
	if err != nil {
		return fmt.Errorf("could not get absolute path for destination: %w", err)
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading tar header: %w", err)
		}

		targetPath := filepath.Join(absDestDir, header.Name)

		if !strings.HasPrefix(targetPath, absDestDir) {
			return fmt.Errorf("illegal file path in archive: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("could not create directory: %w", err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
				return fmt.Errorf("could not create parent directory for file: %w", err)
			}

			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("could not create file: %w", err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("could not write to file: %w", err)
			}
			outFile.Close()
		default:
			return fmt.Errorf("unsupported type in archive: %c for file %s", header.Typeflag, header.Name)
		}
	}

	return nil
}
