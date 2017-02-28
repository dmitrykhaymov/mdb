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

// FileType is an object representing the database table.
type FileType struct {
	Name    string      `boil:"name" json:"name" toml:"name" yaml:"name"`
	Extlist null.String `boil:"extlist" json:"extlist,omitempty" toml:"extlist" yaml:"extlist,omitempty"`
	Pic     null.String `boil:"pic" json:"pic,omitempty" toml:"pic" yaml:"pic,omitempty"`

	R *fileTypeR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L fileTypeL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// fileTypeR is where relationships are stored.
type fileTypeR struct {
}

// fileTypeL is where Load methods for each relationship are stored.
type fileTypeL struct{}

var (
	fileTypeColumns               = []string{"name", "extlist", "pic"}
	fileTypeColumnsWithoutDefault = []string{"extlist", "pic"}
	fileTypeColumnsWithDefault    = []string{"name"}
	fileTypePrimaryKeyColumns     = []string{"name"}
)

type (
	// FileTypeSlice is an alias for a slice of pointers to FileType.
	// This should generally be used opposed to []FileType.
	FileTypeSlice []*FileType

	fileTypeQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	fileTypeType                 = reflect.TypeOf(&FileType{})
	fileTypeMapping              = queries.MakeStructMapping(fileTypeType)
	fileTypePrimaryKeyMapping, _ = queries.BindMapping(fileTypeType, fileTypeMapping, fileTypePrimaryKeyColumns)
	fileTypeInsertCacheMut       sync.RWMutex
	fileTypeInsertCache          = make(map[string]insertCache)
	fileTypeUpdateCacheMut       sync.RWMutex
	fileTypeUpdateCache          = make(map[string]updateCache)
	fileTypeUpsertCacheMut       sync.RWMutex
	fileTypeUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single fileType record from the query, and panics on error.
func (q fileTypeQuery) OneP() *FileType {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single fileType record from the query.
func (q fileTypeQuery) One() (*FileType, error) {
	o := &FileType{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "kmodels: failed to execute a one query for file_types")
	}

	return o, nil
}

// AllP returns all FileType records from the query, and panics on error.
func (q fileTypeQuery) AllP() FileTypeSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all FileType records from the query.
func (q fileTypeQuery) All() (FileTypeSlice, error) {
	var o FileTypeSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "kmodels: failed to assign all query results to FileType slice")
	}

	return o, nil
}

// CountP returns the count of all FileType records in the query, and panics on error.
func (q fileTypeQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all FileType records in the query.
func (q fileTypeQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "kmodels: failed to count file_types rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q fileTypeQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q fileTypeQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "kmodels: failed to check if file_types exists")
	}

	return count > 0, nil
}

// FileTypesG retrieves all records.
func FileTypesG(mods ...qm.QueryMod) fileTypeQuery {
	return FileTypes(boil.GetDB(), mods...)
}

// FileTypes retrieves all the records using an executor.
func FileTypes(exec boil.Executor, mods ...qm.QueryMod) fileTypeQuery {
	mods = append(mods, qm.From("\"file_types\""))
	return fileTypeQuery{NewQuery(exec, mods...)}
}

// FindFileTypeG retrieves a single record by ID.
func FindFileTypeG(name string, selectCols ...string) (*FileType, error) {
	return FindFileType(boil.GetDB(), name, selectCols...)
}

