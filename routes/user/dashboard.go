package user

import (
	"github.com/hani17/chtq/pkg/context"
	"net/http"
)

const (
	DASHBOARD = "user/dashboard/dashboard"
)

func Dashboard(c *context.Context) {
	c.HTML(http.StatusOK, DASHBOARD)
}
