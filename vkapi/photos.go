package vkapi

// MessagesUploadServer struct.
type MessagesUploadServer struct {
	UploadURL string `json:"upload_url"`
	AlbumID   int    `json:"album_id"`
	UserID    UserID `json:"user_id"`
}

// PhotosGetMessagesUploadServer returns the server address for photo upload in a private message for a user.
// When uploaded successfully, the photo can be saved using the photos.saveMessagesPhoto method.
func (c *Client) PhotosGetMessagesUploadServer(peerID ID) (*MessagesUploadServer, error) {
	var mus MessagesUploadServer
	return &mus, c.method(&mus, "photos.getMessagesUploadServer", vkargs{
		"peer_id": peerID,
	})
}

// PhotosSaveMessagesPhoto saves a photo after being successfully uploaded.
// URL obtained with photos.getMessagesUploadServer method.
func (c *Client) PhotosSaveMessagesPhoto(server int, photo, hash string) ([]*Photo, error) {
	var smp []*Photo
	return smp, c.method(&smp, "photos.saveMessagesPhoto", vkargs{
		"server": server,
		"photo":  photo,
		"hash":   hash,
	})
}
