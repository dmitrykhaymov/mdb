package kmodels

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

// Language is an object representing the database table.
type Language struct {
	ID       int         `boil:"id" json:"id" toml:"id" yaml:"id"`
	Locale   null.String `boil:"locale" json:"locale,omitempty" toml:"locale" yaml:"locale,omitempty"`
	Code3    null.String `boil:"code3" json:"code3,omitempty" toml:"code3" yaml:"code3,omitempty"`
	Language null.String `boil:"language" json:"language,omitempty" toml:"language" yaml:"language,omitempty"`

	R *languageR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L languageL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// languageR is where relationships are stored.
type languageR struct {
	LangCatalogDescriptions   CatalogDescriptionSlice
	LangContainerDescriptions ContainerDescriptionSlice
	LangContainers            ContainerSlice
	LangFileAssets            FileAssetSlice
}

// languageL is where Load methods for each relationship are stored.
type languageL struct{}

var (
	languageColumns               = []string{"id", "locale", "code3", "language"}
	languageColumnsWithoutDefault = []string{"locale", "code3", "language"}
	languageColumnsWithDefault    = []string{"id"}
	languagePrimaryKeyColumns     = []string{"id"}
)

type (
	// LanguageSlice is an alias for a slice of pointers to Language.
	// This should generally be used opposed to []Language.
	LanguageSlice []*Language

	languageQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	languageType                 = reflect.TypeOf(&Language{})
	languageMapping              = queries.MakeStructMapping(languageType)
	languagePrimaryKeyMapping, _ = queries.BindMapping(languageType, languageMapping, languagePrimaryKeyColumns)
	languageInsertCacheMut       sync.RWMutex
	languageInsertCache          = make(map[string]insertCache)
	languageUpdateCacheMut       sync.RWMutex
	languageUpdateCache          = make(map[string]updateCache)
	languageUpsertCacheMut       sync.RWMutex
	languageUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single language record from the query, and panics on error.
func (q languageQuery) OneP() *Language {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single language record from the query.
func (q languageQuery) One() (*Language, error) {
	o := &Language{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "kmodels: failed to execute a one query for languages")
	}

	return o, nil
}

// AllP returns all Language records from the query, and panics on error.
func (q languageQuery) AllP() LanguageSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Language records from the query.
func (q languageQuery) All() (LanguageSlice, error) {
	var o LanguageSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "kmodels: failed to assign all query results to Language slice")
	}

	return o, nil
}

// CountP returns the count of all Language records in the query, and panics on error.
func (q languageQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Language records in the query.
func (q languageQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "kmodels: failed to count languages rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q languageQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q languageQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "kmodels: failed to check if languages exists")
	}

	return count > 0, nil
}

// LangCatalogDescriptionsG retrieves all the catalog_description's catalog descriptions via lang_id column.
func (o *Language) LangCatalogDescriptionsG(mods ...qm.QueryMod) catalogDescriptionQuery {
	return o.LangCatalogDescriptions(boil.GetDB(), mods...)
}

// LangCatalogDescriptions retrieves all the catalog_description's catalog descriptions with an executor via lang_id column.
func (o *Language) LangCatalogDescriptions(exec boil.Executor, mods ...qm.QueryMod) catalogDescriptionQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"lang_id\"=?", o.Code3),
	)

	query := CatalogDescriptions(exec, queryMods...)
	queries.SetFrom(query.Query, "\"catalog_descriptions\" as \"a\"")
	return query
}

// LangContainerDescriptionsG retrieves all the container_description's container descriptions via lang_id column.
func (o *Language) LangContainerDescriptionsG(mods ...qm.QueryMod) containerDescriptionQuery {
	return o.LangContainerDescriptions(boil.GetDB(), mods...)
}

// LangContainerDescriptions retrieves all the container_description's container descriptions with an executor via lang_id column.
func (o *Language) LangContainerDescriptions(exec boil.Executor, mods ...qm.QueryMod) containerDescriptionQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"lang_id\"=?", o.Code3),
	)

	query := ContainerDescriptions(exec, queryMods...)
	queries.SetFrom(query.Query, "\"container_descriptions\" as \"a\"")
	return query
}

// LangContainersG retrieves all the container's containers via lang_id column.
func (o *Language) LangContainersG(mods ...qm.QueryMod) containerQuery {
	return o.LangContainers(boil.GetDB(), mods...)
}

// LangContainers retrieves all the container's containers with an executor via lang_id column.
func (o *Language) LangContainers(exec boil.Executor, mods ...qm.QueryMod) containerQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"lang_id\"=?", o.Code3),
	)

	query := Containers(exec, queryMods...)
	queries.SetFrom(query.Query, "\"containers\" as \"a\"")
	return query
}

