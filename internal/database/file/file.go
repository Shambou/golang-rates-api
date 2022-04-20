package database

import (
	"context"
)

type File struct {
	BaseCurrency string
	FxPath       string
	Ext          string
}

// NewFile - returns a pointer to a file struct
func NewFile(baseCurrency string, fxPath string, ext string) *File {
	return &File{
		BaseCurrency: baseCurrency,
		FxPath:       fxPath,
		Ext:          ext,
	}
}

// Ping - pings the db
func (f *File) Ping(ctx context.Context) error {
	return nil
}
