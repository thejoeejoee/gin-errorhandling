package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/JosephWoodward/gin-errorhandling/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	NotFoundError = fmt.Errorf("this is an error")
)

type ValidationError struct {
	customError string
}

func (e *ValidationError) Error() string {
	return "Invalid request"
}

func TestMapSimpleErrorToStatusCode(t *testing.T) {
	// Arrange
	router := gin.Default()
	router.Use(
		ErrorHandler(
			Map(NotFoundError).ToStatusCode(http.StatusNotFound),
		))

	// Act
	router.GET("/", func(c *gin.Context) {
		_ = c.Error(NotFoundError)
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest("GET", "/", nil))

	// Assert
	assert.Equal(t, recorder.Result().StatusCode, http.StatusNotFound)
}

func TestMapWrappedErrorToStatusCode(t *testing.T) {
	wrapped := fmt.Errorf("wrapped error: %w", NotFoundError)

	// Arrange
	router := gin.Default()
	router.Use(
		ErrorHandler(
			Map(NotFoundError).ToStatusCode(http.StatusNotFound),
		))

	// Act
	router.GET("/", func(c *gin.Context) {
		_ = c.Error(wrapped)
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest("GET", "/", nil))

	// Assert
	assert.Equal(t, recorder.Result().StatusCode, http.StatusNotFound)
}

func TestMapErrorStructToStatusCode(t *testing.T) {
	// Arrange
	router := gin.Default()
	router.Use(
		ErrorHandler(
			Map(&ValidationError{}).ToStatusCode(http.StatusBadRequest),
		))

	// Act
	router.GET("/", func(c *gin.Context) {
		_ = c.Error(&ValidationError{})
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest("GET", "/", nil))

	// Assert
	assert.Equal(t, recorder.Result().StatusCode, http.StatusBadRequest)
}

func TestMapErrorResponseFunc(t *testing.T) {
	// Arrange
	router := gin.Default()
	router.Use(
		ErrorHandler(
			Map(NotFoundError).ToResponse(func(c *gin.Context, err error) {
				c.Status(http.StatusNotFound)
				c.Writer.Write([]byte(err.Error()))
			}),
		))

	// Act
	router.GET("/", func(c *gin.Context) {
		_ = c.Error(NotFoundError)
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest("GET", "/", nil))

	// Assert
	assert.Equal(t, http.StatusNotFound, recorder.Result().StatusCode)
	assert.Equal(t, NotFoundError.Error(), recorder.Body.String())
}
