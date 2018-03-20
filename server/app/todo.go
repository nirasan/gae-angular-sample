package app

import (
	"github.com/labstack/echo"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
	"errors"
)

type Todo struct {
	ID      int64 `json:"id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
	Done    bool `json:"done"`
}

func TodoListHandler(e echo.Context) error {
	ctx := appengine.NewContext(e.Request())

	u, err := GetUser(e)
	if err != nil {
		return err
	}

	q := datastore.NewQuery("Todo").Filter("UserID =", u.ID)

	var list []Todo
	if _, err := q.GetAll(ctx, &list); err != nil {
		return err
	}

	e.JSON(http.StatusOK, list)

	return nil
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
	if t.UserID != u.ID {
		return errors.New("invalid owner")
	}

	key := datastore.NewKey(ctx, "Todo", "", t.ID, nil)
	if _, err := datastore.Put(ctx, key, t); err != nil {
		return err
	}

	e.JSON(http.StatusOK, t)

	return nil
}

func TodoDeleteHandler(e echo.Context) error {
	ctx := appengine.NewContext(e.Request())

	u, err := GetUser(e)
	if err != nil {
		return err
	}

	t := new(Todo)
	if err := e.Bind(t); err != nil {
		return err
	}
	if t.UserID != u.ID {
		return errors.New("invalid owner")
	}

	key := datastore.NewKey(ctx, "Todo", "", t.ID, nil)
	if err := datastore.Delete(ctx, key); err != nil {
		return err
	}

	e.JSON(http.StatusOK, t)

	return nil
}
