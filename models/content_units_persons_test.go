package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testContentUnitsPersons(t *testing.T) {
	t.Parallel()

	query := ContentUnitsPersons(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testContentUnitsPersonsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnitsPerson := &ContentUnitsPerson{}
	if err = randomize.Struct(seed, contentUnitsPerson, contentUnitsPersonDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitsPerson.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = contentUnitsPerson.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentUnitsPersons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testContentUnitsPersonsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnitsPerson := &ContentUnitsPerson{}
	if err = randomize.Struct(seed, contentUnitsPerson, contentUnitsPersonDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitsPerson.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = ContentUnitsPersons(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := ContentUnitsPersons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testContentUnitsPersonsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnitsPerson := &ContentUnitsPerson{}
	if err = randomize.Struct(seed, contentUnitsPerson, contentUnitsPersonDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitsPerson.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ContentUnitsPersonSlice{contentUnitsPerson}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentUnitsPersons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testContentUnitsPersonsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnitsPerson := &ContentUnitsPerson{}
	if err = randomize.Struct(seed, contentUnitsPerson, contentUnitsPersonDBTypes, true, contentUnitsPersonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitsPerson.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := ContentUnitsPersonExists(tx, contentUnitsPerson.ContentUnitID, contentUnitsPerson.PersonID)
	if err != nil {
		t.Errorf("Unable to check if ContentUnitsPerson exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ContentUnitsPersonExistsG to return true, but got false.")
	}
}
func testContentUnitsPersonsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnitsPerson := &ContentUnitsPerson{}
	if err = randomize.Struct(seed, contentUnitsPerson, contentUnitsPersonDBTypes, true, contentUnitsPersonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitsPerson.Insert(tx); err != nil {
		t.Error(err)
	}

	contentUnitsPersonFound, err := FindContentUnitsPerson(tx, contentUnitsPerson.ContentUnitID, contentUnitsPerson.PersonID)
	if err != nil {
		t.Error(err)
	}

	if contentUnitsPersonFound == nil {
		t.Error("want a record, got nil")
	}
}
func testContentUnitsPersonsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnitsPerson := &ContentUnitsPerson{}
	if err = randomize.Struct(seed, contentUnitsPerson, contentUnitsPersonDBTypes, true, contentUnitsPersonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitsPerson.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = ContentUnitsPersons(tx).Bind(contentUnitsPerson); err != nil {
		t.Error(err)
	}
}

func testContentUnitsPersonsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnitsPerson := &ContentUnitsPerson{}
	if err = randomize.Struct(seed, contentUnitsPerson, contentUnitsPersonDBTypes, true, contentUnitsPersonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitsPerson.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := ContentUnitsPersons(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testContentUnitsPersonsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnitsPersonOne := &ContentUnitsPerson{}
	contentUnitsPersonTwo := &ContentUnitsPerson{}
	if err = randomize.Struct(seed, contentUnitsPersonOne, contentUnitsPersonDBTypes, false, contentUnitsPersonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}
	if err = randomize.Struct(seed, contentUnitsPersonTwo, contentUnitsPersonDBTypes, false, contentUnitsPersonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitsPersonOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = contentUnitsPersonTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := ContentUnitsPersons(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testContentUnitsPersonsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	contentUnitsPersonOne := &ContentUnitsPerson{}
	contentUnitsPersonTwo := &ContentUnitsPerson{}
	if err = randomize.Struct(seed, contentUnitsPersonOne, contentUnitsPersonDBTypes, false, contentUnitsPersonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}
	if err = randomize.Struct(seed, contentUnitsPersonTwo, contentUnitsPersonDBTypes, false, contentUnitsPersonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitsPersonOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = contentUnitsPersonTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentUnitsPersons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testContentUnitsPersonsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnitsPerson := &ContentUnitsPerson{}
	if err = randomize.Struct(seed, contentUnitsPerson, contentUnitsPersonDBTypes, true, contentUnitsPersonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitsPerson.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentUnitsPersons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testContentUnitsPersonsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnitsPerson := &ContentUnitsPerson{}
	if err = randomize.Struct(seed, contentUnitsPerson, contentUnitsPersonDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitsPerson.Insert(tx, contentUnitsPersonColumns...); err != nil {
		t.Error(err)
	}

	count, err := ContentUnitsPersons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testContentUnitsPersonToOneContentUnitUsingContentUnit(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local ContentUnitsPerson
	var foreign ContentUnit

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, contentUnitsPersonDBTypes, true, contentUnitsPersonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, contentUnitDBTypes, true, contentUnitColumnsWithDefault...); err != nil {
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

	slice := ContentUnitsPersonSlice{&local}
	if err = local.L.LoadContentUnit(tx, false, &slice); err != nil {
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

func testContentUnitsPersonToOnePersonUsingPerson(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local ContentUnitsPerson
	var foreign Person

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, contentUnitsPersonDBTypes, true, contentUnitsPersonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, personDBTypes, true, personColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Person struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.PersonID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Person(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ContentUnitsPersonSlice{&local}
	if err = local.L.LoadPerson(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Person == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Person = nil
	if err = local.L.LoadPerson(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Person == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testContentUnitsPersonToOneContentRoleUsingRole(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local ContentUnitsPerson
	var foreign ContentRole

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, contentUnitsPersonDBTypes, true, contentUnitsPersonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, contentRoleDBTypes, true, contentRoleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentRole struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.RoleID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Role(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ContentUnitsPersonSlice{&local}
	if err = local.L.LoadRole(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Role == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Role = nil
	if err = local.L.LoadRole(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Role == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testContentUnitsPersonToOneSetOpContentUnitUsingContentUnit(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentUnitsPerson
	var b, c ContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentUnitsPersonDBTypes, false, strmangle.SetComplement(contentUnitsPersonPrimaryKeyColumns, contentUnitsPersonColumnsWithoutDefault)...); err != nil {
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

		if x.R.ContentUnitsPersons[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.ContentUnitID != x.ID {
			t.Error("foreign key was wrong value", a.ContentUnitID)
		}

		if exists, err := ContentUnitsPersonExists(tx, a.ContentUnitID, a.PersonID); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}

	}
}
func testContentUnitsPersonToOneSetOpPersonUsingPerson(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentUnitsPerson
	var b, c Person

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentUnitsPersonDBTypes, false, strmangle.SetComplement(contentUnitsPersonPrimaryKeyColumns, contentUnitsPersonColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, personDBTypes, false, strmangle.SetComplement(personPrimaryKeyColumns, personColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, personDBTypes, false, strmangle.SetComplement(personPrimaryKeyColumns, personColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Person{&b, &c} {
		err = a.SetPerson(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Person != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.ContentUnitsPersons[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.PersonID != x.ID {
			t.Error("foreign key was wrong value", a.PersonID)
		}

		if exists, err := ContentUnitsPersonExists(tx, a.ContentUnitID, a.PersonID); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}

	}
}
func testContentUnitsPersonToOneSetOpContentRoleUsingRole(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentUnitsPerson
	var b, c ContentRole

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentUnitsPersonDBTypes, false, strmangle.SetComplement(contentUnitsPersonPrimaryKeyColumns, contentUnitsPersonColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, contentRoleDBTypes, false, strmangle.SetComplement(contentRolePrimaryKeyColumns, contentRoleColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, contentRoleDBTypes, false, strmangle.SetComplement(contentRolePrimaryKeyColumns, contentRoleColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*ContentRole{&b, &c} {
		err = a.SetRole(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Role != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.RoleContentUnitsPersons[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.RoleID != x.ID {
			t.Error("foreign key was wrong value", a.RoleID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.RoleID))
		reflect.Indirect(reflect.ValueOf(&a.RoleID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.RoleID != x.ID {
			t.Error("foreign key was wrong value", a.RoleID, x.ID)
		}
	}
}
func testContentUnitsPersonsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnitsPerson := &ContentUnitsPerson{}
	if err = randomize.Struct(seed, contentUnitsPerson, contentUnitsPersonDBTypes, true, contentUnitsPersonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitsPerson.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = contentUnitsPerson.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testContentUnitsPersonsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnitsPerson := &ContentUnitsPerson{}
	if err = randomize.Struct(seed, contentUnitsPerson, contentUnitsPersonDBTypes, true, contentUnitsPersonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitsPerson.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ContentUnitsPersonSlice{contentUnitsPerson}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testContentUnitsPersonsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnitsPerson := &ContentUnitsPerson{}
	if err = randomize.Struct(seed, contentUnitsPerson, contentUnitsPersonDBTypes, true, contentUnitsPersonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitsPerson.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := ContentUnitsPersons(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	contentUnitsPersonDBTypes = map[string]string{`ContentUnitID`: `bigint`, `PersonID`: `bigint`, `RoleID`: `bigint`}
	_                         = bytes.MinRead
)

func testContentUnitsPersonsUpdate(t *testing.T) {
	t.Parallel()

	if len(contentUnitsPersonColumns) == len(contentUnitsPersonPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	contentUnitsPerson := &ContentUnitsPerson{}
	if err = randomize.Struct(seed, contentUnitsPerson, contentUnitsPersonDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitsPerson.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentUnitsPersons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, contentUnitsPerson, contentUnitsPersonDBTypes, true, contentUnitsPersonColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	if err = contentUnitsPerson.Update(tx); err != nil {
		t.Error(err)
	}
}

func testContentUnitsPersonsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(contentUnitsPersonColumns) == len(contentUnitsPersonPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	contentUnitsPerson := &ContentUnitsPerson{}
	if err = randomize.Struct(seed, contentUnitsPerson, contentUnitsPersonDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitsPerson.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentUnitsPersons(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, contentUnitsPerson, contentUnitsPersonDBTypes, true, contentUnitsPersonPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(contentUnitsPersonColumns, contentUnitsPersonPrimaryKeyColumns) {
		fields = contentUnitsPersonColumns
	} else {
		fields = strmangle.SetComplement(
			contentUnitsPersonColumns,
			contentUnitsPersonPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(contentUnitsPerson))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := ContentUnitsPersonSlice{contentUnitsPerson}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testContentUnitsPersonsUpsert(t *testing.T) {
	t.Parallel()

	if len(contentUnitsPersonColumns) == len(contentUnitsPersonPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	contentUnitsPerson := ContentUnitsPerson{}
	if err = randomize.Struct(seed, &contentUnitsPerson, contentUnitsPersonDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitsPerson.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert ContentUnitsPerson: %s", err)
	}

	count, err := ContentUnitsPersons(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &contentUnitsPerson, contentUnitsPersonDBTypes, false, contentUnitsPersonPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ContentUnitsPerson struct: %s", err)
	}

	if err = contentUnitsPerson.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert ContentUnitsPerson: %s", err)
	}

	count, err = ContentUnitsPersons(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
