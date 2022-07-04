package importers

import "os"

type DocumentImporter interface {
	ImportFile(filename string, file *os.File) error
}
