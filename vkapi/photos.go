package vkapi

// ProvidePhotos makes photos provider.
func ProvidePhotos(c *Client) *PhotosProvider {
	return &PhotosProvider{
		APIProvider: APIProvider{c},
	}
}

// PhotosProvider provides messages api.
type PhotosProvider struct {
	APIProvider
}

// MessagesUploadServer struct.
type MessagesUploadServer struct {
	UploadURL string `json:"upload_url"`
	AlbumID   int    `json:"album_id"`
	UserID    UserID `json:"user_id"`
}

// GetMessagesUploadServer returns the server address for photo upload in a private message for a user.
// When uploaded successfully, the photo can be saved using the photos.saveMessagesPhoto method.
func (p *PhotosProvider) GetMessagesUploadServer(peerID ID) (*MessagesUploadServer, error) {
	args := NewArgs().Set("peer_id", itoa(int(peerID)))
	var mus MessagesUploadServer
	err := p.client.Method("photos.getMessagesUploadServer", args, &mus)
	return &mus, err
}

// SaveMessagesPhoto saves a photo after being successfully uploaded.
// URL obtained with photos.getMessagesUploadServer method.
func (p *PhotosProvider) SaveMessagesPhoto(server int, photo, hash string) ([]Photo, error) {
	args := NewArgs()
	args.Set("server", itoa(server))
	args.Set("photo", photo)
	args.Set("hash", hash)
	var smp []Photo
	err := p.client.Method("photos.saveMessagesPhoto", args, &smp)
	return smp, err
}
