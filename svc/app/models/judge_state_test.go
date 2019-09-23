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

func testJudgeStates(t *testing.T) {
	t.Parallel()

	query := JudgeStates()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testJudgeStatesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &JudgeState{}
	if err = randomize.Struct(seed, o, judgeStateDBTypes, true, judgeStateColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
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

	count, err := JudgeStates().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testJudgeStatesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &JudgeState{}
	if err = randomize.Struct(seed, o, judgeStateDBTypes, true, judgeStateColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := JudgeStates().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := JudgeStates().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testJudgeStatesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &JudgeState{}
	if err = randomize.Struct(seed, o, judgeStateDBTypes, true, judgeStateColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := JudgeStateSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := JudgeStates().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testJudgeStatesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &JudgeState{}
	if err = randomize.Struct(seed, o, judgeStateDBTypes, true, judgeStateColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := JudgeStateExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if JudgeState exists: %s", err)
	}
	if !e {
		t.Errorf("Expected JudgeStateExists to return true, but got false.")
	}
}

func testJudgeStatesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &JudgeState{}
	if err = randomize.Struct(seed, o, judgeStateDBTypes, true, judgeStateColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	judgeStateFound, err := FindJudgeState(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if judgeStateFound == nil {
		t.Error("want a record, got nil")
	}
}

func testJudgeStatesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &JudgeState{}
	if err = randomize.Struct(seed, o, judgeStateDBTypes, true, judgeStateColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = JudgeStates().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testJudgeStatesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &JudgeState{}
	if err = randomize.Struct(seed, o, judgeStateDBTypes, true, judgeStateColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := JudgeStates().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testJudgeStatesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	judgeStateOne := &JudgeState{}
	judgeStateTwo := &JudgeState{}
	if err = randomize.Struct(seed, judgeStateOne, judgeStateDBTypes, false, judgeStateColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}
	if err = randomize.Struct(seed, judgeStateTwo, judgeStateDBTypes, false, judgeStateColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = judgeStateOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = judgeStateTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := JudgeStates().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testJudgeStatesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	judgeStateOne := &JudgeState{}
	judgeStateTwo := &JudgeState{}
	if err = randomize.Struct(seed, judgeStateOne, judgeStateDBTypes, false, judgeStateColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}
	if err = randomize.Struct(seed, judgeStateTwo, judgeStateDBTypes, false, judgeStateColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = judgeStateOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = judgeStateTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := JudgeStates().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func judgeStateBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *JudgeState) error {
	*o = JudgeState{}
	return nil
}

func judgeStateAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *JudgeState) error {
	*o = JudgeState{}
	return nil
}

func judgeStateAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *JudgeState) error {
	*o = JudgeState{}
	return nil
}

func judgeStateBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *JudgeState) error {
	*o = JudgeState{}
	return nil
}

func judgeStateAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *JudgeState) error {
	*o = JudgeState{}
	return nil
}

func judgeStateBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *JudgeState) error {
	*o = JudgeState{}
	return nil
}

func judgeStateAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *JudgeState) error {
	*o = JudgeState{}
	return nil
}

func judgeStateBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *JudgeState) error {
	*o = JudgeState{}
	return nil
}

func judgeStateAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *JudgeState) error {
	*o = JudgeState{}
	return nil
}

func testJudgeStatesHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &JudgeState{}
	o := &JudgeState{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, judgeStateDBTypes, false); err != nil {
		t.Errorf("Unable to randomize JudgeState object: %s", err)
	}

	AddJudgeStateHook(boil.BeforeInsertHook, judgeStateBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	judgeStateBeforeInsertHooks = []JudgeStateHook{}

	AddJudgeStateHook(boil.AfterInsertHook, judgeStateAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	judgeStateAfterInsertHooks = []JudgeStateHook{}

	AddJudgeStateHook(boil.AfterSelectHook, judgeStateAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	judgeStateAfterSelectHooks = []JudgeStateHook{}

	AddJudgeStateHook(boil.BeforeUpdateHook, judgeStateBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	judgeStateBeforeUpdateHooks = []JudgeStateHook{}

	AddJudgeStateHook(boil.AfterUpdateHook, judgeStateAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	judgeStateAfterUpdateHooks = []JudgeStateHook{}

	AddJudgeStateHook(boil.BeforeDeleteHook, judgeStateBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	judgeStateBeforeDeleteHooks = []JudgeStateHook{}

	AddJudgeStateHook(boil.AfterDeleteHook, judgeStateAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	judgeStateAfterDeleteHooks = []JudgeStateHook{}

	AddJudgeStateHook(boil.BeforeUpsertHook, judgeStateBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	judgeStateBeforeUpsertHooks = []JudgeStateHook{}

	AddJudgeStateHook(boil.AfterUpsertHook, judgeStateAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	judgeStateAfterUpsertHooks = []JudgeStateHook{}
}

func testJudgeStatesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &JudgeState{}
	if err = randomize.Struct(seed, o, judgeStateDBTypes, true, judgeStateColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := JudgeStates().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testJudgeStatesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &JudgeState{}
	if err = randomize.Struct(seed, o, judgeStateDBTypes, true); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(judgeStateColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := JudgeStates().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testJudgeStatesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &JudgeState{}
	if err = randomize.Struct(seed, o, judgeStateDBTypes, true, judgeStateColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
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

func testJudgeStatesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &JudgeState{}
	if err = randomize.Struct(seed, o, judgeStateDBTypes, true, judgeStateColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := JudgeStateSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testJudgeStatesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &JudgeState{}
	if err = randomize.Struct(seed, o, judgeStateDBTypes, true, judgeStateColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := JudgeStates().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	judgeStateDBTypes = map[string]string{`ID`: `int`, `Code`: `mediumtext`, `Language`: `varchar`, `Status`: `enum('Accepted','Compile Error','File Error','Invalid Interaction','Judgement Failed','Memory Limit Exceeded','No Testdata','Output Limit Exceeded','Partially Correct','Runtime Error','System Error','Time Limit Exceeded','Unknown','Wrong Answer','Waiting')`, `TaskID`: `varchar`, `Score`: `int`, `TotalTime`: `int`, `CodeLength`: `int`, `Pending`: `tinyint`, `MaxMemory`: `int`, `Compilation`: `longtext`, `Result`: `longtext`, `UserID`: `int`, `ProblemID`: `int`, `SubmitTime`: `int`, `Type`: `int`, `TypeInfo`: `int`, `IsPublic`: `tinyint`}
	_                 = bytes.MinRead
)

func testJudgeStatesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(judgeStatePrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(judgeStateAllColumns) == len(judgeStatePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &JudgeState{}
	if err = randomize.Struct(seed, o, judgeStateDBTypes, true, judgeStateColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := JudgeStates().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, judgeStateDBTypes, true, judgeStatePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testJudgeStatesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(judgeStateAllColumns) == len(judgeStatePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &JudgeState{}
	if err = randomize.Struct(seed, o, judgeStateDBTypes, true, judgeStateColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := JudgeStates().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, judgeStateDBTypes, true, judgeStatePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(judgeStateAllColumns, judgeStatePrimaryKeyColumns) {
		fields = judgeStateAllColumns
	} else {
		fields = strmangle.SetComplement(
			judgeStateAllColumns,
			judgeStatePrimaryKeyColumns,
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

	slice := JudgeStateSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testJudgeStatesUpsert(t *testing.T) {
	t.Parallel()

	if len(judgeStateAllColumns) == len(judgeStatePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}
	if len(mySQLJudgeStateUniqueColumns) == 0 {
		t.Skip("Skipping table with no unique columns to conflict on")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := JudgeState{}
	if err = randomize.Struct(seed, &o, judgeStateDBTypes, false); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert JudgeState: %s", err)
	}

	count, err := JudgeStates().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, judgeStateDBTypes, false, judgeStatePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize JudgeState struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert JudgeState: %s", err)
	}

	count, err = JudgeStates().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
