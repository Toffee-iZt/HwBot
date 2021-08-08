package vkapi

import (
	"HwBot/common/storage"
	"HwBot/shttp"
	"encoding/json"
	"io/fs"
)

// uploadMultipart uploads multipart data from file to uploadURL.
func (c *Client) uploadMultipart(dst interface{}, uploadURL, field string, f fs.File) error {
	req, resp := shttp.New(shttp.POSTStr, &shttp.URIBuffer{B: []byte(uploadURL)})
	defer shttp.Release(req, resp)

	err := shttp.WriteMultipartRequest(req, field, f)
	if err != nil {
		return err
	}

	err = c.client.Do(req, resp)
	if err != nil {
		return err
	}

	return json.Unmarshal(resp.Body(), dst)
}

// PhotoCache contains photos as vk string.
var PhotoCache = storage.RuntimeStorage()

// UploadMessagesPhoto uploads photo from source and returns vk string.
func (c *Client) UploadMessagesPhoto(peerID int, f fs.File) (string, error) {
	//if cached := PhotoCache.Get(info.Name()); cached != nil {
	//	return cached.(string), nil
	//}

	mus, vkerr := c.Photos.GetMessagesUploadServer(peerID)
	if vkerr != nil {
		return "", vkerr
	}

	var res struct {
		Photo  string `json:"photo"`
		Hash   string `json:"hash"`
		Server int    `json:"server"`
	}

	err := c.uploadMultipart(&res, mus.UploadURL, "photo", f)
	if err != nil {
		return "", err
	}

	saved, vkerr := c.Photos.SaveMessagesPhoto(res.Server, res.Photo, res.Hash)
	if vkerr != nil {
		return "", vkerr
	}

	p := saved[0].String()
	//PhotoCache.Store(src, p)

	return p, nil
}
