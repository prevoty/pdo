package pdo

import (
	"bytes"
	stdlib_sql "database/sql"
	"fmt"
	"reflect"
	"strings"
)

type (
	DBO struct {
		DB *stdlib_sql.DB
	}
)

var (
	EmptySkipList = []string{}
)

func (d *DBO) Create(record_ptr interface{}) (int64, error) {

	switch t := reflect.TypeOf(record_ptr); {
	case t.Kind() != reflect.Ptr,
		t.Elem().Kind() != reflect.Struct:
		panic("only *Struct value allowed")
	}

	var sql bytes.Buffer

	sql.WriteString("INSERT INTO `")
	sql.WriteString(Table(record_ptr))
	sql.WriteString("` (")

	for _, col := range Columns(record_ptr, []string{"id"}) {
		sql.WriteString("`")
		sql.WriteString(col)
		sql.WriteString("`,")
	}

	// trim off last comma
	sql.Truncate(sql.Len() - 1)

	sql.WriteString(") VALUES (")

	values := FieldPointers(record_ptr, []string{"id"})

	for i := 0; i < len(values); i++ {
		sql.WriteString("?,")
	}

	// trim off last comma
	sql.Truncate(sql.Len() - 1)

	sql.WriteString(")")

	stmt, err := d.DB.Prepare(sql.String())
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	r, err := stmt.Exec(values...)
	if err != nil {
		return 0, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil

}

func (d *DBO) Update(record_ptr interface{}) error {

	switch t := reflect.TypeOf(record_ptr); {
	case t.Kind() != reflect.Ptr,
		t.Elem().Kind() != reflect.Struct:
		panic("only *Struct value allowed")
	}

	var sql bytes.Buffer

	sql.WriteString("UPDATE `")
	sql.WriteString(Table(record_ptr))
	sql.WriteString("` SET ")

	for _, col := range Columns(record_ptr, []string{"id"}) {
		sql.WriteString("`")
		sql.WriteString(col)
		sql.WriteString("` = ?,")
	}

	// trim off last comma
	sql.Truncate(sql.Len() - 1)

	sql.WriteString(" WHERE `id` = ?")

	stmt, err := d.DB.Prepare(sql.String())
	if err != nil {
		return err
	}
	defer stmt.Close()

	// grab id value
	id := reflect.ValueOf(record_ptr).Elem().FieldByName("Id").Interface()

	// add id to end of values list
	values := append(FieldPointers(record_ptr, []string{"id"}), id)

	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}

	return nil

}

func (d *DBO) Delete(record_ptr interface{}) error {

	switch t := reflect.TypeOf(record_ptr); {
	case t.Kind() != reflect.Ptr,
		t.Elem().Kind() != reflect.Struct:
		panic("only *Struct value allowed")
	}

	var sql bytes.Buffer

	sql.WriteString("DELETE FROM `")
	sql.WriteString(Table(record_ptr))
	sql.WriteString("` WHERE `id` = ?")

	stmt, err := d.DB.Prepare(sql.String())
	if err != nil {
		return err
	}
	defer stmt.Close()

	id := reflect.ValueOf(record_ptr).Elem().FieldByName("Id").Interface()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil

}

func (d *DBO) Find(record_ptr interface{}, where string, params ...interface{}) error {

	switch t := reflect.TypeOf(record_ptr); {
	case t.Kind() != reflect.Ptr,
		t.Elem().Kind() != reflect.Struct:
		panic("only *Struct value allowed")
	}

	var sql bytes.Buffer

	sql.WriteString("SELECT ")

	for _, col := range Columns(record_ptr, EmptySkipList) {
		sql.WriteString("`")
		sql.WriteString(col)
		sql.WriteString("`,")
	}

	// trim off last comma
	sql.Truncate(sql.Len() - 1)

	sql.WriteString(" FROM `")
	sql.WriteString(Table(record_ptr))
	sql.WriteString("` ")
	sql.WriteString(where)

	return d.DB.QueryRow(sql.String(), params...).Scan(FieldPointers(record_ptr, EmptySkipList)...)

}

func (d *DBO) FindAll(record_set_ptr interface{}, where string, params ...interface{}) error {

	var (
		sql bytes.Buffer
		t   = reflect.TypeOf(record_set_ptr)
		v   = reflect.ValueOf(record_set_ptr)
	)

	switch {
	case t.Kind() != reflect.Ptr,
		t.Elem().Kind() != reflect.Slice,
		t.Elem().Elem().Kind() != reflect.Ptr,
		t.Elem().Elem().Elem().Kind() != reflect.Struct:
		panic("only *[]*Struct value allowed")
	}

	struct_val := t.Elem().Elem().Elem()
	struct_clone := reflect.New(struct_val).Interface()

	sql.WriteString("SELECT ")

	for _, col := range Columns(struct_clone, EmptySkipList) {
		sql.WriteString("`")
		sql.WriteString(col)
		sql.WriteString("`,")
	}

	// trim off last comma
	sql.Truncate(sql.Len() - 1)

	sql.WriteString(" FROM `")
	sql.WriteString(Table(struct_clone))
	sql.WriteString("` ")
	sql.WriteString(where)

	rows, err := d.DB.Query(sql.String(), params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// dereference the *[]*Type array
	rowset := v.Elem()

	for rows.Next() {
		// new struct per row
		new_struct := reflect.New(struct_val).Interface()

		err = rows.Scan(FieldPointers(new_struct, EmptySkipList)...)
		if err != nil {
			return err
		}

		rowset = reflect.Append(rowset, reflect.ValueOf(new_struct))
	}

	v.Elem().Set(rowset)

	return nil

}

func (d *DBO) Close() error {

	return d.DB.Close()

}

func Table(record_ptr interface{}) string {
	field, ok := reflect.ValueOf(record_ptr).Elem().Type().FieldByName("_meta")
	if !ok {
		panic("_meta member missing on struct definition for: " + fmt.Sprintf("%#v\n", record_ptr))
	}
	table := field.Tag.Get("table")
	if table == "" {
		panic("table tag on _meta missing on struct definition for: " + fmt.Sprintf("%#v\n", record_ptr))
	}
	return table
}

func Columns(record_ptr interface{}, skip_list []string) []string {

	var (
		v    = reflect.ValueOf(record_ptr).Elem()
		t    = v.Type()
		cols = make([]string, 0, v.NumField())
	)

	for i := 0; i < v.NumField(); i++ {

		field := t.Field(i)

		if inSlice(skip_list, field.Name) {
			continue
		}

		col := field.Tag.Get("column")
		if col == "" {
			continue
		}

		cols = append(cols, col)

	}

	return cols

}

func FieldPointers(record_ptr interface{}, skip_list []string) []interface{} {

	var (
		v        = reflect.ValueOf(record_ptr).Elem()
		t        = v.Type()
		pointers = make([]interface{}, 0, v.NumField())
	)

	for i := 0; i < v.NumField(); i++ {

		field := t.Field(i)

		if inSlice(skip_list, field.Name) {
			continue
		}

		col := field.Tag.Get("column")
		if col == "" {
			continue
		}

		pointers = append(pointers, v.Field(i).Addr().Interface())

	}

	return pointers

}

// case insensitive search
func inSlice(haystack []string, needle string) bool {
	for _, v := range haystack {
		if strings.EqualFold(v, needle) {
			return true
		}
	}
	return false
}
