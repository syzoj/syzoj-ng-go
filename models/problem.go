// Code generated by SQLBoiler 3.5.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// Problem is an object representing the database table.
type Problem struct {
	ID               int         `boil:"id" json:"id" toml:"id" yaml:"id"`
	Title            null.String `boil:"title" json:"title,omitempty" toml:"title" yaml:"title,omitempty"`
	UserID           null.Int    `boil:"user_id" json:"user_id,omitempty" toml:"user_id" yaml:"user_id,omitempty"`
	PublicizerID     null.Int    `boil:"publicizer_id" json:"publicizer_id,omitempty" toml:"publicizer_id" yaml:"publicizer_id,omitempty"`
	IsAnonymous      null.Int8   `boil:"is_anonymous" json:"is_anonymous,omitempty" toml:"is_anonymous" yaml:"is_anonymous,omitempty"`
	Description      null.String `boil:"description" json:"description,omitempty" toml:"description" yaml:"description,omitempty"`
	InputFormat      null.String `boil:"input_format" json:"input_format,omitempty" toml:"input_format" yaml:"input_format,omitempty"`
	OutputFormat     null.String `boil:"output_format" json:"output_format,omitempty" toml:"output_format" yaml:"output_format,omitempty"`
	Example          null.String `boil:"example" json:"example,omitempty" toml:"example" yaml:"example,omitempty"`
	LimitAndHint     null.String `boil:"limit_and_hint" json:"limit_and_hint,omitempty" toml:"limit_and_hint" yaml:"limit_and_hint,omitempty"`
	TimeLimit        null.Int    `boil:"time_limit" json:"time_limit,omitempty" toml:"time_limit" yaml:"time_limit,omitempty"`
	MemoryLimit      null.Int    `boil:"memory_limit" json:"memory_limit,omitempty" toml:"memory_limit" yaml:"memory_limit,omitempty"`
	AdditionalFileID null.Int    `boil:"additional_file_id" json:"additional_file_id,omitempty" toml:"additional_file_id" yaml:"additional_file_id,omitempty"`
	AcNum            null.Int    `boil:"ac_num" json:"ac_num,omitempty" toml:"ac_num" yaml:"ac_num,omitempty"`
	SubmitNum        null.Int    `boil:"submit_num" json:"submit_num,omitempty" toml:"submit_num" yaml:"submit_num,omitempty"`
	IsPublic         null.Int8   `boil:"is_public" json:"is_public,omitempty" toml:"is_public" yaml:"is_public,omitempty"`
	FileIo           null.Int8   `boil:"file_io" json:"file_io,omitempty" toml:"file_io" yaml:"file_io,omitempty"`
	FileIoInputName  null.String `boil:"file_io_input_name" json:"file_io_input_name,omitempty" toml:"file_io_input_name" yaml:"file_io_input_name,omitempty"`
	FileIoOutputName null.String `boil:"file_io_output_name" json:"file_io_output_name,omitempty" toml:"file_io_output_name" yaml:"file_io_output_name,omitempty"`
	PublicizeTime    null.Time   `boil:"publicize_time" json:"publicize_time,omitempty" toml:"publicize_time" yaml:"publicize_time,omitempty"`
	Type             null.String `boil:"type" json:"type,omitempty" toml:"type" yaml:"type,omitempty"`
	Tags             null.String `boil:"tags" json:"tags,omitempty" toml:"tags" yaml:"tags,omitempty"`

	R *problemR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L problemL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ProblemColumns = struct {
	ID               string
	Title            string
	UserID           string
	PublicizerID     string
	IsAnonymous      string
	Description      string
	InputFormat      string
	OutputFormat     string
	Example          string
	LimitAndHint     string
	TimeLimit        string
	MemoryLimit      string
	AdditionalFileID string
	AcNum            string
	SubmitNum        string
	IsPublic         string
	FileIo           string
	FileIoInputName  string
	FileIoOutputName string
	PublicizeTime    string
	Type             string
	Tags             string
}{
	ID:               "id",
	Title:            "title",
	UserID:           "user_id",
	PublicizerID:     "publicizer_id",
	IsAnonymous:      "is_anonymous",
	Description:      "description",
	InputFormat:      "input_format",
	OutputFormat:     "output_format",
	Example:          "example",
	LimitAndHint:     "limit_and_hint",
	TimeLimit:        "time_limit",
	MemoryLimit:      "memory_limit",
	AdditionalFileID: "additional_file_id",
	AcNum:            "ac_num",
	SubmitNum:        "submit_num",
	IsPublic:         "is_public",
	FileIo:           "file_io",
	FileIoInputName:  "file_io_input_name",
	FileIoOutputName: "file_io_output_name",
	PublicizeTime:    "publicize_time",
	Type:             "type",
	Tags:             "tags",
}

// Generated where

type whereHelpernull_Time struct{ field string }

func (w whereHelpernull_Time) EQ(x null.Time) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Time) NEQ(x null.Time) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Time) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Time) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }
func (w whereHelpernull_Time) LT(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Time) LTE(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Time) GT(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Time) GTE(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var ProblemWhere = struct {
	ID               whereHelperint
	Title            whereHelpernull_String
	UserID           whereHelpernull_Int
	PublicizerID     whereHelpernull_Int
	IsAnonymous      whereHelpernull_Int8
	Description      whereHelpernull_String
	InputFormat      whereHelpernull_String
	OutputFormat     whereHelpernull_String
	Example          whereHelpernull_String
	LimitAndHint     whereHelpernull_String
	TimeLimit        whereHelpernull_Int
	MemoryLimit      whereHelpernull_Int
	AdditionalFileID whereHelpernull_Int
	AcNum            whereHelpernull_Int
	SubmitNum        whereHelpernull_Int
	IsPublic         whereHelpernull_Int8
	FileIo           whereHelpernull_Int8
	FileIoInputName  whereHelpernull_String
	FileIoOutputName whereHelpernull_String
	PublicizeTime    whereHelpernull_Time
	Type             whereHelpernull_String
	Tags             whereHelpernull_String
}{
	ID:               whereHelperint{field: "`problem`.`id`"},
	Title:            whereHelpernull_String{field: "`problem`.`title`"},
	UserID:           whereHelpernull_Int{field: "`problem`.`user_id`"},
	PublicizerID:     whereHelpernull_Int{field: "`problem`.`publicizer_id`"},
	IsAnonymous:      whereHelpernull_Int8{field: "`problem`.`is_anonymous`"},
	Description:      whereHelpernull_String{field: "`problem`.`description`"},
	InputFormat:      whereHelpernull_String{field: "`problem`.`input_format`"},
	OutputFormat:     whereHelpernull_String{field: "`problem`.`output_format`"},
	Example:          whereHelpernull_String{field: "`problem`.`example`"},
	LimitAndHint:     whereHelpernull_String{field: "`problem`.`limit_and_hint`"},
	TimeLimit:        whereHelpernull_Int{field: "`problem`.`time_limit`"},
	MemoryLimit:      whereHelpernull_Int{field: "`problem`.`memory_limit`"},
	AdditionalFileID: whereHelpernull_Int{field: "`problem`.`additional_file_id`"},
	AcNum:            whereHelpernull_Int{field: "`problem`.`ac_num`"},
	SubmitNum:        whereHelpernull_Int{field: "`problem`.`submit_num`"},
	IsPublic:         whereHelpernull_Int8{field: "`problem`.`is_public`"},
	FileIo:           whereHelpernull_Int8{field: "`problem`.`file_io`"},
	FileIoInputName:  whereHelpernull_String{field: "`problem`.`file_io_input_name`"},
	FileIoOutputName: whereHelpernull_String{field: "`problem`.`file_io_output_name`"},
	PublicizeTime:    whereHelpernull_Time{field: "`problem`.`publicize_time`"},
	Type:             whereHelpernull_String{field: "`problem`.`type`"},
	Tags:             whereHelpernull_String{field: "`problem`.`tags`"},
}

// ProblemRels is where relationship names are stored.
var ProblemRels = struct {
}{}

// problemR is where relationships are stored.
type problemR struct {
}

// NewStruct creates a new relationship struct
func (*problemR) NewStruct() *problemR {
	return &problemR{}
}

// problemL is where Load methods for each relationship are stored.
type problemL struct{}

var (
	problemAllColumns            = []string{"id", "title", "user_id", "publicizer_id", "is_anonymous", "description", "input_format", "output_format", "example", "limit_and_hint", "time_limit", "memory_limit", "additional_file_id", "ac_num", "submit_num", "is_public", "file_io", "file_io_input_name", "file_io_output_name", "publicize_time", "type", "tags"}
	problemColumnsWithoutDefault = []string{"title", "user_id", "publicizer_id", "is_anonymous", "description", "input_format", "output_format", "example", "limit_and_hint", "time_limit", "memory_limit", "additional_file_id", "ac_num", "submit_num", "is_public", "file_io", "file_io_input_name", "file_io_output_name", "publicize_time", "tags"}
	problemColumnsWithDefault    = []string{"id", "type"}
	problemPrimaryKeyColumns     = []string{"id"}
)

type (
	// ProblemSlice is an alias for a slice of pointers to Problem.
	// This should generally be used opposed to []Problem.
	ProblemSlice []*Problem
	// ProblemHook is the signature for custom Problem hook methods
	ProblemHook func(context.Context, boil.ContextExecutor, *Problem) error

	problemQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	problemType                 = reflect.TypeOf(&Problem{})
	problemMapping              = queries.MakeStructMapping(problemType)
	problemPrimaryKeyMapping, _ = queries.BindMapping(problemType, problemMapping, problemPrimaryKeyColumns)
	problemInsertCacheMut       sync.RWMutex
	problemInsertCache          = make(map[string]insertCache)
	problemUpdateCacheMut       sync.RWMutex
	problemUpdateCache          = make(map[string]updateCache)
	problemUpsertCacheMut       sync.RWMutex
	problemUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var problemBeforeInsertHooks []ProblemHook
var problemBeforeUpdateHooks []ProblemHook
var problemBeforeDeleteHooks []ProblemHook
var problemBeforeUpsertHooks []ProblemHook

var problemAfterInsertHooks []ProblemHook
var problemAfterSelectHooks []ProblemHook
var problemAfterUpdateHooks []ProblemHook
var problemAfterDeleteHooks []ProblemHook
var problemAfterUpsertHooks []ProblemHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Problem) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range problemBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Problem) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range problemBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Problem) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range problemBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Problem) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range problemBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Problem) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range problemAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Problem) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range problemAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Problem) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range problemAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Problem) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range problemAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Problem) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range problemAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddProblemHook registers your hook function for all future operations.
func AddProblemHook(hookPoint boil.HookPoint, problemHook ProblemHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		problemBeforeInsertHooks = append(problemBeforeInsertHooks, problemHook)
	case boil.BeforeUpdateHook:
		problemBeforeUpdateHooks = append(problemBeforeUpdateHooks, problemHook)
	case boil.BeforeDeleteHook:
		problemBeforeDeleteHooks = append(problemBeforeDeleteHooks, problemHook)
	case boil.BeforeUpsertHook:
		problemBeforeUpsertHooks = append(problemBeforeUpsertHooks, problemHook)
	case boil.AfterInsertHook:
		problemAfterInsertHooks = append(problemAfterInsertHooks, problemHook)
	case boil.AfterSelectHook:
		problemAfterSelectHooks = append(problemAfterSelectHooks, problemHook)
	case boil.AfterUpdateHook:
		problemAfterUpdateHooks = append(problemAfterUpdateHooks, problemHook)
	case boil.AfterDeleteHook:
		problemAfterDeleteHooks = append(problemAfterDeleteHooks, problemHook)
	case boil.AfterUpsertHook:
		problemAfterUpsertHooks = append(problemAfterUpsertHooks, problemHook)
	}
}

