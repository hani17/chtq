package routes

import (
	"github.com/hani17/chtq/pkg/config"
	"github.com/hani17/chtq/pkg/context"
	"net/http"
)

const (
	HOME = "home"
)

func Home(c *context.Context) {

	if c.IsLogged {
		c.HTML(http.StatusOK, HOME)
		return
	}

	uname := c.GetCookie(config.CookieUsername)
	if len(uname) != 0 {
		c.Redirect("/user/login")
		return
	}

	c.Data["PageIsHome"] = true
	c.HTML(http.StatusOK, HOME)
}
