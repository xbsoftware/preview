package main

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"io/ioutil"
	"log"
)

func initFont() *truetype.Font {
	fontBytes, err := ioutil.ReadFile(Config.Text.FontFile)
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
