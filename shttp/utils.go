package shttp

import (
	"io"
	"io/fs"
	"mime/multipart"
)

// WriteMultipartRequest writes multipart data from file to request
func WriteMultipartRequest(req *Request, field string, f fs.File) error {
	writer := multipart.NewWriter(req.BodyWriter())
	defer writer.Close()

	s, err := f.Stat()
	if err != nil {
		return err
	}

	part, _ := writer.CreateFormFile(field, s.Name())
	_, err = io.Copy(part, f)
	if err != nil {
		return err
	}

	req.Header.SetContentType(writer.FormDataContentType())

	return nil
}
