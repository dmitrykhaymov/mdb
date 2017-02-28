package gmodels

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/vattle/sqlboiler/strmangle"
)

// ContentUnitsPerson is an object representing the database table.
type ContentUnitsPerson struct {
	ContentUnitID int64 `boil:"content_unit_id" json:"content_unit_id" toml:"content_unit_id" yaml:"content_unit_id"`
	PersonID      int64 `boil:"person_id" json:"person_id" toml:"person_id" yaml:"person_id"`
	RoleID        int64 `boil:"role_id" json:"role_id" toml:"role_id" yaml:"role_id"`

	R *contentUnitsPersonR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L contentUnitsPersonL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// contentUnitsPersonR is where relationships are stored.
type contentUnitsPersonR struct {
	ContentUnit *ContentUnit
	Person      *Person
	Role        *ContentRole
}

// contentUnitsPersonL is where Load methods for each relationship are stored.
type contentUnitsPersonL struct{}

var (
	contentUnitsPersonColumns               = []string{"content_unit_id", "person_id", "role_id"}
	contentUnitsPersonColumnsWithoutDefault = []string{"content_unit_id", "person_id", "role_id"}
	contentUnitsPersonColumnsWithDefault    = []string{}
	contentUnitsPersonPrimaryKeyColumns     = []string{"content_unit_id", "person_id"}
)

type (
	// ContentUnitsPersonSlice is an alias for a slice of pointers to ContentUnitsPerson.
	// This should generally be used opposed to []ContentUnitsPerson.
	ContentUnitsPersonSlice []*ContentUnitsPerson

	contentUnitsPersonQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	contentUnitsPersonType                 = reflect.TypeOf(&ContentUnitsPerson{})
	contentUnitsPersonMapping              = queries.MakeStructMapping(contentUnitsPersonType)
	contentUnitsPersonPrimaryKeyMapping, _ = queries.BindMapping(contentUnitsPersonType, contentUnitsPersonMapping, contentUnitsPersonPrimaryKeyColumns)
	contentUnitsPersonInsertCacheMut       sync.RWMutex
	contentUnitsPersonInsertCache          = make(map[string]insertCache)
	contentUnitsPersonUpdateCacheMut       sync.RWMutex
	contentUnitsPersonUpdateCache          = make(map[string]updateCache)
	contentUnitsPersonUpsertCacheMut       sync.RWMutex
	contentUnitsPersonUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single contentUnitsPerson record from the query, and panics on error.
func (q contentUnitsPersonQuery) OneP() *ContentUnitsPerson {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single contentUnitsPerson record from the query.
func (q contentUnitsPersonQuery) One() (*ContentUnitsPerson, error) {
	o := &ContentUnitsPerson{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "gmodels: failed to execute a one query for content_units_persons")
	}

	return o, nil
}

// AllP returns all ContentUnitsPerson records from the query, and panics on error.
func (q contentUnitsPersonQuery) AllP() ContentUnitsPersonSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all ContentUnitsPerson records from the query.
func (q contentUnitsPersonQuery) All() (ContentUnitsPersonSlice, error) {
	var o ContentUnitsPersonSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "gmodels: failed to assign all query results to ContentUnitsPerson slice")
	}

	return o, nil
}

// CountP returns the count of all ContentUnitsPerson records in the query, and panics on error.
func (q contentUnitsPersonQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all ContentUnitsPerson records in the query.
func (q contentUnitsPersonQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "gmodels: failed to count content_units_persons rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q contentUnitsPersonQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q contentUnitsPersonQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "gmodels: failed to check if content_units_persons exists")
	}

	return count > 0, nil
}

// ContentUnitG pointed to by the foreign key.
func (o *ContentUnitsPerson) ContentUnitG(mods ...qm.QueryMod) contentUnitQuery {
	return o.ContentUnit(boil.GetDB(), mods...)
}

