package context

import (
	"fmt"
	"gopkg.in/macaron.v1"
)

type ToggleOptions struct {
	SignOutRequired bool
	SignInRequired  bool
}

func Toggle(options *ToggleOptions) macaron.Handler {
	return func(c *Context) {
		if options.SignOutRequired && c.IsLogged && c.Req.RequestURI != "/" {
			fmt.Println(c.Req.RequestURI)
			c.Redirect("/")
		}
	}
}
