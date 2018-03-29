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
	"github.com/satori/go.uuid"
	"fmt"
)

func TestTodoCreateHandler(t *testing.T) {
	// create aetest instance, request, responce
	uid := uuid.NewV4().String()
	instance, req, res, err := NewRequests("POST", "/todo/", strings.NewReader(fmt.Sprintf(`{"user_id":"%s", "content":"content1"}`, uid)))
	if err != nil {
		t.Fatal(err)
	}
	defer instance.Close()

	// create echo context
	e := echo.New().NewContext(req, res)
	e.Set("User", &User{ID: uid})

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

func TestTodoListHandler(t *testing.T) {
	// create aetest instance, request, responce
	instance, req, res, err := NewRequests("GET", "/todo/", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer instance.Close()

	// create echo context
	e := echo.New().NewContext(req, res)
	uid := uuid.NewV4().String()
	e.Set("User", &User{ID: uid})

	// create todo list
	ctx := appengine.NewContext(req)
	for i := 0; i < 3; i++ {
		key := datastore.NewIncompleteKey(ctx, "Todo", nil)
		td := Todo{
			UserID: uid,
			Content: "content",
		}
		if _, err := datastore.Put(ctx, key, &td); err != nil {
			t.Error(err)
		}
	}

	// exec create handler
	if err := TodoListHandler(e); err != nil {
		t.Fatal(err)
	}

	// check response code
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %v", res)
	}
	// check response body
	list := []Todo{}
	b := res.Body.String()
	t.Log(b)
	if err := json.Unmarshal([]byte(b), &list); err != nil {
		t.Error(err)
	}
	if len(list) != 3 {
		t.Errorf("invalid response: %v", list)
	}
	t.Logf("%#v", list)
}

func TestTodoUpdateHandler(t *testing.T) {
	// create aetest instance, request, responce
	uid := uuid.NewV4().String()
	id := int64(1)
	instance, req, res, err := NewRequests("PUT", "/todo/", strings.NewReader(fmt.Sprintf(`{"id":%d, "user_id":"%s", "content":"updated content"}`, id, uid)))
	if err != nil {
		t.Fatal(err)
	}
	defer instance.Close()

	// create echo context
	e := echo.New().NewContext(req, res)
	e.Set("User", &User{ID: uid})

	// create record
	ctx := appengine.NewContext(req)
	key := datastore.NewKey(ctx, "Todo", "", id,nil)
	td := Todo{
		ID: id,
		UserID: uid,
		Content: "content",
	}
	if _, err := datastore.Put(ctx, key, &td); err != nil {
		t.Error(err)
	}

	// exec update handler
	if err := TodoUpdateHandler(e); err != nil {
		t.Fatal(err)
	}

	// check response code
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %v", res)
	}
	// check response body
	if err := json.NewDecoder(res.Body).Decode(&td); err != nil {
		t.Error(err)
	}
	if td.Content != "updated content" {
		t.Errorf("invalid content: %v", td)
	}
	t.Logf("%#v", td)

	// is todo record updated?
	td2 := Todo{}
	if err := datastore.Get(ctx, key, &td2); err != nil {
		t.Error(err)
	}
	if td2.Content != "updated content" {
		t.Errorf("invalid content: %v", td2)
	}
	t.Logf("%#v", td2)
}

func TestTodoDeleteHandler(t *testing.T) {
	// create aetest instance, request, responce
	uid := uuid.NewV4().String()
	id := int64(1)
	instance, req, res, err := NewRequests("DELETE", fmt.Sprintf("/todo/%d", id), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer instance.Close()

	// create echo context
	e := echo.New().NewContext(req, res)
	e.Set("User", &User{ID: uid})
	e.SetParamNames("id")
	e.SetParamValues("1")

	// create record
	ctx := appengine.NewContext(req)
	key := datastore.NewKey(ctx, "Todo", "", id,nil)
	td := Todo{
		ID: id,
		UserID: uid,
		Content: "content",
	}
	if _, err := datastore.Put(ctx, key, &td); err != nil {
		t.Error(err)
	}

	// exec delete handler
	if err := TodoDeleteHandler(e); err != nil {
		t.Fatal(err)
	}

	// check response code
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %v", res)
	}

	// is todo record not exists?
	if err := datastore.Get(ctx, key, &Todo{}); err != datastore.ErrNoSuchEntity {
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
