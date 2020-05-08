package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-backend-example/server/models"
	"github.com/hwangseonu/gin-backend-example/server/requests"
	"github.com/hwangseonu/gin-backend-example/server/responses"
	"github.com/hwangseonu/gin-backend-example/server/security"
	"net/http"
	"time"
)

const DAY = 24 * time.Hour

func SignIn(c *gin.Context) {
	body, _ := c.Get("body")
	req := body.(*requests.SignInRequest)

	user := models.FindUserByUsername(req.Username)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "cannot find user by username"})
		return
	} else if user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "password is mismatch"})
		return
	}

	access, err1 := security.GenerateToken(security.ACCESS, user.Username)
	refresh, err2 := security.GenerateToken(security.REFRESH, user.Username)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err1.Error() + "\n" + err2.Error()})
		return
	}
	c.JSON(http.StatusOK, responses.AuthResponse{
		Access:  access,
		Refresh: refresh,
	})
	return
}

func Refresh(c *gin.Context) {
	u, _ := c.Get("user")
	ex, _ := c.Get("exp")
	user := u.(*models.User)
	exp := ex.(int64)

	access, err1 := security.GenerateToken(security.ACCESS, user.Username)
	refresh, err2 := security.GenerateToken(security.REFRESH, user.Username)

	if err1 != nil || err2 != nil {
		c.JSON(500, gin.H{"message": err1.Error() + err2.Error()})
	} else {
		if time.Unix(exp, 0).Before(time.Now().Add(7 * DAY)) {
			c.JSON(200, gin.H{"access": access, "refresh": refresh})
		} else {
			c.JSON(200, gin.H{"access": access})
		}
	}
}