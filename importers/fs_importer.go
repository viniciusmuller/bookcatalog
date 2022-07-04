package importers

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/arcticlimer/bookcatalog/business"
	log "github.com/sirupsen/logrus"
)

var ErrWarning = errors.New("warning")

type FsImporterConfig struct {
	LibraryPath string
	ImagesPath  string
}

type FsImporter struct {
	Config              FsImporterConfig
	DocumentsRepository business.DocumentsRepository
}

type ImportSummary map[string]*FiletypeImportResult

type FiletypeImportResult struct {
	Failures  int
	Warnings  int
	Successes int
}

func NewFsImporter(cfg FsImporterConfig, documentsRepository business.DocumentsRepository) FsImporter {
	return FsImporter{Config: cfg, DocumentsRepository: documentsRepository}
}

func (i *FsImporter) ImportFiles(directory string) (ImportSummary, error) {
	supportedTypes := ImportSummary{
		".pdf":  &FiletypeImportResult{},
		".epub": &FiletypeImportResult{},
	}

	var paths []string
	err := filepath.Walk(directory, func(path string, fi os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		_, ok := supportedTypes[ext]
		if ok {
			paths = append(paths, path)
		}

		return nil
	})
	if err != nil {
		return ImportSummary{}, fmt.Errorf("couldn't walk directory %s: %s", directory, err)
	}

	// TODO: move code
	err = os.Mkdir(i.Config.LibraryPath, os.ModePerm)
	err = os.Mkdir(i.Config.ImagesPath, os.ModePerm)

	for _, path := range paths {
		filename := filepath.Base(path)
		ext := filepath.Ext(path)
		log.Printf("Importing file %s\n", filename)
		file, err := os.Open(path)

		defer file.Close()
		if err != nil {
			log.Errorf("could not open file %s: %s\n", filename, err)
		}

		importResult := supportedTypes[ext]
		err = i.ImportFile(filename, file)
		if err != nil {
			if errors.Is(err, ErrWarning) {
				log.Warnln(err)
				importResult.Warnings++
			} else {
				log.Errorf("could not import file %s: %s\n", filename, err)
				importResult.Failures++
				continue
			}
		}

		_, err = i.DocumentsRepository.CreateDocument(filename)
		if err != nil {
			// TODO: what to do with the file that was already imported in that case?
			// TODO: Maybe create the document in the database first
			log.Errorf("could not persist file to the database %s: %s\n", filename, err)
		}

		importResult.Successes++
	}

	return supportedTypes, nil
}

func (i *FsImporter) ImportFile(filename string, file *os.File) error {
	_, err := os.Stat(filepath.Join(i.Config.LibraryPath, filename))
	if err == nil {
		return fmt.Errorf("file already imported: %s", filename)
	}

	importedFilePath := filepath.Join(i.Config.LibraryPath, filename)
	destFile, err := os.Create(importedFilePath)
	defer destFile.Close()
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}

	_, err = io.Copy(destFile, file)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	err = destFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	if isPdf(filename) {
		err = importPdfImage(importedFilePath, i.Config.ImagesPath)
		if err != nil {
			return fmt.Errorf("%w: couldn't get pdf front page: %s", ErrWarning, err)
		}
	}

	return nil
}

func importPdfImage(path, imagesPath string) error {
	fileName := filepath.Base(path)
	err := exec.Command("/bin/sh", "-c", "command -v convert").Run()
	if err != nil {
		return fmt.Errorf("command 'convert' is not available. Please install imagemagick")
	}

	pdfFirstPage := fmt.Sprintf("%s[0]", path)
	pageDestionation := filepath.Join(imagesPath, fmt.Sprintf("%s.jpg", fileNameWithoutExtTrimSuffix(fileName)))
	err = exec.Command("convert", pdfFirstPage, pageDestionation).Run()
	if err != nil {
		return fmt.Errorf("command 'convert' errored out: %w", err)
	}

	return nil
}

func fileNameWithoutExtTrimSuffix(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func isPdf(path string) bool {
	return filepath.Ext(path) == ".pdf"
}
