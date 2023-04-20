package courses

import (
	"net/http"

	"github.com/gin-gonic/gin"

	mooc "github.com/AguilaMike/go-gin-hexagonal/internal"
)

type getRequest struct {
	ID string `json:"id" binding:"required"`
}

// GetAllHandler returns an HTTP handler for all courses.
func GetAllHandler(courseRepository mooc.CourseRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		courses, err := courseRepository.GetAll(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, courses)
	}
}

// GetByIDHandler returns an HTTP handler for a course by id.
func GetByIDHandler(courseRepository mooc.CourseRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req getRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		id, err := mooc.NewCourseID(req.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		course, err := courseRepository.GetByID(ctx, id.String())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, course)
	}
}
