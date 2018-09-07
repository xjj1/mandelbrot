// mandelbrot project main.go
package main

import (
	"fmt"
	"log"

	//	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 640
	maxCol       = 256
)

var offscreen *ebiten.Image
var offscreenPix []byte

var m *Mandel

var leftButton bool
var rightButton bool

type Color struct {
	Red, Green, Blue, Alfa byte
}

var palette []Color

func InitPalette() {
	palette = make([]Color, 512)
	for i := 0; i < 256; i++ {
		palette[i].Red = byte(i)
		palette[i].Blue = 255 - byte(i)
	}
	for i := 0; i < 256; i++ {
		palette[i+256].Red = 255 - byte(i)
		palette[i+256].Blue = byte(i)
	}
}

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.DrawImage(offscreen, nil)
	x, y := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !leftButton {
			m.Update(x, y, 0.7)
			m.DrawMandel()
			screen.DrawImage(offscreen, nil)
			leftButton = true
		}
	} else {
		leftButton = false
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		if !rightButton {
			m.Update(x, y, 1.3)
			m.DrawMandel()
			screen.DrawImage(offscreen, nil)
			rightButton = true
		}
	} else {
		rightButton = false
	}
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("%d,%d", x, y))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%f,%f,%f,%f", m.centerX, m.centerY, m.size, 0.05/m.size))
	return nil
}

func main() {
	offscreenPix = make([]byte, screenWidth*screenHeight*4)
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Mandelbrot (Ebiten Demo)"); err != nil {
		log.Fatal(err)
	}
}

type Mandel struct {
	centerX, centerY, size float64
}

func (m *Mandel) Update(x, y int, scale float64) {
	newCentX := float64(x)*m.size/screenWidth - m.size/2 + m.centerX
	newCentY := (screenHeight-float64(y))*m.size/screenHeight - m.size/2 + m.centerY
	m.size *= scale
	m.centerX, m.centerY = newCentX, newCentY
	fmt.Println(newCentX, newCentY, m.size)
}

func (m *Mandel) Set(centerX, centerY, size float64) {
	m.centerX, m.centerY, m.size = centerX, centerY, size
}

func (m *Mandel) DrawMandel() {
	maxIter := 512
	for j := 0; j < screenHeight; j++ {
		for i := 0; i < screenHeight; i++ {
			x := float64(i)*m.size/screenWidth - m.size/2 + m.centerX
			y := (screenHeight-float64(j))*m.size/screenHeight - m.size/2 + m.centerY
			c := complex(x, y)
			z := complex(0, 0)
			it := 0
			for ; it < maxIter; it++ {
				z = z*z + c
				if real(z)*real(z)+imag(z)*imag(z) > 4 {
					break
				}
			}
			r, g, b := color(it, maxIter)
			p := 4 * (i + j*screenWidth)
			offscreenPix[p] = r
			offscreenPix[p+1] = g
			offscreenPix[p+2] = b
			offscreenPix[p+3] = 0xff
		}
	}
	offscreen.ReplacePixels(offscreenPix)
}

func init() {
	offscreen, _ = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterDefault)
	offscreenPix = make([]byte, screenWidth*screenHeight*4)
	//	for i := range palette {
	//		palette[i] = byte(math.Sqrt(float64(i)/float64(len(palette))) * 0x80)
	//	}
	InitPalette()
	leftButton, rightButton = false, false
}

func init() {
	// Now it is not feasible to call updateOffscreen every frame due to performance.
	m = new(Mandel)
	m.Set(0.0, 1.0, 0.5)
	m.DrawMandel()
}

func color(it, max int) (r, g, b byte) {
	if it == max {
		return 0xff, 0x22, 0xff
	}
	c := palette[it%maxCol]
	return c.Red, c.Green, c.Blue
}
