package downloader

import (
	"archive/zip"
	"bytes"
	"context"
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

type ArchiveResult struct {
	Ok  []byte
	Err error
}

type GoArchiver struct {
	in       chan *MIMEFile
	out      chan ArchiveResult
	inClosed bool
}

func NewGoArchiver() *GoArchiver {
	ga := &GoArchiver{
		in:  make(chan *MIMEFile),
		out: make(chan ArchiveResult),
	}

	go func() {
		defer close(ga.out)
		archiver := NewArchiver()

		for file := range ga.in {
			if err := archiver.AddFile(file.Name, file.Content); err != nil {
				ga.out <- ArchiveResult{Err: err}
				return
			}
		}

		result, err := archiver.FinishArchive()
		if err != nil {
			ga.out <- ArchiveResult{Err: err}
		} else {
			ga.out <- ArchiveResult{Ok: result}
		}
	}()

	return ga
}

func (ga *GoArchiver) AddFile(name string, content []byte) {
	ga.in <- &MIMEFile{
		Name:    name,
		Content: content,
	}
}

func (ga *GoArchiver) Close() {
	if !ga.inClosed {
		close(ga.in)
	}
	ga.inClosed = true
}

func (ga *GoArchiver) Finish(ctx context.Context) ([]byte, error) {
	ga.Close()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-ga.out:
		if result.Err != nil {
			return nil, result.Err
		}
		return result.Ok, nil
	}
}
