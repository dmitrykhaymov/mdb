package models

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
	"gopkg.in/nullbio/null.v6"
)

// Author is an object representing the database table.
type Author struct {
	ID        int64       `boil:"id" json:"id" toml:"id" yaml:"id"`
	Code      string      `boil:"code" json:"code" toml:"code" yaml:"code"`
	Name      string      `boil:"name" json:"name" toml:"name" yaml:"name"`
	FullName  null.String `boil:"full_name" json:"full_name,omitempty" toml:"full_name" yaml:"full_name,omitempty"`
	CreatedAt time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *authorR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L authorL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// authorR is where relationships are stored.
type authorR struct {
	Sources     SourceSlice
	AuthorI18ns AuthorI18nSlice
}

// authorL is where Load methods for each relationship are stored.
type authorL struct{}

var (
	authorColumns               = []string{"id", "code", "name", "full_name", "created_at"}
	authorColumnsWithoutDefault = []string{"code", "name", "full_name"}
	authorColumnsWithDefault    = []string{"id", "created_at"}
	authorPrimaryKeyColumns     = []string{"id"}
)

type (
	// AuthorSlice is an alias for a slice of pointers to Author.
	// This should generally be used opposed to []Author.
	AuthorSlice []*Author

	authorQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	authorType                 = reflect.TypeOf(&Author{})
	authorMapping              = queries.MakeStructMapping(authorType)
	authorPrimaryKeyMapping, _ = queries.BindMapping(authorType, authorMapping, authorPrimaryKeyColumns)
	authorInsertCacheMut       sync.RWMutex
	authorInsertCache          = make(map[string]insertCache)
	authorUpdateCacheMut       sync.RWMutex
	authorUpdateCache          = make(map[string]updateCache)
	authorUpsertCacheMut       sync.RWMutex
	authorUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single author record from the query, and panics on error.
func (q authorQuery) OneP() *Author {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single author record from the query.
func (q authorQuery) One() (*Author, error) {
	o := &Author{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for authors")
	}

	return o, nil
}

// AllP returns all Author records from the query, and panics on error.
func (q authorQuery) AllP() AuthorSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Author records from the query.
func (q authorQuery) All() (AuthorSlice, error) {
	var o AuthorSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Author slice")
	}

	return o, nil
}

// CountP returns the count of all Author records in the query, and panics on error.
func (q authorQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Author records in the query.
func (q authorQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count authors rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q authorQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q authorQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if authors exists")
	}

	return count > 0, nil
}

// SourcesG retrieves all the source's sources.
func (o *Author) SourcesG(mods ...qm.QueryMod) sourceQuery {
	return o.Sources(boil.GetDB(), mods...)
}

// Sources retrieves all the source's sources with an executor.
func (o *Author) Sources(exec boil.Executor, mods ...qm.QueryMod) sourceQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.InnerJoin("\"authors_sources\" as \"b\" on \"a\".\"id\" = \"b\".\"source_id\""),
		qm.Where("\"b\".\"author_id\"=?", o.ID),
	)

	query := Sources(exec, queryMods...)
	queries.SetFrom(query.Query, "\"sources\" as \"a\"")
	return query
}

// AuthorI18nsG retrieves all the author_i18n's author i18n.
func (o *Author) AuthorI18nsG(mods ...qm.QueryMod) authorI18nQuery {
	return o.AuthorI18ns(boil.GetDB(), mods...)
}

// AuthorI18ns retrieves all the author_i18n's author i18n with an executor.
func (o *Author) AuthorI18ns(exec boil.Executor, mods ...qm.QueryMod) authorI18nQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"author_id\"=?", o.ID),
	)

	query := AuthorI18ns(exec, queryMods...)
	queries.SetFrom(query.Query, "\"author_i18n\" as \"a\"")
	return query
}

// LoadSources allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (authorL) LoadSources(e boil.Executor, singular bool, maybeAuthor interface{}) error {
	var slice []*Author
	var object *Author

	count := 1
	if singular {
		object = maybeAuthor.(*Author)
	} else {
		slice = *maybeAuthor.(*AuthorSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &authorR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &authorR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select \"a\".*, \"b\".\"author_id\" from \"sources\" as \"a\" inner join \"authors_sources\" as \"b\" on \"a\".\"id\" = \"b\".\"source_id\" where \"b\".\"author_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load sources")
	}
	defer results.Close()

	var resultSlice []*Source

	var localJoinCols []int64
	for results.Next() {
		one := new(Source)
		var localJoinCol int64

		err = results.Scan(&one.ID, &one.UID, &one.ParentID, &one.Pattern, &one.TypeID, &one.Position, &one.Name, &one.Description, &one.CreatedAt, &one.Properties, &localJoinCol)
		if err = results.Err(); err != nil {
			return errors.Wrap(err, "failed to plebian-bind eager loaded slice sources")
		}

		resultSlice = append(resultSlice, one)
		localJoinCols = append(localJoinCols, localJoinCol)
	}

	if err = results.Err(); err != nil {
		return errors.Wrap(err, "failed to plebian-bind eager loaded slice sources")
	}

	if singular {
		object.R.Sources = resultSlice
		return nil
	}

	for i, foreign := range resultSlice {
		localJoinCol := localJoinCols[i]
		for _, local := range slice {
			if local.ID == localJoinCol {
				local.R.Sources = append(local.R.Sources, foreign)
				break
			}
		}
	}

	return nil
}

