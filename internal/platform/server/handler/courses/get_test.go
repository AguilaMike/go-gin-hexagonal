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

	mooc "github.com/AguilaMike/go-gin-hexagonal/internal"
	"github.com/AguilaMike/go-gin-hexagonal/internal/platform/storage/storagemocks"
)

func TestGetAllHandler(t *testing.T) {
	setMockRepository := func(data []mooc.Course, err error) *gin.Engine {
		courseRepository := new(storagemocks.CourseRepository)
		courseRepository.On("GetAll", mock.Anything).Return(data, err)

		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.GET("/courses", GetAllHandler(courseRepository))
		return r
	}

	courseDB, _ := mooc.NewCourse("8a1c5cdc-ba57-445a-994d-aa412d23723f", "Demo Course", "10 months")

	type args struct {
		mockReturnValue []mooc.Course
		mockReturnError error
	}
	type want struct {
		status int
		data   []mooc.Course
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "given a valid request it returns 200 and zero courses",
			args: args{
				mockReturnValue: []mooc.Course{},
				mockReturnError: nil,
			},
			want: want{
				status: http.StatusOK,
				data:   []mooc.Course{},
			},
		},
		{
			name: "given a valid request it returns 200 with courses",
			args: args{
				mockReturnValue: []mooc.Course{
					courseDB,
				},
				mockReturnError: nil,
			},
			want: want{
				status: http.StatusOK,
				data: []mooc.Course{
					courseDB,
				},
			},
		},
		{
			name: "given a invalid request it returns 500",
			args: args{
				mockReturnValue: nil,
				mockReturnError: errors.New("error"),
			},
			want: want{
				status: http.StatusInternalServerError,
				data:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setMockRepository(tt.args.mockReturnValue, tt.args.mockReturnError)

			req, err := http.NewRequest(http.MethodGet, "/courses", nil)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.status, res.StatusCode)

			dataRes := []mooc.Course{}
			json.NewDecoder(res.Body).Decode(&dataRes)
			assert.Equal(t, len(tt.want.data), len(dataRes))
		})
	}
}

func TestGetByIDHandler(t *testing.T) {
	setMockRepository := func(data *mooc.Course, err error) *gin.Engine {
		courseRepository := new(storagemocks.CourseRepository)
		courseRepository.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(data, err)

		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.GET("/courses/:id", GetByIDHandler(courseRepository))
		return r
	}

	courseDB, _ := mooc.NewCourse("8a1c5cdc-ba57-445a-994d-aa412d23723f", "Demo Course", "10 months")

	type args struct {
		courseRepository getRequest
		mockReturnValue  *mooc.Course
		mockReturnError  error
	}
	type want struct {
		status int
		data   *mooc.Course
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "given a valid request it returns 200",
			args: args{
				courseRepository: getRequest{
					ID: "8a1c5cdc-ba57-445a-994d-aa412d23723f",
				},
				mockReturnValue: &courseDB,
				mockReturnError: nil,
			},
			want: want{
				status: http.StatusOK,
				data:   &courseDB,
			},
		},
		{
			name: "given a invalid request it returns 500",
			args: args{
				courseRepository: getRequest{
					ID: "8a1c5cdc-ba57-445a-994d-aa412d23723f",
				},
				mockReturnValue: nil,
				mockReturnError: errors.New("error"),
			},
			want: want{
				status: http.StatusInternalServerError,
				data:   nil,
			},
		},
		{
			name: "given a invalid request it returns 400",
			args: args{
				courseRepository: getRequest{
					ID: "token-invalid",
				},
				mockReturnValue: nil,
				mockReturnError: errors.New("error"),
			},
			want: want{
				status: http.StatusBadRequest,
				data:   nil,
			},
		},
		{
			name: "given a invalid request it returns 400",
			args: args{
				courseRepository: getRequest{
					ID: "",
				},
				mockReturnValue: nil,
				mockReturnError: errors.New("error"),
			},
			want: want{
				status: http.StatusBadRequest,
				data:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setMockRepository(tt.args.mockReturnValue, tt.args.mockReturnError)

			b, err := json.Marshal(tt.args.courseRepository)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodGet, "/courses/:id", bytes.NewBuffer(b))
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.status, res.StatusCode)

			// var dataRes mooc.Course
			// body, err := ioutil.ReadAll(res.Body)
			// require.NoError(t, err)
			// err = json.Unmarshal(body, &dataRes)
			// require.NoError(t, err)
			// assert.Equal(t, tt.want.data.ID().String(), dataRes.ID().String())
		})
	}
}
