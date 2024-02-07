package http

import (
	"fmt"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) registerMiddlewares() {
	s.Rtr.Use(cORSMiddleware())
	s.Rtr.Use(requestIDMiddleware())
	s.Rtr.Use(gzip.Gzip(gzip.DefaultCompression))
}

// CORSMiddleware ...
func cORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

// RequestIDMiddleware ...
// Generate a unique ID and attach it to each request for future reference or use
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		suuid, _ := uuid.NewUUID() //uuid.NewV4()
		c.Writer.Header().Set("X-Request-Id", suuid.String())
		c.Next()
	}
}
