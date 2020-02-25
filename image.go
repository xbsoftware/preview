// +build extralibs

package main

import (
	"gopkg.in/h2non/bimg.v1"
)

func genImagePreview(source, target string, width, height int) error {
	buffer, err := bimg.Read(source)
	if err != nil {
		return err
	}

	newImage, err := bimg.NewImage(buffer).Process(bimg.Options{
		Type:       bimg.PNG,
		Width:      width,
		Height:     height,
		Background: bimg.Color{255, 255, 255},
		Embed:      true,
	})
	if err != nil {
		return err
	}

	bimg.Write(target, newImage)
	return nil
}
