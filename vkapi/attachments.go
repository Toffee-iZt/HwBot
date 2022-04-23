package vkapi

import "strconv"

// vk types
const (
	AttachmentTypePhoto        = "photo"
	AttachmentTypeAudioMessage = "audio_message"
)

type uploaded struct {
	AccessKey string `json:"access_key"`
	ID        int    `json:"id"`
	OwnerID   ID     `json:"owner_id"`
}

func (a *uploaded) string(typ string) string {
	s := string(typ) + strconv.Itoa(int(a.OwnerID)) + "_" + strconv.Itoa(a.ID)
	if a.AccessKey != "" {
		s += "_" + a.AccessKey
	}
	return s
}

// Photo struct.
type Photo struct {
	uploaded

	AlbumID int    `json:"album_id"`
	UserID  UserID `json:"user_id"`
	Date    int64  `json:"date"`
	HasTag  bool   `json:"has_tag"`
	Text    string `json:"text"`

	Sizes []struct {
		URL    string `json:"url"`
		Type   string `json:"type"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"sizes"`
}

func (p *Photo) String() string {
	return p.string(AttachmentTypePhoto)
}

// AudioMessage struct.
type AudioMessage struct {
	uploaded

	Duration int    `json:"duration"`
	LinkOGG  string `json:"link_ogg"`
	LinkMP3  string `json:"link_mp3"`
	// Waveform []int `json:"waveform"`
}

func (a *AudioMessage) String() string {
	return a.string(AttachmentTypeAudioMessage)
}
