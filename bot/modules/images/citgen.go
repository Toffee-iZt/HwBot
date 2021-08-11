package images

import (
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"strings"
	"time"

	"github.com/Toffee-iZt/HwBot/bot"
	"github.com/Toffee-iZt/HwBot/vkapi"
	"github.com/Toffee-iZt/HwBot/vkapi/vkutils"
	"github.com/Toffee-iZt/workfs"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

const copyrightSymbol = "©"

var citgen = bot.Command{
	Cmd:  "citgen",
	Desc: "Генерация цитаты",
	Help: "/citen и ответить или переслать сообщение",
	Priv: true,
	Chat: true,
	Run: func(b *bot.Bot, m *bot.IncomingMessage, _ []string) {
		var fromID int
		var text string
		switch {
		case m.Message.Reply != nil:
			r := m.Message.Reply
			fromID, text = r.FromID, r.Text
		case len(m.Message.Forward) > 0:
			f := m.Message.Forward[0]
			fromID, text = f.FromID, f.Text
		default:
			b.SimpleReply(m, "Ответьте или перешлите сообщение")
			return
		}
		if text == "" {
			b.SimpleReply(m, "Сообщение не содержит текста")
			return
		}

		api := b.API()

		name, photo, err := getNamePhoto(api, fromID)
		if err != nil {
			log.Error("citgen: %s", err.Error())
			return
		}

		f, err := generateQuote(photo, name, text, fromID == m.Message.FromID, color.Black, color.White)
		if err != nil {
			log.Error("citgen generate: %s", err.Error())
			return
		}
		defer f.Close()

		s, err := api.UploadMessagesPhoto(m.Message.PeerID, f)
		if err != nil {
			log.Error("citgen upload: %s", err.Error())
			return
		}

		_, vkerr := api.Messages.Send(vkapi.MessagePeer{PeerID: m.Message.PeerID}, vkapi.MessageContent{
			Attachment: []string{s},
		})
		if vkerr != nil {
			log.Error("citgen send: %s", vkerr.Error())
			return
		}
	},
}

func getNamePhoto(api *vkapi.Client, from int) (string, image.Image, error) {
	var name string
	var photoCrop string
	if from < 0 {
		g, err := api.Groups.GetByID([]int{vkutils.PeerToGroup(from)}, "photo_200")
		if err != nil || len(g) == 0 {
			return "", nil, err
		}
		name = g[0].Name
		photoCrop = g[0].Photo200
	} else {
		u, err := api.Users.Get([]int{vkutils.PeerToUser(from)}, "photo_200")
		if err != nil || len(u) == 0 {
			return "", nil, err
		}

		name = u[0].FirstName + " " + u[0].LastName
		photoCrop = *u[0].Photo200
	}

	//size := photoCrop.Photo.Sizes[len(photoCrop.Photo.Sizes)-1]
	//img, err := dl(size.URL)
	img, err := dl(photoCrop)
	if err != nil {
		return "", nil, err
	}

	//c := photoCrop.Crop
	//img = crop(img, image.Rect(
	//	int(float64(size.Width)*(c.X/100)),
	//	int(float64(size.Height)*(c.Y/100)),
	//	int(float64(size.Width)*(c.X2/100)),
	//	int(float64(size.Height)*(c.Y2/100)),
	//))

	return name, img, nil
}

const (
	fontSize = 20

	width     = 700
	minHeight = 400

	photoSize = 200

	textPointX = width / 3

	padding = 15
)

var fontFace = getGoFontFace()

func getGoFontFace() font.Face {
	gottf, _ := truetype.Parse(goregular.TTF)
	return truetype.NewFace(gottf, &truetype.Options{Size: fontSize})
}

func generateQuote(photo image.Image, name, quote string, self bool, bg, fg color.Color) (fs.File, error) {
	lines, height := generateLines(quote, textPointX)

	name += " " + copyrightSymbol
	if self {
		name += " (self)"
	}

	dc := gg.NewContext(width, height)

	dc.SetFontFace(fontFace)
	dc.SetColor(bg)
	dc.Clear()
	dc.SetColor(fg)

	// Draw quote text
	for i, line := range lines {
		y := height/2 - fontSize*len(lines)/2 + i*fontSize
		dc.DrawString(line, textPointX, float64(y))
	}

	// Draw name
	dc.DrawString(name, float64(padding), float64(height-padding))

	// Draw time
	t := time.Now().UTC().Format("02.01.2006 15:04")
	dc.DrawString(t, float64(width-getStringWidth(t)-padding), float64(height-padding))

	// Draw photo and make it round
	px, py := width/6, height/2
	dc.DrawEllipse(float64(px), float64(py), 100, 100)
	dc.Clip()
	dc.DrawImageAnchored( /*scaleQuad(photo, photoSize)*/ photo, px, py, 0.5, 0.5)

	f := workfs.OpenAtOnce("citgen.png")
	err := png.Encode(f, dc.Image())
	if err != nil {
		f.Close()
		return nil, err
	}
	return f, nil
}

func generateLines(s string, w int) ([]string, int) {
	var lines []string

	for _, line := range strings.Split(s, "\n") {
		var newLine string
		for _, word := range strings.Split(line, " ") {
			if getStringWidth(newLine+" "+word) > (width - w - 10) {
				lines = append(lines, newLine)
				newLine = word
			} else {
				newLine += " " + word
			}
		}

		if newLine != "" {
			lines = append(lines, strings.TrimSpace(newLine))
		}
	}

	lines[0] = "«" + lines[0]
	lines[len(lines)-1] += "»"

	h := len(lines) * (fontFace.Metrics().Height.Ceil() + 2)
	if h < minHeight {
		h = minHeight
	}

	return lines, h
}

func getStringWidth(s string) int {
	w := 0
	for _, r := range s {
		_, a, _ := fontFace.GlyphBounds(r)
		w += a.Round()
	}
	return w
}
