// +build extralibs

package main

import (
	"gopkg.in/h2non/bimg.v1"
	"io"
	"io/ioutil"
)

func genImagePreview(source io.Reader, target io.Writer, width, height int) error {
	buffer, err := ioutil.ReadAll(source)
	if err != nil {
		return err
	}

	newImage, err := bimg.NewImage(buffer).Process(bimg.Options{
		Type:       bimg.PNG,
		Width:      width,
		Height:     height,
		Background: bimg.Color{255, 255, 255},
		Embed:      true,
		Crop:       true,
	})
	if err != nil {
		return err
	}

	target.Write(newImage)
	return nil
}
