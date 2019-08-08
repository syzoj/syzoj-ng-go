package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	httplib "github.com/syzoj/syzoj-ng-go/lib/fasthttp"
	"github.com/syzoj/syzoj-ng-go/lib/testdata"
	"github.com/valyala/fasthttp"
	"golang.org/x/sys/unix"
)

func (app *App) deleteProblem(ctx *fasthttp.RequestCtx) {
	name := ctx.UserValue("name").(string)
	reason, ok := app.tryLock(name, "delete")
	if !ok {
		httplib.SendConflict(ctx, fmt.Errorf("Conflicting operation: %s", reason))
		return
	}
	defer app.unlock(name)
	path, err := app.ensurePath(name)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	if err := os.RemoveAll(path); err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	ctx.SetStatusCode(204)
	go app.automationCli.Trigger(map[string]interface{}{
		"tags": []string{"problem-data/*/delete", fmt.Sprintf("problem-data/%s/delete", name)},
		"problem": map[string]interface{}{
			"uid": name,
		},
	})
}

func (app *App) putProblemData(ctx *fasthttp.RequestCtx) {
	name := ctx.UserValue("name").(string)
	reason, ok := app.tryLock(name, "upload-data")
	if !ok {
		httplib.SendConflict(ctx, fmt.Errorf("Conflicting operation: %s", reason))
		return
	}
	defer app.unlock(name)
	path, err := app.ensurePath(name)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	tempFile, err := app.makeTempFile()
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	if err := ctx.Request.BodyWriteTo(tempFile); err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	destName := filepath.Join(path, "data.zip")
	if err := os.Rename(tempFile.Name(), destName); err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	ctx.SetStatusCode(204)
	go app.automationCli.Trigger(map[string]interface{}{
		"tags": []string{"problem-data/*/upload-data", fmt.Sprintf("problem-data/%s/upload-data", name)},
		"problem": map[string]interface{}{
			"uid": name,
		},
	})
}

func (app *App) postProblemExtract(ctx *fasthttp.RequestCtx) {
	name := ctx.UserValue("name").(string)
	reason, ok := app.tryLock(name, "extract")
	if !ok {
		httplib.SendConflict(ctx, fmt.Errorf("Conflicting operation: %s", reason))
		return
	}
	defer app.unlock(name)
	path, err := app.ensurePath(name)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	zipPath := filepath.Join(path, "data.zip")
	zipFile, err := zip.OpenReader(zipPath)
	if err != nil {
		if os.IsNotExist(err) {
			httplib.SendError(ctx, "zip file not exist")
			return
		} else {
			httplib.SendInternalError(ctx, err)
			return
		}
	}
	defer zipFile.Close()
	if len(zipFile.File) > 500 {
		httplib.SendError(ctx, "too many files")
		return
	}
	dir, err := app.makeTempDir()
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	defer os.RemoveAll(dir) // ignore errors
	var totalSize uint64
	for _, f := range zipFile.File {
		curSize := f.UncompressedSize64
		totalSize += curSize
		fpath := filepath.Join(dir, f.Name)
		if !strings.HasPrefix(fpath, dir+string(os.PathSeparator)) {
			httplib.SendError(ctx, "invalid zip filename")
			return
		}
		inFile, err := f.Open()
		if err != nil {
			httplib.SendInternalError(ctx, err)
			return
		}
		outFile, err := os.Create(fpath) // Ignore f.Mode()
		if err != nil {
			httplib.SendInternalError(ctx, err)
			return
		}
		_, err = io.CopyN(outFile, inFile, int64(curSize))
		if err != nil {
			httplib.SendInternalError(ctx, err)
			return
		}
		outFile.Close()
		inFile.Close()
	}
	if totalSize >= 50*1024*1024 {
		log.Warningf("ZIP file for problem %s has size %d", name, totalSize)
	}
	// try to swap dirs atomically if possible
	targetDir := filepath.Join(path, "data")
	err = unix.Renameat2(unix.AT_FDCWD, dir, unix.AT_FDCWD, targetDir, unix.RENAME_EXCHANGE)
	if err != nil {
		log.WithError(err).Warning("Failed to use renameat2")
		tempDir := filepath.Join(app.dataPath, "temp", randomHex())
		if os.Rename(targetDir, tempDir) != nil {
			go os.RemoveAll(tempDir)
		}
		err = os.Rename(dir, targetDir)
	}
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	ctx.SetStatusCode(204)
	go app.automationCli.Trigger(map[string]interface{}{
		"tags": []string{"problem-data/*/extract", fmt.Sprintf("problem-data/%s/extract", name)},
		"problem": map[string]interface{}{
			"uid": name,
		},
	})
}

func (app *App) getProblemParseData(ctx *fasthttp.RequestCtx) {
	name := ctx.UserValue("name").(string)
	reason, ok := app.tryLock(name, "parse-data")
	if !ok {
		httplib.SendConflict(ctx, fmt.Errorf("Conflicting operation: %s", reason))
		return
	}
	defer app.unlock(name)
	path, err := app.ensurePath(name)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	info, err := testdata.ParseTestdata(filepath.Join(path, "data"))
	if err != nil {
		httplib.SendError(ctx, err.Error())
		return
	}
	httplib.SendJSON(ctx, map[string]interface{}{
		"info": info,
	})
}
