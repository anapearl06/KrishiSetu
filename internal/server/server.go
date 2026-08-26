package server

import ("github.com/gin-gonic/gin"
	"github.com/raaj2493/KrishiSetu/internal/server/response"
)

func New() *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		response.Success(c, 200, gin.H{
			"status": "ok",
		})
	})

	return router
}