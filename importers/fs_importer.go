package importers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"image/png"
	"io/ioutil"

	"github.com/arcticlimer/bookcatalog/business"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/single_threaded"
	log "github.com/sirupsen/logrus"
)

type FsImporterConfig struct {
	LibraryPath string
	ImagesPath  string
}

type DocumentMetadata struct {
	Pages  int
	Author string
	Title  string
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

// Insert the single/multi-threaded init() here.
var pool pdfium.Pool
var instance pdfium.Pdfium

func init() {
	// Init the PDFium library and return the instance to open documents.
	pool = single_threaded.Init(single_threaded.Config{})

	var err error
	instance, err = pool.GetInstance(time.Second * 30)
	if err != nil {
		log.Fatal(err)
	}
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
	if _, err := os.Stat(i.Config.LibraryPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(i.Config.LibraryPath, os.ModePerm)
		if err != nil {
			return ImportSummary{}, fmt.Errorf("couldn't create library directory: %w", err)
		}
	}

	if _, err := os.Stat(i.Config.ImagesPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(i.Config.ImagesPath, os.ModePerm)
		if err != nil {
			return ImportSummary{}, fmt.Errorf("couldn't create images directory: %w", err)
		}
	}

	for _, path := range paths {
		filename := filepath.Base(path)
		ext := filepath.Ext(path)
		log.Printf("Importing file %s\n", filename)
		file, err := os.Open(path)
		if err != nil {
			log.Errorf("could not open file %s: %s\n", filename, err)
		}
		defer file.Close()

		_, ok := supportedTypes[ext]
		if !ok {
			log.Errorf("unsupported file extension: %s", ext)
			continue
		}

		_, err = i.ImportFile(filename, file)
		if err != nil {
			log.Errorf("could not import file %s: %s\n", filename, err)
			supportedTypes[ext].Failures++
			continue
		}

		supportedTypes[ext].Successes++
	}

	return supportedTypes, nil
}

func (i *FsImporter) ImportFile(filename string, file io.Reader) (business.Document, error) {
	var (
		meta DocumentMetadata
		doc  business.Document
		buf  bytes.Buffer
	)

	tee := io.TeeReader(file, &buf)

	_, err := os.Stat(filepath.Join(i.Config.LibraryPath, filename))
	// TODO: Check in database if file was already imported instead of locally
	if err == nil {
		return doc, fmt.Errorf("file already imported: %s", filename)
	}

	if isPdf(filename) {
		handle, err := openPdf(tee)
		if err != nil {
			return doc, fmt.Errorf("couldn't open pdf file: %w", err)
		}
		// Always close the document, this will release its resources.
		defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
			Document: handle,
		})

		meta = getMetaData(handle)

		output := filepath.Join(i.Config.ImagesPath, fmt.Sprintf("%s.png", removeExt(filename)))
		err = renderPdf(handle, output)
		if err != nil {
			// return doc, fmt.Errorf("%w: couldn't get pdf front page: %s", ErrWarning, err)
			log.Warnf("couldn't get pdf front page: %s", err)
		}
	}

	importedFilePath := filepath.Join(i.Config.LibraryPath, filename)
	destFile, err := os.Create(importedFilePath)
	if err != nil {
		return doc, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, &buf)
	if err != nil {
		return doc, fmt.Errorf("failed to copy file: %w", err)
	}

	err = destFile.Sync()
	if err != nil {
		return doc, fmt.Errorf("failed to sync file: %w", err)
	}

	data := business.CreateDocumentRequest{
		Filename: filename,
		Pages:    meta.Pages,
		Author:   meta.Author,
		Title:    meta.Title,
	}
	doc, err = i.DocumentsRepository.CreateDocument(data)
	if err != nil {
		return doc, fmt.Errorf("could not persist file to the database %s: %s\n", filename, err)
	}

	return doc, nil
}

func removeExt(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func isPdf(path string) bool {
	return filepath.Ext(path) == ".pdf"
}

func openPdf(file io.Reader) (references.FPDF_DOCUMENT, error) {
	pdfBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Open the PDF using PDFium (and claim a worker)
	doc, err := instance.OpenDocument(&requests.OpenDocument{
		File: &pdfBytes,
	})
	if err != nil {
		return "", err
	}

	return doc.Document, nil
}

func getMetaData(doc references.FPDF_DOCUMENT) DocumentMetadata {
	var meta DocumentMetadata

	author, err := getMetaText(doc, "Author")
	if err == nil {
		meta.Author = author
	}

	title, err := getMetaText(doc, "Title")
	if err == nil {
		meta.Title = title
	}

	pageCount, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
		Document: doc,
	})
	if err == nil {
		meta.Pages = pageCount.PageCount
	}

	return meta
}

func getMetaText(doc references.FPDF_DOCUMENT, tag string) (string, error) {
	restag, err := instance.FPDF_GetMetaText(&requests.FPDF_GetMetaText{
		Document: doc,
		Tag:      tag,
	})
	if err != nil {
		return "", err
	}
	return restag.Value, nil
}

func renderPdf(doc references.FPDF_DOCUMENT, output string) error {
	// Render the page in DPI 200.
	pageRender, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{
		DPI: 200, // The DPI to render the page in.
		Page: requests.Page{
			ByIndex: &requests.PageByIndex{
				Document: doc,
				Index:    0,
			},
		}, // The page to render, 0-indexed.
	})
	if err != nil {
		return err
	}

	// Write the output to a file.
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, pageRender.Result.Image)
	if err != nil {
		return err
	}

	return nil
}
