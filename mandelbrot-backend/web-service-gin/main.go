package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

// object Mandelbrot
var mandelbrot = Mandelbrot{
	PosX: -2, PosY: -1.2, Height: 2.5, ImgWidth: 1024, ImgHeight: 1024, MaxIter: 100, Samples: 200, NumBlocks: 64, NumThreads: 16, ShowProgress: true, CloseOnEnd: false,
}

// router
func main() {
	router := gin.Default()

	router.GET("/mandelbrot", getMandelbrot)
	// router.GET("/albums/:id", getAlbumByID)
	// router.POST("/albums", postAlbums)

	router.Run("localhost:8000")
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
