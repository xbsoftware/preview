// +build !extralibs

package main

import (
	"errors"
	"io"
)

func genOfficeDocPreview(source io.ReadCloser, target io.Writer, name string, width, height int) error {
	return errors.New("Not supported without extra libraries")
}

func genOfficeOtherPreview(source io.Reader, target io.Writer, name string, width, height int) error {
	return errors.New("Not supported without extra libraries")
}

func genImagePreview(source io.Reader, target io.Writer, width, height int) error {
	return errors.New("Not supported without extra libraries")
}

func convertOffice(source io.Reader, writer io.Writer, name, outType string) error {
	return errors.New("Not supported without extra libraries")
}
