// +build extralibs

package main

import (
	"image"
	"unsafe"

	"github.com/docsbox/go-libreofficekit"
)

var office, _ = libreofficekit.NewOffice("/usr/lib/libreoffice/program")

func genOfficeDocPreview(source, target string, width, height int) error {
	document, err := office.LoadDocument(source)
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

	return saveText(target, m)
}

func genOfficeOtherPreview(source, target string, width, height int) error {
	document, err := office.LoadDocument(source)
	if err != nil {
		return err
	}

	w, h := document.GetSize()

	pageWidth := libreofficekit.TwipsToPixels(w, 120)
	pageHeight := libreofficekit.TwipsToPixels(h, 120)

	m := image.NewRGBA(image.Rect(0, 0, pageWidth, pageHeight))

	document.PaintTile(unsafe.Pointer(&m.Pix[0]), pageWidth, pageHeight, 0, 0, pageWidth, pageHeight)
	libreofficekit.BGRA(m.Pix)
	document.Close()

	return saveText(target, m)
}