// One returns a single problem record from the query.
func (q problemQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Problem, error) {
	o := &Problem{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for problem")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Problem records from the query.
func (q problemQuery) All(ctx context.Context, exec boil.ContextExecutor) (ProblemSlice, error) {
	var o []*Problem

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Problem slice")
	}

	if len(problemAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Problem records in the query.
func (q problemQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count problem rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q problemQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if problem exists")
	}

	return count > 0, nil
}

// Problems retrieves all the records using an executor.
func Problems(mods ...qm.QueryMod) problemQuery {
	mods = append(mods, qm.From("`problem`"))
	return problemQuery{NewQuery(mods...)}
}

// FindProblem retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindProblem(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*Problem, error) {
	problemObj := &Problem{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `problem` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, problemObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from problem")
	}

	return problemObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Problem) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no problem provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(problemColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	problemInsertCacheMut.RLock()
	cache, cached := problemInsertCache[key]
	problemInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			problemAllColumns,
			problemColumnsWithDefault,
			problemColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(problemType, problemMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(problemType, problemMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `problem` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `problem` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `problem` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, problemPrimaryKeyColumns))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	result, err := exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into problem")
	}

	var lastID int64
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = int(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == problemMapping["ID"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.ID,
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, identifierCols...)
	}

	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for problem")
	}

