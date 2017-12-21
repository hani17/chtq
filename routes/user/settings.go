package user

import (
	"github.com/hani17/chtq/pkg/context"
	"net/http"
)

const (
	SETTINGS = "user/settings/settings"
)

func Settings(c *context.Context) {
	c.Data["Title"] = "Settings"
	c.HTML(http.StatusOK, SETTINGS)
}
