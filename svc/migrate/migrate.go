package migrate

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/syzoj/syzoj-ng-go/lib/rediskey"
	"github.com/syzoj/syzoj-ng-go/models"
	"github.com/syzoj/syzoj-ng-go/svc/redis"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

var log = logrus.StandardLogger()

type MigrateService struct {
	Db    *sql.DB
	Redis *redis.RedisService
}

func DefaultMigrateService(db *sql.DB, redis *redis.RedisService) *MigrateService {
	return &MigrateService{
		Db:    db,
		Redis: redis,
	}
}

func (m *MigrateService) MigrateProblemTags(ctx context.Context) error {
	log.Info("Computing problem tags")
	m.Db.ExecContext(ctx, "ALTER TABLE `problem` ADD COLUMN `tags` JSON DEFAULT NULL")
	tagMaps, err := models.ProblemTagMaps().All(ctx, m.Db)
	if err != nil {
		return err
	}
	tags, err := models.ProblemTags().All(ctx, m.Db)
	if err != nil {
		return err
	}
	tagsMap := make(map[int]string)
	for _, tag := range tags {
		if tag.Name.Valid {
			tagsMap[tag.ID] = tag.Name.String
		}
	}
	problemTags := make(map[int][]string)
	for _, tagMap := range tagMaps {
		pid := tagMap.ProblemID
		tid := tagMap.TagID
		problemTags[pid] = append(problemTags[pid], tagsMap[tid])
	}
	for pid, tags := range problemTags {
		b, err := json.Marshal(tags)
		if err != nil {
			return err
		}
		if _, err := m.Db.ExecContext(ctx, "UPDATE `problem` SET `tags`=? WHERE `id`=?", b, pid); err != nil {
			return err
		}
	}
	return nil
}

func (m *MigrateService) MigrateProblemCounter(ctx context.Context) error {
	log.Info("Creating problem counters")
	problems, err := models.Problems().All(ctx, m.Db)
	if err != nil {
		return err
	}
	pipeline, err := m.Redis.NewPipeline(ctx)
	if err != nil {
		return err
	}
	for _, problem := range problems {
		if problem.AcNum.Valid {
			pipeline.Do(nil, "SET", rediskey.MAIN_PROBLEM_ACCEPTS.Format(strconv.Itoa(problem.ID)), problem.AcNum.Int)
		}
		if problem.SubmitNum.Valid {
			pipeline.Do(nil, "SET", rediskey.MAIN_PROBLEM_SUBMITS.Format(strconv.Itoa(problem.ID)), problem.SubmitNum.Int)
		}
	}
	return pipeline.Flush(ctx)
}

func (m *MigrateService) MigrateUserSubmissions(ctx context.Context) error {
	log.Info("Migrating user submissions")
	query, err := models.JudgeStates(qm.OrderBy("submit_time")).QueryContext(ctx, m.Db)
	if err != nil {
		return err
	}
	pipeline, err := m.Redis.NewPipeline(ctx)
	if err != nil {
		return err
	}
	for query.Next() {
		submission := &models.JudgeState{}
		if err := queries.Bind(query, submission); err != nil {
			log.WithError(err).Error("failed to bind model")
			continue
		}
		if !submission.UserID.Valid || !submission.ProblemID.Valid {
			continue
		}
		userId := submission.UserID.Int
		problemId := submission.ProblemID.Int
		if submission.Status == null.StringFrom("Accepted") {
			pipeline.Do(nil, "HSET", rediskey.MAIN_USER_LAST_ACCEPT.Format(strconv.Itoa(userId)), problemId, submission.ID)
		}
		pipeline.Do(nil, "HSET", rediskey.MAIN_USER_LAST_SUBMISSION.Format(strconv.Itoa(userId)), problemId, submission.ID)
	}
	if err := query.Err(); err != nil {
		return err
	}
	return pipeline.Flush(ctx)
}

func (m *MigrateService) All(ctx context.Context) error {
	if err := m.MigrateProblemTags(ctx); err != nil {
		return err
	}
	if err := m.MigrateProblemCounter(ctx); err != nil {
		return err
	}
	if err := m.MigrateUserSubmissions(ctx); err != nil {
		return err
	}
	return nil
}
