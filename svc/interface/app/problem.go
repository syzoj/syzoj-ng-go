package app

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
	"github.com/syzoj/syzoj-ng-go/lib/util"
	"github.com/syzoj/syzoj-ng-go/svc/interface/problem"
)

// The limit for problem body.
const ProblemBodySizeLimit = 1000 * 1024

// POST /api/problem/new
func (a *App) HandleProblemNew(c *gin.Context) {
	ctx := c.Request.Context()
	inf := a.auth.GetInfo(c)
	if inf.UserId == "" {
		panic("user id is empty")
	}

	reader := io.LimitReader(c.Request.Body, ProblemBodySizeLimit)
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, reader); err != nil {
		c.AbortWithError(500, err)
		return
	}
	probId := util.RandomHex(16)
	if err := a.prob.CreateProblem(ctx, probId, buf.Bytes()); err != nil {
		c.AbortWithError(500, err)
		return
	}

	const SQLSetProblemOwner = "UPDATE `problems` SET `owner_uid` = ? WHERE `uid` = ?"
	if _, err := a.Db.ExecContext(ctx, SQLSetProblemOwner, inf.UserId, probId); err != nil {
		c.AbortWithError(500, err)
		return
	}
	a.stats.Inc(a.ctx, "user.problems:"+inf.UserId, 1)

	c.JSON(200, gin.H{
		"success":    true,
		"problem_id": probId,
	})
}

// PUT /api/problem/id/:problem_id
func (a *App) HandleProblemPut(c *gin.Context) {
	ctx := c.Request.Context()
	inf := a.auth.GetInfo(c)
	if inf.UserId == "" {
		panic("user id is empty")
		return
	}

	const SQLGetProblem = "SELECT `uid` FROM `problems` WHERE `uid`=? AND `owner_uid`=?"
	var probId string
	if err := a.Db.QueryRowContext(ctx, SQLGetProblem, c.Param("problem_id"), inf.UserId).Scan(&probId); err != nil {
		c.JSON(404, gin.H{
			"success": false,
			"message": "Not found",
		})
		return
	}

	reader := io.LimitReader(c.Request.Body, ProblemBodySizeLimit)
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, reader); err != nil {
		c.AbortWithError(500, err)
		return
	}
	if err := a.prob.UpdateProblem(ctx, probId, buf.Bytes()); err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.Status(204)
}

// GET /api/problem/id/:problem_id/statement
func (a *App) GetProblemStatement(c *gin.Context) {
	ctx := c.Request.Context()
	probId := c.Param("problem_id")
	v, err := a.prob.GetStatementHTML(ctx, probId)
	if err == problem.ErrNotFound {
		c.JSON(404, gin.H{
			"success": false,
			"message": "Not found",
		})
		return
	}
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.Data(200, "text/html; charset=UTF-8", v)
}

// GET /api/problem/id/:problem_id/judge-info
func (a *App) GetProblemJudgeInfo(c *gin.Context) {
	ctx := c.Request.Context()
	probId := c.Param("problem_id")
	data, err := a.prob.GetJudgeInfo(ctx, probId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.Data(200, "text/xml; charset=UTF-8", data)
}

// POST /api/problem/upload-temp
func (a *App) HandleProblemUploadTemp(c *gin.Context) {
	ctx := c.Request.Context()
	id := util.RandomHex(16)
	objectName := "temp/" + id
	// Minio handles ContentLength==-1 case
	_, err := a.Minio.PutObjectWithContext(ctx, a.TestdataBucket, objectName, c.Request.Body, c.Request.ContentLength, minio.PutObjectOptions{})
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"file_id": id,
	})
}

type UploadZipRequest struct {
	ZipFileId string `json:"zip_file_id"`
}

// POST /api/problem/id/:problem_id/upload-data
func (a *App) HandleProblemUploadData(c *gin.Context) {
	ctx := c.Request.Context()
	req := &UploadZipRequest{}
	if err := c.BindJSON(req); err != nil {
		c.AbortWithError(400, err)
		return
	}
	inf := a.auth.GetInfo(c)
	if inf.UserId == "" {
		panic("user id is empty")
		return
	}

	const SQLGetProblem = "SELECT `uid` FROM `problems` WHERE `uid`=? AND `owner_uid`=?"
	var probId string
	if err := a.Db.QueryRowContext(ctx, SQLGetProblem, c.Param("problem_id"), inf.UserId).Scan(&probId); err != nil {
		c.JSON(404, gin.H{
			"success": false,
			"message": "Not found",
		})
		return
	}

	objectName := "temp/" + req.ZipFileId
	log.Info(objectName)
	file, err := a.Minio.GetObjectWithContext(ctx, a.TestdataBucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	size, err := file.Seek(0, 2)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if err := a.prob.UnzipTestdata(ctx, probId, file, size); err != nil {
		c.AbortWithError(500, err)
		return
	}
	if err := a.Minio.RemoveObject(a.TestdataBucket, objectName); err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.Status(204)
}

// POST /api/problem/id/:problem_id/delete-data
func (a *App) HandleProblemDeleteData(c *gin.Context) {
	ctx := c.Request.Context()
	inf := a.auth.GetInfo(c)
	if inf.UserId == "" {
		panic("user id is empty")
		return
	}

	const SQLGetProblem = "SELECT `uid` FROM `problems` WHERE `uid`=? AND `owner_uid`=?"
	var probId string
	if err := a.Db.QueryRowContext(ctx, SQLGetProblem, c.Param("problem_id"), inf.UserId).Scan(&probId); err != nil {
		c.JSON(404, gin.H{
			"success": false,
			"message": "Not found",
		})
		return
	}

	if err := a.prob.DeleteTestdata(ctx, probId); err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.Status(204)
}
