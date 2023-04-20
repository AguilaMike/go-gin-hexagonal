package courses

import (
	"github.com/gin-gonic/gin"

	mooc "github.com/AguilaMike/go-gin-hexagonal/internal"
)

func Router(route *gin.Engine, r mooc.CourseRepository) {
	courses := route.Group("/courses")
	courses.POST("", CreateHandler(r))
	courses.GET("", GetAllHandler(r))
	courses.GET("/:id", GetByIDHandler(r))
}
