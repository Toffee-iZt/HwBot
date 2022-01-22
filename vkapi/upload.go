package vkapi

import (
	"io"
	"mime/multipart"
)

// UploadMessagesPhoto uploads photo from source and returns vk string.
func (c *Client) UploadMessagesPhoto(peerID ID, fname string, data io.Reader) (string, *Error) {
	mus, err := c.PhotosGetMessagesUploadServer(peerID)
	if err != nil {
		return "", err
	}

	var res struct {
		Photo  string `json:"photo"`
		Hash   string `json:"hash"`
		Server int    `json:"server"`
	}

	c.uploadMultipart(&res, mus.UploadURL, "photo", fname, data)

	saved, err := c.PhotosSaveMessagesPhoto(res.Server, res.Photo, res.Hash)
	if err != nil {
		return "", err
	}

	return saved[0].String(), nil
}

// uploadMultipart uploads multipart data from file to uploadURL.
func (c *Client) uploadMultipart(dst interface{}, uploadURL, field string, fname string, data io.Reader) {
	req := NewRequest(uploadURL, "POST")

	writer := multipart.NewWriter(req.BodyWriter())
	part, _ := writer.CreateFormFile(field, fname)
	io.Copy(part, data)
	writer.Close()

	req.Header.SetContentType(writer.FormDataContentType())

	c.Do(req, dst)
}
