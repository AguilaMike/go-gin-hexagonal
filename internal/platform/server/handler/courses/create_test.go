package courses

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/AguilaMike/go-gin-hexagonal/internal/platform/storage/storagemocks"
)

func TestCreateHandler(t *testing.T) {
	setMockRepository := func(err error) *gin.Engine {
		courseRepository := new(storagemocks.CourseRepository)
		courseRepository.On("Save", mock.Anything, mock.AnythingOfType("mooc.Course")).Return(err)

		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.POST("/courses", CreateHandler(courseRepository))
		return r
	}

	type args struct {
		courseRepository createRequest
		mockReturnValue  error
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "given an invalid request it returns 400",
			args: args{
				courseRepository: createRequest{
					ID:   "8a1c5cdc-ba57-445a-994d-aa412d23723f",
					Name: "Demo Course",
				},
				mockReturnValue: nil,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "given a valid request it returns 201",
			args: args{
				courseRepository: createRequest{
					ID:       "8a1c5cdc-ba57-445a-994d-aa412d23723f",
					Name:     "Demo Course",
					Duration: "10 months",
				},
				mockReturnValue: nil,
			},
			want: http.StatusCreated,
		},
		{
			name: "given a invalid token it returns 400",
			args: args{
				courseRepository: createRequest{
					ID:       "NO-VALID-UUID",
					Name:     "Demo Course",
					Duration: "10 months",
				},
				mockReturnValue: nil,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "given a empty name it returns 400",
			args: args{
				courseRepository: createRequest{
					ID:       "8a1c5cdc-ba57-445a-994d-aa412d23723f",
					Name:     "",
					Duration: "10 months",
				},
				mockReturnValue: nil,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "given a empty duration it returns 400",
			args: args{
				courseRepository: createRequest{
					ID:       "8a1c5cdc-ba57-445a-994d-aa412d23723f",
					Name:     "Demo Course",
					Duration: "",
				},
				mockReturnValue: nil,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "given a invalid duration it returns 406",
			args: args{
				courseRepository: createRequest{
					ID:       "8a1c5cdc-ba57-445a-994d-aa412d23723f",
					Name:     "Demo Course",
					Duration: "Duration",
				},
				mockReturnValue: nil,
			},
			want: http.StatusNotAcceptable,
		},
		{
			name: "given a error to save it returns 500",
			args: args{
				courseRepository: createRequest{
					ID:       "8a1c5cdc-ba57-445a-994d-aa412d23723f",
					Name:     "Demo Course",
					Duration: "10 months",
				},
				mockReturnValue: errors.New("error to save"),
			},
			want: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setMockRepository(tt.args.mockReturnValue)

			b, err := json.Marshal(tt.args.courseRepository)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/courses", bytes.NewBuffer(b))
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want, res.StatusCode)
		})
	}
}
