package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testContentUnits(t *testing.T) {
	t.Parallel()

	query := ContentUnits(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testContentUnitsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnit := &ContentUnit{}
	if err = randomize.Struct(seed, contentUnit, contentUnitDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = contentUnit.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testContentUnitsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnit := &ContentUnit{}
	if err = randomize.Struct(seed, contentUnit, contentUnitDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = ContentUnits(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := ContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testContentUnitsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnit := &ContentUnit{}
	if err = randomize.Struct(seed, contentUnit, contentUnitDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ContentUnitSlice{contentUnit}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testContentUnitsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnit := &ContentUnit{}
	if err = randomize.Struct(seed, contentUnit, contentUnitDBTypes, true, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := ContentUnitExists(tx, contentUnit.ID)
	if err != nil {
		t.Errorf("Unable to check if ContentUnit exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ContentUnitExistsG to return true, but got false.")
	}
}
func testContentUnitsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnit := &ContentUnit{}
	if err = randomize.Struct(seed, contentUnit, contentUnitDBTypes, true, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	contentUnitFound, err := FindContentUnit(tx, contentUnit.ID)
	if err != nil {
		t.Error(err)
	}

	if contentUnitFound == nil {
		t.Error("want a record, got nil")
	}
}
func testContentUnitsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnit := &ContentUnit{}
	if err = randomize.Struct(seed, contentUnit, contentUnitDBTypes, true, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = ContentUnits(tx).Bind(contentUnit); err != nil {
		t.Error(err)
	}
}

func testContentUnitsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnit := &ContentUnit{}
	if err = randomize.Struct(seed, contentUnit, contentUnitDBTypes, true, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := ContentUnits(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testContentUnitsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnitOne := &ContentUnit{}
	contentUnitTwo := &ContentUnit{}
	if err = randomize.Struct(seed, contentUnitOne, contentUnitDBTypes, false, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}
	if err = randomize.Struct(seed, contentUnitTwo, contentUnitDBTypes, false, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = contentUnitTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := ContentUnits(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testContentUnitsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	contentUnitOne := &ContentUnit{}
	contentUnitTwo := &ContentUnit{}
	if err = randomize.Struct(seed, contentUnitOne, contentUnitDBTypes, false, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}
	if err = randomize.Struct(seed, contentUnitTwo, contentUnitDBTypes, false, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnitOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = contentUnitTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testContentUnitsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnit := &ContentUnit{}
	if err = randomize.Struct(seed, contentUnit, contentUnitDBTypes, true, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testContentUnitsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnit := &ContentUnit{}
	if err = randomize.Struct(seed, contentUnit, contentUnitDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnit.Insert(tx, contentUnitColumns...); err != nil {
		t.Error(err)
	}

	count, err := ContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testContentUnitToManyContentUnitsPersons(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentUnit
	var b, c ContentUnitsPerson

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentUnitDBTypes, true, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, contentUnitsPersonDBTypes, false, contentUnitsPersonColumnsWithDefault...)
	randomize.Struct(seed, &c, contentUnitsPersonDBTypes, false, contentUnitsPersonColumnsWithDefault...)

	b.ContentUnitID = a.ID
	c.ContentUnitID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	contentUnitsPerson, err := a.ContentUnitsPersons(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range contentUnitsPerson {
		if v.ContentUnitID == b.ContentUnitID {
			bFound = true
		}
		if v.ContentUnitID == c.ContentUnitID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := ContentUnitSlice{&a}
	if err = a.L.LoadContentUnitsPersons(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ContentUnitsPersons); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.ContentUnitsPersons = nil
	if err = a.L.LoadContentUnitsPersons(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ContentUnitsPersons); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", contentUnitsPerson)
	}
}

func testContentUnitToManyCollectionsContentUnits(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentUnit
	var b, c CollectionsContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentUnitDBTypes, true, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, collectionsContentUnitDBTypes, false, collectionsContentUnitColumnsWithDefault...)
	randomize.Struct(seed, &c, collectionsContentUnitDBTypes, false, collectionsContentUnitColumnsWithDefault...)

	b.ContentUnitID = a.ID
	c.ContentUnitID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	collectionsContentUnit, err := a.CollectionsContentUnits(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range collectionsContentUnit {
		if v.ContentUnitID == b.ContentUnitID {
			bFound = true
		}
		if v.ContentUnitID == c.ContentUnitID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := ContentUnitSlice{&a}
	if err = a.L.LoadCollectionsContentUnits(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.CollectionsContentUnits); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.CollectionsContentUnits = nil
	if err = a.L.LoadCollectionsContentUnits(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.CollectionsContentUnits); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", collectionsContentUnit)
	}
}

func testContentUnitToManyFiles(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentUnit
	var b, c File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentUnitDBTypes, true, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, fileDBTypes, false, fileColumnsWithDefault...)
	randomize.Struct(seed, &c, fileDBTypes, false, fileColumnsWithDefault...)

	b.ContentUnitID.Valid = true
	c.ContentUnitID.Valid = true
	b.ContentUnitID.Int64 = a.ID
	c.ContentUnitID.Int64 = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	file, err := a.Files(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range file {
		if v.ContentUnitID.Int64 == b.ContentUnitID.Int64 {
			bFound = true
		}
		if v.ContentUnitID.Int64 == c.ContentUnitID.Int64 {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := ContentUnitSlice{&a}
	if err = a.L.LoadFiles(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Files); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Files = nil
	if err = a.L.LoadFiles(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Files); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", file)
	}
}

func testContentUnitToManyAddOpContentUnitsPersons(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentUnit
	var b, c, d, e ContentUnitsPerson

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ContentUnitsPerson{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, contentUnitsPersonDBTypes, false, strmangle.SetComplement(contentUnitsPersonPrimaryKeyColumns, contentUnitsPersonColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*ContentUnitsPerson{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddContentUnitsPersons(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.ContentUnitID {
			t.Error("foreign key was wrong value", a.ID, first.ContentUnitID)
		}
		if a.ID != second.ContentUnitID {
			t.Error("foreign key was wrong value", a.ID, second.ContentUnitID)
		}

		if first.R.ContentUnit != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.ContentUnit != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.ContentUnitsPersons[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.ContentUnitsPersons[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.ContentUnitsPersons(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testContentUnitToManyAddOpCollectionsContentUnits(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentUnit
	var b, c, d, e CollectionsContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*CollectionsContentUnit{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, collectionsContentUnitDBTypes, false, strmangle.SetComplement(collectionsContentUnitPrimaryKeyColumns, collectionsContentUnitColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*CollectionsContentUnit{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddCollectionsContentUnits(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.ContentUnitID {
			t.Error("foreign key was wrong value", a.ID, first.ContentUnitID)
		}
		if a.ID != second.ContentUnitID {
			t.Error("foreign key was wrong value", a.ID, second.ContentUnitID)
		}

		if first.R.ContentUnit != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.ContentUnit != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.CollectionsContentUnits[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.CollectionsContentUnits[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.CollectionsContentUnits(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testContentUnitToManyAddOpFiles(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentUnit
	var b, c, d, e File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*File{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*File{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddFiles(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.ContentUnitID.Int64 {
			t.Error("foreign key was wrong value", a.ID, first.ContentUnitID.Int64)
		}
		if a.ID != second.ContentUnitID.Int64 {
			t.Error("foreign key was wrong value", a.ID, second.ContentUnitID.Int64)
		}

		if first.R.ContentUnit != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.ContentUnit != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Files[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Files[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Files(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testContentUnitToManySetOpFiles(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentUnit
	var b, c, d, e File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*File{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.SetFiles(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Files(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetFiles(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Files(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.ContentUnitID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.ContentUnitID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.ID != d.ContentUnitID.Int64 {
		t.Error("foreign key was wrong value", a.ID, d.ContentUnitID.Int64)
	}
	if a.ID != e.ContentUnitID.Int64 {
		t.Error("foreign key was wrong value", a.ID, e.ContentUnitID.Int64)
	}

	if b.R.ContentUnit != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.ContentUnit != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.ContentUnit != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.ContentUnit != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.Files[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Files[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testContentUnitToManyRemoveOpFiles(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentUnit
	var b, c, d, e File

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*File{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, fileDBTypes, false, strmangle.SetComplement(filePrimaryKeyColumns, fileColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddFiles(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Files(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveFiles(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Files(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.ContentUnitID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.ContentUnitID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.ContentUnit != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.ContentUnit != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.ContentUnit != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.ContentUnit != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.Files) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.Files[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.Files[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testContentUnitToOneStringTranslationUsingDescription(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local ContentUnit
	var foreign StringTranslation

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, contentUnitDBTypes, true, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	local.DescriptionID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.DescriptionID.Int64 = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Description(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ContentUnitSlice{&local}
	if err = local.L.LoadDescription(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Description == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Description = nil
	if err = local.L.LoadDescription(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Description == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testContentUnitToOneStringTranslationUsingName(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local ContentUnit
	var foreign StringTranslation

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, contentUnitDBTypes, true, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.NameID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Name(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ContentUnitSlice{&local}
	if err = local.L.LoadName(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Name == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Name = nil
	if err = local.L.LoadName(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Name == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testContentUnitToOneContentTypeUsingType(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local ContentUnit
	var foreign ContentType

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, contentUnitDBTypes, true, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.TypeID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Type(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ContentUnitSlice{&local}
	if err = local.L.LoadType(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Type == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Type = nil
	if err = local.L.LoadType(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Type == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testContentUnitToOneSetOpStringTranslationUsingDescription(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentUnit
	var b, c StringTranslation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*StringTranslation{&b, &c} {
		err = a.SetDescription(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Description != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.DescriptionContentUnits[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.DescriptionID.Int64 != x.ID {
			t.Error("foreign key was wrong value", a.DescriptionID.Int64)
		}

		zero := reflect.Zero(reflect.TypeOf(a.DescriptionID.Int64))
		reflect.Indirect(reflect.ValueOf(&a.DescriptionID.Int64)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.DescriptionID.Int64 != x.ID {
			t.Error("foreign key was wrong value", a.DescriptionID.Int64, x.ID)
		}
	}
}

func testContentUnitToOneRemoveOpStringTranslationUsingDescription(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentUnit
	var b StringTranslation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.SetDescription(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveDescription(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.Description(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.Description != nil {
		t.Error("R struct entry should be nil")
	}

	if a.DescriptionID.Valid {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.DescriptionContentUnits) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testContentUnitToOneSetOpStringTranslationUsingName(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentUnit
	var b, c StringTranslation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*StringTranslation{&b, &c} {
		err = a.SetName(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Name != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.NameContentUnits[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.NameID != x.ID {
			t.Error("foreign key was wrong value", a.NameID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.NameID))
		reflect.Indirect(reflect.ValueOf(&a.NameID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.NameID != x.ID {
			t.Error("foreign key was wrong value", a.NameID, x.ID)
		}
	}
}
func testContentUnitToOneSetOpContentTypeUsingType(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a ContentUnit
	var b, c ContentType

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, contentUnitDBTypes, false, strmangle.SetComplement(contentUnitPrimaryKeyColumns, contentUnitColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, contentTypeDBTypes, false, strmangle.SetComplement(contentTypePrimaryKeyColumns, contentTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, contentTypeDBTypes, false, strmangle.SetComplement(contentTypePrimaryKeyColumns, contentTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*ContentType{&b, &c} {
		err = a.SetType(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Type != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.TypeContentUnits[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.TypeID != x.ID {
			t.Error("foreign key was wrong value", a.TypeID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.TypeID))
		reflect.Indirect(reflect.ValueOf(&a.TypeID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.TypeID != x.ID {
			t.Error("foreign key was wrong value", a.TypeID, x.ID)
		}
	}
}
func testContentUnitsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnit := &ContentUnit{}
	if err = randomize.Struct(seed, contentUnit, contentUnitDBTypes, true, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = contentUnit.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testContentUnitsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnit := &ContentUnit{}
	if err = randomize.Struct(seed, contentUnit, contentUnitDBTypes, true, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ContentUnitSlice{contentUnit}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testContentUnitsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	contentUnit := &ContentUnit{}
	if err = randomize.Struct(seed, contentUnit, contentUnitDBTypes, true, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := ContentUnits(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	contentUnitDBTypes = map[string]string{`CreatedAt`: `timestamp with time zone`, `DescriptionID`: `bigint`, `ID`: `bigint`, `NameID`: `bigint`, `Properties`: `jsonb`, `TypeID`: `bigint`, `UID`: `character`}
	_                  = bytes.MinRead
)

func testContentUnitsUpdate(t *testing.T) {
	t.Parallel()

	if len(contentUnitColumns) == len(contentUnitPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	contentUnit := &ContentUnit{}
	if err = randomize.Struct(seed, contentUnit, contentUnitDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, contentUnit, contentUnitDBTypes, true, contentUnitColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	if err = contentUnit.Update(tx); err != nil {
		t.Error(err)
	}
}

func testContentUnitsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(contentUnitColumns) == len(contentUnitPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	contentUnit := &ContentUnit{}
	if err = randomize.Struct(seed, contentUnit, contentUnitDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnit.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := ContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, contentUnit, contentUnitDBTypes, true, contentUnitPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(contentUnitColumns, contentUnitPrimaryKeyColumns) {
		fields = contentUnitColumns
	} else {
		fields = strmangle.SetComplement(
			contentUnitColumns,
			contentUnitPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(contentUnit))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := ContentUnitSlice{contentUnit}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testContentUnitsUpsert(t *testing.T) {
	t.Parallel()

	if len(contentUnitColumns) == len(contentUnitPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	contentUnit := ContentUnit{}
	if err = randomize.Struct(seed, &contentUnit, contentUnitDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = contentUnit.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert ContentUnit: %s", err)
	}

	count, err := ContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &contentUnit, contentUnitDBTypes, false, contentUnitPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ContentUnit struct: %s", err)
	}

	if err = contentUnit.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert ContentUnit: %s", err)
	}

	count, err = ContentUnits(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}