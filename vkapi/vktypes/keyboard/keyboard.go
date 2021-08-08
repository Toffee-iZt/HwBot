package keyboard

import "encoding/json"

//
const (
	MaxButtonsOnLine = 5
	MaxDefaultLines  = 10
	MaxInlineLines   = 6
)

//
const (
	// VkBlue
	ColorPrimary = "primary"
	// White
	ColorSecondary = "secondary"
	// Red
	ColorNegative = "negative"
	// Green
	ColorPositive = "positive"
)

//
const (
	ButtonTypeText     = "text"
	ButtonTypeLocation = "location"
	ButtonTypeVkPay    = "vkpay"
	ButtonTypeOpenApp  = "open_app"
	ButtonTypeOpenLink = "open_link"
	ButtonTypeCallback = "callback"
)

// Empty creates empty keyboard to remove an existing.
func Empty() *Keyboard {
	return New(true, false)
}

// New creates new keyboard.
func New(oneTime bool, inline bool) *Keyboard {
	return &Keyboard{
		OneTime: oneTime,
		Inline:  inline,
	}
}

// Keyboard struct.
type Keyboard struct {
	Buttons [][]Button `json:"buttons"`
	OneTime bool       `json:"one_time"`
	Inline  bool       `json:"inline"`
}

func (k *Keyboard) String() string {
	if k == nil {
		return ""
	}
	b, _ := json.Marshal(k)
	return string(b)
}

// AddRow adds line of buttons.
func (k *Keyboard) AddRow() bool {
	if k.Inline && len(k.Buttons) >= MaxInlineLines {
		return false
	}
	if !k.Inline && len(k.Buttons) >= MaxDefaultLines {
		return false
	}
	k.Buttons = append(k.Buttons, make([]Button, 0, 1))
	return true
}

func (k *Keyboard) add(b Button) {
	r := k.Buttons[len(k.Buttons)-1]
	l := len(r)
	max := MaxButtonsOnLine
	for i := 0; i < l; i++ {
		switch r[i].Action.Type {
		case ButtonTypeVkPay:
			max = 1
		case ButtonTypeLocation, ButtonTypeOpenLink, ButtonTypeOpenApp:
			if max > 2 {
				max = 2
			}
		}
	}
	if l < max {
		k.Buttons[len(k.Buttons)-1] = append(r, b)
	}
}

// AddText adds a text button to the last row.
func (k *Keyboard) AddText(payload string, label string, color string) {
	k.add(Button{
		Color: color,
		Action: Action{
			Type:    ButtonTypeText,
			Payload: payload,
			Label:   label,
		},
	})
}

// AddLocation adds a location button to the last row.
func (k *Keyboard) AddLocation(payload string) {
	k.add(Button{
		Action: Action{
			Type:    ButtonTypeLocation,
			Payload: payload,
		},
	})
}

// AddVkPay adds a VKPay button to the last row.
func (k *Keyboard) AddVkPay(payload string, hash string) {
	k.add(Button{
		Action: Action{
			Type:    ButtonTypeVkPay,
			Payload: payload,
			Hash:    hash,
		},
	})
}

// AddOpenApp adds a button with link to the vkapp to the last row.
func (k *Keyboard) AddOpenApp(payload string, appID, ownerID int, hash string) {
	k.add(Button{
		Action: Action{
			Type:    ButtonTypeOpenApp,
			Payload: payload,
			AppID:   appID,
			OwnerID: ownerID,
			Hash:    hash,
		},
	})
}

// AddOpenLink adds a button with external link to the last row.
func (k *Keyboard) AddOpenLink(payload string, link string) {
	k.add(Button{
		Action: Action{
			Type:    ButtonTypeOpenLink,
			Payload: payload,
			Link:    link,
		},
	})
}

// AddCallback adds a callback text button to the last row.
func (k *Keyboard) AddCallback(payload string, label string, color string) {
	k.add(Button{
		Color: color,
		Action: Action{
			Type:    ButtonTypeCallback,
			Payload: payload,
			Label:   label,
		},
	})
}

// Button struct.
type Button struct {
	Color  string `json:"color,omitempty"`
	Action Action `json:"action"`
}

// Action struct.
type Action struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
	Label   string `json:"label,omitempty"`
	Link    string `json:"link,omitempty"`
	Hash    string `json:"hash,omitempty"`
	AppID   int    `json:"app_id,omitempty"`
	OwnerID int    `json:"owner_id,omitempty"`
}
