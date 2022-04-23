package upload

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/Toffee-iZt/HwBot/vkapi"
)

// MessagesPhoto uploads photo from source and returns vk string.
func MessagesPhoto(c *vkapi.Client, peerID vkapi.ID, fname string, data io.Reader) (string, *vkapi.Error) {
	mus, err := c.PhotosGetMessagesUploadServer(peerID)
	if err != nil {
		return "", err
	}

	var res struct {
		Photo  string `json:"photo"`
		Hash   string `json:"hash"`
		Server int    `json:"server"`
	}

	upload(c, &res, mus.UploadURL, "photo", fname, data)

	saved, err := c.PhotosSaveMessagesPhoto(res.Server, res.Photo, res.Hash)
	if err != nil {
		return "", err
	}

	return saved[0].String(), nil
}

func upload(c *vkapi.Client, dst interface{}, uploadURL, field string, fname string, data io.Reader) {
	var body bytes.Buffer
	req, _ := http.NewRequest(http.MethodPost, uploadURL, &body)

	writer := multipart.NewWriter(&body)
	part, _ := writer.CreateFormFile(field, fname)
	io.Copy(part, data)
	writer.Close()

	req.Header.Set("Content-Type", writer.FormDataContentType())

	c.Do(req, dst)
}
