package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	ucMocks "github.com/Abenuterefe/a2sv-project/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCommentsByBlog_BadID(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	uc := ucMocks.NewCommentUseCaseInterface(t)
	h := NewCommentHandler(uc)

	r := gin.New()
	r.GET("/blogs/:id/comments", h.GetCommentsByBlog)

	w := httptest.NewRecorder()
	// We'll use a normal ID and UC returns error 500 to exercise error path
	uc.On("GetCommentsByBlogID", mock.Anything, "blog-1").Return(nil, assert.AnError)
	req := httptest.NewRequest(http.MethodGet, "/blogs/blog-1/comments", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateComment_Unauthorized(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	uc := ucMocks.NewCommentUseCaseInterface(t)
	h := NewCommentHandler(uc)

	r := gin.New()
	r.POST("/blogs/:id/comments", h.CreateComment)

	body, _ := json.Marshal(&entities.Comment{Content: "hi"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/blogs/blog-1/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateComment_Happy(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	uc := ucMocks.NewCommentUseCaseInterface(t)
	h := NewCommentHandler(uc)

	uc.On("CreateComment", mock.Anything, mock.AnythingOfType("*entities.Comment"), "user-1", "blog-1").Return(nil)

	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("userID", "user-1") })
	r.POST("/blogs/:id/comments", h.CreateComment)

	body, _ := json.Marshal(&entities.Comment{Content: "nice"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/blogs/blog-1/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestUpdateComment_ForbiddenWhenNotOwner(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	uc := ucMocks.NewCommentUseCaseInterface(t)
	h := NewCommentHandler(uc)

	// existing comment owned by someone else
	uc.On("GetCommentByID", mock.Anything, "c-1").Return(&entities.Comment{UserID: "owner-2"}, nil)

	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("userID", "owner-1") })
	r.PUT("/comments/:id", h.UpdateComment)

	body, _ := json.Marshal(&entities.Comment{Content: "edit"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/comments/c-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteComment_ForbiddenWhenNotOwner(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	uc := ucMocks.NewCommentUseCaseInterface(t)
	h := NewCommentHandler(uc)

	// existing comment owned by someone else
	uc.On("GetCommentByID", mock.Anything, "c-1").Return(&entities.Comment{UserID: "owner-2"}, nil)

	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("userID", "owner-1") })
	r.DELETE("/comments/:id", h.DeleteComment)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/comments/c-1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
