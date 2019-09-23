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

// Contest is an object representing the database table.
type Contest struct {
	ID             int         `boil:"id" json:"id" toml:"id" yaml:"id"`
	Title          null.String `boil:"title" json:"title,omitempty" toml:"title" yaml:"title,omitempty"`
	Subtitle       null.String `boil:"subtitle" json:"subtitle,omitempty" toml:"subtitle" yaml:"subtitle,omitempty"`
	StartTime      null.Int    `boil:"start_time" json:"start_time,omitempty" toml:"start_time" yaml:"start_time,omitempty"`
	EndTime        null.Int    `boil:"end_time" json:"end_time,omitempty" toml:"end_time" yaml:"end_time,omitempty"`
	HolderID       null.Int    `boil:"holder_id" json:"holder_id,omitempty" toml:"holder_id" yaml:"holder_id,omitempty"`
	Type           null.String `boil:"type" json:"type,omitempty" toml:"type" yaml:"type,omitempty"`
	Information    null.String `boil:"information" json:"information,omitempty" toml:"information" yaml:"information,omitempty"`
	Problems       null.String `boil:"problems" json:"problems,omitempty" toml:"problems" yaml:"problems,omitempty"`
	Admins         null.String `boil:"admins" json:"admins,omitempty" toml:"admins" yaml:"admins,omitempty"`
	RanklistID     null.Int    `boil:"ranklist_id" json:"ranklist_id,omitempty" toml:"ranklist_id" yaml:"ranklist_id,omitempty"`
	IsPublic       null.Int8   `boil:"is_public" json:"is_public,omitempty" toml:"is_public" yaml:"is_public,omitempty"`
	HideStatistics null.Int8   `boil:"hide_statistics" json:"hide_statistics,omitempty" toml:"hide_statistics" yaml:"hide_statistics,omitempty"`

	R *contestR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L contestL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ContestColumns = struct {
	ID             string
	Title          string
	Subtitle       string
	StartTime      string
	EndTime        string
	HolderID       string
	Type           string
	Information    string
	Problems       string
	Admins         string
	RanklistID     string
	IsPublic       string
	HideStatistics string
}{
	ID:             "id",
	Title:          "title",
	Subtitle:       "subtitle",
	StartTime:      "start_time",
	EndTime:        "end_time",
	HolderID:       "holder_id",
	Type:           "type",
	Information:    "information",
	Problems:       "problems",
	Admins:         "admins",
	RanklistID:     "ranklist_id",
	IsPublic:       "is_public",
	HideStatistics: "hide_statistics",
}

// Generated where

var ContestWhere = struct {
	ID             whereHelperint
	Title          whereHelpernull_String
	Subtitle       whereHelpernull_String
	StartTime      whereHelpernull_Int
	EndTime        whereHelpernull_Int
	HolderID       whereHelpernull_Int
	Type           whereHelpernull_String
	Information    whereHelpernull_String
	Problems       whereHelpernull_String
	Admins         whereHelpernull_String
	RanklistID     whereHelpernull_Int
	IsPublic       whereHelpernull_Int8
	HideStatistics whereHelpernull_Int8
}{
	ID:             whereHelperint{field: "`contest`.`id`"},
	Title:          whereHelpernull_String{field: "`contest`.`title`"},
	Subtitle:       whereHelpernull_String{field: "`contest`.`subtitle`"},
	StartTime:      whereHelpernull_Int{field: "`contest`.`start_time`"},
	EndTime:        whereHelpernull_Int{field: "`contest`.`end_time`"},
	HolderID:       whereHelpernull_Int{field: "`contest`.`holder_id`"},
	Type:           whereHelpernull_String{field: "`contest`.`type`"},
	Information:    whereHelpernull_String{field: "`contest`.`information`"},
	Problems:       whereHelpernull_String{field: "`contest`.`problems`"},
	Admins:         whereHelpernull_String{field: "`contest`.`admins`"},
	RanklistID:     whereHelpernull_Int{field: "`contest`.`ranklist_id`"},
	IsPublic:       whereHelpernull_Int8{field: "`contest`.`is_public`"},
	HideStatistics: whereHelpernull_Int8{field: "`contest`.`hide_statistics`"},
}

// ContestRels is where relationship names are stored.
var ContestRels = struct {
}{}

// contestR is where relationships are stored.
type contestR struct {
}

// NewStruct creates a new relationship struct
func (*contestR) NewStruct() *contestR {
	return &contestR{}
}

// contestL is where Load methods for each relationship are stored.
type contestL struct{}

var (
	contestAllColumns            = []string{"id", "title", "subtitle", "start_time", "end_time", "holder_id", "type", "information", "problems", "admins", "ranklist_id", "is_public", "hide_statistics"}
	contestColumnsWithoutDefault = []string{"title", "subtitle", "start_time", "end_time", "holder_id", "type", "information", "problems", "admins", "ranklist_id", "is_public", "hide_statistics"}
	contestColumnsWithDefault    = []string{"id"}
	contestPrimaryKeyColumns     = []string{"id"}
)

type (
	// ContestSlice is an alias for a slice of pointers to Contest.
	// This should generally be used opposed to []Contest.
	ContestSlice []*Contest
	// ContestHook is the signature for custom Contest hook methods
	ContestHook func(context.Context, boil.ContextExecutor, *Contest) error

	contestQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	contestType                 = reflect.TypeOf(&Contest{})
	contestMapping              = queries.MakeStructMapping(contestType)
	contestPrimaryKeyMapping, _ = queries.BindMapping(contestType, contestMapping, contestPrimaryKeyColumns)
	contestInsertCacheMut       sync.RWMutex
	contestInsertCache          = make(map[string]insertCache)
	contestUpdateCacheMut       sync.RWMutex
	contestUpdateCache          = make(map[string]updateCache)
	contestUpsertCacheMut       sync.RWMutex
	contestUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var contestBeforeInsertHooks []ContestHook
var contestBeforeUpdateHooks []ContestHook
var contestBeforeDeleteHooks []ContestHook
var contestBeforeUpsertHooks []ContestHook

var contestAfterInsertHooks []ContestHook
var contestAfterSelectHooks []ContestHook
var contestAfterUpdateHooks []ContestHook
var contestAfterDeleteHooks []ContestHook
var contestAfterUpsertHooks []ContestHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Contest) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range contestBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Contest) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range contestBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Contest) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range contestBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Contest) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range contestBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Contest) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range contestAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Contest) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range contestAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Contest) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range contestAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Contest) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range contestAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Contest) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range contestAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddContestHook registers your hook function for all future operations.
func AddContestHook(hookPoint boil.HookPoint, contestHook ContestHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		contestBeforeInsertHooks = append(contestBeforeInsertHooks, contestHook)
	case boil.BeforeUpdateHook:
		contestBeforeUpdateHooks = append(contestBeforeUpdateHooks, contestHook)
	case boil.BeforeDeleteHook:
		contestBeforeDeleteHooks = append(contestBeforeDeleteHooks, contestHook)
	case boil.BeforeUpsertHook:
		contestBeforeUpsertHooks = append(contestBeforeUpsertHooks, contestHook)
	case boil.AfterInsertHook:
		contestAfterInsertHooks = append(contestAfterInsertHooks, contestHook)
	case boil.AfterSelectHook:
		contestAfterSelectHooks = append(contestAfterSelectHooks, contestHook)
	case boil.AfterUpdateHook:
		contestAfterUpdateHooks = append(contestAfterUpdateHooks, contestHook)
	case boil.AfterDeleteHook:
		contestAfterDeleteHooks = append(contestAfterDeleteHooks, contestHook)
	case boil.AfterUpsertHook:
		contestAfterUpsertHooks = append(contestAfterUpsertHooks, contestHook)
	}
}

