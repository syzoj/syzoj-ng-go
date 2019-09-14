package app

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/syzoj/syzoj-ng-go/lib/util"
	"github.com/syzoj/syzoj-ng-go/svc/interface/model"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// POST /api/register
func (a *App) HandleRegister(c *gin.Context) {
	ctx := c.Request.Context()
	req := &RegisterRequest{}
	err := c.BindJSON(req)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	// I know there is a race condition here but we actually rely on database unique constraint
	const SQLCheckDuplicate = "SELECT `uid` FROM `users` WHERE `username` = ?"
	var x interface{}
	if err := a.Db.QueryRowxContext(ctx, SQLCheckDuplicate, req.UserName).Scan(&x); err == nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "Duplicate username",
		})
		return
	} else if err != sql.ErrNoRows {
		c.AbortWithError(500, err)
		return
	}

	const SQLCreateUser = "INSERT INTO `users` (`uid`, `username`, `password`, `register_time`) VALUES (:uid, :username, :password, :register_time)"
	pass, err := bcrypt.GenerateFromPassword([]byte(req.Password), 0)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	user := &model.User{
		Uid:          util.RandomHex(16),
		UserName:     req.UserName,
		Password:     pass,
		RegisterTime: time.Now(),
	}
	if _, err := a.Db.NamedExecContext(ctx, SQLCreateUser, user); err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})
}

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// POST /api/login
func (a *App) HandleLogin(c *gin.Context) {
	ctx := c.Request.Context()
	req := &LoginRequest{}
	if err := c.BindJSON(req); err != nil {
		c.AbortWithError(400, err)
		return
	}

	user := &model.User{}
	const SQLGetUserPassword = "SELECT `uid`, `password` FROM `users` WHERE `username` = ?"
	if err := a.Db.GetContext(ctx, user, SQLGetUserPassword, req.UserName); err == sql.ErrNoRows {
		c.JSON(200, gin.H{
			"success": false,
			"message": "Unknown username",
		})
		return
	} else if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password)); err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": "Invalid password",
		})
		return
	}

	if err := a.auth.SetInfo(c, user.Uid, time.Hour*24*365); err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})
}
