package vkapi

import (
	"crypto/md5"
	"encoding/json"
	"hash"
	"io"
	"io/fs"
	"mime/multipart"
	"sync"

	"github.com/Toffee-iZt/HwBot/shttp"
)

var photoCache = newCache()

// UploadMessagesPhoto uploads photo from source and returns vk string.
func (c *Client) UploadMessagesPhoto(peerID int, f fs.File) (string, error) {
	data, h, err := readHash(f, photoCache)
	if err != nil {
		return "", err
	}

	cachedData := photoCache.get(h)
	cachedData.mut.Lock()
	defer cachedData.mut.Unlock()

	s := cachedData.val[peerID]
	if s != "" {
		return s, nil
	}

	mus, vkerr := c.Photos.GetMessagesUploadServer(peerID)
	if vkerr != nil {
		return "", vkerr
	}

	var res struct {
		Photo  string `json:"photo"`
		Hash   string `json:"hash"`
		Server int    `json:"server"`
	}

	err = c.uploadMultipart(&res, mus.UploadURL, "photo", data)
	if err != nil {
		return "", err
	}

	saved, vkerr := c.Photos.SaveMessagesPhoto(res.Server, res.Photo, res.Hash)
	if vkerr != nil {
		return "", vkerr
	}

	s = saved[0].String()
	cachedData.val[peerID] = s

	return s, nil
}

// uploadMultipart uploads multipart data from file to uploadURL.
func (c *Client) uploadMultipart(dst interface{}, uploadURL, field string, data []byte) error {
	req, resp := shttp.New(shttp.POSTStr, shttp.URIFromString(uploadURL))
	defer shttp.Release(req, resp)

	writer := multipart.NewWriter(req.BodyWriter())
	part, _ := writer.CreateFormField(field)
	_, err := part.Write(data)
	if err != nil {
		return err
	}
	writer.Close()

	req.Header.SetContentType(writer.FormDataContentType())

	err = c.client.Do(req, resp)
	if err != nil {
		return err
	}

	return json.Unmarshal(resp.Body(), dst)
}

func readHash(f fs.File, c *cache) ([]byte, md5hash, error) {
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, md5hash{}, err
	}

	return data, c.hash(data), nil
}

type md5hash = [md5.Size]byte

func newCache() *cache {
	return &cache{
		h: md5.New(),
		c: make(map[md5hash]*cached),
	}
}

type cached struct {
	val map[int]string
	mut sync.Mutex
}

type cache struct {
	h hash.Hash
	c map[md5hash]*cached
	m sync.Mutex
}

func (c *cache) hash(data []byte) md5hash {
	c.m.Lock()
	c.h.Reset()
	c.h.Write(data)
	var sum md5hash
	c.h.Sum(sum[:0])
	c.m.Unlock()
	return sum
}

func (c *cache) get(hash md5hash) *cached {
	c.m.Lock()
	v := c.c[hash]
	if v == nil {
		v = &cached{
			val: make(map[int]string),
		}
		c.c[hash] = v
	}
	c.m.Unlock()
	return v
}