// FindFileTypeGP retrieves a single record by ID, and panics on error.
func FindFileTypeGP(name string, selectCols ...string) *FileType {
	retobj, err := FindFileType(boil.GetDB(), name, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindFileType retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindFileType(exec boil.Executor, name string, selectCols ...string) (*FileType, error) {
	fileTypeObj := &FileType{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"file_types\" where \"name\"=$1", sel,
	)

	q := queries.Raw(exec, query, name)

	err := q.Bind(fileTypeObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "kmodels: unable to select from file_types")
	}

	return fileTypeObj, nil
}

// FindFileTypeP retrieves a single record by ID with an executor, and panics on error.
func FindFileTypeP(exec boil.Executor, name string, selectCols ...string) *FileType {
	retobj, err := FindFileType(exec, name, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *FileType) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *FileType) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *FileType) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *FileType) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("kmodels: no file_types provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(fileTypeColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	fileTypeInsertCacheMut.RLock()
	cache, cached := fileTypeInsertCache[key]
	fileTypeInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			fileTypeColumns,
			fileTypeColumnsWithDefault,
			fileTypeColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(fileTypeType, fileTypeMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(fileTypeType, fileTypeMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"file_types\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "kmodels: unable to insert into file_types")
	}

	if !cached {
		fileTypeInsertCacheMut.Lock()
		fileTypeInsertCache[key] = cache
		fileTypeInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single FileType record. See Update for
// whitelist behavior description.
func (o *FileType) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single FileType record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *FileType) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the FileType, and panics on error.
// See Update for whitelist behavior description.
func (o *FileType) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the FileType.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *FileType) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	key := makeCacheKey(whitelist, nil)
	fileTypeUpdateCacheMut.RLock()
	cache, cached := fileTypeUpdateCache[key]
	fileTypeUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(fileTypeColumns, fileTypePrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("kmodels: unable to update file_types, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"file_types\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, fileTypePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(fileTypeType, fileTypeMapping, append(wl, fileTypePrimaryKeyColumns...))
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
		return errors.Wrap(err, "kmodels: unable to update file_types row")
	}

	if !cached {
		fileTypeUpdateCacheMut.Lock()
		fileTypeUpdateCache[key] = cache
		fileTypeUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q fileTypeQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q fileTypeQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to update all for file_types")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o FileTypeSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o FileTypeSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o FileTypeSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o FileTypeSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), fileTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"file_types\" SET %s WHERE (\"name\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(fileTypePrimaryKeyColumns), len(colNames)+1, len(fileTypePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to update all in fileType slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *FileType) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *FileType) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *FileType) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *FileType) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("kmodels: no file_types provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(fileTypeColumnsWithDefault, o)

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

	fileTypeUpsertCacheMut.RLock()
	cache, cached := fileTypeUpsertCache[key]
	fileTypeUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			fileTypeColumns,
			fileTypeColumnsWithDefault,
			fileTypeColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			fileTypeColumns,
			fileTypePrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("kmodels: unable to upsert file_types, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(fileTypePrimaryKeyColumns))
			copy(conflict, fileTypePrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"file_types\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(fileTypeType, fileTypeMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(fileTypeType, fileTypeMapping, ret)
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
		return errors.Wrap(err, "kmodels: unable to upsert file_types")
	}

	if !cached {
		fileTypeUpsertCacheMut.Lock()
		fileTypeUpsertCache[key] = cache
		fileTypeUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single FileType record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *FileType) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single FileType record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *FileType) DeleteG() error {
	if o == nil {
		return errors.New("kmodels: no FileType provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single FileType record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *FileType) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single FileType record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *FileType) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("kmodels: no FileType provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), fileTypePrimaryKeyMapping)
	sql := "DELETE FROM \"file_types\" WHERE \"name\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to delete from file_types")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q fileTypeQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q fileTypeQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("kmodels: no fileTypeQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to delete all from file_types")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o FileTypeSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o FileTypeSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("kmodels: no FileType slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o FileTypeSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o FileTypeSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("kmodels: no FileType slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), fileTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"file_types\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, fileTypePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(fileTypePrimaryKeyColumns), 1, len(fileTypePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to delete all from fileType slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *FileType) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *FileType) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *FileType) ReloadG() error {
	if o == nil {
		return errors.New("kmodels: no FileType provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *FileType) Reload(exec boil.Executor) error {
	ret, err := FindFileType(exec, o.Name)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *FileTypeSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *FileTypeSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *FileTypeSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("kmodels: empty FileTypeSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *FileTypeSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	fileTypes := FileTypeSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), fileTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"file_types\".* FROM \"file_types\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, fileTypePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(fileTypePrimaryKeyColumns), 1, len(fileTypePrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&fileTypes)
	if err != nil {
		return errors.Wrap(err, "kmodels: unable to reload all in FileTypeSlice")
	}

	*o = fileTypes

	return nil
}

// FileTypeExists checks if the FileType row exists.
func FileTypeExists(exec boil.Executor, name string) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"file_types\" where \"name\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, name)
	}

	row := exec.QueryRow(sql, name)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "kmodels: unable to check if file_types exists")
	}

	return exists, nil
}

// FileTypeExistsG checks if the FileType row exists.
func FileTypeExistsG(name string) (bool, error) {
	return FileTypeExists(boil.GetDB(), name)
}

// FileTypeExistsGP checks if the FileType row exists. Panics on error.
func FileTypeExistsGP(name string) bool {
	e, err := FileTypeExists(boil.GetDB(), name)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// FileTypeExistsP checks if the FileType row exists. Panics on error.
func FileTypeExistsP(exec boil.Executor, name string) bool {
	e, err := FileTypeExists(exec, name)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
