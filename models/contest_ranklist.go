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

// ContestRanklist is an object representing the database table.
type ContestRanklist struct {
	ID            int         `boil:"id" json:"id" toml:"id" yaml:"id"`
	RankingParams null.String `boil:"ranking_params" json:"ranking_params,omitempty" toml:"ranking_params" yaml:"ranking_params,omitempty"`
	Ranklist      string      `boil:"ranklist" json:"ranklist" toml:"ranklist" yaml:"ranklist"`

	R *contestRanklistR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L contestRanklistL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ContestRanklistColumns = struct {
	ID            string
	RankingParams string
	Ranklist      string
}{
	ID:            "id",
	RankingParams: "ranking_params",
	Ranklist:      "ranklist",
}

// Generated where

var ContestRanklistWhere = struct {
	ID            whereHelperint
	RankingParams whereHelpernull_String
	Ranklist      whereHelperstring
}{
	ID:            whereHelperint{field: "`contest_ranklist`.`id`"},
	RankingParams: whereHelpernull_String{field: "`contest_ranklist`.`ranking_params`"},
	Ranklist:      whereHelperstring{field: "`contest_ranklist`.`ranklist`"},
}

// ContestRanklistRels is where relationship names are stored.
var ContestRanklistRels = struct {
}{}

// contestRanklistR is where relationships are stored.
type contestRanklistR struct {
}

// NewStruct creates a new relationship struct
func (*contestRanklistR) NewStruct() *contestRanklistR {
	return &contestRanklistR{}
}

// contestRanklistL is where Load methods for each relationship are stored.
type contestRanklistL struct{}

var (
	contestRanklistAllColumns            = []string{"id", "ranking_params", "ranklist"}
	contestRanklistColumnsWithoutDefault = []string{"ranking_params"}
	contestRanklistColumnsWithDefault    = []string{"id", "ranklist"}
	contestRanklistPrimaryKeyColumns     = []string{"id"}
)

type (
	// ContestRanklistSlice is an alias for a slice of pointers to ContestRanklist.
	// This should generally be used opposed to []ContestRanklist.
	ContestRanklistSlice []*ContestRanklist
	// ContestRanklistHook is the signature for custom ContestRanklist hook methods
	ContestRanklistHook func(context.Context, boil.ContextExecutor, *ContestRanklist) error

	contestRanklistQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	contestRanklistType                 = reflect.TypeOf(&ContestRanklist{})
	contestRanklistMapping              = queries.MakeStructMapping(contestRanklistType)
	contestRanklistPrimaryKeyMapping, _ = queries.BindMapping(contestRanklistType, contestRanklistMapping, contestRanklistPrimaryKeyColumns)
	contestRanklistInsertCacheMut       sync.RWMutex
	contestRanklistInsertCache          = make(map[string]insertCache)
	contestRanklistUpdateCacheMut       sync.RWMutex
	contestRanklistUpdateCache          = make(map[string]updateCache)
	contestRanklistUpsertCacheMut       sync.RWMutex
	contestRanklistUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var contestRanklistBeforeInsertHooks []ContestRanklistHook
var contestRanklistBeforeUpdateHooks []ContestRanklistHook
var contestRanklistBeforeDeleteHooks []ContestRanklistHook
var contestRanklistBeforeUpsertHooks []ContestRanklistHook

var contestRanklistAfterInsertHooks []ContestRanklistHook
var contestRanklistAfterSelectHooks []ContestRanklistHook
var contestRanklistAfterUpdateHooks []ContestRanklistHook
var contestRanklistAfterDeleteHooks []ContestRanklistHook
var contestRanklistAfterUpsertHooks []ContestRanklistHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *ContestRanklist) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range contestRanklistBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *ContestRanklist) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range contestRanklistBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *ContestRanklist) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range contestRanklistBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *ContestRanklist) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range contestRanklistBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *ContestRanklist) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range contestRanklistAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *ContestRanklist) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range contestRanklistAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *ContestRanklist) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range contestRanklistAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *ContestRanklist) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range contestRanklistAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *ContestRanklist) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range contestRanklistAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddContestRanklistHook registers your hook function for all future operations.
func AddContestRanklistHook(hookPoint boil.HookPoint, contestRanklistHook ContestRanklistHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		contestRanklistBeforeInsertHooks = append(contestRanklistBeforeInsertHooks, contestRanklistHook)
	case boil.BeforeUpdateHook:
		contestRanklistBeforeUpdateHooks = append(contestRanklistBeforeUpdateHooks, contestRanklistHook)
	case boil.BeforeDeleteHook:
		contestRanklistBeforeDeleteHooks = append(contestRanklistBeforeDeleteHooks, contestRanklistHook)
	case boil.BeforeUpsertHook:
		contestRanklistBeforeUpsertHooks = append(contestRanklistBeforeUpsertHooks, contestRanklistHook)
	case boil.AfterInsertHook:
		contestRanklistAfterInsertHooks = append(contestRanklistAfterInsertHooks, contestRanklistHook)
	case boil.AfterSelectHook:
		contestRanklistAfterSelectHooks = append(contestRanklistAfterSelectHooks, contestRanklistHook)
	case boil.AfterUpdateHook:
		contestRanklistAfterUpdateHooks = append(contestRanklistAfterUpdateHooks, contestRanklistHook)
	case boil.AfterDeleteHook:
		contestRanklistAfterDeleteHooks = append(contestRanklistAfterDeleteHooks, contestRanklistHook)
	case boil.AfterUpsertHook:
		contestRanklistAfterUpsertHooks = append(contestRanklistAfterUpsertHooks, contestRanklistHook)
	}
}

