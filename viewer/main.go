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
	shader    *ebiten.Shader
	setImages [4]*ebiten.Image
	images    map[string]*ebiten.Image
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
	init bool
	idx  int
	time int
}

func (g *Game) Update() error {
	g.time++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !g.init {
		g.init = true
		images["screen"] = screen
	}
	if shader == nil {
		msg := "No shader program loaded."
		ebitenutil.DebugPrint(screen, msg)
		return
	}

	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	cx, cy := ebiten.CursorPosition()

	op := &ebiten.DrawRectShaderOptions{}
	op.Uniforms = map[string]any{
		"Time":       float32(g.time) / 60,
		"Cursor":     []float32{float32(cx), float32(cy)},
		"Resolution": []float32{float32(screenWidth), float32(screenHeight)},
	}
	op.Images = setImages
	screen.DrawRectShader(w, h, shader, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
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
	setImages[0] = images[args[1].String()]
	setImages[1] = images[args[2].String()]
	setImages[2] = images[args[3].String()]
	setImages[3] = images[args[4].String()]
	return nil
}

func main() {
	js.Global().Set("compileShader", js.FuncOf(compileShader))
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Shader (Ebitengine Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
