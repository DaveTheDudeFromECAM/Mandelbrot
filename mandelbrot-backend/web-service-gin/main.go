package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"time"
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
	pixelgl.Run(run)
	/*
	   router := gin.Default()

	   router.GET("/mandelbrot", getMandelbrot)
	   // router.GET("/albums/:id", getAlbumByID)
	   // router.POST("/albums", postAlbums)

	   router.Run("localhost:8001")
	*/
}

// // getAlbums responds with the list of all albums as JSON.
//
//	func getAlbums(c *gin.Context) {
//		c.IndentedJSON(http.StatusOK, albums)
//	}
func getMandelbrot(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, mandelbrot)

}

// // postAlbums adds an album from JSON received in the request body.
// func postAlbums(c *gin.Context) {
// 	var newAlbum album

// 	// Call BindJSON to bind the received JSON to
// 	// newAlbum.
// 	if err := c.BindJSON(&newAlbum); err != nil {
// 		return
// 	}

// 	// Add the new album to the slice.
// 	albums = append(albums, newAlbum)
// 	c.IndentedJSON(http.StatusCreated, newAlbum)
// }

// // getAlbumByID locates the album whose ID value matches the id
// // parameter sent by the client, then returns that album as a response.
// func getAlbumByID(c *gin.Context) {
// 	id := c.Param("id")

// 	// Loop over the list of albums, looking for
// 	// an album whose ID value matches the parameter.
// 	for _, a := range albums {
// 		if a.ID == id {
// 			c.IndentedJSON(http.StatusOK, a)
// 			return
// 		}
// 	}
// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
// }

func run() {
	log.Println("Initial processing...")
	pixelCount = 0
	img = image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	cfg := pixelgl.WindowConfig{
		Title:  "Parallel Mandelbrot in Go",
		Bounds: pixel.R(0, 0, imgWidth, imgHeight),
		VSync:  true,
		// Invisible: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	log.Println("Rendering...")
	start := time.Now()
	workBuffer := make(chan WorkItem, numBlocks)
	threadBuffer := make(chan bool, numThreads)
	drawBuffer := make(chan Pix, pixelTotal)

	workBufferInit(workBuffer)
	go workersInit(drawBuffer, workBuffer, threadBuffer)
	go drawThread(drawBuffer, win)

	for !win.Closed() {
		pic := pixel.PictureDataFromImage(img)
		sprite := pixel.NewSprite(pic, pic.Bounds())
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()

		if showProgress {
			fmt.Printf("\r%d/%d (%d%%)", pixelCount, pixelTotal, int(100*(float64(pixelCount)/float64(pixelTotal))))
		}

		if pixelCount == pixelTotal {
			end := time.Now()
			fmt.Println("\nFinished with time = ", end.Sub(start))
			pixelCount++

			if closeOnEnd {
				break
			}
		}
	}
}

func workBufferInit(workBuffer chan WorkItem) {
	var sqrt = int(math.Sqrt(numBlocks))

	for i := sqrt - 1; i >= 0; i-- {
		for j := 0; j < sqrt; j++ {
			workBuffer <- WorkItem{
				initialX: i * (imgWidth / sqrt),
				finalX:   (i + 1) * (imgWidth / sqrt),
				initialY: j * (imgHeight / sqrt),
				finalY:   (j + 1) * (imgHeight / sqrt),
			}
		}
	}
}

func workersInit(drawBuffer chan Pix, workBuffer chan WorkItem, threadBuffer chan bool) {
	for i := 1; i <= numThreads; i++ {
		threadBuffer <- true
	}

	for range threadBuffer {
		workItem := <-workBuffer

		go workerThread(workItem, drawBuffer, threadBuffer)
	}
}

func workerThread(workItem WorkItem, drawBuffer chan Pix, threadBuffer chan bool) {
	for x := workItem.initialX; x < workItem.finalX; x++ {
		for y := workItem.initialY; y < workItem.finalY; y++ {
			var colorR, colorG, colorB int
			for k := 0; k < samples; k++ {
				a := height*ratio*((float64(x)+RandFloat64())/float64(imgWidth)) + posX
				b := height*((float64(y)+RandFloat64())/float64(imgHeight)) + posY
				c := pixelColor(mandelbrotIteraction(a, b, mandelbrot.MaxIter /*maxIter*/))
				colorR += int(c.R)
				colorG += int(c.G)
				colorB += int(c.B)
			}
			var cr, cg, cb uint8
			cr = uint8(float64(colorR) / float64(samples))
			cg = uint8(float64(colorG) / float64(samples))
			cb = uint8(float64(colorB) / float64(samples))

			drawBuffer <- Pix{
				x, y, cr, cg, cb,
			}

		}
	}
	threadBuffer <- true
}

func drawThread(drawBuffer chan Pix, win *pixelgl.Window) {
	for i := range drawBuffer {
		img.SetRGBA(i.x, i.y, color.RGBA{R: i.cr, G: i.cg, B: i.cb, A: 255})
		pixelCount++
	}
}

func mandelbrotIteraction(a, b float64, maxIter int) (float64, int) {
	var x, y, xx, yy, xy float64

	for i := 0; i < maxIter; i++ {
		xx, yy, xy = x*x, y*y, x*y
		if xx+yy > 4 {
			return xx + yy, i
		}
		// xn+1 = x^2 - y^2 + a
		x = xx - yy + a
		// yn+1 = 2xy + b
		y = 2*xy + b
	}

	return xx + yy, maxIter
}

func pixelColor(r float64, iter int) color.RGBA {
	insideSet := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	// validar se está dentro do conjunto
	// https://pt.wikipedia.org/wiki/Conjunto_de_Mandelbrot
	if r > 4 {
		// return hslToRGB(float64(0.70)-float64(iter)/3500*r, 1, 0.5)
		return hslToRGB(float64(iter)/100*r, 1, 0.5)
	}

	return insideSet
}

// -------------------------------

// xorshift random

var randState = uint64(time.Now().UnixNano())

func RandUint64() uint64 {
	randState = ((randState ^ (randState << 13)) ^ (randState >> 7)) ^ (randState << 17)
	return randState
}

func RandFloat64() float64 {
	return float64(RandUint64()/2) / (1 << 63)
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	switch {
	case t < 1.0/6.0:
		return p + (q-p)*6*t
	case t < 1.0/2.0:
		return q
	case t < 2.0/3.0:
		return p + (q-p)*(2.0/3.0-t)*6
	default:
		return p
	}
}

func hslToRGB(h, s, l float64) color.RGBA {
	var r, g, b float64
	if s == 0 {
		r, g, b = l, l, l
	} else {
		var q, p float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p = 2*l - q
		r = hueToRGB(p, q, h+1.0/3.0)
		g = hueToRGB(p, q, h)
		b = hueToRGB(p, q, h-1.0/3.0)
	}
	return color.RGBA{R: uint8(r * 255), G: uint8(g * 255), B: uint8(b * 255), A: 255}
}
