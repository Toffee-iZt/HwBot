package vkapi

import "encoding/json"

// vk types
const (
	AttTypePhoto        = "photo"
	AttTypeAudioMessage = "audio_message"
)

// Attachment struct.
type Attachment struct {
	Object interface{}
	Type   string
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (a *Attachment) UnmarshalJSON(data []byte) error {
	var fvk map[string]json.RawMessage
	err := unmarshal(data, &fvk)
	if err != nil {
		return err
	}

	a.Type = string(fvk["type"])
	a.Type = a.Type[1 : len(a.Type)-1] // remove quotes

	raw := fvk[a.Type]

	switch a.Type {
	case AttTypePhoto:
		a.Object = new(Photo)
	case AttTypeAudioMessage:
		a.Object = new(AudioMessage)
	default:
		a.Object = raw
		return nil
	}

	return unmarshal(raw, a.Object)
}

type attachment struct {
	AccessKey string `json:"access_key"`
	ID        int    `json:"id"`
	OwnerID   ID     `json:"owner_id"`
}

func (a *attachment) string(typ string) string {
	s := string(typ) + itoa(int(a.OwnerID)) + "_" + itoa(a.ID)
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
	return p.string(AttTypePhoto)
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
	return a.string(AttTypeAudioMessage)
}
