// mandelbrot project main.go
package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	//	"os"

	//	"runtime/trace"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 480
	screenHeight = 480
	maxCol       = 512

//	log_escape   = 0.6931471805599453
)

var offscreen *ebiten.Image

var m *Mandel

var leftButton bool
var rightButton bool

type Color struct {
	Red, Green, Blue byte
}

type Palette []Color

//var palette Palette

//func InitPalette() {
//	palette = make([]Color, 512)
//	for i := 0; i < 128; i++ {
//		palette[i].Red = 255 - 2*byte(i)
//		palette[i].Blue = 2 * byte(i)
//		palette[i+128].Red = 255 - 2*byte(i)
//		palette[i+128].Blue = 2 * byte(i)
//		palette[i+256].Blue = 255 - 2*byte(i)
//		palette[i+256].Green = 2 * byte(i)
//		palette[i+384].Green = 255 - 2*byte(i)
//		palette[i+384].Red = 2 * byte(i)
//	}
//}

//func InitPalette() {
//	palette = make([]Color, 1024)
//	var step = float32(1.0 / 1024)
//	ind := 0
//	for i := float32(0); i < 1; i += step {
//		c := GetColor(i)
//		fmt.Println(c)
//		palette[ind] = c
//		ind++
//	}
//}

//func RotatePalette() {
//	//fmt.Println(palette)
//	temp := make([]Color, 256)
//	copy(temp, palette[1:])
//	temp[255] = palette[0]
//	copy(palette, temp)
//	//fmt.Println(palette)
//}

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

	if string(ebiten.InputChars()) == "r" {
		m.Set(0.0, 1.0, 0.5)
		m.DrawMandel()
	}
	newCentX := float64(x)*m.size/screenWidth - m.size/2 + m.centerX
	newCentY := (screenHeight-float64(y))*m.size/screenHeight - m.size/2 + m.centerY

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%f,%f,%f", newCentX, newCentY, m.size))
	return nil
}

func init() {
	offscreen, _ = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterDefault)
	leftButton, rightButton = false, false

	//fmt.Println(palette)
	s.Init()
	m = new(Mandel)
	m.InitPalette()
	m.Set(0.0, 1.0, 0.5)
	m.DrawMandel()
}

func (m *Mandel) color(it, max int) Color {
	if it == max {
		return Color{0, 0, 0}
	}
	c := m.palette[it%maxCol]
	return c
}

func main() {

	//	f, err := os.Create("cpu.trace")
	//	if err != nil {
	//		fmt.Println("cannot profile", err)
	//		return
	//	}
	//	defer f.Close()
	//	trace.Start(f)
	//	defer trace.Stop()
	//	pprof.StartCPUProfile(f)
	//	defer func() {
	//		pprof.StopCPUProfile()

	//}()
	//s.Init()

	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Mandelbrot (Ebiten Demo)"); err != nil {
		log.Fatal(err)
	}
}

// var offscreenPix []byte

// Screen is a drawing buffer used to feed to ReplacePixels.
var s Sreen

type Sreen struct {
	pix []byte
	m   sync.Mutex
}

func (s *Sreen) Init() {
	s.pix = make([]byte, screenWidth*screenHeight*4)
}

func (s *Sreen) Set(x, y int, c Color) {
	p := 4 * (x + y*screenWidth)
	s.m.Lock()
	s.pix[p] = c.Red
	s.pix[p+1] = c.Green
	s.pix[p+2] = c.Blue
	s.pix[p+3] = 0xff
	s.m.Unlock()
}

type Mandel struct {
	centerX, centerY, size float64
	palette                Palette
}

func (m *Mandel) InitPalette() {
	m.palette = make([]Color, maxCol)
	var step = float32(1.0 / maxCol)
	ind := 0
	for i := float32(0); i < 1; i += step {
		c := GetColor(i)
		m.palette[ind] = c
		ind++
	}
}

func GetColor(pos float32) Color {
	var r, g, b byte
	if pos > 1.0 {
		if pos-float32(int(pos)) == 0.0 {
			pos = 1.0
		} else {
			pos = pos - float32(int(pos))
		}
	}
	nmax := 6
	m := float32(nmax) * pos
	n := int(m)
	f := m - float32(n)
	t := byte(255 * f)

	switch n {
	case 0:
		r, g, b = 255, t, 0
	case 1:
		r, g, b = 255-t, 255, 0
	case 2:
		r, g, b = 0, 255, t
	case 3:
		r, g, b = 0, 255-t, 255
	case 4:
		r, g, b = t, 0, 255
	case 5:
		r, g, b = 255, 0, 255-t
	default:
		r, g, b = 255, 0, 0
	}

	return Color{r, g, b}

}

func (m *Mandel) Update(x, y int, scale float64) {
	newCentX := float64(x)*m.size/screenWidth - m.size/2 + m.centerX
	newCentY := (screenHeight-float64(y))*m.size/screenHeight - m.size/2 + m.centerY
	m.size *= scale
	m.centerX, m.centerY = newCentX, newCentY
	//fmt.Println(newCentX, newCentY, m.size)
}

func (m *Mandel) Set(centerX, centerY, size float64) {
	m.centerX, m.centerY, m.size = centerX, centerY, size
}

func (m *Mandel) DrawMandel() {
	start := time.Now()
	var wg sync.WaitGroup
	maxIter := 512
	for j := 0; j < screenHeight; j++ {
		wg.Add(1)
		go func(j int) {
			for i := 0; i < screenWidth; i++ {
				x := float64(i)*m.size/screenWidth - m.size/2 + m.centerX
				y := (screenHeight-float64(j))*m.size/screenHeight - m.size/2 + m.centerY
				c := complex(x, y)
				z := complex(0, 0)
				it := 0
				for ; it < maxIter; it++ {
					z = z*z + c
					if real(z)*real(z)+imag(z)*imag(z) > 2 {
						break
					}
				}
				col := m.color(it, maxIter)
				s.Set(i, j, col)
			}
			wg.Done()
		}(j)
	}
	wg.Wait()
	fmt.Println("Draw mandelbrot: ", time.Since(start))
	start = time.Now()
	offscreen.ReplacePixels(s.pix)
	fmt.Println("Replace Pixels: ", time.Since(start))
}

