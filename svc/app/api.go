package app

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/syzoj/syzoj-ng-go/lib/rediskey"
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
			c.AbortWithError(500, err)
			return
		}
		resp.User = &models.User{
			Username: user.Username,
			IsAdmin:  user.IsAdmin,
		}
	}
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

type ForgetRequest struct {
	Email string `json:"email" form:"email"`
}

type forgetData struct {
	Username string
	Url      string
}

var forgetTitleTemplate = template.Must(template.New("forget").Parse("{{.Username}} 的密码重置邮件"))
var forgetBodyTemplate = template.Must(template.New("forget").Parse(`<p>请点击该链接来重置密码：</p><p><a href="{{.Url}}">{{.Url}}</a></p><p>链接有效期为 12h。如果您不是 {{.Username}}，请忽略此邮件。</p>`))

func (a *App) postApiForget(c *gin.Context) {
	if a.EmailService == nil {
		c.JSON(500, gin.H{"error": "Email service not configured"})
		return
	}
	ctx := c.Request.Context()
	var req ForgetRequest
	if err := c.Bind(&req); err != nil {
		c.AbortWithError(400, err)
		return
	}
	user, err := models.Users(qm.Where("email=?", req.Email)).One(ctx, a.Db)
	if err == sql.ErrNoRows {
		c.JSON(200, gin.H{"error_code": 1001})
		return
	} else if err != nil {
		c.AbortWithError(500, err)
		return
	}

	_, err = a.RedisCache.DoContext(ctx, "SET", rediskey.MAIN_EMAIL_PASSWORD_RECOVERY_RATELIM.Format(req.Email), "", "EX", 3600, "NX")
	if err == redis.ErrNil {
		c.JSON(200, gin.H{"error": "Only one email allowed per hour"})
		return
	} else if err != nil {
		c.AbortWithError(500, err)
		return
	}

	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		c.AbortWithError(500, err)
		return
	}
	tokenString := base64.URLEncoding.EncodeToString(token)
	_, err = a.RedisCache.DoContext(ctx, "SET", rediskey.MAIN_EMAIL_PASSWORD_RECOVERY_TOKEN.Format(req.Email, tokenString), "", "EX", 43200)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	forgetData := &forgetData{Username: user.Username.String, Url: "/forget_confirm?email=" + url.QueryEscape(req.Email) + "&token=" + tokenString}
	bufTitle := &bytes.Buffer{}
	if err := forgetTitleTemplate.Execute(bufTitle, forgetData); err != nil {
		c.AbortWithError(500, err)
		return
	}
	bufBody := &bytes.Buffer{}
	if err := forgetBodyTemplate.Execute(bufBody, forgetData); err != nil {
		c.AbortWithError(500, err)
		return
	}

	err = a.EmailService.SendEmail(ctx, req.Email, bufTitle.String(), bufBody.String())
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, gin.H{"error_code": 1})
}