// LoadAuthorI18ns allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (authorL) LoadAuthorI18ns(e boil.Executor, singular bool, maybeAuthor interface{}) error {
	var slice []*Author
	var object *Author

	count := 1
	if singular {
		object = maybeAuthor.(*Author)
	} else {
		slice = *maybeAuthor.(*AuthorSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &authorR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &authorR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"author_i18n\" where \"author_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load author_i18n")
	}
	defer results.Close()

	var resultSlice []*AuthorI18n
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice author_i18n")
	}

	if singular {
		object.R.AuthorI18ns = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.AuthorID {
				local.R.AuthorI18ns = append(local.R.AuthorI18ns, foreign)
				break
			}
		}
	}

	return nil
}

// AddSourcesG adds the given related objects to the existing relationships
// of the author, optionally inserting them as new records.
// Appends related to o.R.Sources.
// Sets related.R.Authors appropriately.
// Uses the global database handle.
func (o *Author) AddSourcesG(insert bool, related ...*Source) error {
	return o.AddSources(boil.GetDB(), insert, related...)
}

// AddSourcesP adds the given related objects to the existing relationships
// of the author, optionally inserting them as new records.
// Appends related to o.R.Sources.
// Sets related.R.Authors appropriately.
// Panics on error.
func (o *Author) AddSourcesP(exec boil.Executor, insert bool, related ...*Source) {
	if err := o.AddSources(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddSourcesGP adds the given related objects to the existing relationships
// of the author, optionally inserting them as new records.
// Appends related to o.R.Sources.
// Sets related.R.Authors appropriately.
// Uses the global database handle and panics on error.
func (o *Author) AddSourcesGP(insert bool, related ...*Source) {
	if err := o.AddSources(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddSources adds the given related objects to the existing relationships
// of the author, optionally inserting them as new records.
// Appends related to o.R.Sources.
// Sets related.R.Authors appropriately.
func (o *Author) AddSources(exec boil.Executor, insert bool, related ...*Source) error {
	var err error
	for _, rel := range related {
		if insert {
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		}
	}

	for _, rel := range related {
		query := "insert into \"authors_sources\" (\"author_id\", \"source_id\") values ($1, $2)"
		values := []interface{}{o.ID, rel.ID}

		if boil.DebugMode {
			fmt.Fprintln(boil.DebugWriter, query)
			fmt.Fprintln(boil.DebugWriter, values)
		}

		_, err = exec.Exec(query, values...)
		if err != nil {
			return errors.Wrap(err, "failed to insert into join table")
		}
	}
	if o.R == nil {
		o.R = &authorR{
			Sources: related,
		}
	} else {
		o.R.Sources = append(o.R.Sources, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &sourceR{
				Authors: AuthorSlice{o},
			}
		} else {
			rel.R.Authors = append(rel.R.Authors, o)
		}
	}
	return nil
}

// SetSourcesG removes all previously related items of the
// author replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Authors's Sources accordingly.
// Replaces o.R.Sources with related.
// Sets related.R.Authors's Sources accordingly.
// Uses the global database handle.
func (o *Author) SetSourcesG(insert bool, related ...*Source) error {
	return o.SetSources(boil.GetDB(), insert, related...)
}

// SetSourcesP removes all previously related items of the
// author replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Authors's Sources accordingly.
// Replaces o.R.Sources with related.
// Sets related.R.Authors's Sources accordingly.
// Panics on error.
func (o *Author) SetSourcesP(exec boil.Executor, insert bool, related ...*Source) {
	if err := o.SetSources(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetSourcesGP removes all previously related items of the
// author replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Authors's Sources accordingly.
// Replaces o.R.Sources with related.
// Sets related.R.Authors's Sources accordingly.
// Uses the global database handle and panics on error.
func (o *Author) SetSourcesGP(insert bool, related ...*Source) {
	if err := o.SetSources(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetSources removes all previously related items of the
// author replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Authors's Sources accordingly.
// Replaces o.R.Sources with related.
// Sets related.R.Authors's Sources accordingly.
func (o *Author) SetSources(exec boil.Executor, insert bool, related ...*Source) error {
	query := "delete from \"authors_sources\" where \"author_id\" = $1"
	values := []interface{}{o.ID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	removeSourcesFromAuthorsSlice(o, related)
	o.R.Sources = nil
	return o.AddSources(exec, insert, related...)
}

// RemoveSourcesG relationships from objects passed in.
// Removes related items from R.Sources (uses pointer comparison, removal does not keep order)
// Sets related.R.Authors.
// Uses the global database handle.
func (o *Author) RemoveSourcesG(related ...*Source) error {
	return o.RemoveSources(boil.GetDB(), related...)
}

// RemoveSourcesP relationships from objects passed in.
// Removes related items from R.Sources (uses pointer comparison, removal does not keep order)
// Sets related.R.Authors.
// Panics on error.
func (o *Author) RemoveSourcesP(exec boil.Executor, related ...*Source) {
	if err := o.RemoveSources(exec, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveSourcesGP relationships from objects passed in.
// Removes related items from R.Sources (uses pointer comparison, removal does not keep order)
// Sets related.R.Authors.
// Uses the global database handle and panics on error.
func (o *Author) RemoveSourcesGP(related ...*Source) {
	if err := o.RemoveSources(boil.GetDB(), related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveSources relationships from objects passed in.
// Removes related items from R.Sources (uses pointer comparison, removal does not keep order)
// Sets related.R.Authors.
func (o *Author) RemoveSources(exec boil.Executor, related ...*Source) error {
	var err error
	query := fmt.Sprintf(
		"delete from \"authors_sources\" where \"author_id\" = $1 and \"source_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, len(related), 1, 1),
	)
	values := []interface{}{o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}
	removeSourcesFromAuthorsSlice(o, related)
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.Sources {
			if rel != ri {
				continue
			}

			ln := len(o.R.Sources)
			if ln > 1 && i < ln-1 {
				o.R.Sources[i] = o.R.Sources[ln-1]
			}
			o.R.Sources = o.R.Sources[:ln-1]
			break
		}
	}

	return nil
}

func removeSourcesFromAuthorsSlice(o *Author, related []*Source) {
	for _, rel := range related {
		if rel.R == nil {
			continue
		}
		for i, ri := range rel.R.Authors {
			if o.ID != ri.ID {
				continue
			}

			ln := len(rel.R.Authors)
			if ln > 1 && i < ln-1 {
				rel.R.Authors[i] = rel.R.Authors[ln-1]
			}
			rel.R.Authors = rel.R.Authors[:ln-1]
			break
		}
	}
}

// AddAuthorI18nsG adds the given related objects to the existing relationships
// of the author, optionally inserting them as new records.
// Appends related to o.R.AuthorI18ns.
// Sets related.R.Author appropriately.
// Uses the global database handle.
func (o *Author) AddAuthorI18nsG(insert bool, related ...*AuthorI18n) error {
	return o.AddAuthorI18ns(boil.GetDB(), insert, related...)
}

// AddAuthorI18nsP adds the given related objects to the existing relationships
// of the author, optionally inserting them as new records.
// Appends related to o.R.AuthorI18ns.
// Sets related.R.Author appropriately.
// Panics on error.
func (o *Author) AddAuthorI18nsP(exec boil.Executor, insert bool, related ...*AuthorI18n) {
	if err := o.AddAuthorI18ns(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddAuthorI18nsGP adds the given related objects to the existing relationships
// of the author, optionally inserting them as new records.
// Appends related to o.R.AuthorI18ns.
// Sets related.R.Author appropriately.
// Uses the global database handle and panics on error.
func (o *Author) AddAuthorI18nsGP(insert bool, related ...*AuthorI18n) {
	if err := o.AddAuthorI18ns(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddAuthorI18ns adds the given related objects to the existing relationships
// of the author, optionally inserting them as new records.
// Appends related to o.R.AuthorI18ns.
// Sets related.R.Author appropriately.
func (o *Author) AddAuthorI18ns(exec boil.Executor, insert bool, related ...*AuthorI18n) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.AuthorID = o.ID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"author_i18n\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"author_id"}),
				strmangle.WhereClause("\"", "\"", 2, authorI18nPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.AuthorID, rel.Language}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.AuthorID = o.ID
		}
	}

	if o.R == nil {
		o.R = &authorR{
			AuthorI18ns: related,
		}
	} else {
		o.R.AuthorI18ns = append(o.R.AuthorI18ns, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &authorI18nR{
				Author: o,
			}
		} else {
			rel.R.Author = o
		}
	}
	return nil
}

// AuthorsG retrieves all records.
func AuthorsG(mods ...qm.QueryMod) authorQuery {
	return Authors(boil.GetDB(), mods...)
}

// Authors retrieves all the records using an executor.
func Authors(exec boil.Executor, mods ...qm.QueryMod) authorQuery {
	mods = append(mods, qm.From("\"authors\""))
	return authorQuery{NewQuery(exec, mods...)}
}

// FindAuthorG retrieves a single record by ID.
func FindAuthorG(id int64, selectCols ...string) (*Author, error) {
	return FindAuthor(boil.GetDB(), id, selectCols...)
}

// FindAuthorGP retrieves a single record by ID, and panics on error.
func FindAuthorGP(id int64, selectCols ...string) *Author {
	retobj, err := FindAuthor(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindAuthor retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindAuthor(exec boil.Executor, id int64, selectCols ...string) (*Author, error) {
	authorObj := &Author{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"authors\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(authorObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from authors")
	}

	return authorObj, nil
}

// FindAuthorP retrieves a single record by ID with an executor, and panics on error.
func FindAuthorP(exec boil.Executor, id int64, selectCols ...string) *Author {
	retobj, err := FindAuthor(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Author) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Author) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Author) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Author) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no authors provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(authorColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	authorInsertCacheMut.RLock()
	cache, cached := authorInsertCache[key]
	authorInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			authorColumns,
			authorColumnsWithDefault,
			authorColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(authorType, authorMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(authorType, authorMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"authors\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into authors")
	}

	if !cached {
		authorInsertCacheMut.Lock()
		authorInsertCache[key] = cache
		authorInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single Author record. See Update for
// whitelist behavior description.
func (o *Author) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Author record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Author) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Author, and panics on error.
// See Update for whitelist behavior description.
func (o *Author) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Author.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Author) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	key := makeCacheKey(whitelist, nil)
	authorUpdateCacheMut.RLock()
	cache, cached := authorUpdateCache[key]
	authorUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(authorColumns, authorPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update authors, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"authors\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, authorPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(authorType, authorMapping, append(wl, authorPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update authors row")
	}

	if !cached {
		authorUpdateCacheMut.Lock()
		authorUpdateCache[key] = cache
		authorUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q authorQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q authorQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for authors")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o AuthorSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o AuthorSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o AuthorSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o AuthorSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), authorPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"authors\" SET %s WHERE (\"id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(authorPrimaryKeyColumns), len(colNames)+1, len(authorPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in author slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Author) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Author) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Author) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Author) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no authors provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(authorColumnsWithDefault, o)

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

	authorUpsertCacheMut.RLock()
	cache, cached := authorUpsertCache[key]
	authorUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			authorColumns,
			authorColumnsWithDefault,
			authorColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			authorColumns,
			authorPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert authors, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(authorPrimaryKeyColumns))
			copy(conflict, authorPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"authors\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(authorType, authorMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(authorType, authorMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert authors")
	}

	if !cached {
		authorUpsertCacheMut.Lock()
		authorUpsertCache[key] = cache
		authorUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single Author record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Author) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Author record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Author) DeleteG() error {
	if o == nil {
		return errors.New("models: no Author provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Author record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Author) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Author record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Author) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Author provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), authorPrimaryKeyMapping)
	sql := "DELETE FROM \"authors\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from authors")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q authorQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q authorQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no authorQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from authors")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o AuthorSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o AuthorSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Author slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o AuthorSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o AuthorSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Author slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), authorPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"authors\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, authorPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(authorPrimaryKeyColumns), 1, len(authorPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from author slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Author) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Author) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Author) ReloadG() error {
	if o == nil {
		return errors.New("models: no Author provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Author) Reload(exec boil.Executor) error {
	ret, err := FindAuthor(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *AuthorSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *AuthorSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *AuthorSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty AuthorSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *AuthorSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	authors := AuthorSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), authorPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"authors\".* FROM \"authors\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, authorPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(authorPrimaryKeyColumns), 1, len(authorPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&authors)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in AuthorSlice")
	}

	*o = authors

	return nil
}

// AuthorExists checks if the Author row exists.
func AuthorExists(exec boil.Executor, id int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"authors\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if authors exists")
	}

	return exists, nil
}

// AuthorExistsG checks if the Author row exists.
func AuthorExistsG(id int64) (bool, error) {
	return AuthorExists(boil.GetDB(), id)
}

// AuthorExistsGP checks if the Author row exists. Panics on error.
func AuthorExistsGP(id int64) bool {
	e, err := AuthorExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// AuthorExistsP checks if the Author row exists. Panics on error.
func AuthorExistsP(exec boil.Executor, id int64) bool {
	e, err := AuthorExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
