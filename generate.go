package main

import (
	"log"

	"github.com/fogleman/gg"
)

func generateImage(chartImageName, newImageName string) {
	size := 1440
	c := gg.NewContext(size, size)

	im, err := gg.LoadImage(chartImageName)
	if err != nil {
		log.Panic(err)
	}

	c.DrawImage(im, 0, 0)

	var (
		num    float64 = 10
		startX float64 = 10
		w      float64 = 90
		h      float64 = 40
	)

	var i float64
	for i = 0; i < num; i++ {
		c.DrawRectangle(startX+i*w+i*startX, 900, w, h)
		c.SetRGBA(0, 0, 0, 0.7)
		c.SetLineWidth(1)
		c.Stroke()
	}

	c.SavePNG(newImageName)
}
