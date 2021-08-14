package vkapi

import (
	"io"
	"io/fs"
	"mime/multipart"

	"github.com/Toffee-iZt/HwBot/shttp"
)

// UploadMessagesPhoto uploads photo from source and returns vk string.
func (c *Client) UploadMessagesPhoto(peerID ID, f fs.File) (string, error) {
	mus, err := c.Photos.GetMessagesUploadServer(peerID)
	if err != nil {
		return "", err
	}

	var res struct {
		Photo  string `json:"photo"`
		Hash   string `json:"hash"`
		Server int    `json:"server"`
	}

	err = c.uploadMultipart(&res, mus.UploadURL, "photo", f)
	if err != nil {
		return "", err
	}

	saved, err := c.Photos.SaveMessagesPhoto(res.Server, res.Photo, res.Hash)
	if err != nil {
		return "", err
	}

	return saved[0].String(), nil
}

// uploadMultipart uploads multipart data from file to uploadURL.
func (c *Client) uploadMultipart(dst interface{}, uploadURL, field string, file fs.File) error {
	fstat, err := file.Stat()
	if err != nil {
		return err
	}

	req := shttp.New(shttp.POSTStr, shttp.URIFromString(uploadURL))

	writer := multipart.NewWriter(req.BodyWriter())
	part, _ := writer.CreateFormFile(field, fstat.Name())
	io.Copy(part, file)
	writer.Close()

	req.Header.SetContentType(writer.FormDataContentType())

	body, err := c.client.Do(req)
	if err != nil {
		return err
	}

	return unmarshal(body, dst)
}
