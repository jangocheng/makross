// Package makross is a high productive and modular web framework in Golang.

package makross

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockStore struct {
	*store
	data map[string]interface{}
}

func newMockStore() *mockStore {
	return &mockStore{newStore(), make(map[string]interface{})}
}

func (s *mockStore) Add(key string, data interface{}) int {
	for _, handler := range data.([]Handler) {
		handler(nil)
	}
	return s.store.Add(key, data)
}

func TestRouteNew(t *testing.T) {
	makross := New()
	group := newRouteGroup("/admin", makross, nil)

	r1 := group.newRoute("GET", "/users").Get()
	assert.Equal(t, "", r1.name, "route.name =")
	assert.Equal(t, "/users", r1.path, "route.path =")
	assert.Equal(t, 1, len(makross.Routes()))

	r2 := group.newRoute("GET", "/users/<id:\\d+>/*").Post()
	assert.Equal(t, "", r2.name, "route.name =")
	assert.Equal(t, "/users/<id:\\d+>/*", r2.path, "route.path =")
	assert.Equal(t, "/admin/users/<id>/", r2.template, "route.template =")
	assert.Equal(t, 2, len(makross.Routes()))
}

func TestRouteName(t *testing.T) {
	makross := New()
	group := newRouteGroup("/admin", makross, nil)

	r1 := group.newRoute("GET", "/users")
	assert.Equal(t, "", r1.name, "route.name =")
	r1.Name("user")
	assert.Equal(t, "user", r1.name, "route.name =")
	_, exists := makross.namedRoutes[r1.name]
	assert.True(t, exists)
}

func TestRouteURL(t *testing.T) {
	makross := New()
	group := newRouteGroup("/admin", makross, nil)
	r := group.newRoute("GET", "/users/<id:\\d+>/<action>/*")
	assert.Equal(t, "/admin/users/123/address/", r.URL("id", 123, "action", "address"))
	assert.Equal(t, "/admin/users/123/<action>/", r.URL("id", 123))
	assert.Equal(t, "/admin/users/123//", r.URL("id", 123, "action"))
	assert.Equal(t, "/admin/users/123/profile/", r.URL("id", 123, "action", "profile", ""))
	assert.Equal(t, "/admin/users/123/profile/", r.URL("id", 123, "action", "profile", "", "xyz/abc"))
	assert.Equal(t, "/admin/users/123/a%2C%3C%3E%3F%23/", r.URL("id", 123, "action", "a,<>?#"))
}

func newHandler(tag string, buf *bytes.Buffer) Handler {
	return func(*Context) error {
		fmt.Fprintf(buf, tag)
		return nil
	}
}

func TestRouteAdd(t *testing.T) {
	store := newMockStore()
	makross := New()
	makross.stores["GET"] = store
	assert.Equal(t, 0, store.count, "makross.stores[GET].count =")

	var buf bytes.Buffer

	group := newRouteGroup("/admin", makross, []Handler{newHandler("1.", &buf), newHandler("2.", &buf)})
	group.newRoute("GET", "/users").Get(newHandler("3.", &buf), newHandler("4.", &buf))
	assert.Equal(t, "1.2.3.4.", buf.String(), "buf@1 =")

	buf.Reset()
	group = newRouteGroup("/admin", makross, []Handler{})
	group.newRoute("GET", "/users").Get(newHandler("3.", &buf), newHandler("4.", &buf))
	assert.Equal(t, "3.4.", buf.String(), "buf@2 =")

	buf.Reset()
	group = newRouteGroup("/admin", makross, []Handler{newHandler("1.", &buf), newHandler("2.", &buf)})
	group.newRoute("GET", "/users").Get()
	assert.Equal(t, "1.2.", buf.String(), "buf@3 =")
}

func TestRouteTag(t *testing.T) {
	makross := New()
	makross.Get("/posts").Tag("posts")
	makross.Any("/users").Tag("users")
	makross.To("PUT,PATCH", "/comments").Tag("comments")
	makross.Get("/orders").Tag("GET orders").Post().Tag("POST orders")
	routes := makross.Routes()
	for _, route := range routes {
		if !assert.True(t, len(route.Tags()) > 0, route.method+" "+route.path+" should have a tag") {
			continue
		}
		tag := route.Tags()[0].(string)
		switch route.path {
		case "/posts":
			assert.Equal(t, "posts", tag)
		case "/users":
			assert.Equal(t, "users", tag)
		case "/comments":
			assert.Equal(t, "comments", tag)
		case "/orders":
			if route.method == "GET" {
				assert.Equal(t, "GET orders", tag)
			} else {
				assert.Equal(t, "POST orders", tag)
			}
		}
	}
}

