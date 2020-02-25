package main

import (
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"io/ioutil"
	"log"
)

func initFont() *truetype.Font {
	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return nil
	}

	return f
}

func getImagePreview(source, target string, width, height int) error {
	src, err := imaging.Open(source)
	if err != nil {
		return err
	}

	return getImagePreviewFromRGBA(src, target, width, height)
}

func getImagePreviewFromRGBA(src image.Image, target string, width, height int) error {
	dst := imaging.Thumbnail(src, width, height, imaging.Lanczos)
	err := imaging.Save(dst, target)

	if err != nil {
		return err
	}
	return nil
}
