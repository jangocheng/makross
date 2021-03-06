// Package makross is a high productive and modular web framework in Golang.

package makross

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultDataWriter(t *testing.T) {
	res := httptest.NewRecorder()
	err := DefaultDataWriter.Write(res, "abc")
	assert.Nil(t, err)
	assert.Equal(t, "abc", res.Body.String())

	res = httptest.NewRecorder()
	err = DefaultDataWriter.Write(res, []byte("abc"))
	assert.Nil(t, err)
	assert.Equal(t, "abc", res.Body.String())

	res = httptest.NewRecorder()
	err = DefaultDataWriter.Write(res, 123)
	assert.Nil(t, err)
	assert.Equal(t, "123", res.Body.String())

	res = httptest.NewRecorder()
	err = DefaultDataWriter.Write(res, nil)
	assert.Nil(t, err)
	assert.Equal(t, "", res.Body.String())

	res = httptest.NewRecorder()
	m := New()
	c := m.NewContext(nil, nil)
	c.Reset(res, nil)
	assert.Nil(t, c.Write("abc"))
	assert.Equal(t, "abc", res.Body.String())
}
