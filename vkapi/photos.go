package vkapi

import (
	"HwBot/vkapi/vktypes"
	"strconv"
)

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
	UserID    int    `json:"user_id"`
}

// GetMessagesUploadServer returns the server address for photo upload in a private message for a user.
// When uploaded successfully, the photo can be saved using the photos.saveMessagesPhoto method.
func (p *PhotosProvider) GetMessagesUploadServer(peerID int) (*MessagesUploadServer, *Error) {
	args := NewArgs().Set("peer_id", strconv.Itoa(peerID))
	var mus MessagesUploadServer
	err := p.client.Method("photos.getMessagesUploadServer", args, &mus)
	return &mus, err
}

// SavedMessagesPhoto ...
//type SavedMessagesPhoto struct {
//	ID       int    `json:"id"`
//	PID      int    `json:"pid"`
//	AID      int    `json:"aid"`
//	OwnerID  int    `json:"owner_id"`
//	Src      string `json:"src"`
//	SrcBig   string `json:"src_big"`
//	SrcSmall string `json:"src_small"`
//	Created  int64  `json:"created"`
//	SrcXBig  string `json:"src_xbig"`
//	SrcXXBig string `json:"src_xxbig"`
//}

// SaveMessagesPhoto saves a photo after being successfully uploaded.
// URL obtained with photos.getMessagesUploadServer method.
func (p *PhotosProvider) SaveMessagesPhoto(server int, photo, hash string) ([]vktypes.Photo, *Error) {
	args := NewArgs()
	args.Set("server", strconv.Itoa(server))
	args.Set("photo", photo)
	args.Set("hash", hash)
	var smp []vktypes.Photo
	err := p.client.Method("photos.saveMessagesPhoto", args, &smp)
	return smp, err
}
