package uploader

import (
	"bytes"
	"io"
)

type ContentProvider interface {
	GetReader() (io.Reader, error)
}

type bytesContentProvider struct {
	data []byte
}

func (this bytesContentProvider) GetReader() (io.Reader, error) {
	return bytes.NewReader(this.data), nil
}

func NewBytesContentProvider(data []byte) ContentProvider {
	return bytesContentProvider{data}
}