// One returns a single contestRanklist record from the query.
func (q contestRanklistQuery) One(ctx context.Context, exec boil.ContextExecutor) (*ContestRanklist, error) {
	o := &ContestRanklist{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for contest_ranklist")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all ContestRanklist records from the query.
func (q contestRanklistQuery) All(ctx context.Context, exec boil.ContextExecutor) (ContestRanklistSlice, error) {
	var o []*ContestRanklist

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to ContestRanklist slice")
	}

	if len(contestRanklistAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all ContestRanklist records in the query.
func (q contestRanklistQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count contest_ranklist rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q contestRanklistQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if contest_ranklist exists")
	}

	return count > 0, nil
}

// ContestRanklists retrieves all the records using an executor.
func ContestRanklists(mods ...qm.QueryMod) contestRanklistQuery {
	mods = append(mods, qm.From("`contest_ranklist`"))
	return contestRanklistQuery{NewQuery(mods...)}
}

// FindContestRanklist retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindContestRanklist(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*ContestRanklist, error) {
	contestRanklistObj := &ContestRanklist{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `contest_ranklist` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, contestRanklistObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from contest_ranklist")
	}

	return contestRanklistObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *ContestRanklist) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no contest_ranklist provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(contestRanklistColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	contestRanklistInsertCacheMut.RLock()
	cache, cached := contestRanklistInsertCache[key]
	contestRanklistInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			contestRanklistAllColumns,
			contestRanklistColumnsWithDefault,
			contestRanklistColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(contestRanklistType, contestRanklistMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(contestRanklistType, contestRanklistMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `contest_ranklist` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `contest_ranklist` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `contest_ranklist` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, contestRanklistPrimaryKeyColumns))
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
		return errors.Wrap(err, "models: unable to insert into contest_ranklist")
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
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == contestRanklistMapping["ID"] {
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
		return errors.Wrap(err, "models: unable to populate default values for contest_ranklist")
	}

CacheNoHooks:
	if !cached {
		contestRanklistInsertCacheMut.Lock()
		contestRanklistInsertCache[key] = cache
		contestRanklistInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the ContestRanklist.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *ContestRanklist) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	contestRanklistUpdateCacheMut.RLock()
	cache, cached := contestRanklistUpdateCache[key]
	contestRanklistUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			contestRanklistAllColumns,
			contestRanklistPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update contest_ranklist, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `contest_ranklist` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, contestRanklistPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(contestRanklistType, contestRanklistMapping, append(wl, contestRanklistPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update contest_ranklist row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for contest_ranklist")
	}

	if !cached {
		contestRanklistUpdateCacheMut.Lock()
		contestRanklistUpdateCache[key] = cache
		contestRanklistUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q contestRanklistQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for contest_ranklist")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for contest_ranklist")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ContestRanklistSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contestRanklistPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `contest_ranklist` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, contestRanklistPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in contestRanklist slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all contestRanklist")
	}
	return rowsAff, nil
}

var mySQLContestRanklistUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *ContestRanklist) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no contest_ranklist provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(contestRanklistColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLContestRanklistUniqueColumns, o)

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

	contestRanklistUpsertCacheMut.RLock()
	cache, cached := contestRanklistUpsertCache[key]
	contestRanklistUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			contestRanklistAllColumns,
			contestRanklistColumnsWithDefault,
			contestRanklistColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			contestRanklistAllColumns,
			contestRanklistPrimaryKeyColumns,
		)

		if len(update) == 0 {
			return errors.New("models: unable to upsert contest_ranklist, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "contest_ranklist", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `contest_ranklist` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(contestRanklistType, contestRanklistMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(contestRanklistType, contestRanklistMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert for contest_ranklist")
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
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == contestRanklistMapping["id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(contestRanklistType, contestRanklistMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "models: unable to retrieve unique values for contest_ranklist")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, nzUniqueCols...)
	}

	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for contest_ranklist")
	}

CacheNoHooks:
	if !cached {
		contestRanklistUpsertCacheMut.Lock()
		contestRanklistUpsertCache[key] = cache
		contestRanklistUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single ContestRanklist record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ContestRanklist) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no ContestRanklist provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), contestRanklistPrimaryKeyMapping)
	sql := "DELETE FROM `contest_ranklist` WHERE `id`=?"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from contest_ranklist")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for contest_ranklist")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q contestRanklistQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no contestRanklistQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from contest_ranklist")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for contest_ranklist")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ContestRanklistSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(contestRanklistBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contestRanklistPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `contest_ranklist` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, contestRanklistPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from contestRanklist slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for contest_ranklist")
	}

	if len(contestRanklistAfterDeleteHooks) != 0 {
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
func (o *ContestRanklist) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindContestRanklist(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ContestRanklistSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ContestRanklistSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contestRanklistPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `contest_ranklist`.* FROM `contest_ranklist` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, contestRanklistPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ContestRanklistSlice")
	}

	*o = slice

	return nil
}

// ContestRanklistExists checks if the ContestRanklist row exists.
func ContestRanklistExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `contest_ranklist` where `id`=? limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}

	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if contest_ranklist exists")
	}

	return exists, nil
}
