package main

import (
	"image/color"
	"log"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
)

var (
	userText   = "Begin your masterpiece, or maybe just a shopping list:\n"
	counter    = 0
	tt         *truetype.Font
	normalFont font.Face
	dpi        = 72
)

type Game struct{}

// repeatingKeyPressed return true when key is pressed considering the repeat state.
func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {

	// outsideWidth and outsideHeight are set on line 85 by default
	// they change when the user resizes the window
	var newWidth int = int(float64(outsideWidth) * 0.6)
	var newHeight int = int(float64(outsideHeight) * 0.6)

	return newHeight, newWidth
}

func (g *Game) Update(screen *ebiten.Image) error {

	// Add a string from InputChars, that returns string input by users.
	// Note that InputChars result changes every frame, so you need to call this
	// every frame.
	userText += string(ebiten.InputChars())

	// Adjust the string to be at most 10 lines.
	ss := strings.Split(userText, "\n")
	if len(ss) > 10 {
		userText = strings.Join(ss[len(ss)-10:], "\n")
	}

	// If the enter key is pressed, add a line break.
	if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyKPEnter) {
		userText += "\n"
	}

	// If the backspace key is pressed, remove one character.
	if repeatingKeyPressed(ebiten.KeyBackspace) {
		if len(userText) >= 1 {
			userText = userText[:len(userText)-1]
		}
	}

	counter++

	if ebiten.IsDrawingSkipped() {
		// When the game is running slowly, the rendering result
		// will not be adopted.
		return nil
	}

	// append blinking cursor to text
	t := userText
	if counter%60 < 30 {
		t += "_"
	}

	screen.Fill(color.Black)

	slices := strings.Split(t, "\n")
	for i := 0; i < len(slices); i++ {
		text.Draw(screen, slices[i], normalFont, 15, 15+(20*i), color.White)
	}

	return nil
}

func init() {
	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	normalFont = truetype.NewFace(tt, &truetype.Options{
		Size:    12,
		DPI:     float64(dpi),
		Hinting: font.HintingFull,
	})
}

func main() {
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("carver")

	game := &Game{}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
