package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// router
func main() {
	router := gin.Default()
	router.GET("/mandelbrot", getMandelbrot)
	router.Run("localhost:8001")
}

func getMandelbrot(c *gin.Context) {

	// parameters from the request
	iterations := c.Query("iterations")
	height := c.Query("height")
	width := c.Query("width")
	xmin := c.Query("xmin")
	xmax := c.Query("xmax")
	ymin := c.Query("ymin")
	ymax := c.Query("ymax")

	// type conversion
	widthInt, err := strconv.Atoi(width)
	heightInt, err := strconv.Atoi(height)
	xminFloat, err := strconv.ParseFloat(xmin, 64)
	xmaxFloat, err := strconv.ParseFloat(xmax, 64)
	yminFloat, err := strconv.ParseFloat(ymin, 64)
	ymaxFloat, err := strconv.ParseFloat(ymax, 64)
	iterationsInt, err := strconv.Atoi(iterations)

	//
	img := image.NewRGBA(image.Rect(0, 0, widthInt, heightInt))
	startTime := time.Now()

	// counter for number of goroutines/tasks
	var wg sync.WaitGroup

	for py := 0; py < heightInt; py++ {
		y := float64(py)/float64(heightInt)*(ymaxFloat-yminFloat) + yminFloat
		wg.Add(1)

		// goroutine for each row
		go func(py int, y float64) {
			defer wg.Done()

			// each worker computes his row,pixel by pixel
			for px := 0; px < widthInt; px++ {
				x := float64(px)/float64(widthInt)*(xmaxFloat-xminFloat) + xminFloat
				z := complex(x, y)

				// sets color of a pixel according to computation
				img.Set(px, py, getColor(z, iterationsInt))
			}
		}(py, y)
	}

	// waint untill all routines are done
	wg.Wait()
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)

	// saves mandelbot set to a file
	f, err := os.Create("mandelbrot.png")
	if err != nil {
		panic(err)
	}
	png.Encode(f, img)

	// sends PNG path and computing duration to frontend
	c.JSON(http.StatusOK, ImageResponse{
		ImagePath: "mandelbrot.png",
		Duration:  elapsedTime,
	})
}

func getColor(z complex128, iterations int) color.Color {

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
