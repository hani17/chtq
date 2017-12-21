package context

import (
	"fmt"
	"github.com/go-macaron/session"
	"github.com/hani17/chtq/models"
	"github.com/hani17/chtq/pkg/auth"
	"gopkg.in/macaron.v1"
	"net/http"
)

type Context struct {
	*macaron.Context
	Session session.Store

	User     *models.User
	IsLogged bool
}

// Contexter initializes a classic context for a request.
func Contexter() macaron.Handler {
	return func(ctx *macaron.Context, sess session.Store) {
		c := &Context{
			Context: ctx,
			Session: sess,
		}

		c.User = auth.SignedInUser(ctx, sess)
		if c.User != nil {
			c.IsLogged = true
			c.Data["IsLogged"] = c.IsLogged
			c.Data["LoggedUser"] = c.User
			c.Data["LoggedUserID"] = c.User.ID
			c.Data["LoggedUserName"] = c.User.UserName
			c.Data["IsAdmin"] = c.User.IsAdmin
		} else {
			c.Data["LoggedUserID"] = 0
			c.Data["LoggedUserName"] = ""
		}
		ctx.Map(c)
	}
}

func (c *Context) Handle(status int, title string, err error) {
	switch status {
	case http.StatusNotFound:
		c.Data["Title"] = "Page Not Found"
	case http.StatusInternalServerError:
		c.Data["Title"] = "Internal Server Error"
	}
	c.HTML(status, fmt.Sprintf("status/%d", status))
}

func (c *Context) ServerError(title string, err error) {
	c.Handle(http.StatusInternalServerError, title, err)
}
