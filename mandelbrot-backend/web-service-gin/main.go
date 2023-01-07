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

	maxIter = 50
	samples = 100

	numBlocks  = 64
	numThreads = 32

	ratio = float64(imgWidth) / float64(imgHeight)

	showProgress = true
	closeOnEnd   = false
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
		width, height          = 1024, 1024
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