// ContentUnit pointed to by the foreign key.
func (o *ContentUnitsPerson) ContentUnit(exec boil.Executor, mods ...qm.QueryMod) contentUnitQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.ContentUnitID),
	}

	queryMods = append(queryMods, mods...)

	query := ContentUnits(exec, queryMods...)
	queries.SetFrom(query.Query, "\"content_units\"")

	return query
}

// PersonG pointed to by the foreign key.
func (o *ContentUnitsPerson) PersonG(mods ...qm.QueryMod) personQuery {
	return o.Person(boil.GetDB(), mods...)
}

// Person pointed to by the foreign key.
func (o *ContentUnitsPerson) Person(exec boil.Executor, mods ...qm.QueryMod) personQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.PersonID),
	}

	queryMods = append(queryMods, mods...)

	query := Persons(exec, queryMods...)
	queries.SetFrom(query.Query, "\"persons\"")

	return query
}

// RoleG pointed to by the foreign key.
func (o *ContentUnitsPerson) RoleG(mods ...qm.QueryMod) contentRoleQuery {
	return o.Role(boil.GetDB(), mods...)
}

// Role pointed to by the foreign key.
func (o *ContentUnitsPerson) Role(exec boil.Executor, mods ...qm.QueryMod) contentRoleQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.RoleID),
	}

	queryMods = append(queryMods, mods...)

	query := ContentRoles(exec, queryMods...)
	queries.SetFrom(query.Query, "\"content_roles\"")

	return query
}

// LoadContentUnit allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (contentUnitsPersonL) LoadContentUnit(e boil.Executor, singular bool, maybeContentUnitsPerson interface{}) error {
	var slice []*ContentUnitsPerson
	var object *ContentUnitsPerson

	count := 1
	if singular {
		object = maybeContentUnitsPerson.(*ContentUnitsPerson)
	} else {
		slice = *maybeContentUnitsPerson.(*ContentUnitsPersonSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &contentUnitsPersonR{}
		}
		args[0] = object.ContentUnitID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &contentUnitsPersonR{}
			}
			args[i] = obj.ContentUnitID
		}
	}

	query := fmt.Sprintf(
		"select * from \"content_units\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load ContentUnit")
	}
	defer results.Close()

	var resultSlice []*ContentUnit
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice ContentUnit")
	}

	if singular && len(resultSlice) != 0 {
		object.R.ContentUnit = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ContentUnitID == foreign.ID {
				local.R.ContentUnit = foreign
				break
			}
		}
	}

	return nil
}

// LoadPerson allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (contentUnitsPersonL) LoadPerson(e boil.Executor, singular bool, maybeContentUnitsPerson interface{}) error {
	var slice []*ContentUnitsPerson
	var object *ContentUnitsPerson

	count := 1
	if singular {
		object = maybeContentUnitsPerson.(*ContentUnitsPerson)
	} else {
		slice = *maybeContentUnitsPerson.(*ContentUnitsPersonSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &contentUnitsPersonR{}
		}
		args[0] = object.PersonID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &contentUnitsPersonR{}
			}
			args[i] = obj.PersonID
		}
	}

	query := fmt.Sprintf(
		"select * from \"persons\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Person")
	}
	defer results.Close()

	var resultSlice []*Person
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Person")
	}

	if singular && len(resultSlice) != 0 {
		object.R.Person = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.PersonID == foreign.ID {
				local.R.Person = foreign
				break
			}
		}
	}

	return nil
}

// LoadRole allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (contentUnitsPersonL) LoadRole(e boil.Executor, singular bool, maybeContentUnitsPerson interface{}) error {
	var slice []*ContentUnitsPerson
	var object *ContentUnitsPerson

	count := 1
	if singular {
		object = maybeContentUnitsPerson.(*ContentUnitsPerson)
	} else {
		slice = *maybeContentUnitsPerson.(*ContentUnitsPersonSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &contentUnitsPersonR{}
		}
		args[0] = object.RoleID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &contentUnitsPersonR{}
			}
			args[i] = obj.RoleID
		}
	}

	query := fmt.Sprintf(
		"select * from \"content_roles\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load ContentRole")
	}
	defer results.Close()

	var resultSlice []*ContentRole
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice ContentRole")
	}

	if singular && len(resultSlice) != 0 {
		object.R.Role = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.RoleID == foreign.ID {
				local.R.Role = foreign
				break
			}
		}
	}

	return nil
}

