package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testSourceTypes(t *testing.T) {
	t.Parallel()

	query := SourceTypes(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testSourceTypesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceType := &SourceType{}
	if err = randomize.Struct(seed, sourceType, sourceTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceType.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = sourceType.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := SourceTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testSourceTypesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceType := &SourceType{}
	if err = randomize.Struct(seed, sourceType, sourceTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceType.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = SourceTypes(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := SourceTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testSourceTypesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceType := &SourceType{}
	if err = randomize.Struct(seed, sourceType, sourceTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceType.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := SourceTypeSlice{sourceType}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := SourceTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testSourceTypesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceType := &SourceType{}
	if err = randomize.Struct(seed, sourceType, sourceTypeDBTypes, true, sourceTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceType.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := SourceTypeExists(tx, sourceType.ID)
	if err != nil {
		t.Errorf("Unable to check if SourceType exists: %s", err)
	}
	if !e {
		t.Errorf("Expected SourceTypeExistsG to return true, but got false.")
	}
}
func testSourceTypesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceType := &SourceType{}
	if err = randomize.Struct(seed, sourceType, sourceTypeDBTypes, true, sourceTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceType.Insert(tx); err != nil {
		t.Error(err)
	}

	sourceTypeFound, err := FindSourceType(tx, sourceType.ID)
	if err != nil {
		t.Error(err)
	}

	if sourceTypeFound == nil {
		t.Error("want a record, got nil")
	}
}
func testSourceTypesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceType := &SourceType{}
	if err = randomize.Struct(seed, sourceType, sourceTypeDBTypes, true, sourceTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceType.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = SourceTypes(tx).Bind(sourceType); err != nil {
		t.Error(err)
	}
}

func testSourceTypesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceType := &SourceType{}
	if err = randomize.Struct(seed, sourceType, sourceTypeDBTypes, true, sourceTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceType.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := SourceTypes(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testSourceTypesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceTypeOne := &SourceType{}
	sourceTypeTwo := &SourceType{}
	if err = randomize.Struct(seed, sourceTypeOne, sourceTypeDBTypes, false, sourceTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}
	if err = randomize.Struct(seed, sourceTypeTwo, sourceTypeDBTypes, false, sourceTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceTypeOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = sourceTypeTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := SourceTypes(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testSourceTypesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	sourceTypeOne := &SourceType{}
	sourceTypeTwo := &SourceType{}
	if err = randomize.Struct(seed, sourceTypeOne, sourceTypeDBTypes, false, sourceTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}
	if err = randomize.Struct(seed, sourceTypeTwo, sourceTypeDBTypes, false, sourceTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceTypeOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = sourceTypeTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := SourceTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testSourceTypesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceType := &SourceType{}
	if err = randomize.Struct(seed, sourceType, sourceTypeDBTypes, true, sourceTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceType.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := SourceTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testSourceTypesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceType := &SourceType{}
	if err = randomize.Struct(seed, sourceType, sourceTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceType.Insert(tx, sourceTypeColumns...); err != nil {
		t.Error(err)
	}

	count, err := SourceTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testSourceTypeToManyTypeSources(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a SourceType
	var b, c Source

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceTypeDBTypes, true, sourceTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, sourceDBTypes, false, sourceColumnsWithDefault...)
	randomize.Struct(seed, &c, sourceDBTypes, false, sourceColumnsWithDefault...)

	b.TypeID = a.ID
	c.TypeID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	source, err := a.TypeSources(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range source {
		if v.TypeID == b.TypeID {
			bFound = true
		}
		if v.TypeID == c.TypeID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := SourceTypeSlice{&a}
	if err = a.L.LoadTypeSources(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TypeSources); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.TypeSources = nil
	if err = a.L.LoadTypeSources(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.TypeSources); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", source)
	}
}

func testSourceTypeToManyAddOpTypeSources(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a SourceType
	var b, c, d, e Source

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, sourceTypeDBTypes, false, strmangle.SetComplement(sourceTypePrimaryKeyColumns, sourceTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Source{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, sourceDBTypes, false, strmangle.SetComplement(sourcePrimaryKeyColumns, sourceColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*Source{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddTypeSources(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.TypeID {
			t.Error("foreign key was wrong value", a.ID, first.TypeID)
		}
		if a.ID != second.TypeID {
			t.Error("foreign key was wrong value", a.ID, second.TypeID)
		}

		if first.R.Type != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Type != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.TypeSources[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.TypeSources[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.TypeSources(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testSourceTypesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceType := &SourceType{}
	if err = randomize.Struct(seed, sourceType, sourceTypeDBTypes, true, sourceTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceType.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = sourceType.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testSourceTypesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceType := &SourceType{}
	if err = randomize.Struct(seed, sourceType, sourceTypeDBTypes, true, sourceTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceType.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := SourceTypeSlice{sourceType}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testSourceTypesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	sourceType := &SourceType{}
	if err = randomize.Struct(seed, sourceType, sourceTypeDBTypes, true, sourceTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceType.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := SourceTypes(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	sourceTypeDBTypes = map[string]string{`ID`: `bigint`, `Name`: `character varying`}
	_                 = bytes.MinRead
)

func testSourceTypesUpdate(t *testing.T) {
	t.Parallel()

	if len(sourceTypeColumns) == len(sourceTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	sourceType := &SourceType{}
	if err = randomize.Struct(seed, sourceType, sourceTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceType.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := SourceTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, sourceType, sourceTypeDBTypes, true, sourceTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	if err = sourceType.Update(tx); err != nil {
		t.Error(err)
	}
}

func testSourceTypesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(sourceTypeColumns) == len(sourceTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	sourceType := &SourceType{}
	if err = randomize.Struct(seed, sourceType, sourceTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceType.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := SourceTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, sourceType, sourceTypeDBTypes, true, sourceTypePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(sourceTypeColumns, sourceTypePrimaryKeyColumns) {
		fields = sourceTypeColumns
	} else {
		fields = strmangle.SetComplement(
			sourceTypeColumns,
			sourceTypePrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(sourceType))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := SourceTypeSlice{sourceType}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testSourceTypesUpsert(t *testing.T) {
	t.Parallel()

	if len(sourceTypeColumns) == len(sourceTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	sourceType := SourceType{}
	if err = randomize.Struct(seed, &sourceType, sourceTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = sourceType.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert SourceType: %s", err)
	}

	count, err := SourceTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &sourceType, sourceTypeDBTypes, false, sourceTypePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize SourceType struct: %s", err)
	}

	if err = sourceType.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert SourceType: %s", err)
	}

	count, err = SourceTypes(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
