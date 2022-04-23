package images

import (
	"image"

	"github.com/Toffee-iZt/HwBot/bot"
	"github.com/nfnt/resize"
)

var Module = bot.Module{
	Name:     "images",
	Commands: []*bot.Command{&citgen},
}

func crop(img image.Image, rect image.Rectangle) image.Image {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	simg, ok := img.(subImager)
	if !ok {
		return cropCopy(img, rect)
	}
	return simg.SubImage(rect)
}

func cropCopy(img image.Image, rect image.Rectangle) image.Image {
	rect = rect.Intersect(img.Bounds())
	rgbaImg := &image.RGBA{}
	if !rect.Empty() {
		rgbaImg = image.NewRGBA(rect)
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				rgbaImg.Set(x, y, img.At(x, y))
			}
		}
	}
	return rgbaImg
}

func scaleQuad(img image.Image, needSize int) image.Image {
	return resize.Resize(uint(needSize), uint(needSize), img, resize.Lanczos3)
}
