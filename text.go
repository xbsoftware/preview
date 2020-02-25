package main

import (
	"bufio"
	"github.com/alecthomas/chroma/lexers"
	"github.com/golang/freetype"
	"golang.org/x/image/font"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func genTxtPreview(source, target string, width, height int) error {
	c, rgba, size := textCanvas(width, height)

	sFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sFile.Close()
	scanner := bufio.NewScanner(sFile)

	err = printText(c, scanner, size)
	if err != nil {
		return err
	}

	return saveText(target, rgba)
}

func genCodePreview(source, target string, width, height int) error {
	c, rgba, size := textCanvas(width, height)

	sFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sFile.Close()

	err = printCode(c, path.Base(source), sFile, size)
	if err != nil {
		return err
	}

	return saveText(target, rgba)
}

func textCanvas(width, height int) (*freetype.Context, image.Image, float64) {
	fg := image.Black
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(rgba, rgba.Bounds(), image.White, image.ZP, draw.Src)

	// scale down fonts for small images
	size := *fontsize
	if width < 400 {
		size = size * float64(width) / 400
	}

	c := freetype.NewContext()
	c.SetDPI(*fontdpi)
	c.SetFont(initFont())
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetHinting(font.HintingNone)

	return c, rgba, size
}

func saveText(target string, rgba image.Image) error {
	_, err := os.Stat(target)
	if err == nil {
		os.Remove(target)
	}
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	tExt := path.Ext(target)
	if tExt == ".png" {
		return png.Encode(f, rgba)
	} else {
		return jpeg.Encode(f, rgba, &jpeg.Options{Quality: 80})
	}
}

func printCode(c *freetype.Context, name string, reader io.Reader, size float64) error {
	pt := freetype.Pt(10, 10+int(c.PointToFixed(size)>>6))

	lexer := lexers.Match(name)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	iterator, err := lexer.Tokenise(nil, string(contents))
	if err != nil {
		return err
	}

	left := pt.X
	for _, t := range iterator.Tokens() {
		style := style.Get(t.Type)
		c.SetSrc(colors[style.Colour])

		lines := strings.Split(strings.Replace(strings.Replace(t.Value, "\r", "", -1), "\t", "    ", -1), "\n")
		last := len(lines) - 1
		for i := range lines {
			pt, err = c.DrawString(lines[i], pt)
			if err != nil {
				return err
			}
			if i != last {
				pt.Y += c.PointToFixed(size * 1.5)
				pt.X = left
			}
		}
	}
	return nil
}

func printText(c *freetype.Context, scanner *bufio.Scanner, size float64) error {
	pt := freetype.Pt(10, 10+int(c.PointToFixed(size)>>6))
	for scanner.Scan() {
		_, err := c.DrawString(scanner.Text(), pt)
		if err != nil {
			return err
		}
		pt.Y += c.PointToFixed(size * 1.5)
	}

	return nil
}