// One returns a single contest record from the query.
func (q contestQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Contest, error) {
	o := &Contest{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for contest")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Contest records from the query.
func (q contestQuery) All(ctx context.Context, exec boil.ContextExecutor) (ContestSlice, error) {
	var o []*Contest

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Contest slice")
	}

	if len(contestAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Contest records in the query.
func (q contestQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count contest rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q contestQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if contest exists")
	}

	return count > 0, nil
}

// Contests retrieves all the records using an executor.
func Contests(mods ...qm.QueryMod) contestQuery {
	mods = append(mods, qm.From("`contest`"))
	return contestQuery{NewQuery(mods...)}
}

// FindContest retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindContest(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*Contest, error) {
	contestObj := &Contest{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `contest` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, contestObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from contest")
	}

	return contestObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Contest) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no contest provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(contestColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	contestInsertCacheMut.RLock()
	cache, cached := contestInsertCache[key]
	contestInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			contestAllColumns,
			contestColumnsWithDefault,
			contestColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(contestType, contestMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(contestType, contestMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `contest` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `contest` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `contest` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, contestPrimaryKeyColumns))
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
		return errors.Wrap(err, "models: unable to insert into contest")
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
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == contestMapping["ID"] {
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
		return errors.Wrap(err, "models: unable to populate default values for contest")
	}

CacheNoHooks:
	if !cached {
		contestInsertCacheMut.Lock()
		contestInsertCache[key] = cache
		contestInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Contest.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Contest) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	contestUpdateCacheMut.RLock()
	cache, cached := contestUpdateCache[key]
	contestUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			contestAllColumns,
			contestPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update contest, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `contest` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, contestPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(contestType, contestMapping, append(wl, contestPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update contest row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for contest")
	}

	if !cached {
		contestUpdateCacheMut.Lock()
		contestUpdateCache[key] = cache
		contestUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q contestQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for contest")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for contest")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ContestSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contestPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `contest` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, contestPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in contest slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all contest")
	}
	return rowsAff, nil
}

var mySQLContestUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Contest) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no contest provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(contestColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLContestUniqueColumns, o)

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

	contestUpsertCacheMut.RLock()
	cache, cached := contestUpsertCache[key]
	contestUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			contestAllColumns,
			contestColumnsWithDefault,
			contestColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			contestAllColumns,
			contestPrimaryKeyColumns,
		)

		if len(update) == 0 {
			return errors.New("models: unable to upsert contest, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "contest", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `contest` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(contestType, contestMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(contestType, contestMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert for contest")
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
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == contestMapping["id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(contestType, contestMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "models: unable to retrieve unique values for contest")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, nzUniqueCols...)
	}

	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for contest")
	}

CacheNoHooks:
	if !cached {
		contestUpsertCacheMut.Lock()
		contestUpsertCache[key] = cache
		contestUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Contest record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Contest) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Contest provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), contestPrimaryKeyMapping)
	sql := "DELETE FROM `contest` WHERE `id`=?"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from contest")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for contest")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q contestQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no contestQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from contest")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for contest")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ContestSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(contestBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contestPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `contest` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, contestPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from contest slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for contest")
	}

	if len(contestAfterDeleteHooks) != 0 {
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
func (o *Contest) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindContest(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ContestSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ContestSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contestPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `contest`.* FROM `contest` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, contestPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ContestSlice")
	}

	*o = slice

	return nil
}

// ContestExists checks if the Contest row exists.
func ContestExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `contest` where `id`=? limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}

	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if contest exists")
	}

	return exists, nil
}
