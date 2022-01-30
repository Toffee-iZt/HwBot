package vkapi

import "strconv"

// vk types
const (
	AttachmentTypePhoto        = "photo"
	AttachmentTypeAudioMessage = "audio_message"
)

// Attachment struct.
type Attachment struct {
	Photo        *Photo
	AudioMessage *AudioMessage
	Type         string
}

type attachment struct {
	AccessKey string `json:"access_key"`
	ID        int    `json:"id"`
	OwnerID   ID     `json:"owner_id"`
}

func (a *attachment) string(typ string) string {
	s := string(typ) + strconv.Itoa(int(a.OwnerID)) + "_" + strconv.Itoa(a.ID)
	if a.AccessKey != "" {
		s += "_" + a.AccessKey
	}
	return s
}

// Photo struct.
type Photo struct {
	attachment

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
	attachment

	Duration int    `json:"duration"`
	LinkOGG  string `json:"link_ogg"`
	LinkMP3  string `json:"link_mp3"`
	// Waveform []int `json:"waveform"`
}

func (a *AudioMessage) String() string {
	return a.string(AttachmentTypeAudioMessage)
}
