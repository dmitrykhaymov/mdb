// This file is generated by SQLBoiler (https://github.com/vattle/sqlboiler)
// and is meant to be re-generated in place and/or deleted at any time.
// DO NOT EDIT

package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testCollectionsContentUnits(t *testing.T) {
	t.Parallel()

	query := CollectionsContentUnits(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testCollectionsContentUnitsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionsContentUnit := &CollectionsContentUnit{}
	if err = randomize.Struct(seed, collectionsContentUnit, collectionsContentUnitDBTypes, true, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionsContentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = collectionsContentUnit.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := CollectionsContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCollectionsContentUnitsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionsContentUnit := &CollectionsContentUnit{}
	if err = randomize.Struct(seed, collectionsContentUnit, collectionsContentUnitDBTypes, true, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionsContentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = CollectionsContentUnits(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := CollectionsContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCollectionsContentUnitsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionsContentUnit := &CollectionsContentUnit{}
	if err = randomize.Struct(seed, collectionsContentUnit, collectionsContentUnitDBTypes, true, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionsContentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := CollectionsContentUnitSlice{collectionsContentUnit}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := CollectionsContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testCollectionsContentUnitsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionsContentUnit := &CollectionsContentUnit{}
	if err = randomize.Struct(seed, collectionsContentUnit, collectionsContentUnitDBTypes, true, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionsContentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := CollectionsContentUnitExists(tx, collectionsContentUnit.CollectionID, collectionsContentUnit.ContentUnitID)
	if err != nil {
		t.Errorf("Unable to check if CollectionsContentUnit exists: %s", err)
	}
	if !e {
		t.Errorf("Expected CollectionsContentUnitExistsG to return true, but got false.")
	}
}
func testCollectionsContentUnitsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionsContentUnit := &CollectionsContentUnit{}
	if err = randomize.Struct(seed, collectionsContentUnit, collectionsContentUnitDBTypes, true, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionsContentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	collectionsContentUnitFound, err := FindCollectionsContentUnit(tx, collectionsContentUnit.CollectionID, collectionsContentUnit.ContentUnitID)
	if err != nil {
		t.Error(err)
	}

	if collectionsContentUnitFound == nil {
		t.Error("want a record, got nil")
	}
}
func testCollectionsContentUnitsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionsContentUnit := &CollectionsContentUnit{}
	if err = randomize.Struct(seed, collectionsContentUnit, collectionsContentUnitDBTypes, true, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionsContentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = CollectionsContentUnits(tx).Bind(collectionsContentUnit); err != nil {
		t.Error(err)
	}
}

func testCollectionsContentUnitsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionsContentUnit := &CollectionsContentUnit{}
	if err = randomize.Struct(seed, collectionsContentUnit, collectionsContentUnitDBTypes, true, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionsContentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := CollectionsContentUnits(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testCollectionsContentUnitsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionsContentUnitOne := &CollectionsContentUnit{}
	collectionsContentUnitTwo := &CollectionsContentUnit{}
	if err = randomize.Struct(seed, collectionsContentUnitOne, collectionsContentUnitDBTypes, false, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}
	if err = randomize.Struct(seed, collectionsContentUnitTwo, collectionsContentUnitDBTypes, false, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionsContentUnitOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = collectionsContentUnitTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := CollectionsContentUnits(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testCollectionsContentUnitsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	collectionsContentUnitOne := &CollectionsContentUnit{}
	collectionsContentUnitTwo := &CollectionsContentUnit{}
	if err = randomize.Struct(seed, collectionsContentUnitOne, collectionsContentUnitDBTypes, false, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}
	if err = randomize.Struct(seed, collectionsContentUnitTwo, collectionsContentUnitDBTypes, false, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionsContentUnitOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = collectionsContentUnitTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := CollectionsContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testCollectionsContentUnitsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionsContentUnit := &CollectionsContentUnit{}
	if err = randomize.Struct(seed, collectionsContentUnit, collectionsContentUnitDBTypes, true, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionsContentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := CollectionsContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testCollectionsContentUnitsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionsContentUnit := &CollectionsContentUnit{}
	if err = randomize.Struct(seed, collectionsContentUnit, collectionsContentUnitDBTypes, true); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionsContentUnit.Insert(tx, collectionsContentUnitColumnsWithoutDefault...); err != nil {
		t.Error(err)
	}

	count, err := CollectionsContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testCollectionsContentUnitToOneCollectionUsingCollection(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local CollectionsContentUnit
	var foreign Collection

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, collectionsContentUnitDBTypes, false, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, collectionDBTypes, false, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.CollectionID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Collection(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := CollectionsContentUnitSlice{&local}
	if err = local.L.LoadCollection(tx, false, (*[]*CollectionsContentUnit)(&slice)); err != nil {
		t.Fatal(err)
	}
	if local.R.Collection == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Collection = nil
	if err = local.L.LoadCollection(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Collection == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testCollectionsContentUnitToOneContentUnitUsingContentUnit(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local CollectionsContentUnit
	var foreign ContentUnit

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, collectionsContentUnitDBTypes, false, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, contentUnitDBTypes, false, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.ContentUnitID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.ContentUnit(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := CollectionsContentUnitSlice{&local}
	if err = local.L.LoadContentUnit(tx, false, (*[]*CollectionsContentUnit)(&slice)); err != nil {
		t.Fatal(err)
	}
	if local.R.ContentUnit == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.ContentUnit = nil
	if err = local.L.LoadContentUnit(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.ContentUnit == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testCollectionsContentUnitToOneSetOpCollectionUsingCollection(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a CollectionsContentUnit
	var b, c Collection

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, collectionsContentUnitDBTypes, false, strmangle.SetComplement(collectionsContentUnitPrimaryKeyColumns, collectionsContentUnitColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, collectionDBTypes, false, strmangle.SetComplement(collectionPrimaryKeyColumns, collectionColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, collectionDBTypes, false, strmangle.SetComplement(collectionPrimaryKeyColumns, collectionColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Collection{&b, &c} {
		err = a.SetCollection(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Collection != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.CollectionsContentUnits[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.CollectionID != x.ID {
			t.Error("foreign key was wrong value", a.CollectionID)
		}

		if exists, err := CollectionsContentUnitExists(tx, a.CollectionID, a.ContentUnitID); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}

	}
}
func testCollectionsContentUnitToOneSetOpContentUnitUsingContentUnit(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a CollectionsContentUnit
	var b, c ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, collectionsContentUnitDBTypes, false, strmangle.SetComplement(collectionsContentUnitPrimaryKeyColumns, collectionsContentUnitColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*ContentUnit{&b, &c} {
		err = a.SetContentUnit(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.ContentUnit != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.CollectionsContentUnits[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.ContentUnitID != x.ID {
			t.Error("foreign key was wrong value", a.ContentUnitID)
		}

		if exists, err := CollectionsContentUnitExists(tx, a.CollectionID, a.ContentUnitID); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}

	}
}
func testCollectionsContentUnitsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionsContentUnit := &CollectionsContentUnit{}
	if err = randomize.Struct(seed, collectionsContentUnit, collectionsContentUnitDBTypes, true, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionsContentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = collectionsContentUnit.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testCollectionsContentUnitsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionsContentUnit := &CollectionsContentUnit{}
	if err = randomize.Struct(seed, collectionsContentUnit, collectionsContentUnitDBTypes, true, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionsContentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := CollectionsContentUnitSlice{collectionsContentUnit}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testCollectionsContentUnitsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionsContentUnit := &CollectionsContentUnit{}
	if err = randomize.Struct(seed, collectionsContentUnit, collectionsContentUnitDBTypes, true, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionsContentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := CollectionsContentUnits(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	collectionsContentUnitDBTypes = map[string]string{`CollectionID`: `bigint`, `ContentUnitID`: `bigint`, `Name`: `character varying`}
	_                             = bytes.MinRead
)

func testCollectionsContentUnitsUpdate(t *testing.T) {
	t.Parallel()

	if len(collectionsContentUnitColumns) == len(collectionsContentUnitPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	collectionsContentUnit := &CollectionsContentUnit{}
	if err = randomize.Struct(seed, collectionsContentUnit, collectionsContentUnitDBTypes, true, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionsContentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := CollectionsContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, collectionsContentUnit, collectionsContentUnitDBTypes, true, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	if err = collectionsContentUnit.Update(tx); err != nil {
		t.Error(err)
	}
}

func testCollectionsContentUnitsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(collectionsContentUnitColumns) == len(collectionsContentUnitPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	collectionsContentUnit := &CollectionsContentUnit{}
	if err = randomize.Struct(seed, collectionsContentUnit, collectionsContentUnitDBTypes, true, collectionsContentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionsContentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := CollectionsContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, collectionsContentUnit, collectionsContentUnitDBTypes, true, collectionsContentUnitPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(collectionsContentUnitColumns, collectionsContentUnitPrimaryKeyColumns) {
		fields = collectionsContentUnitColumns
	} else {
		fields = strmangle.SetComplement(
			collectionsContentUnitColumns,
			collectionsContentUnitPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(collectionsContentUnit))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := CollectionsContentUnitSlice{collectionsContentUnit}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testCollectionsContentUnitsUpsert(t *testing.T) {
	t.Parallel()

	if len(collectionsContentUnitColumns) == len(collectionsContentUnitPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	collectionsContentUnit := CollectionsContentUnit{}
	if err = randomize.Struct(seed, &collectionsContentUnit, collectionsContentUnitDBTypes, true); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionsContentUnit.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert CollectionsContentUnit: %s", err)
	}

	count, err := CollectionsContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &collectionsContentUnit, collectionsContentUnitDBTypes, false, collectionsContentUnitPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize CollectionsContentUnit struct: %s", err)
	}

	if err = collectionsContentUnit.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert CollectionsContentUnit: %s", err)
	}

	count, err = CollectionsContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}