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

func testUserPrivileges(t *testing.T) {
	t.Parallel()

	query := UserPrivileges()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testUserPrivilegesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserPrivilege{}
	if err = randomize.Struct(seed, o, userPrivilegeDBTypes, true, userPrivilegeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
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

	count, err := UserPrivileges().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUserPrivilegesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserPrivilege{}
	if err = randomize.Struct(seed, o, userPrivilegeDBTypes, true, userPrivilegeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := UserPrivileges().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := UserPrivileges().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUserPrivilegesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserPrivilege{}
	if err = randomize.Struct(seed, o, userPrivilegeDBTypes, true, userPrivilegeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := UserPrivilegeSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := UserPrivileges().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testUserPrivilegesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserPrivilege{}
	if err = randomize.Struct(seed, o, userPrivilegeDBTypes, true, userPrivilegeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := UserPrivilegeExists(ctx, tx, o.UserID, o.Privilege)
	if err != nil {
		t.Errorf("Unable to check if UserPrivilege exists: %s", err)
	}
	if !e {
		t.Errorf("Expected UserPrivilegeExists to return true, but got false.")
	}
}

func testUserPrivilegesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserPrivilege{}
	if err = randomize.Struct(seed, o, userPrivilegeDBTypes, true, userPrivilegeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	userPrivilegeFound, err := FindUserPrivilege(ctx, tx, o.UserID, o.Privilege)
	if err != nil {
		t.Error(err)
	}

	if userPrivilegeFound == nil {
		t.Error("want a record, got nil")
	}
}

func testUserPrivilegesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserPrivilege{}
	if err = randomize.Struct(seed, o, userPrivilegeDBTypes, true, userPrivilegeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = UserPrivileges().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testUserPrivilegesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserPrivilege{}
	if err = randomize.Struct(seed, o, userPrivilegeDBTypes, true, userPrivilegeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := UserPrivileges().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testUserPrivilegesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	userPrivilegeOne := &UserPrivilege{}
	userPrivilegeTwo := &UserPrivilege{}
	if err = randomize.Struct(seed, userPrivilegeOne, userPrivilegeDBTypes, false, userPrivilegeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}
	if err = randomize.Struct(seed, userPrivilegeTwo, userPrivilegeDBTypes, false, userPrivilegeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = userPrivilegeOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = userPrivilegeTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := UserPrivileges().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testUserPrivilegesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	userPrivilegeOne := &UserPrivilege{}
	userPrivilegeTwo := &UserPrivilege{}
	if err = randomize.Struct(seed, userPrivilegeOne, userPrivilegeDBTypes, false, userPrivilegeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}
	if err = randomize.Struct(seed, userPrivilegeTwo, userPrivilegeDBTypes, false, userPrivilegeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = userPrivilegeOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = userPrivilegeTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UserPrivileges().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func userPrivilegeBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *UserPrivilege) error {
	*o = UserPrivilege{}
	return nil
}

func userPrivilegeAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *UserPrivilege) error {
	*o = UserPrivilege{}
	return nil
}

func userPrivilegeAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *UserPrivilege) error {
	*o = UserPrivilege{}
	return nil
}

func userPrivilegeBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *UserPrivilege) error {
	*o = UserPrivilege{}
	return nil
}

func userPrivilegeAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *UserPrivilege) error {
	*o = UserPrivilege{}
	return nil
}

func userPrivilegeBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *UserPrivilege) error {
	*o = UserPrivilege{}
	return nil
}

func userPrivilegeAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *UserPrivilege) error {
	*o = UserPrivilege{}
	return nil
}

func userPrivilegeBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *UserPrivilege) error {
	*o = UserPrivilege{}
	return nil
}

func userPrivilegeAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *UserPrivilege) error {
	*o = UserPrivilege{}
	return nil
}

func testUserPrivilegesHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &UserPrivilege{}
	o := &UserPrivilege{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, userPrivilegeDBTypes, false); err != nil {
		t.Errorf("Unable to randomize UserPrivilege object: %s", err)
	}

	AddUserPrivilegeHook(boil.BeforeInsertHook, userPrivilegeBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	userPrivilegeBeforeInsertHooks = []UserPrivilegeHook{}

	AddUserPrivilegeHook(boil.AfterInsertHook, userPrivilegeAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	userPrivilegeAfterInsertHooks = []UserPrivilegeHook{}

	AddUserPrivilegeHook(boil.AfterSelectHook, userPrivilegeAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	userPrivilegeAfterSelectHooks = []UserPrivilegeHook{}

	AddUserPrivilegeHook(boil.BeforeUpdateHook, userPrivilegeBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	userPrivilegeBeforeUpdateHooks = []UserPrivilegeHook{}

	AddUserPrivilegeHook(boil.AfterUpdateHook, userPrivilegeAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	userPrivilegeAfterUpdateHooks = []UserPrivilegeHook{}

	AddUserPrivilegeHook(boil.BeforeDeleteHook, userPrivilegeBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	userPrivilegeBeforeDeleteHooks = []UserPrivilegeHook{}

	AddUserPrivilegeHook(boil.AfterDeleteHook, userPrivilegeAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	userPrivilegeAfterDeleteHooks = []UserPrivilegeHook{}

	AddUserPrivilegeHook(boil.BeforeUpsertHook, userPrivilegeBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	userPrivilegeBeforeUpsertHooks = []UserPrivilegeHook{}

	AddUserPrivilegeHook(boil.AfterUpsertHook, userPrivilegeAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	userPrivilegeAfterUpsertHooks = []UserPrivilegeHook{}
}

func testUserPrivilegesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserPrivilege{}
	if err = randomize.Struct(seed, o, userPrivilegeDBTypes, true, userPrivilegeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UserPrivileges().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testUserPrivilegesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserPrivilege{}
	if err = randomize.Struct(seed, o, userPrivilegeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(userPrivilegeColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := UserPrivileges().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testUserPrivilegesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserPrivilege{}
	if err = randomize.Struct(seed, o, userPrivilegeDBTypes, true, userPrivilegeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
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

func testUserPrivilegesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserPrivilege{}
	if err = randomize.Struct(seed, o, userPrivilegeDBTypes, true, userPrivilegeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := UserPrivilegeSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testUserPrivilegesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &UserPrivilege{}
	if err = randomize.Struct(seed, o, userPrivilegeDBTypes, true, userPrivilegeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := UserPrivileges().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	userPrivilegeDBTypes = map[string]string{`UserID`: `int`, `Privilege`: `varchar`}
	_                    = bytes.MinRead
)

func testUserPrivilegesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(userPrivilegePrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(userPrivilegeAllColumns) == len(userPrivilegePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &UserPrivilege{}
	if err = randomize.Struct(seed, o, userPrivilegeDBTypes, true, userPrivilegeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UserPrivileges().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, userPrivilegeDBTypes, true, userPrivilegePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testUserPrivilegesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(userPrivilegeAllColumns) == len(userPrivilegePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &UserPrivilege{}
	if err = randomize.Struct(seed, o, userPrivilegeDBTypes, true, userPrivilegeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := UserPrivileges().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, userPrivilegeDBTypes, true, userPrivilegePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(userPrivilegeAllColumns, userPrivilegePrimaryKeyColumns) {
		fields = userPrivilegeAllColumns
	} else {
		fields = strmangle.SetComplement(
			userPrivilegeAllColumns,
			userPrivilegePrimaryKeyColumns,
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

	slice := UserPrivilegeSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testUserPrivilegesUpsert(t *testing.T) {
	t.Parallel()

	if len(userPrivilegeAllColumns) == len(userPrivilegePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}
	if len(mySQLUserPrivilegeUniqueColumns) == 0 {
		t.Skip("Skipping table with no unique columns to conflict on")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := UserPrivilege{}
	if err = randomize.Struct(seed, &o, userPrivilegeDBTypes, false); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert UserPrivilege: %s", err)
	}

	count, err := UserPrivileges().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, userPrivilegeDBTypes, false, userPrivilegePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize UserPrivilege struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert UserPrivilege: %s", err)
	}

	count, err = UserPrivileges().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
