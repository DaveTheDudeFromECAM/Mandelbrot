package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"sync"
)

type Pix struct {
	x  int
	y  int
	cr uint8
	cg uint8
	cb uint8
}

type WorkItem struct {
	initialX int
	finalX   int
	initialY int
	finalY   int
}

// Mandelbrot class struct.
type Mandelbrot struct {
	ID           string  `json:"id"`
	PosX         float32 `json:"posX"`
	PosY         float32 `json:"posY"`
	Height       float32 `json:"height"`
	ImgWidth     float32 `json:"imgWidth"`
	ImgHeight    float32 `json:"imgHeight"`
	MaxIter      int     `json:"maxIter"`
	Samples      int     `json:"samples"`
	NumBlocks    int     `json:"numBlocks"`
	NumThreads   int     `json:"numThreads"`
	ShowProgress bool    `json:"showProgress"`
	CloseOnEnd   bool    `json:"closeOnEnd"`
}

const (
	posX   = -2
	posY   = -1.2
	height = 2.5

	imgWidth   = 800
	imgHeight  = 600
	pixelTotal = imgWidth * imgHeight
)

// object Mandelbrot
var mandelbrot = Mandelbrot{
	PosX: -2, PosY: -1.2, Height: 2.5, ImgWidth: 1024, ImgHeight: 1024, MaxIter: 1000, Samples: 200, NumBlocks: 64, NumThreads: 16, ShowProgress: true, CloseOnEnd: false,
}

var (
	img        *image.RGBA
	pixelCount int
)

// router
func main() {

	router := gin.Default()

	router.GET("/mandelbrot", getMandelbrot)
	// router.GET("/albums/:id", getAlbumByID)
	// router.POST("/albums", postAlbums)

	router.Run("localhost:8001")
}

func getMandelbrot(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, mandelbrot)
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 6000, 5000
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	var wg sync.WaitGroup
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		wg.Add(1)
		go func(py int, y float64) {
			defer wg.Done()
			for px := 0; px < width; px++ {
				x := float64(px)/width*(xmax-xmin) + xmin
				z := complex(x, y)
				// Image point (px, py) represents complex value z.
				img.Set(px, py, makemandelbrot(z))
			}
		}(py, y)
	}
	wg.Wait()
	f, err := os.Create("mandelbrot.png")
	if err != nil {
		panic(err)
	}
	png.Encode(f, img)
}

/*
//Geaysacel version
func makemandelbrot(z complex128) color.Color {
	const iterations = 30
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - contrast*n}
		}
	}
	return color.Black
}
*/

func makemandelbrot(z complex128) color.Color {
	const iterations = 50
	const contrast = 15

	var v complex128
	for n := uint16(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			// Use the number of iterations as an index into a color palette to
			// determine the color to return.
			return colorPalette[int(n)%len(colorPalette)]

		}
	}
	return color.Black
}

// colorPalette is a slice of colors to use as a color palette.
var colorPalette = []color.Color{
	color.RGBA{66, 30, 15, 255},    // Dark brown
	color.RGBA{25, 7, 26, 255},     // Dark purple
	color.RGBA{9, 1, 47, 255},      // Deep blue
	color.RGBA{4, 4, 73, 255},      // Dark blue
	color.RGBA{0, 7, 100, 255},     // Blue
	color.RGBA{12, 44, 138, 255},   // Light blue
	color.RGBA{24, 82, 177, 255},   // Sky blue
	color.RGBA{57, 125, 209, 255},  // Light sky blue
	color.RGBA{134, 181, 229, 255}, // Very light blue
}
