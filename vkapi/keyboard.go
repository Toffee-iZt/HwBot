package vkapi

//
const (
	KeyboardMaxButtonsOnLine = 5
	KeyboardMaxDefaultLines  = 10
	KeyboardMaxInlineLines   = 6
)

//
const (
	// VkBlue
	KeyboardColorPrimary = "primary"
	// White
	KeyboardColorSecondary = "secondary"
	// Red
	KeyboardColorNegative = "negative"
	// Green
	KeyboardColorPositive = "positive"
)

//
const (
	KeyboardButtonTypeText     = "text"
	KeyboardButtonTypeLocation = "location"
	KeyboardButtonTypeVkPay    = "vkpay"
	KeyboardButtonTypeOpenApp  = "open_app"
	KeyboardButtonTypeOpenLink = "open_link"
	KeyboardButtonTypeCallback = "callback"
)

// EmptyKeyboard creates empty keyboard to remove an existing.
func EmptyKeyboard() *Keyboard {
	return NewKeyboard(true, false)
}

// NewKeyboard creates new keyboard.
func NewKeyboard(oneTime bool, inline bool) *Keyboard {
	if oneTime && inline {
		return nil
	}
	return &Keyboard{
		OneTime: oneTime,
		Inline:  inline,
	}
}

// Keyboard struct.
type Keyboard struct {
	Buttons [][]KeyboardButton `json:"buttons"`
	OneTime bool               `json:"one_time"`
	Inline  bool               `json:"inline"`
}

func (k *Keyboard) String() string {
	if k == nil {
		return ""
	}
	return string(marshal(k))
}

// AddRow adds line of buttons.
func (k *Keyboard) AddRow() bool {
	if k.Inline && len(k.Buttons) >= KeyboardMaxInlineLines {
		return false
	}
	if !k.Inline && len(k.Buttons) >= KeyboardMaxDefaultLines {
		return false
	}
	k.Buttons = append(k.Buttons, make([]KeyboardButton, 0, 1))
	return true
}

func (k *Keyboard) add(b KeyboardButton) bool {
	r := append(k.Buttons[len(k.Buttons)-1], b)
	l := len(r)
	max := KeyboardMaxButtonsOnLine
	for i := 0; i < l; i++ {
		switch r[i].Action.Type {
		case KeyboardButtonTypeVkPay:
			max = 1
		case KeyboardButtonTypeLocation, KeyboardButtonTypeOpenLink, KeyboardButtonTypeOpenApp:
			if max > 2 {
				max = 2
			}
		}
	}
	if l > max {
		return false
	}
	k.Buttons[len(k.Buttons)-1] = append(r, b)
	return true
}

// AddText adds a text button to the last row.
func (k *Keyboard) AddText(payload string, label string, color string) bool {
	return k.add(KeyboardButton{
		Color: color,
		Action: KeyboardAction{
			Type:    KeyboardButtonTypeText,
			Payload: payload,
			Label:   label,
		},
	})
}

// AddLocation adds a location button to the last row.
func (k *Keyboard) AddLocation(payload string) bool {
	return k.add(KeyboardButton{
		Action: KeyboardAction{
			Type:    KeyboardButtonTypeLocation,
			Payload: payload,
		},
	})
}

// AddVkPay adds a VKPay button to the last row.
func (k *Keyboard) AddVkPay(payload string, hash string) bool {
	return k.add(KeyboardButton{
		Action: KeyboardAction{
			Type:    KeyboardButtonTypeVkPay,
			Payload: payload,
			Hash:    hash,
		},
	})
}

// AddOpenApp adds a button with link to the vkapp to the last row.
func (k *Keyboard) AddOpenApp(payload string, appID, ownerID int, hash string) bool {
	return k.add(KeyboardButton{
		Action: KeyboardAction{
			Type:    KeyboardButtonTypeOpenApp,
			Payload: payload,
			AppID:   appID,
			OwnerID: ownerID,
			Hash:    hash,
		},
	})
}

// AddOpenLink adds a button with external link to the last row.
func (k *Keyboard) AddOpenLink(payload string, label string, link string) bool {
	return k.add(KeyboardButton{
		Action: KeyboardAction{
			Type:    KeyboardButtonTypeOpenLink,
			Payload: payload,
			Label:   label,
			Link:    link,
		},
	})
}

// AddCallback adds a callback text button to the last row.
func (k *Keyboard) AddCallback(payload string, label string, color string) bool {
	return k.add(KeyboardButton{
		Color: color,
		Action: KeyboardAction{
			Type:    KeyboardButtonTypeCallback,
			Payload: payload,
			Label:   label,
		},
	})
}

// KeyboardButton struct.
type KeyboardButton struct {
	Color  string         `json:"color,omitempty"`
	Action KeyboardAction `json:"action"`
}

// KeyboardAction struct.
type KeyboardAction struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
	Label   string `json:"label,omitempty"`
	Link    string `json:"link,omitempty"`
	Hash    string `json:"hash,omitempty"`
	AppID   int    `json:"app_id,omitempty"`
	OwnerID int    `json:"owner_id,omitempty"`
}
