// +build extralibs

package main

import (
	"image"
	"io"
	"io/ioutil"
	"os"
	"unsafe"
	"sync"

	"github.com/docsbox/go-libreofficekit"
)

var office, _ = libreofficekit.NewOffice("/usr/lib/libreoffice/program")
var officeMux = sync.Mutex{}

func genOfficeDocPreview(source io.Reader, target io.Writer, name string, width, height int) error {
	officeMux.Lock()
	defer officeMux.Unlock()

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
	officeMux.Lock()
	defer officeMux.Unlock()

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
