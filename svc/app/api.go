package app

import (
	"crypto/subtle"
	"database/sql"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/syzoj/syzoj-ng-go/models"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

type HeaderResponse struct {
	User *models.User `json:"user"`
}

func (a *App) getHeader(c *gin.Context) {
	ctx := c.Request.Context()
	resp := &HeaderResponse{}
	userId := c.GetInt(GIN_USER_ID)
	if userId != 0 {
		user, err := models.Users(qm.Select("username", "is_admin"), qm.Where("id=?", userId)).One(ctx, a.Db)
		if err != nil {
			c.Error(err)
			goto label1
		}
		resp.User = &models.User{
			Username: user.Username,
			IsAdmin:  user.IsAdmin,
		}
	}
label1:
	c.JSON(200, resp)
}

type LoginRequest struct {
	UserName string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (a *App) postApiLogin(c *gin.Context) {
	ctx := c.Request.Context()
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithError(400, err)
		return
	}

	user, err := models.Users(qm.Where("username=?", req.UserName)).One(ctx, a.Db)
	if err == sql.ErrNoRows {
		c.JSON(200, gin.H{"error_code": 1001})
		return
	} else if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if !user.Password.Valid || user.Password.String == "" {
		c.JSON(200, gin.H{"error_code": 1003})
		return
	}
	if subtle.ConstantTimeCompare([]byte(user.Password.String), []byte(req.Password)) == 0 {
		c.JSON(200, gin.H{"error_code": 1002})
		return
	}
	data, _ := json.Marshal([]string{req.UserName, req.Password})
	c.SetCookie("login", string(data), 86400*31, "/", "", false, true)
	c.JSON(200, gin.H{"error_code": 1})
}
