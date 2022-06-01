package main

import (
	"fmt"
	"net/http"

	. "github.com/JosephWoodward/gin-errorhandling/middleware"
	"github.com/gin-gonic/gin"
)

var (
	NotFoundError = fmt.Errorf("resource could not be found")
)

func main() {
	r := gin.Default()
	r.Use(
		ErrorHandler(
			Map(NotFoundError).ToResponse(func(c *gin.Context, err error) {
				c.Status(http.StatusNotFound)
				c.Writer.Write([]byte(err.Error()))
			}),
		))

	r.GET("/ping", func(c *gin.Context) {
		c.Error(NotFoundError)
	})

	r.Run()
}
