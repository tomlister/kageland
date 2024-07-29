package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"
	"syscall/js"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	resources "github.com/hajimehoshi/ebiten/v2/examples/resources/images/shader"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	shader        *ebiten.Shader
	setImagesKeys [4]string
	setImages     [4]*ebiten.Image
	images        map[string]*ebiten.Image
	game          *Game
)

func init() {
	images = map[string]*ebiten.Image{}
	setImages = [4]*ebiten.Image{}
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(resources.Gopher_png))
	if err != nil {
		log.Fatal(err)
	}
	images["gopher"] = ebiten.NewImageFromImage(img)
}

func init() {
	img, _, err := image.Decode(bytes.NewReader(resources.GopherBg_png))
	if err != nil {
		log.Fatal(err)
	}
	images["gopher_background"] = ebiten.NewImageFromImage(img)
}

func init() {
	img, _, err := image.Decode(bytes.NewReader(resources.Normal_png))
	if err != nil {
		log.Fatal(err)
	}
	images["gopher_normal"] = ebiten.NewImageFromImage(img)
}

func init() {
	img, _, err := image.Decode(bytes.NewReader(resources.Noise_png))
	if err != nil {
		log.Fatal(err)
	}
	images["noise"] = ebiten.NewImageFromImage(img)
}

type Game struct {
	init   bool
	idx    int
	time   int
	lw, lh int
}

func (g *Game) Update() error {
	g.time++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !g.init {
		g.init = true
		images["screen"] = ebiten.NewImageFromImage(screen)
	}

	if shader == nil {
		msg := "No shader program loaded."
		ebitenutil.DebugPrint(screen, msg)
		return
	}

	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()

	if g.lw != w || g.lh != h {
		// resize detected
		g.lw = w
		g.lh = h
		images["screen"].Deallocate()
		images["screen"] = ebiten.NewImageFromImage(screen)
		for i, _ := range setImages {
			setImages[i] = images[setImagesKeys[i]]
		}
	}
	cx, cy := ebiten.CursorPosition()

	op := &ebiten.DrawRectShaderOptions{}
	op.Uniforms = map[string]any{
		"Time":       float32(g.time) / 60,
		"Cursor":     []float32{float32(cx), float32(cy)},
		"Resolution": []float32{float32(screenWidth), float32(screenHeight)},
	}
	op.Images = setImages
	screen.DrawRectShader(w, h, shader, op)
	images["screen"].DrawImage(screen, &ebiten.DrawImageOptions{})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) setTime(v int) {
	g.time = v
}

func resetTime(this js.Value, args []js.Value) any {
	game.setTime(0)
	return nil
}

func compileShader(this js.Value, args []js.Value) any {
	if shader != nil {
		shader.Deallocate()
	}
	s, err := ebiten.NewShader([]byte(args[0].String()))
	if err != nil {
		return err.Error()
	}
	shader = s
	// sketchy
	setImagesKeys[0] = args[1].String()
	setImagesKeys[1] = args[2].String()
	setImagesKeys[2] = args[3].String()
	setImagesKeys[3] = args[4].String()

	setImages[0] = images[args[1].String()]
	setImages[1] = images[args[2].String()]
	setImages[2] = images[args[3].String()]
	setImages[3] = images[args[4].String()]
	return nil
}

func main() {
	js.Global().Set("compileShader", js.FuncOf(compileShader))
	js.Global().Set("resetTime", js.FuncOf(resetTime))
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Shader (Ebitengine Demo)")

	game = &Game{}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
