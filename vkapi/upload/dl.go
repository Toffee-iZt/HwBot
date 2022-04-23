package upload

import (
	"context"
	"image"
	"net/http"

	"github.com/Toffee-iZt/HwBot/vkapi"
)

// DownloadImage downloads and encodes image.
func DownloadImage(ctx context.Context, c *vkapi.Client, url string) (image.Image, error) {
	req, err := http.NewRequestWithContext(ctx, url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTP(req)
	if err != nil {
		return nil, err
	}

	photo, _, err := image.Decode(resp.Body)
	resp.Body.Close()
	return photo, err
}