// LangFileAssetsG retrieves all the file_asset's file assets via lang_id column.
func (o *Language) LangFileAssetsG(mods ...qm.QueryMod) fileAssetQuery {
	return o.LangFileAssets(boil.GetDB(), mods...)
}

// LangFileAssets retrieves all the file_asset's file assets with an executor via lang_id column.
func (o *Language) LangFileAssets(exec boil.Executor, mods ...qm.QueryMod) fileAssetQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"lang_id\"=?", o.Code3),
	)

	query := FileAssets(exec, queryMods...)
	queries.SetFrom(query.Query, "\"file_assets\" as \"a\"")
	return query
}

// LoadLangCatalogDescriptions allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (languageL) LoadLangCatalogDescriptions(e boil.Executor, singular bool, maybeLanguage interface{}) error {
	var slice []*Language
	var object *Language

	count := 1
	if singular {
		object = maybeLanguage.(*Language)
	} else {
		slice = *maybeLanguage.(*LanguageSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &languageR{}
		}
		args[0] = object.Code3
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &languageR{}
			}
			args[i] = obj.Code3
		}
	}

	query := fmt.Sprintf(
		"select * from \"catalog_descriptions\" where \"lang_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load catalog_descriptions")
	}
	defer results.Close()

	var resultSlice []*CatalogDescription
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice catalog_descriptions")
	}

	if singular {
		object.R.LangCatalogDescriptions = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.Code3.String == foreign.LangID.String {
				local.R.LangCatalogDescriptions = append(local.R.LangCatalogDescriptions, foreign)
				break
			}
		}
	}

	return nil
}

// LoadLangContainerDescriptions allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (languageL) LoadLangContainerDescriptions(e boil.Executor, singular bool, maybeLanguage interface{}) error {
	var slice []*Language
	var object *Language

	count := 1
	if singular {
		object = maybeLanguage.(*Language)
	} else {
		slice = *maybeLanguage.(*LanguageSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &languageR{}
		}
		args[0] = object.Code3
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &languageR{}
			}
			args[i] = obj.Code3
		}
	}

	query := fmt.Sprintf(
		"select * from \"container_descriptions\" where \"lang_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load container_descriptions")
	}
	defer results.Close()

	var resultSlice []*ContainerDescription
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice container_descriptions")
	}

	if singular {
		object.R.LangContainerDescriptions = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.Code3.String == foreign.LangID.String {
				local.R.LangContainerDescriptions = append(local.R.LangContainerDescriptions, foreign)
				break
			}
		}
	}

	return nil
}

// LoadLangContainers allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (languageL) LoadLangContainers(e boil.Executor, singular bool, maybeLanguage interface{}) error {
	var slice []*Language
	var object *Language

	count := 1
	if singular {
		object = maybeLanguage.(*Language)
	} else {
		slice = *maybeLanguage.(*LanguageSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &languageR{}
		}
		args[0] = object.Code3
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &languageR{}
			}
			args[i] = obj.Code3
		}
	}

	query := fmt.Sprintf(
		"select * from \"containers\" where \"lang_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load containers")
	}
	defer results.Close()

	var resultSlice []*Container
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice containers")
	}

	if singular {
		object.R.LangContainers = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.Code3.String == foreign.LangID.String {
				local.R.LangContainers = append(local.R.LangContainers, foreign)
				break
			}
		}
	}

	return nil
}

// LoadLangFileAssets allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (languageL) LoadLangFileAssets(e boil.Executor, singular bool, maybeLanguage interface{}) error {
	var slice []*Language
	var object *Language

	count := 1
	if singular {
		object = maybeLanguage.(*Language)
	} else {
		slice = *maybeLanguage.(*LanguageSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &languageR{}
		}
		args[0] = object.Code3
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &languageR{}
			}
			args[i] = obj.Code3
		}
	}

	query := fmt.Sprintf(
		"select * from \"file_assets\" where \"lang_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load file_assets")
	}
	defer results.Close()

	var resultSlice []*FileAsset
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice file_assets")
	}

	if singular {
		object.R.LangFileAssets = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.Code3.String == foreign.LangID.String {
				local.R.LangFileAssets = append(local.R.LangFileAssets, foreign)
				break
			}
		}
	}

	return nil
}

// AddLangCatalogDescriptionsG adds the given related objects to the existing relationships
// of the language, optionally inserting them as new records.
// Appends related to o.R.LangCatalogDescriptions.
// Sets related.R.Lang appropriately.
// Uses the global database handle.
func (o *Language) AddLangCatalogDescriptionsG(insert bool, related ...*CatalogDescription) error {
	return o.AddLangCatalogDescriptions(boil.GetDB(), insert, related...)
}

