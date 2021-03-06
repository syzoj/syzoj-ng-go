// Code generated by SQLBoiler 3.5.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/randomize"
	"github.com/volatiletech/sqlboiler/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testProblems(t *testing.T) {
	t.Parallel()

	query := Problems()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testProblemsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Problem{}
	if err = randomize.Struct(seed, o, problemDBTypes, true, problemColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Problems().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testProblemsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Problem{}
	if err = randomize.Struct(seed, o, problemDBTypes, true, problemColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Problems().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Problems().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testProblemsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Problem{}
	if err = randomize.Struct(seed, o, problemDBTypes, true, problemColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ProblemSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Problems().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testProblemsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Problem{}
	if err = randomize.Struct(seed, o, problemDBTypes, true, problemColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := ProblemExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if Problem exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ProblemExists to return true, but got false.")
	}
}

func testProblemsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Problem{}
	if err = randomize.Struct(seed, o, problemDBTypes, true, problemColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	problemFound, err := FindProblem(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if problemFound == nil {
		t.Error("want a record, got nil")
	}
}

func testProblemsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Problem{}
	if err = randomize.Struct(seed, o, problemDBTypes, true, problemColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Problems().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testProblemsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Problem{}
	if err = randomize.Struct(seed, o, problemDBTypes, true, problemColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Problems().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testProblemsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	problemOne := &Problem{}
	problemTwo := &Problem{}
	if err = randomize.Struct(seed, problemOne, problemDBTypes, false, problemColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}
	if err = randomize.Struct(seed, problemTwo, problemDBTypes, false, problemColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = problemOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = problemTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Problems().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testProblemsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	problemOne := &Problem{}
	problemTwo := &Problem{}
	if err = randomize.Struct(seed, problemOne, problemDBTypes, false, problemColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}
	if err = randomize.Struct(seed, problemTwo, problemDBTypes, false, problemColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = problemOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = problemTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Problems().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func problemBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *Problem) error {
	*o = Problem{}
	return nil
}

func problemAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *Problem) error {
	*o = Problem{}
	return nil
}

func problemAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *Problem) error {
	*o = Problem{}
	return nil
}

func problemBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Problem) error {
	*o = Problem{}
	return nil
}

func problemAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Problem) error {
	*o = Problem{}
	return nil
}

func problemBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Problem) error {
	*o = Problem{}
	return nil
}

func problemAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Problem) error {
	*o = Problem{}
	return nil
}

func problemBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Problem) error {
	*o = Problem{}
	return nil
}

func problemAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Problem) error {
	*o = Problem{}
	return nil
}

func testProblemsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &Problem{}
	o := &Problem{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, problemDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Problem object: %s", err)
	}

	AddProblemHook(boil.BeforeInsertHook, problemBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	problemBeforeInsertHooks = []ProblemHook{}

	AddProblemHook(boil.AfterInsertHook, problemAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	problemAfterInsertHooks = []ProblemHook{}

	AddProblemHook(boil.AfterSelectHook, problemAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	problemAfterSelectHooks = []ProblemHook{}

	AddProblemHook(boil.BeforeUpdateHook, problemBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	problemBeforeUpdateHooks = []ProblemHook{}

	AddProblemHook(boil.AfterUpdateHook, problemAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	problemAfterUpdateHooks = []ProblemHook{}

	AddProblemHook(boil.BeforeDeleteHook, problemBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	problemBeforeDeleteHooks = []ProblemHook{}

	AddProblemHook(boil.AfterDeleteHook, problemAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	problemAfterDeleteHooks = []ProblemHook{}

	AddProblemHook(boil.BeforeUpsertHook, problemBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	problemBeforeUpsertHooks = []ProblemHook{}

	AddProblemHook(boil.AfterUpsertHook, problemAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	problemAfterUpsertHooks = []ProblemHook{}
}

func testProblemsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Problem{}
	if err = randomize.Struct(seed, o, problemDBTypes, true, problemColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Problems().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testProblemsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Problem{}
	if err = randomize.Struct(seed, o, problemDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(problemColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Problems().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testProblemsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Problem{}
	if err = randomize.Struct(seed, o, problemDBTypes, true, problemColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testProblemsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Problem{}
	if err = randomize.Struct(seed, o, problemDBTypes, true, problemColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ProblemSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testProblemsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Problem{}
	if err = randomize.Struct(seed, o, problemDBTypes, true, problemColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Problems().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	problemDBTypes = map[string]string{`ID`: `int`, `Title`: `varchar`, `UserID`: `int`, `PublicizerID`: `int`, `IsAnonymous`: `tinyint`, `Description`: `text`, `InputFormat`: `text`, `OutputFormat`: `text`, `Example`: `text`, `LimitAndHint`: `text`, `TimeLimit`: `int`, `MemoryLimit`: `int`, `AdditionalFileID`: `int`, `AcNum`: `int`, `SubmitNum`: `int`, `IsPublic`: `tinyint`, `FileIo`: `tinyint`, `FileIoInputName`: `text`, `FileIoOutputName`: `text`, `PublicizeTime`: `datetime`, `Type`: `enum('traditional','submit-answer','interaction')`, `Tags`: `longtext`}
	_              = bytes.MinRead
)

func testProblemsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(problemPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(problemAllColumns) == len(problemPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Problem{}
	if err = randomize.Struct(seed, o, problemDBTypes, true, problemColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Problems().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, problemDBTypes, true, problemPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testProblemsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(problemAllColumns) == len(problemPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Problem{}
	if err = randomize.Struct(seed, o, problemDBTypes, true, problemColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Problems().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, problemDBTypes, true, problemPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(problemAllColumns, problemPrimaryKeyColumns) {
		fields = problemAllColumns
	} else {
		fields = strmangle.SetComplement(
			problemAllColumns,
			problemPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := ProblemSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testProblemsUpsert(t *testing.T) {
	t.Parallel()

	if len(problemAllColumns) == len(problemPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}
	if len(mySQLProblemUniqueColumns) == 0 {
		t.Skip("Skipping table with no unique columns to conflict on")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Problem{}
	if err = randomize.Struct(seed, &o, problemDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Problem: %s", err)
	}

	count, err := Problems().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, problemDBTypes, false, problemPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Problem struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Problem: %s", err)
	}

	count, err = Problems().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
