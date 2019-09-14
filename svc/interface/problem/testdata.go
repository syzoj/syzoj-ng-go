package problem

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"time"

	"github.com/minio/minio-go"
	"github.com/syzoj/syzoj-ng-go/lib/xml"
)

var ErrLimitExceeded = errors.New("Limit exceeded")
var ErrInvalidData = errors.New("Invalid data")

// Unzip testdata.
func (s *ProblemService) UnzipTestdata(ctx context.Context, problemId string, file io.ReaderAt, size int64) error {
	prefix := "problem/" + problemId + "/testdata/"
	r, err := zip.NewReader(file, size)
	if err != nil {
		return err
	}
	const FileLimit = 1000
	const SizeLimit uint64 = 1024 * 1024 * 500
	if len(r.File) > FileLimit {
		return ErrLimitExceeded
	}
	quota := SizeLimit
	for _, file := range r.File {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		size := file.UncompressedSize64
		if size > quota {
			return ErrLimitExceeded
		}
		quota -= size
		cfile, err := file.Open()
		if err != nil {
			return err
		}
		hasher := sha256.New()
		objectName := prefix + file.Name
		if _, err := s.Minio.PutObjectWithContext(ctx, s.TestdataBucket, objectName, io.TeeReader(io.LimitReader(cfile, int64(size)), hasher), int64(size), minio.PutObjectOptions{}); err != nil {
			return err
		}
		shasum := hex.EncodeToString(hasher.Sum(nil))
		const SQLInsertMetadata = "INSERT INTO `testdata_meta` (`problem_id`, `object_name`, `sha256`) VALUES (?, ?, ?)"
		if _, err := s.Db.ExecContext(ctx, SQLInsertMetadata, problemId, objectName, shasum); err != nil {
			return err
		}
	}
	return nil
}

// Delete testdata.
func (s *ProblemService) DeleteTestdata(ctx context.Context, problemId string) error {
	const SQLDeleteMetadata = "DELETE FROM `testdata_meta` WHERE `problem_id` = ?"
	if _, err := s.Db.ExecContext(ctx, SQLDeleteMetadata, problemId); err != nil {
		return err
	}
	prefix := "problem/" + problemId + "/testdata/"
	c, cancel := context.WithCancel(ctx)
	defer cancel()
	for info := range s.Minio.ListObjects(s.TestdataBucket, prefix, true, c.Done()) {
		if err := s.Minio.RemoveObject(s.TestdataBucket, info.Key); err != nil {
			return err
		}
	}
	return nil
}

// Prepare info for judging, in XML format. The root node will be <Judge>.
func (s *ProblemService) GetJudgeInfo(ctx context.Context, problemId string) ([]byte, error) {
	return s.doCache(ctx, problemId, "judge", func() ([]byte, error) {
		root, err := s.getProblemXML(ctx, problemId)
		if err != nil {
			return nil, err
		}
		xml.EncodeTo(root, os.Stdout)
		buf := &bytes.Buffer{}
		w := new(WarningList) // TODO: The warnings are silently ignored, output them somehow
		if root != nil && root.Name.Local == "Problem" {
			judgeNode := root.SelectElement("Judge")
			if judgeNode != nil {
				if err := s.processJudgeNode(ctx, problemId, w, judgeNode); err != nil {
					return nil, err
				}
			} else {
				return nil, ErrInvalidData
			}
			if err := xml.EncodeTo(judgeNode, buf); err != nil {
				panic(err)
			}
		} else {
			w.Warningf("Root node doesn't exist or is not <Problem>")
		}
		return buf.Bytes(), nil
	})
}

// Scans the whole subtree for <File> nodes and append metadata.
func (s *ProblemService) processJudgeNode(ctx context.Context, problemId string, w Warner, tok xml.Token) error {
	prefix := "problem/" + problemId + "/testdata/"
	switch obj := tok.(type) {
	case *xml.Element:
		if obj.Name.Local == "File" {
			attr := obj.SelectAttrDefault("filename", "")
			if attr == "" {
				w.Warningf("File tag without filename attr")
				goto done
			}
			objectName := prefix + attr
			const SQLGetSHA256 = "SELECT `sha256` FROM `testdata_meta` WHERE `problem_id` = ? AND `object_name` = ?"
			var sha256 string
			err := s.Db.QueryRowContext(ctx, SQLGetSHA256, problemId, objectName).Scan(&sha256)
			if err != nil {
				if err == sql.ErrNoRows {
					w.Warningf("File %s doesn't exist", attr)
				} else {
					return err
				}
			} else {
				obj.CreateAttr("sha256", sha256)
			}

			// TODO: hardcoded one hour expiration
			url, err := s.Minio.PresignedGetObject(s.TestdataBucket, objectName, time.Hour, nil)
			if err != nil {
				return err
			} else {
				obj.CreateAttr("url", url.String())
			}
		}
	done:
		for _, child := range obj.Child {
			if err := s.processJudgeNode(ctx, problemId, w, child); err != nil {
				return err
			}
		}
	}
	return nil
}
