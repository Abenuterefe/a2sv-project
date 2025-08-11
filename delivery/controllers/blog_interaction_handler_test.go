package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ucMocks "github.com/Abenuterefe/a2sv-project/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLikeBlog_Unauthorized(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	uc := ucMocks.NewBlogInteractionUseCaseInterface(t)
	h := NewBlogInteractionHandler(uc)

	r := gin.New()
	r.POST("/blogs/:id/like", h.LikeBlog)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/blogs/abc123/like", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLikeBlog_Happy(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	uc := ucMocks.NewBlogInteractionUseCaseInterface(t)
	h := NewBlogInteractionHandler(uc)

	uc.On("LikeBlog", mock.Anything, "abc123", "user-1").Return(nil)

	r := gin.New()
	// simple middleware to inject userID
	r.Use(func(c *gin.Context) { c.Set("userID", "user-1") })
	r.POST("/blogs/:id/like", h.LikeBlog)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/blogs/abc123/like", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestViewBlog_Anonymous_OK(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	uc := ucMocks.NewBlogInteractionUseCaseInterface(t)
	h := NewBlogInteractionHandler(uc)

	uc.On("ViewBlog", mock.Anything, "abc123", "", mock.Anything, mock.Anything).Return(nil)

	r := gin.New()
	r.POST("/blogs/:id/view", h.ViewBlog)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/blogs/abc123/view", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
