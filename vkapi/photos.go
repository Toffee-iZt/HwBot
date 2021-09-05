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
	args := newArgs()
	args.Set("peer_id", itoa(int(peerID)))
	var mus MessagesUploadServer
	return &mus, c.method("photos.getMessagesUploadServer", args, &mus)
}

// PhotosSaveMessagesPhoto saves a photo after being successfully uploaded.
// URL obtained with photos.getMessagesUploadServer method.
func (c *Client) PhotosSaveMessagesPhoto(server int, photo, hash string) ([]*Photo, error) {
	args := newArgs()
	args.Set("server", itoa(server))
	args.Set("photo", photo)
	args.Set("hash", hash)
	var smp []*Photo
	return smp, c.method("photos.saveMessagesPhoto", args, &smp)
}