// SetContentUnit of the content_units_person to the related item.
// Sets o.R.ContentUnit to related.
// Adds o to related.R.ContentUnitsPersons.
func (o *ContentUnitsPerson) SetContentUnit(exec boil.Executor, insert bool, related *ContentUnit) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"content_units_persons\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"content_unit_id"}),
		strmangle.WhereClause("\"", "\"", 2, contentUnitsPersonPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ContentUnitID, o.PersonID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ContentUnitID = related.ID

	if o.R == nil {
		o.R = &contentUnitsPersonR{
			ContentUnit: related,
		}
	} else {
		o.R.ContentUnit = related
	}

	if related.R == nil {
		related.R = &contentUnitR{
			ContentUnitsPersons: ContentUnitsPersonSlice{o},
		}
	} else {
		related.R.ContentUnitsPersons = append(related.R.ContentUnitsPersons, o)
	}

	return nil
}

// SetPerson of the content_units_person to the related item.
// Sets o.R.Person to related.
// Adds o to related.R.ContentUnitsPersons.
func (o *ContentUnitsPerson) SetPerson(exec boil.Executor, insert bool, related *Person) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"content_units_persons\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"person_id"}),
		strmangle.WhereClause("\"", "\"", 2, contentUnitsPersonPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ContentUnitID, o.PersonID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.PersonID = related.ID

	if o.R == nil {
		o.R = &contentUnitsPersonR{
			Person: related,
		}
	} else {
		o.R.Person = related
	}

	if related.R == nil {
		related.R = &personR{
			ContentUnitsPersons: ContentUnitsPersonSlice{o},
		}
	} else {
		related.R.ContentUnitsPersons = append(related.R.ContentUnitsPersons, o)
	}

	return nil
}

// SetRole of the content_units_person to the related item.
// Sets o.R.Role to related.
// Adds o to related.R.RoleContentUnitsPersons.
func (o *ContentUnitsPerson) SetRole(exec boil.Executor, insert bool, related *ContentRole) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"content_units_persons\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"role_id"}),
		strmangle.WhereClause("\"", "\"", 2, contentUnitsPersonPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ContentUnitID, o.PersonID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.RoleID = related.ID

	if o.R == nil {
		o.R = &contentUnitsPersonR{
			Role: related,
		}
	} else {
		o.R.Role = related
	}

	if related.R == nil {
		related.R = &contentRoleR{
			RoleContentUnitsPersons: ContentUnitsPersonSlice{o},
		}
	} else {
		related.R.RoleContentUnitsPersons = append(related.R.RoleContentUnitsPersons, o)
	}

	return nil
}

// ContentUnitsPersonsG retrieves all records.
func ContentUnitsPersonsG(mods ...qm.QueryMod) contentUnitsPersonQuery {
	return ContentUnitsPersons(boil.GetDB(), mods...)
}

// ContentUnitsPersons retrieves all the records using an executor.
func ContentUnitsPersons(exec boil.Executor, mods ...qm.QueryMod) contentUnitsPersonQuery {
	mods = append(mods, qm.From("\"content_units_persons\""))
	return contentUnitsPersonQuery{NewQuery(exec, mods...)}
}

// FindContentUnitsPersonG retrieves a single record by ID.
func FindContentUnitsPersonG(contentUnitID int64, personID int64, selectCols ...string) (*ContentUnitsPerson, error) {
	return FindContentUnitsPerson(boil.GetDB(), contentUnitID, personID, selectCols...)
}

