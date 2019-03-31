package main

import "image"

type sixelImage struct {
	source string
	image  *image.Image
	sixel  []byte
}

func newSixelImage(source string) *sixelImage {
	return nil
}
