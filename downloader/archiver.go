package downloader

import (
	"archive/zip"
	"bytes"
	"fmt"
)

type Archiver struct {
	zip *zip.Writer
	buf bytes.Buffer
}

func NewArchiver() *Archiver {
	archiver := &Archiver{buf: bytes.Buffer{}}
	archiver.zip = zip.NewWriter(&archiver.buf)
	return archiver
}

func (a *Archiver) AddFile(name string, content []byte) error {
	f, err := a.zip.Create(name)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	if _, err := f.Write(content); err != nil {
		return fmt.Errorf("write file %v: %w", name, err)
	}
	return nil
}

func (a *Archiver) FinishArchive() ([]byte, error) {
	if err := a.zip.Close(); err != nil {
		return nil, fmt.Errorf("archive finish: %w", err)
	}
	return a.buf.Bytes(), nil
}
