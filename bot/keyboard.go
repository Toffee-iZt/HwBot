package bot

import "github.com/Toffee-iZt/HwBot/vkapi"

// Keyboard generator.
type Keyboard struct {
	kb  *vkapi.Keyboard
	mod *Module
}

// NextRow adds line of buttons.
func (k *Keyboard) NextRow() bool {
	return k.kb.NextRow()
}

// AddText adds a text button to the last row.
func (k *Keyboard) AddText(label string, color string, payload interface{}) bool {
	return k.kb.AddText(wrapPayload(k.mod, payload), label, color)
}

// AddOpenLink adds a button with external link to the last row.
func (k *Keyboard) AddOpenLink(label string, link string) bool {
	return k.kb.AddOpenLink(label, link)
}

// AddCallback adds a callback text button to the last row.
func (k *Keyboard) AddCallback(label string, color string, payload interface{}) bool {
	return k.kb.AddCallback(wrapPayload(k.mod, payload), label, color)
}

// Build builds keyboard object.
func (k *Keyboard) Build() vkapi.JSONData {
	return k.kb.Data()
}
