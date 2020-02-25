package main

import (
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/styles"
	"image"
	"image/color"
)

var colors map[chroma.Colour]image.Image
var style *chroma.Style

func init() {
	colors = make(map[chroma.Colour]image.Image)
	style = styles.MonokaiLight

	for _, t := range style.Types() {
		c := style.Get(t).Colour
		cOrigin := c

		var r, g, b uint8
		b = uint8(c % 256)
		c = c / 256
		g = uint8(c % 256)
		c = c / 256
		r = uint8(c % 256)

		colors[cOrigin] = image.NewUniform(color.RGBA{r, g, b, 255})
	}
}
