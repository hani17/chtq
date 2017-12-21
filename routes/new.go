package routes

import (
	"github.com/hani17/chtq/pkg/context"
	"net/http"
)

const (
	NEW = "new"
)

func NewGet(c *context.Context) {
	c.Data["PageIsNew"] = true
	c.HTML(http.StatusOK, NEW)
}