CacheNoHooks:
	if !cached {
		problemInsertCacheMut.Lock()
		problemInsertCache[key] = cache
		problemInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Problem.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Problem) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	problemUpdateCacheMut.RLock()
	cache, cached := problemUpdateCache[key]
	problemUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			problemAllColumns,
			problemPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update problem, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `problem` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, problemPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(problemType, problemMapping, append(wl, problemPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update problem row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for problem")
	}

	if !cached {
		problemUpdateCacheMut.Lock()
		problemUpdateCache[key] = cache
		problemUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q problemQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for problem")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for problem")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ProblemSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), problemPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `problem` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, problemPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in problem slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all problem")
	}
	return rowsAff, nil
}

var mySQLProblemUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Problem) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no problem provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(problemColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLProblemUniqueColumns, o)

	if len(nzUniques) == 0 {
		return errors.New("cannot upsert with a table that cannot conflict on a unique column")
	}

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzUniques {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	problemUpsertCacheMut.RLock()
	cache, cached := problemUpsertCache[key]
	problemUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			problemAllColumns,
			problemColumnsWithDefault,
			problemColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			problemAllColumns,
			problemPrimaryKeyColumns,
		)

		if len(update) == 0 {
			return errors.New("models: unable to upsert problem, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "problem", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `problem` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(problemType, problemMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(problemType, problemMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	result, err := exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "models: unable to upsert for problem")
	}

	var lastID int64
	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = int(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == problemMapping["id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(problemType, problemMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "models: unable to retrieve unique values for problem")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, nzUniqueCols...)
	}

	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for problem")
	}

CacheNoHooks:
	if !cached {
		problemUpsertCacheMut.Lock()
		problemUpsertCache[key] = cache
		problemUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Problem record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Problem) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Problem provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), problemPrimaryKeyMapping)
	sql := "DELETE FROM `problem` WHERE `id`=?"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from problem")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for problem")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q problemQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no problemQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from problem")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for problem")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ProblemSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(problemBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), problemPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `problem` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, problemPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from problem slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for problem")
	}

	if len(problemAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Problem) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindProblem(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ProblemSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ProblemSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), problemPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `problem`.* FROM `problem` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, problemPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ProblemSlice")
	}

	*o = slice

	return nil
}

// ProblemExists checks if the Problem row exists.
func ProblemExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `problem` where `id`=? limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}

	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if problem exists")
	}

	return exists, nil
}