// AddLangCatalogDescriptionsP adds the given related objects to the existing relationships
// of the language, optionally inserting them as new records.
// Appends related to o.R.LangCatalogDescriptions.
// Sets related.R.Lang appropriately.
// Panics on error.
func (o *Language) AddLangCatalogDescriptionsP(exec boil.Executor, insert bool, related ...*CatalogDescription) {
	if err := o.AddLangCatalogDescriptions(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddLangCatalogDescriptionsGP adds the given related objects to the existing relationships
// of the language, optionally inserting them as new records.
// Appends related to o.R.LangCatalogDescriptions.
// Sets related.R.Lang appropriately.
// Uses the global database handle and panics on error.
func (o *Language) AddLangCatalogDescriptionsGP(insert bool, related ...*CatalogDescription) {
	if err := o.AddLangCatalogDescriptions(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddLangCatalogDescriptions adds the given related objects to the existing relationships
// of the language, optionally inserting them as new records.
// Appends related to o.R.LangCatalogDescriptions.
// Sets related.R.Lang appropriately.
func (o *Language) AddLangCatalogDescriptions(exec boil.Executor, insert bool, related ...*CatalogDescription) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.LangID.String = o.Code3.String
			rel.LangID.Valid = true
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"catalog_descriptions\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"lang_id"}),
				strmangle.WhereClause("\"", "\"", 2, catalogDescriptionPrimaryKeyColumns),
			)
			values := []interface{}{o.Code3, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.LangID.String = o.Code3.String
			rel.LangID.Valid = true
		}
	}

	if o.R == nil {
		o.R = &languageR{
			LangCatalogDescriptions: related,
		}
	} else {
		o.R.LangCatalogDescriptions = append(o.R.LangCatalogDescriptions, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &catalogDescriptionR{
				Lang: o,
			}
		} else {
			rel.R.Lang = o
		}
	}
	return nil
}

// SetLangCatalogDescriptionsG removes all previously related items of the
// language replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Lang's LangCatalogDescriptions accordingly.
// Replaces o.R.LangCatalogDescriptions with related.
// Sets related.R.Lang's LangCatalogDescriptions accordingly.
// Uses the global database handle.
func (o *Language) SetLangCatalogDescriptionsG(insert bool, related ...*CatalogDescription) error {
	return o.SetLangCatalogDescriptions(boil.GetDB(), insert, related...)
}

// SetLangCatalogDescriptionsP removes all previously related items of the
// language replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Lang's LangCatalogDescriptions accordingly.
// Replaces o.R.LangCatalogDescriptions with related.
// Sets related.R.Lang's LangCatalogDescriptions accordingly.
// Panics on error.
func (o *Language) SetLangCatalogDescriptionsP(exec boil.Executor, insert bool, related ...*CatalogDescription) {
	if err := o.SetLangCatalogDescriptions(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetLangCatalogDescriptionsGP removes all previously related items of the
// language replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Lang's LangCatalogDescriptions accordingly.
// Replaces o.R.LangCatalogDescriptions with related.
// Sets related.R.Lang's LangCatalogDescriptions accordingly.
// Uses the global database handle and panics on error.
func (o *Language) SetLangCatalogDescriptionsGP(insert bool, related ...*CatalogDescription) {
	if err := o.SetLangCatalogDescriptions(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetLangCatalogDescriptions removes all previously related items of the
// language replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Lang's LangCatalogDescriptions accordingly.
// Replaces o.R.LangCatalogDescriptions with related.
// Sets related.R.Lang's LangCatalogDescriptions accordingly.
func (o *Language) SetLangCatalogDescriptions(exec boil.Executor, insert bool, related ...*CatalogDescription) error {
	query := "update \"catalog_descriptions\" set \"lang_id\" = null where \"lang_id\" = $1"
	values := []interface{}{o.Code3}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.LangCatalogDescriptions {
			rel.LangID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Lang = nil
		}

		o.R.LangCatalogDescriptions = nil
	}
	return o.AddLangCatalogDescriptions(exec, insert, related...)
}

// RemoveLangCatalogDescriptionsG relationships from objects passed in.
// Removes related items from R.LangCatalogDescriptions (uses pointer comparison, removal does not keep order)
// Sets related.R.Lang.
// Uses the global database handle.
func (o *Language) RemoveLangCatalogDescriptionsG(related ...*CatalogDescription) error {
	return o.RemoveLangCatalogDescriptions(boil.GetDB(), related...)
}

// RemoveLangCatalogDescriptionsP relationships from objects passed in.
// Removes related items from R.LangCatalogDescriptions (uses pointer comparison, removal does not keep order)
// Sets related.R.Lang.
// Panics on error.
func (o *Language) RemoveLangCatalogDescriptionsP(exec boil.Executor, related ...*CatalogDescription) {
	if err := o.RemoveLangCatalogDescriptions(exec, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveLangCatalogDescriptionsGP relationships from objects passed in.
// Removes related items from R.LangCatalogDescriptions (uses pointer comparison, removal does not keep order)
// Sets related.R.Lang.
// Uses the global database handle and panics on error.
func (o *Language) RemoveLangCatalogDescriptionsGP(related ...*CatalogDescription) {
	if err := o.RemoveLangCatalogDescriptions(boil.GetDB(), related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveLangCatalogDescriptions relationships from objects passed in.
// Removes related items from R.LangCatalogDescriptions (uses pointer comparison, removal does not keep order)
// Sets related.R.Lang.
func (o *Language) RemoveLangCatalogDescriptions(exec boil.Executor, related ...*CatalogDescription) error {
	var err error
	for _, rel := range related {
		rel.LangID.Valid = false
		if rel.R != nil {
			rel.R.Lang = nil
		}
		if err = rel.Update(exec, "lang_id"); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.LangCatalogDescriptions {
			if rel != ri {
				continue
			}

			ln := len(o.R.LangCatalogDescriptions)
			if ln > 1 && i < ln-1 {
				o.R.LangCatalogDescriptions[i] = o.R.LangCatalogDescriptions[ln-1]
			}
			o.R.LangCatalogDescriptions = o.R.LangCatalogDescriptions[:ln-1]
			break
		}
	}

	return nil
}

// AddLangContainerDescriptionsG adds the given related objects to the existing relationships
// of the language, optionally inserting them as new records.
// Appends related to o.R.LangContainerDescriptions.
// Sets related.R.Lang appropriately.
// Uses the global database handle.
func (o *Language) AddLangContainerDescriptionsG(insert bool, related ...*ContainerDescription) error {
	return o.AddLangContainerDescriptions(boil.GetDB(), insert, related...)
}

// AddLangContainerDescriptionsP adds the given related objects to the existing relationships
// of the language, optionally inserting them as new records.
// Appends related to o.R.LangContainerDescriptions.
// Sets related.R.Lang appropriately.
// Panics on error.
func (o *Language) AddLangContainerDescriptionsP(exec boil.Executor, insert bool, related ...*ContainerDescription) {
	if err := o.AddLangContainerDescriptions(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddLangContainerDescriptionsGP adds the given related objects to the existing relationships
// of the language, optionally inserting them as new records.
// Appends related to o.R.LangContainerDescriptions.
// Sets related.R.Lang appropriately.
// Uses the global database handle and panics on error.
func (o *Language) AddLangContainerDescriptionsGP(insert bool, related ...*ContainerDescription) {
	if err := o.AddLangContainerDescriptions(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddLangContainerDescriptions adds the given related objects to the existing relationships
// of the language, optionally inserting them as new records.
// Appends related to o.R.LangContainerDescriptions.
// Sets related.R.Lang appropriately.
func (o *Language) AddLangContainerDescriptions(exec boil.Executor, insert bool, related ...*ContainerDescription) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.LangID.String = o.Code3.String
			rel.LangID.Valid = true
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"container_descriptions\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"lang_id"}),
				strmangle.WhereClause("\"", "\"", 2, containerDescriptionPrimaryKeyColumns),
			)
			values := []interface{}{o.Code3, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.LangID.String = o.Code3.String
			rel.LangID.Valid = true
		}
	}

	if o.R == nil {
		o.R = &languageR{
			LangContainerDescriptions: related,
		}
	} else {
		o.R.LangContainerDescriptions = append(o.R.LangContainerDescriptions, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &containerDescriptionR{
				Lang: o,
			}
		} else {
			rel.R.Lang = o
		}
	}
	return nil
}

// SetLangContainerDescriptionsG removes all previously related items of the
// language replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Lang's LangContainerDescriptions accordingly.
// Replaces o.R.LangContainerDescriptions with related.
// Sets related.R.Lang's LangContainerDescriptions accordingly.
// Uses the global database handle.
func (o *Language) SetLangContainerDescriptionsG(insert bool, related ...*ContainerDescription) error {
	return o.SetLangContainerDescriptions(boil.GetDB(), insert, related...)
}

// SetLangContainerDescriptionsP removes all previously related items of the
// language replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Lang's LangContainerDescriptions accordingly.
// Replaces o.R.LangContainerDescriptions with related.
// Sets related.R.Lang's LangContainerDescriptions accordingly.
// Panics on error.
func (o *Language) SetLangContainerDescriptionsP(exec boil.Executor, insert bool, related ...*ContainerDescription) {
	if err := o.SetLangContainerDescriptions(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetLangContainerDescriptionsGP removes all previously related items of the
// language replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Lang's LangContainerDescriptions accordingly.
// Replaces o.R.LangContainerDescriptions with related.
// Sets related.R.Lang's LangContainerDescriptions accordingly.
// Uses the global database handle and panics on error.
func (o *Language) SetLangContainerDescriptionsGP(insert bool, related ...*ContainerDescription) {
	if err := o.SetLangContainerDescriptions(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetLangContainerDescriptions removes all previously related items of the
// language replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Lang's LangContainerDescriptions accordingly.
// Replaces o.R.LangContainerDescriptions with related.
// Sets related.R.Lang's LangContainerDescriptions accordingly.
func (o *Language) SetLangContainerDescriptions(exec boil.Executor, insert bool, related ...*ContainerDescription) error {
	query := "update \"container_descriptions\" set \"lang_id\" = null where \"lang_id\" = $1"
	values := []interface{}{o.Code3}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.LangContainerDescriptions {
			rel.LangID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Lang = nil
		}

		o.R.LangContainerDescriptions = nil
	}
	return o.AddLangContainerDescriptions(exec, insert, related...)
}

// RemoveLangContainerDescriptionsG relationships from objects passed in.
// Removes related items from R.LangContainerDescriptions (uses pointer comparison, removal does not keep order)
// Sets related.R.Lang.
// Uses the global database handle.
func (o *Language) RemoveLangContainerDescriptionsG(related ...*ContainerDescription) error {
	return o.RemoveLangContainerDescriptions(boil.GetDB(), related...)
}

// RemoveLangContainerDescriptionsP relationships from objects passed in.
// Removes related items from R.LangContainerDescriptions (uses pointer comparison, removal does not keep order)
// Sets related.R.Lang.
// Panics on error.
func (o *Language) RemoveLangContainerDescriptionsP(exec boil.Executor, related ...*ContainerDescription) {
	if err := o.RemoveLangContainerDescriptions(exec, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveLangContainerDescriptionsGP relationships from objects passed in.
// Removes related items from R.LangContainerDescriptions (uses pointer comparison, removal does not keep order)
// Sets related.R.Lang.
// Uses the global database handle and panics on error.
func (o *Language) RemoveLangContainerDescriptionsGP(related ...*ContainerDescription) {
	if err := o.RemoveLangContainerDescriptions(boil.GetDB(), related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveLangContainerDescriptions relationships from objects passed in.
// Removes related items from R.LangContainerDescriptions (uses pointer comparison, removal does not keep order)
// Sets related.R.Lang.
func (o *Language) RemoveLangContainerDescriptions(exec boil.Executor, related ...*ContainerDescription) error {
	var err error
	for _, rel := range related {
		rel.LangID.Valid = false
		if rel.R != nil {
			rel.R.Lang = nil
		}
		if err = rel.Update(exec, "lang_id"); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.LangContainerDescriptions {
			if rel != ri {
				continue
			}

			ln := len(o.R.LangContainerDescriptions)
			if ln > 1 && i < ln-1 {
				o.R.LangContainerDescriptions[i] = o.R.LangContainerDescriptions[ln-1]
			}
			o.R.LangContainerDescriptions = o.R.LangContainerDescriptions[:ln-1]
			break
		}
	}

	return nil
}

// AddLangContainersG adds the given related objects to the existing relationships
// of the language, optionally inserting them as new records.
// Appends related to o.R.LangContainers.
// Sets related.R.Lang appropriately.
// Uses the global database handle.
func (o *Language) AddLangContainersG(insert bool, related ...*Container) error {
	return o.AddLangContainers(boil.GetDB(), insert, related...)
}

// AddLangContainersP adds the given related objects to the existing relationships
// of the language, optionally inserting them as new records.
// Appends related to o.R.LangContainers.
// Sets related.R.Lang appropriately.
// Panics on error.
func (o *Language) AddLangContainersP(exec boil.Executor, insert bool, related ...*Container) {
	if err := o.AddLangContainers(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddLangContainersGP adds the given related objects to the existing relationships
// of the language, optionally inserting them as new records.
// Appends related to o.R.LangContainers.
// Sets related.R.Lang appropriately.
// Uses the global database handle and panics on error.
func (o *Language) AddLangContainersGP(insert bool, related ...*Container) {
	if err := o.AddLangContainers(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddLangContainers adds the given related objects to the existing relationships
// of the language, optionally inserting them as new records.
// Appends related to o.R.LangContainers.
// Sets related.R.Lang appropriately.
func (o *Language) AddLangContainers(exec boil.Executor, insert bool, related ...*Container) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.LangID.String = o.Code3.String
			rel.LangID.Valid = true
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"containers\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"lang_id"}),
				strmangle.WhereClause("\"", "\"", 2, containerPrimaryKeyColumns),
			)
			values := []interface{}{o.Code3, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.LangID.String = o.Code3.String
			rel.LangID.Valid = true
		}
	}

	if o.R == nil {
		o.R = &languageR{
			LangContainers: related,
		}
	} else {
		o.R.LangContainers = append(o.R.LangContainers, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &containerR{
				Lang: o,
			}
		} else {
			rel.R.Lang = o
		}
	}
	return nil
}

// SetLangContainersG removes all previously related items of the
// language replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Lang's LangContainers accordingly.
// Replaces o.R.LangContainers with related.
// Sets related.R.Lang's LangContainers accordingly.
// Uses the global database handle.
func (o *Language) SetLangContainersG(insert bool, related ...*Container) error {
	return o.SetLangContainers(boil.GetDB(), insert, related...)
}

// SetLangContainersP removes all previously related items of the
// language replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Lang's LangContainers accordingly.
// Replaces o.R.LangContainers with related.
// Sets related.R.Lang's LangContainers accordingly.
// Panics on error.
func (o *Language) SetLangContainersP(exec boil.Executor, insert bool, related ...*Container) {
	if err := o.SetLangContainers(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetLangContainersGP removes all previously related items of the
// language replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Lang's LangContainers accordingly.
// Replaces o.R.LangContainers with related.
// Sets related.R.Lang's LangContainers accordingly.
// Uses the global database handle and panics on error.
func (o *Language) SetLangContainersGP(insert bool, related ...*Container) {
	if err := o.SetLangContainers(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetLangContainers removes all previously related items of the
// language replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Lang's LangContainers accordingly.
// Replaces o.R.LangContainers with related.
// Sets related.R.Lang's LangContainers accordingly.
func (o *Language) SetLangContainers(exec boil.Executor, insert bool, related ...*Container) error {
	query := "update \"containers\" set \"lang_id\" = null where \"lang_id\" = $1"
	values := []interface{}{o.Code3}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.LangContainers {
			rel.LangID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Lang = nil
		}

		o.R.LangContainers = nil
	}
	return o.AddLangContainers(exec, insert, related...)
}

// RemoveLangContainersG relationships from objects passed in.
// Removes related items from R.LangContainers (uses pointer comparison, removal does not keep order)
// Sets related.R.Lang.
// Uses the global database handle.
func (o *Language) RemoveLangContainersG(related ...*Container) error {
	return o.RemoveLangContainers(boil.GetDB(), related...)
}

// RemoveLangContainersP relationships from objects passed in.
// Removes related items from R.LangContainers (uses pointer comparison, removal does not keep order)
// Sets related.R.Lang.
// Panics on error.
func (o *Language) RemoveLangContainersP(exec boil.Executor, related ...*Container) {
	if err := o.RemoveLangContainers(exec, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveLangContainersGP relationships from objects passed in.
// Removes related items from R.LangContainers (uses pointer comparison, removal does not keep order)
// Sets related.R.Lang.
// Uses the global database handle and panics on error.
func (o *Language) RemoveLangContainersGP(related ...*Container) {
	if err := o.RemoveLangContainers(boil.GetDB(), related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveLangContainers relationships from objects passed in.
// Removes related items from R.LangContainers (uses pointer comparison, removal does not keep order)
// Sets related.R.Lang.
func (o *Language) RemoveLangContainers(exec boil.Executor, related ...*Container) error {
	var err error
	for _, rel := range related {
		rel.LangID.Valid = false
		if rel.R != nil {
			rel.R.Lang = nil
		}
		if err = rel.Update(exec, "lang_id"); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.LangContainers {
			if rel != ri {
				continue
			}

			ln := len(o.R.LangContainers)
			if ln > 1 && i < ln-1 {
				o.R.LangContainers[i] = o.R.LangContainers[ln-1]
			}
			o.R.LangContainers = o.R.LangContainers[:ln-1]
			break
		}
	}

	return nil
}

// AddLangFileAssetsG adds the given related objects to the existing relationships
// of the language, optionally inserting them as new records.
// Appends related to o.R.LangFileAssets.
// Sets related.R.Lang appropriately.
// Uses the global database handle.
func (o *Language) AddLangFileAssetsG(insert bool, related ...*FileAsset) error {
	return o.AddLangFileAssets(boil.GetDB(), insert, related...)
}

// AddLangFileAssetsP adds the given related objects to the existing relationships
// of the language, optionally inserting them as new records.
// Appends related to o.R.LangFileAssets.
// Sets related.R.Lang appropriately.
// Panics on error.
func (o *Language) AddLangFileAssetsP(exec boil.Executor, insert bool, related ...*FileAsset) {
	if err := o.AddLangFileAssets(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddLangFileAssetsGP adds the given related objects to the existing relationships
// of the language, optionally inserting them as new records.
// Appends related to o.R.LangFileAssets.
// Sets related.R.Lang appropriately.
// Uses the global database handle and panics on error.
func (o *Language) AddLangFileAssetsGP(insert bool, related ...*FileAsset) {
	if err := o.AddLangFileAssets(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddLangFileAssets adds the given related objects to the existing relationships
// of the language, optionally inserting them as new records.
// Appends related to o.R.LangFileAssets.
// Sets related.R.Lang appropriately.
func (o *Language) AddLangFileAssets(exec boil.Executor, insert bool, related ...*FileAsset) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.LangID.String = o.Code3.String
			rel.LangID.Valid = true
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"file_assets\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"lang_id"}),
				strmangle.WhereClause("\"", "\"", 2, fileAssetPrimaryKeyColumns),
			)
			values := []interface{}{o.Code3, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.LangID.String = o.Code3.String
			rel.LangID.Valid = true
		}
	}

	if o.R == nil {
		o.R = &languageR{
			LangFileAssets: related,
		}
	} else {
		o.R.LangFileAssets = append(o.R.LangFileAssets, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &fileAssetR{
				Lang: o,
			}
		} else {
			rel.R.Lang = o
		}
	}
	return nil
}

// SetLangFileAssetsG removes all previously related items of the
// language replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Lang's LangFileAssets accordingly.
// Replaces o.R.LangFileAssets with related.
// Sets related.R.Lang's LangFileAssets accordingly.
// Uses the global database handle.
func (o *Language) SetLangFileAssetsG(insert bool, related ...*FileAsset) error {
	return o.SetLangFileAssets(boil.GetDB(), insert, related...)
}

// SetLangFileAssetsP removes all previously related items of the
// language replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Lang's LangFileAssets accordingly.
// Replaces o.R.LangFileAssets with related.
// Sets related.R.Lang's LangFileAssets accordingly.
// Panics on error.
func (o *Language) SetLangFileAssetsP(exec boil.Executor, insert bool, related ...*FileAsset) {
	if err := o.SetLangFileAssets(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetLangFileAssetsGP removes all previously related items of the
// language replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Lang's LangFileAssets accordingly.
// Replaces o.R.LangFileAssets with related.
// Sets related.R.Lang's LangFileAssets accordingly.
// Uses the global database handle and panics on error.
func (o *Language) SetLangFileAssetsGP(insert bool, related ...*FileAsset) {
	if err := o.SetLangFileAssets(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetLangFileAssets removes all previously related items of the
// language replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Lang's LangFileAssets accordingly.
// Replaces o.R.LangFileAssets with related.
// Sets related.R.Lang's LangFileAssets accordingly.
func (o *Language) SetLangFileAssets(exec boil.Executor, insert bool, related ...*FileAsset) error {
	query := "update \"file_assets\" set \"lang_id\" = null where \"lang_id\" = $1"
	values := []interface{}{o.Code3}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.LangFileAssets {
			rel.LangID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Lang = nil
		}

		o.R.LangFileAssets = nil
	}
	return o.AddLangFileAssets(exec, insert, related...)
}

// RemoveLangFileAssetsG relationships from objects passed in.
// Removes related items from R.LangFileAssets (uses pointer comparison, removal does not keep order)
// Sets related.R.Lang.
// Uses the global database handle.
func (o *Language) RemoveLangFileAssetsG(related ...*FileAsset) error {
	return o.RemoveLangFileAssets(boil.GetDB(), related...)
}

// RemoveLangFileAssetsP relationships from objects passed in.
// Removes related items from R.LangFileAssets (uses pointer comparison, removal does not keep order)
// Sets related.R.Lang.
// Panics on error.
func (o *Language) RemoveLangFileAssetsP(exec boil.Executor, related ...*FileAsset) {
	if err := o.RemoveLangFileAssets(exec, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveLangFileAssetsGP relationships from objects passed in.
// Removes related items from R.LangFileAssets (uses pointer comparison, removal does not keep order)
// Sets related.R.Lang.
// Uses the global database handle and panics on error.
func (o *Language) RemoveLangFileAssetsGP(related ...*FileAsset) {
	if err := o.RemoveLangFileAssets(boil.GetDB(), related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveLangFileAssets relationships from objects passed in.
// Removes related items from R.LangFileAssets (uses pointer comparison, removal does not keep order)
// Sets related.R.Lang.
func (o *Language) RemoveLangFileAssets(exec boil.Executor, related ...*FileAsset) error {
	var err error
	for _, rel := range related {
		rel.LangID.Valid = false
		if rel.R != nil {
			rel.R.Lang = nil
		}
		if err = rel.Update(exec, "lang_id"); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.LangFileAssets {
			if rel != ri {
				continue
			}

			ln := len(o.R.LangFileAssets)
			if ln > 1 && i < ln-1 {
				o.R.LangFileAssets[i] = o.R.LangFileAssets[ln-1]
			}
			o.R.LangFileAssets = o.R.LangFileAssets[:ln-1]
			break
		}
	}

	return nil
}

// LanguagesG retrieves all records.
func LanguagesG(mods ...qm.QueryMod) languageQuery {
	return Languages(boil.GetDB(), mods...)
}

// Languages retrieves all the records using an executor.
func Languages(exec boil.Executor, mods ...qm.QueryMod) languageQuery {
	mods = append(mods, qm.From("\"languages\""))
	return languageQuery{NewQuery(exec, mods...)}
}

// FindLanguageG retrieves a single record by ID.
func FindLanguageG(id int, selectCols ...string) (*Language, error) {
	return FindLanguage(boil.GetDB(), id, selectCols...)
}

// FindLanguageGP retrieves a single record by ID, and panics on error.
func FindLanguageGP(id int, selectCols ...string) *Language {
	retobj, err := FindLanguage(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindLanguage retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindLanguage(exec boil.Executor, id int, selectCols ...string) (*Language, error) {
	languageObj := &Language{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"languages\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(languageObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "kmodels: unable to select from languages")
	}

	return languageObj, nil
}

// FindLanguageP retrieves a single record by ID with an executor, and panics on error.
func FindLanguageP(exec boil.Executor, id int, selectCols ...string) *Language {
	retobj, err := FindLanguage(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Language) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Language) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Language) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Language) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("kmodels: no languages provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(languageColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	languageInsertCacheMut.RLock()
	cache, cached := languageInsertCache[key]
	languageInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			languageColumns,
			languageColumnsWithDefault,
			languageColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(languageType, languageMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(languageType, languageMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"languages\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "kmodels: unable to insert into languages")
	}

	if !cached {
		languageInsertCacheMut.Lock()
		languageInsertCache[key] = cache
		languageInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single Language record. See Update for
// whitelist behavior description.
func (o *Language) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Language record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Language) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Language, and panics on error.
// See Update for whitelist behavior description.
func (o *Language) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Language.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Language) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	key := makeCacheKey(whitelist, nil)
	languageUpdateCacheMut.RLock()
	cache, cached := languageUpdateCache[key]
	languageUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(languageColumns, languagePrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("kmodels: unable to update languages, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"languages\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, languagePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(languageType, languageMapping, append(wl, languagePrimaryKeyColumns...))
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
		return errors.Wrap(err, "kmodels: unable to update languages row")
	}

	if !cached {
		languageUpdateCacheMut.Lock()
		languageUpdateCache[key] = cache
		languageUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q languageQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q languageQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to update all for languages")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o LanguageSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o LanguageSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o LanguageSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o LanguageSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("kmodels: update all requires at least one column argument")
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), languagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"languages\" SET %s WHERE (\"id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(languagePrimaryKeyColumns), len(colNames)+1, len(languagePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to update all in language slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Language) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Language) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Language) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Language) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("kmodels: no languages provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(languageColumnsWithDefault, o)

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

	languageUpsertCacheMut.RLock()
	cache, cached := languageUpsertCache[key]
	languageUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			languageColumns,
			languageColumnsWithDefault,
			languageColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			languageColumns,
			languagePrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("kmodels: unable to upsert languages, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(languagePrimaryKeyColumns))
			copy(conflict, languagePrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"languages\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(languageType, languageMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(languageType, languageMapping, ret)
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
		return errors.Wrap(err, "kmodels: unable to upsert languages")
	}

	if !cached {
		languageUpsertCacheMut.Lock()
		languageUpsertCache[key] = cache
		languageUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single Language record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Language) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Language record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Language) DeleteG() error {
	if o == nil {
		return errors.New("kmodels: no Language provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Language record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Language) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Language record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Language) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("kmodels: no Language provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), languagePrimaryKeyMapping)
	sql := "DELETE FROM \"languages\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to delete from languages")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q languageQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q languageQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("kmodels: no languageQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to delete all from languages")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o LanguageSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o LanguageSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("kmodels: no Language slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o LanguageSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o LanguageSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("kmodels: no Language slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), languagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"languages\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, languagePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(languagePrimaryKeyColumns), 1, len(languagePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to delete all from language slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Language) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Language) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Language) ReloadG() error {
	if o == nil {
		return errors.New("kmodels: no Language provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Language) Reload(exec boil.Executor) error {
	ret, err := FindLanguage(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *LanguageSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *LanguageSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *LanguageSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("kmodels: empty LanguageSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *LanguageSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	languages := LanguageSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), languagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"languages\".* FROM \"languages\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, languagePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(languagePrimaryKeyColumns), 1, len(languagePrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&languages)
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to reload all in LanguageSlice")
	}

	*o = languages

	return nil
}

// LanguageExists checks if the Language row exists.
func LanguageExists(exec boil.Executor, id int) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"languages\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "kmodels: unable to check if languages exists")
	}

	return exists, nil
}

// LanguageExistsG checks if the Language row exists.
func LanguageExistsG(id int) (bool, error) {
	return LanguageExists(boil.GetDB(), id)
}

// LanguageExistsGP checks if the Language row exists. Panics on error.
func LanguageExistsGP(id int) bool {
	e, err := LanguageExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// LanguageExistsP checks if the Language row exists. Panics on error.
func LanguageExistsP(exec boil.Executor, id int) bool {
	e, err := LanguageExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
