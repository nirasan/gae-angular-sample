package app

import (
	"testing"
	"strings"
	"net/http/httptest"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
	"github.com/labstack/echo"
	"net/http"
	"encoding/json"
	"io"
)

func TestTodoCreateHandler(t *testing.T) {
	// create aetest instance, request, responce
	instance, req, res, err := NewRequests("POST", "/todo/", strings.NewReader(`{"user_id":"user1", "content":"content1"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer instance.Close()

	// create echo context
	e := echo.New().NewContext(req, res)
	e.Set("User", &User{ID: "user1"})

	// exec create handler
	if err := TodoCreateHandler(e); err != nil {
		t.Fatal(err)
	}

	// check response code
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %v", res)
	}
	// check response body
	td := new(Todo)
	if err := json.NewDecoder(res.Body).Decode(td); err != nil {
		t.Error(err)
	}
	if td.Content != "content1" {
		t.Errorf("invalid content: %v", td)
	}
	t.Logf("%#v", td)

	// is todo record exists?
	ctx := appengine.NewContext(req)
	key := datastore.NewKey(ctx, "Todo", "", td.ID, nil)
	if err := datastore.Get(ctx, key, &Todo{}); err != nil {
		t.Error(err)
	}
}

func NewRequests(method, url string, body io.Reader) (instance aetest.Instance, req *http.Request, res *httptest.ResponseRecorder, err error) {
	instance, err = aetest.NewInstance(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		return
	}

	req, err = instance.NewRequest(method, url, body)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	res = httptest.NewRecorder()

	return
}
