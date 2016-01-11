package main

import (
	"fmt"
	"image"
	"image/color"
	"math/cmplx"
	"sync"

	"github.com/disintegration/imaging"
)

const (
	maxIterations = 50
	size          = 1024
	animate       = true
)

func mandelbrot(z, c complex128) complex128 {
	return z*z + c
}

func iteration(c complex128, iterations int) int {
	var z complex128
	for i := 0; i < iterations; i++ {
		z = mandelbrot(z, c)
		r, _ := cmplx.Polar(z)
		if r >= 2 {
			return i
		}
	}
	return 0
}

func colorCode(iteration, iterations int) color.Color {
	shade := uint8(float64(iteration) / float64(iterations) * 255)
	return color.NRGBA{shade, shade, shade, 255}
}

func pixelToReal(w, h, x, y int, offsetX, offsetY, zoom float64) (float64, float64) {
	rx := ((float64(x)/float64(w)+offsetX)*3 - 2) / zoom
	ry := ((float64(y)/float64(h)+offsetY)*3 - 1.5) / zoom
	return rx, ry
}

func fractal(w, h int, offsetX, offsetY, zoom float64, iterations int) image.Image {
	img := imaging.New(w, h, color.White)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			sx, sy := pixelToReal(w, h, x, y, offsetX, offsetY, zoom)
			j := complex(sx, sy)
			i := iteration(j, iterations)
			c := colorCode(i, iterations)
			img.Set(x, y, c)
		}
	}
	return img
}

func main() {
	const zoom = 1.0
	const x, y = 0.0, 0.0

	img := fractal(size, size, x, y, zoom, maxIterations)
	imaging.Save(img, fmt.Sprintf("%.3f_%.3f_%.3f.png", x, y, zoom))

	if !animate {
		return
	}
	var wg sync.WaitGroup
	for i := 0; i < 300; i++ {
		wg.Add(1)
		go func(i int) {
			frame := fractal(size, size, x-0.03*float64(i), y-0.006*float64(i), zoom+0.12*float64(i), maxIterations)
			imaging.Save(frame, fmt.Sprintf("frames/frame_%03d.png", i))
			println(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