func TestRouteMethods(t *testing.T) {
	makross := New()
	for _, method := range Methods {
		store := newMockStore()
		makross.stores[method] = store
		assert.Equal(t, 0, store.count, "makross.stores["+method+"].count =")
	}
	group := newRouteGroup("/admin", makross, nil)

	group.newRoute("GET", "/users").Get()
	assert.Equal(t, 1, makross.stores["GET"].(*mockStore).count, "makross.stores[GET].count =")
	group.newRoute("GET", "/users").Post()
	assert.Equal(t, 1, makross.stores["POST"].(*mockStore).count, "makross.stores[POST].count =")
	group.newRoute("GET", "/users").Patch()
	assert.Equal(t, 1, makross.stores["PATCH"].(*mockStore).count, "makross.stores[PATCH].count =")
	group.newRoute("GET", "/users").Put()
	assert.Equal(t, 1, makross.stores["PUT"].(*mockStore).count, "makross.stores[PUT].count =")
	group.newRoute("GET", "/users").Delete()
	assert.Equal(t, 1, makross.stores["DELETE"].(*mockStore).count, "makross.stores[DELETE].count =")
	group.newRoute("GET", "/users").Connect()
	assert.Equal(t, 1, makross.stores["CONNECT"].(*mockStore).count, "makross.stores[CONNECT].count =")
	group.newRoute("GET", "/users").Head()
	assert.Equal(t, 1, makross.stores["HEAD"].(*mockStore).count, "makross.stores[HEAD].count =")
	group.newRoute("GET", "/users").Options()
	assert.Equal(t, 1, makross.stores["OPTIONS"].(*mockStore).count, "makross.stores[OPTIONS].count =")
	group.newRoute("GET", "/users").Trace()
	assert.Equal(t, 1, makross.stores["TRACE"].(*mockStore).count, "makross.stores[TRACE].count =")

	group.newRoute("GET", "/posts").To("GET,POST")
	assert.Equal(t, 2, makross.stores["GET"].(*mockStore).count, "makross.stores[GET].count =")
	assert.Equal(t, 2, makross.stores["POST"].(*mockStore).count, "makross.stores[POST].count =")
	assert.Equal(t, 1, makross.stores["PUT"].(*mockStore).count, "makross.stores[PUT].count =")

	group.newRoute("GET", "/posts").To("POST")
	assert.Equal(t, 2, makross.stores["GET"].(*mockStore).count, "makross.stores[GET].count =")
	assert.Equal(t, 3, makross.stores["POST"].(*mockStore).count, "makross.stores[POST].count =")
	assert.Equal(t, 1, makross.stores["PUT"].(*mockStore).count, "makross.stores[PUT].count =")
}

func TestBuildURLTemplate(t *testing.T) {
	tests := []struct {
		path, expected string
	}{
		{"", ""},
		{"/users", "/users"},
		{"<id>", "<id>"},
		{"<id", "<id"},
		{"/users/<id>", "/users/<id>"},
		{"/users/<id:\\d+>", "/users/<id>"},
		{"/users/<:\\d+>", "/users/<>"},
		{"/users/<id>/xyz", "/users/<id>/xyz"},
		{"/users/<id:\\d+>/xyz", "/users/<id>/xyz"},
		{"/users/<id:\\d+>/<test>", "/users/<id>/<test>"},
		{"/users/<id:\\d+>/<test>/", "/users/<id>/<test>/"},
		{"/users/<id:\\d+><test>", "/users/<id><test>"},
		{"/users/<id:\\d+><test>/", "/users/<id><test>/"},
	}
	for _, test := range tests {
		actual := buildURLTemplate(test.path)
		assert.Equal(t, test.expected, actual, "buildURLTemplate("+test.path+") =")
	}
}

func TestRouteString(t *testing.T) {
	makross := New()
	makross.Get("/users/<id>")
	makross.To("GET,POST", "/users/<id>/profile")
	group := makross.Group("/admin")
	group.Post("/users")
	s := ""
	for _, route := range makross.Routes() {
		s += fmt.Sprintln(route)
	}

	assert.Equal(t, `GET /users/<id>
GET /users/<id>/profile
POST /users/<id>/profile
POST /admin/users
`, s)
}
