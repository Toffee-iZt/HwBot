package keyboard

import "github.com/Toffee-iZt/HwBot/vkapi"

//
const (
	MaxButtonsOnLine = 5
	MaxDefaultLines  = 10
	MaxInlineLines   = 6
)

// Color type.
type Color string

//
const (
	// Blue
	ColorPrimary Color = "primary"
	// White
	ColorSecondary Color = "secondary"
	// Red
	ColorNegative Color = "negative"
	// Green
	ColorPositive Color = "positive"
)

//
const (
	buttonTypeText     = "text"
	buttonTypeLocation = "location"
	buttonTypeVkPay    = "vkpay"
	buttonTypeOpenApp  = "open_app"
	buttonTypeOpenLink = "open_link"
	buttonTypeCallback = "callback"
)

// Empty creates empty keyboard to remove an existing.
func Empty() *Keyboard {
	kb := New(true, false)
	kb.Buttons = make([][]Button, 0)
	return kb
}

// New creates new keyboard.
func New(oneTime bool, inline bool) *Keyboard {
	if oneTime && inline {
		oneTime = false
	}
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

// Data converts the keyboard into a sandable object.
func (k *Keyboard) Data() vkapi.JSONData {
	d, _ := vkapi.NewJSONData(k)
	return d
}

// NextRow adds line of buttons.
func (k *Keyboard) NextRow() bool {
	if k.Inline && len(k.Buttons) >= MaxInlineLines {
		return false
	}
	if !k.Inline && len(k.Buttons) >= MaxDefaultLines {
		return false
	}
	k.Buttons = append(k.Buttons, make([]Button, 0, 1))
	return true
}

func (k *Keyboard) add(b Button) bool {
	r := append(k.Buttons[len(k.Buttons)-1], b)
	l := len(r)
	max := MaxButtonsOnLine
	for i := 0; i < l; i++ {
		switch r[i].Action.Type {
		case buttonTypeVkPay, buttonTypeLocation, buttonTypeOpenApp:
			max = 1
		case buttonTypeOpenLink:
			max = 2
		}
	}
	if l > max {
		return false
	}
	k.Buttons[len(k.Buttons)-1] = r
	return true
}

// AddText adds a text button to the last row.
func (k *Keyboard) AddText(label string, color Color, payload vkapi.JSONData) bool {
	return k.add(Button{
		Color: color,
		Action: Action{
			Type:    buttonTypeText,
			Payload: payload,
			Label:   label,
		},
	})
}

// AddLocation adds a location button to the last row.
func (k *Keyboard) AddLocation(payload vkapi.JSONData) bool {
	return k.add(Button{
		Action: Action{
			Type:    buttonTypeLocation,
			Payload: payload,
		},
	})
}

// AddVkPay adds a VKPay button to the last row.
func (k *Keyboard) AddVkPay(hash string) bool {
	return k.add(Button{
		Action: Action{
			Type: buttonTypeVkPay,
			Hash: hash,
		},
	})
}

// AddOpenApp adds a button with link to the vkapp to the last row.
func (k *Keyboard) AddOpenApp(appID, ownerID int, hash string) bool {
	return k.add(Button{
		Action: Action{
			Type:    buttonTypeOpenApp,
			AppID:   appID,
			OwnerID: ownerID,
			Hash:    hash,
		},
	})
}

// AddOpenLink adds a button with external link to the last row.
func (k *Keyboard) AddOpenLink(label string, link string) bool {
	return k.add(Button{
		Action: Action{
			Type:  buttonTypeOpenLink,
			Label: label,
			Link:  link,
		},
	})
}

// AddCallback adds a callback text button to the last row.
func (k *Keyboard) AddCallback(label string, color Color, payload vkapi.JSONData) bool {
	return k.add(Button{
		Color: color,
		Action: Action{
			Type:    buttonTypeCallback,
			Payload: payload,
			Label:   label,
		},
	})
}

// Button type.
type Button struct {
	Color  Color  `json:"color,omitempty"`
	Action Action `json:"action"`
}

// Action type.
type Action struct {
	Type    string         `json:"type"`
	Payload vkapi.JSONData `json:"payload,omitempty"`
	Label   string         `json:"label,omitempty"`
	Link    string         `json:"link,omitempty"`
	Hash    string         `json:"hash,omitempty"`
	AppID   int            `json:"app_id,omitempty"`
	OwnerID int            `json:"owner_id,omitempty"`
}
