package app

import (
	"github.com/labstack/echo"
	"appengine/datastore"
	"appengine"
	"net/http"
)

type Todo struct {
	ID      int64
	UserID  string
	Content string
	Done    bool
}

func TodoCreateHandler(e echo.Context) error {
	ctx := appengine.NewContext(e.Request())

	u, err := GetUser(e)
	if err != nil {
		return err
	}

	t := new(Todo)
	if err := e.Bind(t); err != nil {
		return err
	}
	t.UserID = u.ID
	t.Done = false

	key := datastore.NewIncompleteKey(ctx, "Todo", nil)
	if newkey, err := datastore.Put(ctx, key, t); err != nil {
		return err
	} else {
		t.ID = newkey.IntID()
		if _, err := datastore.Put(ctx, newkey, t); err != nil {
			return err
		}
	}

	e.JSON(http.StatusOK, t)

	return nil
}

func TodoUpdateHandler(e echo.Context) error {
	ctx := appengine.NewContext(e.Request())

	u, err := GetUser(e)
	if err != nil {
		return err
	}

	t := new(Todo)
	if err := e.Bind(t); err != nil {
		return err
	}
	t.UserID = u.ID

	key := datastore.NewKey(ctx, "Todo", "", t.ID, nil)
	if _, err := datastore.Put(ctx, key, t); err != nil {
		return err
	}

	e.JSON(http.StatusOK, t)

	return nil
}
