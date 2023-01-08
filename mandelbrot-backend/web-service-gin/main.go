package main

import (
	"fmt" //console log
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time" //Time taken to generate the set

	//"math/rand"
	"github.com/gin-gonic/gin"
)

/*
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
*/

// router
func main() {

	router := gin.Default()

	router.GET("/mandelbrot", getMandelbrot)
	// router.GET("/albums/:id", getAlbumByID)
	// router.POST("/albums", postAlbums)

	router.Run("localhost:8001")
}

func getMandelbrot(c *gin.Context) {

	// Get the parameters from the request.
	xmin := c.Query("xmin")
	xmax := c.Query("xmax")
	ymin := c.Query("ymin")
	ymax := c.Query("ymax")
	iterations := c.Query("iterations")
	width := c.Query("width")
	height := c.Query("height")

	// Convert the parameters to the appropriate types.
	widthInt, err := strconv.Atoi(width)
	if err != nil {
		// Handle the error.
	}
	heightInt, err := strconv.Atoi(height)
	if err != nil {
		// Handle the error.
	}
	xminFloat, err := strconv.ParseFloat(xmin, 64)
	if err != nil {
		// Handle the error.
	}
	xmaxFloat, err := strconv.ParseFloat(xmax, 64)
	if err != nil {
		// Handle the error.
	}
	yminFloat, err := strconv.ParseFloat(ymin, 64)
	if err != nil {
		// Handle the error.
	}
	ymaxFloat, err := strconv.ParseFloat(ymax, 64)
	if err != nil {
		// Handle the error.
	}
	iterationsInt, err := strconv.Atoi(iterations)
	if err != nil {
		// Handle the error.
	}

	// Use the parameters to generate the Mandelbrot set.
	//const width, height = 1000, 1000
	img := image.NewRGBA(image.Rect(0, 0, widthInt, heightInt))
	startTime := time.Now()
	var wg sync.WaitGroup
	for py := 0; py < heightInt; py++ {
		y := float64(py)/float64(heightInt)*(ymaxFloat-yminFloat) + yminFloat
		wg.Add(1)
		go func(py int, y float64) {
			defer wg.Done()
			for px := 0; px < widthInt; px++ {
				x := float64(px)/float64(widthInt)*(xmaxFloat-xminFloat) + xminFloat
				z := complex(x, y)
				// Image point (px, py) represents complex value z.
				img.Set(px, py, makemandelbrot(z, iterationsInt))
			}
		}(py, y)

	}
	wg.Wait()
	endTime := time.Now()
	f, err := os.Create("mandelbrot.png")
	if err != nil {
		panic(err)
	}
	png.Encode(f, img)
	elapsedTime := endTime.Sub(startTime)
	// Print the elapsed time to the console.
	fmt.Println("Mandelbrot set generated in:", elapsedTime)
	c.JSON(http.StatusOK, ImageResponse{
		ImagePath: "mandelbrot.png",
		Duration:  elapsedTime,
	})
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

func makemandelbrot(z complex128, iterations int) color.Color {
	//const iterations = 15
	const contrast = 255

	var v complex128
	for n := int(0); n < iterations; n++ {
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

// ImageResponse represents the response to the /image endpoint.
type ImageResponse struct {
	ImagePath string        `json:"imagePath"`
	Duration  time.Duration `json:"duration"`
}
