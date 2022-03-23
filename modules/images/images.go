package images

import (
	"bytes"
	"image"

	"github.com/Toffee-iZt/HwBot/bot"
	"github.com/Toffee-iZt/HwBot/logger"
	"github.com/nfnt/resize"
	"github.com/valyala/fasthttp"
)

var Module = bot.Module{
	Name: "images",
	Init: func(l *logger.Logger) bool {
		log = l
		return true
	},
	Terminate: nil,
	Commands:  []*bot.Command{&citgen},
}

var log *logger.Logger

var dlClient fasthttp.Client

func dl(url string) (image.Image, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	err := dlClient.Do(req, resp)
	if err != nil {
		return nil, err
	}

	photo, _, err := image.Decode(bytes.NewReader(resp.Body()))
	return photo, err
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
