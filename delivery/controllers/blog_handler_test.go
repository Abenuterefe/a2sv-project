package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	ucMocks "github.com/Abenuterefe/a2sv-project/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFilterBlogs_InvalidDateParam(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	uc := ucMocks.NewBlogUseCaseInterface(t)
	h := NewBlogHandler(uc)

	// invalid date_from
	req := httptest.NewRequest(http.MethodGet, "/blogs/filter?date_from=2024-13-01", nil)
	w := httptest.NewRecorder()
	r := gin.New()
	r.GET("/blogs/filter", h.FilterBlogs)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSearchBlogs_MissingParams_400(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	uc := ucMocks.NewBlogUseCaseInterface(t)
	h := NewBlogHandler(uc)

	// When UC is called with empty params, handler should still call UC and get error.
	uc.On("SearchBlogs", mock.Anything, &entities.BlogSearch{}).Return((*entities.SearchResponse)(nil), assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/blogs/search", nil)
	w := httptest.NewRecorder()
	r := gin.New()
	r.GET("/blogs/search", h.SearchBlogs)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSearchBlogs_HappyPath_PageToSkip(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	uc := ucMocks.NewBlogUseCaseInterface(t)
	h := NewBlogHandler(uc)

	// page=3, limit=10 => skip should be 20
	uc.On("SearchBlogs", mock.Anything, mock.MatchedBy(func(s *entities.BlogSearch) bool {
		return s.Title == "go" && s.Limit == 10 && s.Skip == 20
	})).Return(&entities.SearchResponse{Blogs: []*entities.BlogWithAuthor{}, Count: 0, TotalCount: 0}, nil)

	req := httptest.NewRequest(http.MethodGet, "/blogs/search?title=go&page=3&limit=10", nil)
	w := httptest.NewRecorder()
	r := gin.New()
	r.GET("/blogs/search", h.SearchBlogs)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateBlog_UnauthorizedAndBadJSON(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	uc := ucMocks.NewBlogUseCaseInterface(t)
	h := NewBlogHandler(uc)

	r := gin.New()
	r.POST("/blogs", h.CreateBlog)

	// Unauthorized
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/blogs", strings.NewReader(`{"title":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Bad JSON
	r = gin.New()
	r.Use(func(c *gin.Context) { c.Set("userID", "u1") })
	r.POST("/blogs", h.CreateBlog)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/blogs", strings.NewReader(`{"title":`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateBlog_Happy(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	uc := ucMocks.NewBlogUseCaseInterface(t)
	h := NewBlogHandler(uc)

	uc.On("CreateBlog", mock.Anything, mock.AnythingOfType("*entities.Blog"), "u1").Return(nil)

	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("userID", "u1") })
	r.POST("/blogs", h.CreateBlog)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/blogs", strings.NewReader(`{"title":"t","content":"c"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetBlogsByUser_UnauthorizedAndCapLimit(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	uc := ucMocks.NewBlogUseCaseInterface(t)
	h := NewBlogHandler(uc)

	// Unauthorized
	r := gin.New()
	r.GET("/blogs", h.GetBlogsByUser)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/blogs", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// With user, limit capped to 5
	uc.On("GetBlogsByUserID", mock.Anything, "u1", int64(1), int64(5)).Return([]*entities.Blog{}, nil)
	r = gin.New()
	r.Use(func(c *gin.Context) { c.Set("userID", "u1") })
	r.GET("/blogs", h.GetBlogsByUser)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/blogs?limit=100", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetBlogByID_404And200(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	uc := ucMocks.NewBlogUseCaseInterface(t)
	h := NewBlogHandler(uc)

	uc.On("GetBlogByID", mock.Anything, "missing").Return((*entities.Blog)(nil), assert.AnError)
	r := gin.New()
	r.GET("/blogs/:id", h.GetBlogByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/blogs/missing", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	blog := &entities.Blog{Title: "ok"}
	uc.ExpectedCalls = nil // reset expectations
	uc.On("GetBlogByID", mock.Anything, "507f1f77bcf86cd799439011").Return(blog, nil)
	r = gin.New()
	r.GET("/blogs/:id", h.GetBlogByID)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/blogs/507f1f77bcf86cd799439011", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateBlog_InvalidIDAndSuccess(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	uc := ucMocks.NewBlogUseCaseInterface(t)
	h := NewBlogHandler(uc)

	// First call: GetBlogByID returns a blog, but then invalid hex triggers 400
	uc.On("GetBlogByID", mock.Anything, "badid").Return(&entities.Blog{Title: "t"}, nil)
	r := gin.New()
	r.PUT("/blogs/:id", h.UpdateBlog)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/blogs/badid", strings.NewReader(`{"content":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Success path
	uc.ExpectedCalls = nil // reset
	uc.On("GetBlogByID", mock.Anything, "507f1f77bcf86cd799439011").Return(&entities.Blog{Title: "t"}, nil)
	uc.On("UpdateBlog", mock.Anything, mock.AnythingOfType("*entities.Blog")).Return(nil)
	r = gin.New()
	r.PUT("/blogs/:id", h.UpdateBlog)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/blogs/507f1f77bcf86cd799439011", strings.NewReader(`{"content":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteBlog_204(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	uc := ucMocks.NewBlogUseCaseInterface(t)
	h := NewBlogHandler(uc)

	uc.On("DeleteBlog", mock.Anything, "507f1f77bcf86cd799439011").Return(nil)
	r := gin.New()
	r.DELETE("/blogs/:id", h.DeleteBlog)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/blogs/507f1f77bcf86cd799439011", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestGetPopularBlogs_DefaultAndCustomLimit(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	uc := ucMocks.NewBlogUseCaseInterface(t)
	h := NewBlogHandler(uc)

	// default limit 10
	uc.On("GetPopularBlogs", mock.Anything, int64(10)).Return([]*entities.BlogWithPopularity{}, nil)
	r := gin.New()
	r.GET("/blogs/popular", h.GetPopularBlogs)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/blogs/popular", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// custom limit 3
	uc.ExpectedCalls = nil
	uc.On("GetPopularBlogs", mock.Anything, int64(3)).Return([]*entities.BlogWithPopularity{}, nil)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/blogs/popular?limit=3", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFilterBlogs_PageToSkip(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	uc := ucMocks.NewBlogUseCaseInterface(t)
	h := NewBlogHandler(uc)

	uc.On("FilterBlogs", mock.Anything, mock.Anything).Return(&entities.FilterResponse{Blogs: []*entities.Blog{}}, nil)

	r := gin.New()
	r.GET("/blogs/filter", h.FilterBlogs)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/blogs/filter?page=2", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
