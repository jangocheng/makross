package static

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/insionng/makross"
	"github.com/stretchr/testify/assert"
)

func TestStatic(t *testing.T) {
	e := makross.New()
	req := httptest.NewRequest(makross.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec, makross.NotFoundHandler)
	config := StaticConfig{
		Root: "../public",
	}

	// Directory
	h := StaticWithConfig(config)
	if assert.NoError(t, h(c)) {
		assert.Contains(t, rec.Body.String(), "Makross")
	}

	// File found
	req = httptest.NewRequest(makross.GET, "/images/makross.jpg", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	//m := makross.New()
	//m.Use(StaticWithConfig(config))
	//m.ServeHTTP(rec, req)
	if assert.NoError(t, h(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotZero(t, rec.Body.Len())
		assert.Equal(t, fmt.Sprintf("%v", rec.Body.Len()), rec.Header().Get(makross.HeaderContentLength))
	}

	// File not found
	req = httptest.NewRequest(makross.GET, "/none", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec, makross.NotFoundHandler)
	he := h(c).(makross.HTTPError)
	assert.Equal(t, http.StatusNotFound, he.StatusCode())

	// HTML5
	req = httptest.NewRequest(makross.GET, "/random", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec, makross.NotFoundHandler)
	config.HTML5 = true
	static := StaticWithConfig(config)
	err := static(c)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Makross")
	}

	// Browse
	req = httptest.NewRequest(makross.GET, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec, makross.NotFoundHandler)
	config.Root = "../public/certs"
	config.Browse = true
	static = StaticWithConfig(config)
	err = static(c)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "cert.pem")
	}
}
