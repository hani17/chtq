package auth

import (
	"github.com/go-macaron/session"
	"github.com/hani17/chtq/models"
	"github.com/hani17/chtq/models/errors"
	log "gopkg.in/clog.v1"
	"gopkg.in/macaron.v1"
)

func SignedInUser(ctx *macaron.Context, sess session.Store) *models.User {
	uid := SignedInID(ctx, sess)
	u, err := models.GetUserByID(uid)
	if err != nil {
		log.Error(4, "GetUserById: %v", err)
		return nil
	}
	return u
}

func SignedInID(c *macaron.Context, sess session.Store) int64 {
	uid := sess.Get("uid")
	if uid == nil {
		return 0
	}

	if id, ok := uid.(int64); ok {
		if _, err := models.GetUserByID(id); err != nil {
			if !errors.IsUserNotExist(err) {
				log.Error(2, "GetUserByID: %v", err)
			}
			return 0
		}
		return id
	}
	return 0
}
