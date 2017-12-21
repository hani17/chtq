package user

import (
	"github.com/hani17/chtq/pkg/context"
	"net/http"
)

const (
	PROFILE = "user/profile"
)

func Profile(c *context.Context) {
	c.Data["Title"] = c.User.UserName
	c.HTML(http.StatusOK, PROFILE)
}
