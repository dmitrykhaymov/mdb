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
	"gopkg.in/nullbio/null.v6"
)

// OperationType is an object representing the database table.
type OperationType struct {
	ID          int64       `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name        string      `boil:"name" json:"name" toml:"name" yaml:"name"`
	Description null.String `boil:"description" json:"description,omitempty" toml:"description" yaml:"description,omitempty"`

	R *operationTypeR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L operationTypeL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// operationTypeR is where relationships are stored.
type operationTypeR struct {
	TypeOperations OperationSlice
}

// operationTypeL is where Load methods for each relationship are stored.
type operationTypeL struct{}

var (
	operationTypeColumns               = []string{"id", "name", "description"}
	operationTypeColumnsWithoutDefault = []string{"name", "description"}
	operationTypeColumnsWithDefault    = []string{"id"}
	operationTypePrimaryKeyColumns     = []string{"id"}
)

type (
	// OperationTypeSlice is an alias for a slice of pointers to OperationType.
	// This should generally be used opposed to []OperationType.
	OperationTypeSlice []*OperationType

	operationTypeQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	operationTypeType                 = reflect.TypeOf(&OperationType{})
	operationTypeMapping              = queries.MakeStructMapping(operationTypeType)
	operationTypePrimaryKeyMapping, _ = queries.BindMapping(operationTypeType, operationTypeMapping, operationTypePrimaryKeyColumns)
	operationTypeInsertCacheMut       sync.RWMutex
	operationTypeInsertCache          = make(map[string]insertCache)
	operationTypeUpdateCacheMut       sync.RWMutex
	operationTypeUpdateCache          = make(map[string]updateCache)
	operationTypeUpsertCacheMut       sync.RWMutex
	operationTypeUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single operationType record from the query, and panics on error.
func (q operationTypeQuery) OneP() *OperationType {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single operationType record from the query.
func (q operationTypeQuery) One() (*OperationType, error) {
	o := &OperationType{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "gmodels: failed to execute a one query for operation_types")
	}

	return o, nil
}

// AllP returns all OperationType records from the query, and panics on error.
func (q operationTypeQuery) AllP() OperationTypeSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all OperationType records from the query.
func (q operationTypeQuery) All() (OperationTypeSlice, error) {
	var o OperationTypeSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "gmodels: failed to assign all query results to OperationType slice")
	}

	return o, nil
}

// CountP returns the count of all OperationType records in the query, and panics on error.
func (q operationTypeQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all OperationType records in the query.
func (q operationTypeQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "gmodels: failed to count operation_types rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q operationTypeQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q operationTypeQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "gmodels: failed to check if operation_types exists")
	}

	return count > 0, nil
}

// TypeOperationsG retrieves all the operation's operations via type_id column.
func (o *OperationType) TypeOperationsG(mods ...qm.QueryMod) operationQuery {
	return o.TypeOperations(boil.GetDB(), mods...)
}

// TypeOperations retrieves all the operation's operations with an executor via type_id column.
func (o *OperationType) TypeOperations(exec boil.Executor, mods ...qm.QueryMod) operationQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"type_id\"=?", o.ID),
	)

	query := Operations(exec, queryMods...)
	queries.SetFrom(query.Query, "\"operations\" as \"a\"")
	return query
}

// LoadTypeOperations allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (operationTypeL) LoadTypeOperations(e boil.Executor, singular bool, maybeOperationType interface{}) error {
	var slice []*OperationType
	var object *OperationType

	count := 1
	if singular {
		object = maybeOperationType.(*OperationType)
	} else {
		slice = *maybeOperationType.(*OperationTypeSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &operationTypeR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &operationTypeR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"operations\" where \"type_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load operations")
	}
	defer results.Close()

	var resultSlice []*Operation
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice operations")
	}

	if singular {
		object.R.TypeOperations = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.TypeID {
				local.R.TypeOperations = append(local.R.TypeOperations, foreign)
				break
			}
		}
	}

	return nil
}

// AddTypeOperations adds the given related objects to the existing relationships
// of the operation_type, optionally inserting them as new records.
// Appends related to o.R.TypeOperations.
// Sets related.R.Type appropriately.
func (o *OperationType) AddTypeOperations(exec boil.Executor, insert bool, related ...*Operation) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.TypeID = o.ID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"operations\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"type_id"}),
				strmangle.WhereClause("\"", "\"", 2, operationPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.TypeID = o.ID
		}
	}

	if o.R == nil {
		o.R = &operationTypeR{
			TypeOperations: related,
		}
	} else {
		o.R.TypeOperations = append(o.R.TypeOperations, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &operationR{
				Type: o,
			}
		} else {
			rel.R.Type = o
		}
	}
	return nil
}

// OperationTypesG retrieves all records.
func OperationTypesG(mods ...qm.QueryMod) operationTypeQuery {
	return OperationTypes(boil.GetDB(), mods...)
}

// OperationTypes retrieves all the records using an executor.
func OperationTypes(exec boil.Executor, mods ...qm.QueryMod) operationTypeQuery {
	mods = append(mods, qm.From("\"operation_types\""))
	return operationTypeQuery{NewQuery(exec, mods...)}
}

// FindOperationTypeG retrieves a single record by ID.
func FindOperationTypeG(id int64, selectCols ...string) (*OperationType, error) {
	return FindOperationType(boil.GetDB(), id, selectCols...)
}

// FindOperationTypeGP retrieves a single record by ID, and panics on error.
func FindOperationTypeGP(id int64, selectCols ...string) *OperationType {
	retobj, err := FindOperationType(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindOperationType retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindOperationType(exec boil.Executor, id int64, selectCols ...string) (*OperationType, error) {
	operationTypeObj := &OperationType{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"operation_types\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(operationTypeObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "gmodels: unable to select from operation_types")
	}

	return operationTypeObj, nil
}

// FindOperationTypeP retrieves a single record by ID with an executor, and panics on error.
func FindOperationTypeP(exec boil.Executor, id int64, selectCols ...string) *OperationType {
	retobj, err := FindOperationType(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *OperationType) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *OperationType) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *OperationType) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *OperationType) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("gmodels: no operation_types provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(operationTypeColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	operationTypeInsertCacheMut.RLock()
	cache, cached := operationTypeInsertCache[key]
	operationTypeInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			operationTypeColumns,
			operationTypeColumnsWithDefault,
			operationTypeColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(operationTypeType, operationTypeMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(operationTypeType, operationTypeMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"operation_types\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "gmodels: unable to insert into operation_types")
	}

	if !cached {
		operationTypeInsertCacheMut.Lock()
		operationTypeInsertCache[key] = cache
		operationTypeInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single OperationType record. See Update for
// whitelist behavior description.
func (o *OperationType) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single OperationType record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *OperationType) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the OperationType, and panics on error.
// See Update for whitelist behavior description.
func (o *OperationType) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the OperationType.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *OperationType) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	key := makeCacheKey(whitelist, nil)
	operationTypeUpdateCacheMut.RLock()
	cache, cached := operationTypeUpdateCache[key]
	operationTypeUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(operationTypeColumns, operationTypePrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("gmodels: unable to update operation_types, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"operation_types\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, operationTypePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(operationTypeType, operationTypeMapping, append(wl, operationTypePrimaryKeyColumns...))
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
		return errors.Wrap(err, "gmodels: unable to update operation_types row")
	}

	if !cached {
		operationTypeUpdateCacheMut.Lock()
		operationTypeUpdateCache[key] = cache
		operationTypeUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q operationTypeQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q operationTypeQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to update all for operation_types")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o OperationTypeSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o OperationTypeSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o OperationTypeSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o OperationTypeSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), operationTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"operation_types\" SET %s WHERE (\"id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(operationTypePrimaryKeyColumns), len(colNames)+1, len(operationTypePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to update all in operationType slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *OperationType) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *OperationType) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *OperationType) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *OperationType) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("gmodels: no operation_types provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(operationTypeColumnsWithDefault, o)

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

	operationTypeUpsertCacheMut.RLock()
	cache, cached := operationTypeUpsertCache[key]
	operationTypeUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			operationTypeColumns,
			operationTypeColumnsWithDefault,
			operationTypeColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			operationTypeColumns,
			operationTypePrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("gmodels: unable to upsert operation_types, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(operationTypePrimaryKeyColumns))
			copy(conflict, operationTypePrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"operation_types\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(operationTypeType, operationTypeMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(operationTypeType, operationTypeMapping, ret)
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
		return errors.Wrap(err, "gmodels: unable to upsert operation_types")
	}

	if !cached {
		operationTypeUpsertCacheMut.Lock()
		operationTypeUpsertCache[key] = cache
		operationTypeUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single OperationType record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *OperationType) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single OperationType record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *OperationType) DeleteG() error {
	if o == nil {
		return errors.New("gmodels: no OperationType provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single OperationType record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *OperationType) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single OperationType record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *OperationType) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("gmodels: no OperationType provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), operationTypePrimaryKeyMapping)
	sql := "DELETE FROM \"operation_types\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to delete from operation_types")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q operationTypeQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q operationTypeQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("gmodels: no operationTypeQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to delete all from operation_types")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o OperationTypeSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o OperationTypeSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("gmodels: no OperationType slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o OperationTypeSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o OperationTypeSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("gmodels: no OperationType slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), operationTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"operation_types\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, operationTypePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(operationTypePrimaryKeyColumns), 1, len(operationTypePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to delete all from operationType slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *OperationType) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *OperationType) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *OperationType) ReloadG() error {
	if o == nil {
		return errors.New("gmodels: no OperationType provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *OperationType) Reload(exec boil.Executor) error {
	ret, err := FindOperationType(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *OperationTypeSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *OperationTypeSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *OperationTypeSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("gmodels: empty OperationTypeSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *OperationTypeSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	operationTypes := OperationTypeSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), operationTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"operation_types\".* FROM \"operation_types\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, operationTypePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(operationTypePrimaryKeyColumns), 1, len(operationTypePrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&operationTypes)
	if err != nil {
		return errors.Wrap(err, "gmodels: unable to reload all in OperationTypeSlice")
	}

	*o = operationTypes

	return nil
}

// OperationTypeExists checks if the OperationType row exists.
func OperationTypeExists(exec boil.Executor, id int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"operation_types\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "gmodels: unable to check if operation_types exists")
	}

	return exists, nil
}

// OperationTypeExistsG checks if the OperationType row exists.
func OperationTypeExistsG(id int64) (bool, error) {
	return OperationTypeExists(boil.GetDB(), id)
}

// OperationTypeExistsGP checks if the OperationType row exists. Panics on error.
func OperationTypeExistsGP(id int64) bool {
	e, err := OperationTypeExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// OperationTypeExistsP checks if the OperationType row exists. Panics on error.
func OperationTypeExistsP(exec boil.Executor, id int64) bool {
	e, err := OperationTypeExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
