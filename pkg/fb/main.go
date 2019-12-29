// Copyright 2019-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fb

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"

	framebuffer "github.com/orangecms/go-framebuffer"
)

const fbdev = "/dev/fb0"

func DrawImageAt(img image.Image, posx int, posy int) error {
	fb, err := framebuffer.Init(fbdev)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}
	width, height := fb.Size()
	stride := fb.Stride()
	bpp := fb.Bpp()
	fmt.Fprintf(os.Stdout, "resolution: %v %v %v %v\n", width, height, stride, bpp)
	buf := make([]byte, width*height*bpp)
	DrawOnBufAt(buf, img, posx, posy, stride, bpp)
	err = ioutil.WriteFile(fbdev, buf, 0600)
	if err != nil {
		return fmt.Errorf("Error writing to framebuffer: %v", err)
	}
	return nil
}

func DrawOnBufAt(
	buf []byte,
	img image.Image,
	posx int,
	posy int,
	stride int,
	bpp int,
) {
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			offset := bpp * ((posy+y)*stride + posx + x)
			// framebuffer is BGR(A)
			buf[offset+0] = byte(b)
			buf[offset+1] = byte(g)
			buf[offset+2] = byte(r)
			if bpp >= 4 {
				buf[offset+3] = byte(a)
			}
		}
	}
}

func DrawScaledOnBufAt(
	buf []byte,
	img image.Image,
	posx int,
	posy int,
	factor int,
	stride int,
	bpp int,
) {
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			for sx := 1; sx <= factor; sx++ {
				for sy := 1; sy <= factor; sy++ {
					offset := bpp * ((posy+y*factor+sy)*stride + posx + x*factor + sx)
					buf[offset+0] = byte(b)
					buf[offset+1] = byte(g)
					buf[offset+2] = byte(r)
					if bpp == 4 {
						buf[offset+3] = byte(a)
					}
				}
			}
		}
	}
}

func DrawDigits(posx int, posy int, size int) error {
	fb, err := framebuffer.Init(fbdev)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}
	width, height := fb.Size()
	stride := fb.Stride()
	bpp := fb.Bpp()
	fmt.Fprintf(os.Stdout, "resolution: %v %v\n", width, height)
	buf := make([]byte, width*height*bpp)
	for digit := 0; digit <= 9; digit++ {
		DrawDigitAt(buf, digit, 205+digit*15, 130, stride, bpp, size)
	}
	err = ioutil.WriteFile(fbdev, buf, 0600)
	if err != nil {
		return fmt.Errorf("Error writing to framebuffer: %v", err)
	}
	return nil
}

func DrawScaledImageAt(img image.Image, posx int, posy int, factor int) error {
	fb, err := framebuffer.Init(fbdev)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}
	width, height := fb.Size()
	stride := fb.Stride()
	bpp := fb.Bpp()
	fmt.Fprintf(os.Stdout, "resolution: %v %v %v %v\n", width, height, stride, bpp)
	buf := make([]byte, width*height*bpp)
	DrawScaledOnBufAt(buf, img, posx, posy, factor, stride, bpp)
	/*
	 */
	size := 3
	digits := []int{4, 7, 1, 1}
	for i, digit := range digits {
		DrawDigitAt(buf, digit, posx+i*15, posy-40, stride, bpp, size)
	}
	err = ioutil.WriteFile(fbdev, buf, 0600)
	if err != nil {
		return fmt.Errorf("Error writing to framebuffer: %v", err)
	}
	return nil
}