// FindContentUnitsPersonGP retrieves a single record by ID, and panics on error.
func FindContentUnitsPersonGP(contentUnitID int64, personID int64, selectCols ...string) *ContentUnitsPerson {
	retobj, err := FindContentUnitsPerson(boil.GetDB(), contentUnitID, personID, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindContentUnitsPerson retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindContentUnitsPerson(exec boil.Executor, contentUnitID int64, personID int64, selectCols ...string) (*ContentUnitsPerson, error) {
	contentUnitsPersonObj := &ContentUnitsPerson{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"content_units_persons\" where \"content_unit_id\"=$1 AND \"person_id\"=$2", sel,
	)

	q := queries.Raw(exec, query, contentUnitID, personID)

	err := q.Bind(contentUnitsPersonObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "gmodels: unable to select from content_units_persons")
	}

	return contentUnitsPersonObj, nil
}

// FindContentUnitsPersonP retrieves a single record by ID with an executor, and panics on error.
func FindContentUnitsPersonP(exec boil.Executor, contentUnitID int64, personID int64, selectCols ...string) *ContentUnitsPerson {
	retobj, err := FindContentUnitsPerson(exec, contentUnitID, personID, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *ContentUnitsPerson) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *ContentUnitsPerson) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *ContentUnitsPerson) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *ContentUnitsPerson) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("gmodels: no content_units_persons provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(contentUnitsPersonColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	contentUnitsPersonInsertCacheMut.RLock()
	cache, cached := contentUnitsPersonInsertCache[key]
	contentUnitsPersonInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			contentUnitsPersonColumns,
			contentUnitsPersonColumnsWithDefault,
			contentUnitsPersonColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(contentUnitsPersonType, contentUnitsPersonMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(contentUnitsPersonType, contentUnitsPersonMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"content_units_persons\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

		if len(cache.retMapping) != 0 {
			cache.query += fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "gmodels: unable to insert into content_units_persons")
	}

	if !cached {
		contentUnitsPersonInsertCacheMut.Lock()
		contentUnitsPersonInsertCache[key] = cache
		contentUnitsPersonInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single ContentUnitsPerson record. See Update for
// whitelist behavior description.
func (o *ContentUnitsPerson) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single ContentUnitsPerson record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *ContentUnitsPerson) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the ContentUnitsPerson, and panics on error.
// See Update for whitelist behavior description.
func (o *ContentUnitsPerson) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the ContentUnitsPerson.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *ContentUnitsPerson) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	key := makeCacheKey(whitelist, nil)
	contentUnitsPersonUpdateCacheMut.RLock()
	cache, cached := contentUnitsPersonUpdateCache[key]
	contentUnitsPersonUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(contentUnitsPersonColumns, contentUnitsPersonPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("gmodels: unable to update content_units_persons, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"content_units_persons\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, contentUnitsPersonPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(contentUnitsPersonType, contentUnitsPersonMapping, append(wl, contentUnitsPersonPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to update content_units_persons row")
	}

	if !cached {
		contentUnitsPersonUpdateCacheMut.Lock()
		contentUnitsPersonUpdateCache[key] = cache
		contentUnitsPersonUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q contentUnitsPersonQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q contentUnitsPersonQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to update all for content_units_persons")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o ContentUnitsPersonSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o ContentUnitsPersonSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o ContentUnitsPersonSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ContentUnitsPersonSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("gmodels: update all requires at least one column argument")
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contentUnitsPersonPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"content_units_persons\" SET %s WHERE (\"content_unit_id\",\"person_id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(contentUnitsPersonPrimaryKeyColumns), len(colNames)+1, len(contentUnitsPersonPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to update all in contentUnitsPerson slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *ContentUnitsPerson) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *ContentUnitsPerson) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *ContentUnitsPerson) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *ContentUnitsPerson) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("gmodels: no content_units_persons provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(contentUnitsPersonColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs postgres problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range updateColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range whitelist {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	contentUnitsPersonUpsertCacheMut.RLock()
	cache, cached := contentUnitsPersonUpsertCache[key]
	contentUnitsPersonUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			contentUnitsPersonColumns,
			contentUnitsPersonColumnsWithDefault,
			contentUnitsPersonColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			contentUnitsPersonColumns,
			contentUnitsPersonPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("gmodels: unable to upsert content_units_persons, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(contentUnitsPersonPrimaryKeyColumns))
			copy(conflict, contentUnitsPersonPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"content_units_persons\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(contentUnitsPersonType, contentUnitsPersonMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(contentUnitsPersonType, contentUnitsPersonMapping, ret)
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

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to upsert content_units_persons")
	}

	if !cached {
		contentUnitsPersonUpsertCacheMut.Lock()
		contentUnitsPersonUpsertCache[key] = cache
		contentUnitsPersonUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single ContentUnitsPerson record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *ContentUnitsPerson) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single ContentUnitsPerson record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *ContentUnitsPerson) DeleteG() error {
	if o == nil {
		return errors.New("gmodels: no ContentUnitsPerson provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single ContentUnitsPerson record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *ContentUnitsPerson) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single ContentUnitsPerson record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ContentUnitsPerson) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("gmodels: no ContentUnitsPerson provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), contentUnitsPersonPrimaryKeyMapping)
	sql := "DELETE FROM \"content_units_persons\" WHERE \"content_unit_id\"=$1 AND \"person_id\"=$2"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to delete from content_units_persons")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q contentUnitsPersonQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q contentUnitsPersonQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("gmodels: no contentUnitsPersonQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to delete all from content_units_persons")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o ContentUnitsPersonSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o ContentUnitsPersonSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("gmodels: no ContentUnitsPerson slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o ContentUnitsPersonSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ContentUnitsPersonSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("gmodels: no ContentUnitsPerson slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contentUnitsPersonPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"content_units_persons\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, contentUnitsPersonPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(contentUnitsPersonPrimaryKeyColumns), 1, len(contentUnitsPersonPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to delete all from contentUnitsPerson slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *ContentUnitsPerson) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *ContentUnitsPerson) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *ContentUnitsPerson) ReloadG() error {
	if o == nil {
		return errors.New("gmodels: no ContentUnitsPerson provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ContentUnitsPerson) Reload(exec boil.Executor) error {
	ret, err := FindContentUnitsPerson(exec, o.ContentUnitID, o.PersonID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *ContentUnitsPersonSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *ContentUnitsPersonSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ContentUnitsPersonSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("gmodels: empty ContentUnitsPersonSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ContentUnitsPersonSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	contentUnitsPersons := ContentUnitsPersonSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contentUnitsPersonPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"content_units_persons\".* FROM \"content_units_persons\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, contentUnitsPersonPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(contentUnitsPersonPrimaryKeyColumns), 1, len(contentUnitsPersonPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&contentUnitsPersons)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to reload all in ContentUnitsPersonSlice")
	}

	*o = contentUnitsPersons

	return nil
}

// ContentUnitsPersonExists checks if the ContentUnitsPerson row exists.
func ContentUnitsPersonExists(exec boil.Executor, contentUnitID int64, personID int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"content_units_persons\" where \"content_unit_id\"=$1 AND \"person_id\"=$2 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, contentUnitID, personID)
	}

	row := exec.QueryRow(sql, contentUnitID, personID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "gmodels: unable to check if content_units_persons exists")
	}

	return exists, nil
}

// ContentUnitsPersonExistsG checks if the ContentUnitsPerson row exists.
func ContentUnitsPersonExistsG(contentUnitID int64, personID int64) (bool, error) {
	return ContentUnitsPersonExists(boil.GetDB(), contentUnitID, personID)
}

// ContentUnitsPersonExistsGP checks if the ContentUnitsPerson row exists. Panics on error.
func ContentUnitsPersonExistsGP(contentUnitID int64, personID int64) bool {
	e, err := ContentUnitsPersonExists(boil.GetDB(), contentUnitID, personID)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// ContentUnitsPersonExistsP checks if the ContentUnitsPerson row exists. Panics on error.
func ContentUnitsPersonExistsP(exec boil.Executor, contentUnitID int64, personID int64) bool {
	e, err := ContentUnitsPersonExists(exec, contentUnitID, personID)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
