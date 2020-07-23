// +build extralibs

package main

import (
	"image"
	"io"
	"io/ioutil"
	"os"
	"unsafe"

	"github.com/docsbox/go-libreofficekit"
)

var office, _ = libreofficekit.NewOffice("/usr/lib/libreoffice/program")

func genOfficeDocPreview(source io.Reader, target io.Writer, name string, width, height int) error {
	file, _ := ioutil.TempFile(os.TempDir(), "*"+name)
	io.Copy(file, source)
	file.Close()

	document, err := office.LoadDocument(file.Name())
	if err != nil {
		return err
	}

	m := image.NewRGBA(image.Rect(0, 0, width, height))

	rectangles := document.GetPartPageRectangles()
	if len(rectangles) > 0 {
		r := rectangles[0]
		document.PaintTile(unsafe.Pointer(&m.Pix[0]), width, height, r.Min.X, r.Min.Y, r.Dx(), r.Dy())
		libreofficekit.BGRA(m.Pix)
	}
	document.Close()
	os.Remove(file.Name())

	return saveText(target, m)
}

func genOfficeOtherPreview(source io.Reader, target io.Writer, name string, width, height int) error {
	file, _ := ioutil.TempFile(os.TempDir(), "*"+name)
	io.Copy(file, source)
	file.Close()

	document, err := office.LoadDocument(file.Name())
	if err != nil {
		return err
	}

	w, h := document.GetSize()

	pageWidth := libreofficekit.TwipsToPixels(w, 120)
	pageHeight := libreofficekit.TwipsToPixels(h, 120)

	dx := float32(pageWidth) / float32(width)
	dy := float32(pageHeight) / float32(height)

	min := dy
	if dy > dx {
		min = dx
	}

	if min > 4 {
		min = 4
	}

	pageWidth = int(float32(width) * min)
	pageHeight = int(float32(height) * min)
	w = libreofficekit.PixelsToTwips(pageWidth, 120)
	h = libreofficekit.PixelsToTwips(pageHeight, 120)

	m := image.NewRGBA(image.Rect(0, 0, width, height))

	document.PaintTile(unsafe.Pointer(&m.Pix[0]), width, height, 0, 0, w, h)
	libreofficekit.BGRA(m.Pix)
	document.Close()
	os.Remove(file.Name())

	return saveText(target, m)
}

func convertOffice(source io.Reader, writer io.Writer, name, outType string) error {
	file, _ := ioutil.TempFile(os.TempDir(), "*"+name)
	io.Copy(file, source)
	file.Close()
	inName := file.Name()

	document, err := office.LoadDocument(inName)
	if err != nil {
		return err
	}
	defer document.Close()

	outName := inName + "." + outType
	err = document.SaveAs(outName, outType, "")
	if err != nil {
		return err
	}

	outFile, err := os.Open(outName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(writer, outFile)

	os.Remove(inName)
	os.Remove(outName)
	return err
}
