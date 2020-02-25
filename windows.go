// +build !extralibs

package main

import "errors"

func genOfficeDocPreview(source, target string, width, height int) error {
	return errors.New("Not supported without extra libraries")
}

func genOfficeOtherPreview(source, target string, width, height int) error {
	return errors.New("Not supported without extra libraries")
}

func genImagePreview(source, target string, width, height int) error {
	return errors.New("Not supported without extra libraries")
}
