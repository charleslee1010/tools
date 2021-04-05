package db

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
	//	"reflect"
	"errors"
	toolkit "github.com/charles/toolkit"
	"io"
	"strconv"
)

type Database struct {
	DB *sql.DB
}

// for sql builder, field name of updatetime
var updatetime string = "updatetime"

//var createtime string = "cupdatetime"

type RecordWriter interface {
	Write(m map[string]interface{})
}

func NewDatabase(dbstr string) (*Database, error) {
	db, err := sql.Open("mysql", dbstr) //"root:Daily.2017@tcp(10.0.10.87:3306)/xxxc")

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Second)

	return &Database{DB: db}, nil
}

func (this *Database) Close() {
	if this.DB != nil {
		this.DB.Close()
	}
}

type AllRecords struct {
	a []map[string]interface{}
}

func (w *AllRecords) Write(m map[string]interface{}) {
	w.a = append(w.a, m)
}

func (this *Database) QueryDBReturnAll(sqlString string, args ...interface{}) ([]map[string]interface{}, error) {

	d := &AllRecords{a: make([]map[string]interface{}, 0)}

	err := this.QueryDBCallBack(d, sqlString, args...)

	return d.a, err
}

func (this *Database) QueryDBCallBack(rw RecordWriter, sqlString string, args ...interface{}) error {

	//	fmt.Println(sqlString)
	//	fmt.Println(args)

	rows, err := this.DB.Query(sqlString, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	//	fmt.Println("query ends")

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	count := len(columns)

	//
	//	ctypes, err := rows.ColumnTypes()
	//
	//	for j:= 0 ; j < count; j++ {
	//		fmt.Println("column name:", columns[j])
	//		fmt.Println("column type:", ctypes[j].DatabaseTypeName())
	//	}
	//

	//	fmt.Println("query return ", count)
	//	tableData := make([]map[string]interface{}, 0)

	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	var lines int = 0

	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]

			//			if val != nil {
			//				fmt.Println("col=", col, ", value=", val, ", type=", reflect.TypeOf(val).String(),
			//				   " type=", )
			//			}

			//			switch ctypes[i].DatabaseTypeName()

			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}

			//			if val != nil {
			//				fmt.Println("col=", col, ", value=", val, ", type=", reflect.TypeOf(val).String())
			//			}
			//
			entry[col] = v
		}

		lines++
		// we have a row in entry
		rw.Write(entry)
		//		if callback(entry) == false {
		//			log.Error("User aborted...")
		//			return false
		//		}
		//		fmt.Println(entry)
	}
	//	fmt.Println("return ", lines, " records")
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

func (this *Database) QueryDBReturnAllNew(sqlString string, args ...interface{}) ([]map[string]interface{}, error) {

	d := &AllRecords{a: make([]map[string]interface{}, 0)}

	err := this.QueryDBCallBackNew(d, sqlString, args...)

	return d.a, err
}

func (this *Database) QueryDBCallBackNew(rw RecordWriter, sqlString string, args ...interface{}) error {

	//	fmt.Printf("sql=%s, args=%+v\n", sqlString, args)

	rows, err := this.DB.Query(sqlString, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	//	fmt.Println("query ends")

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	count := len(columns)

	ctypes, err := rows.ColumnTypes()
	if err != nil {
		return err
	}

	//	for j:= 0 ; j < count; j++ {
	//		fmt.Println("column name:", columns[j])
	//		fmt.Println("column type:", ctypes[j].DatabaseTypeName())
	//	}

	//	fmt.Println("query return ", count)
	//	tableData := make([]map[string]interface{}, 0)

	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	var lines int = 0

	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]

			//			if val != nil {
			//				fmt.Println("col=", col, ", value=", val, ", type=", reflect.TypeOf(val).String(),
			//				   " type=", )
			//			}

			b, ok := val.([]byte)
			if ok {
				switch ctypes[i].DatabaseTypeName() {
				case "DECIMAL":
					v, _ = strconv.Atoi(string(b))
				case "INT":
					v, _ = strconv.Atoi(string(b))
				default:
					v = string(b)
				}
			} else {
				v = val
			}

			//			if val != nil {
			//				fmt.Println("col=", col, ", value=", val, ", type=", reflect.TypeOf(val).String())
			//			}

			entry[col] = v
		}

		lines++
		// we have a row in entry
		rw.Write(entry)
		//		if callback(entry) == false {
		//			log.Error("User aborted...")
		//			return false
		//		}
		//		fmt.Println(entry)
	}
	//	fmt.Println("return ", lines, " records")
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

// turn map into insert sql
// insert into tbl (k1,k2,k3...) values (v1,v2,v3...)
func BuildInsertSql(m map[string]interface{}, tbl string) string {

	m[updatetime] = time.Now().Format("2006-01-02 15:04:05")

	var kb bytes.Buffer
	var vb bytes.Buffer

	cnt := 0
	for k, v := range m {
		if cnt > 0 {
			kb.WriteString(",")
			vb.WriteString(",")
		} else {
			kb.WriteString("(")
			vb.WriteString("(")
		}
		kb.WriteString(k)
		vb.WriteString(printValue(v))
		cnt++
	}
	kb.WriteString(")")
	vb.WriteString(")")

	var sql bytes.Buffer
	sql.WriteString("insert into ")
	sql.WriteString(tbl)
	sql.Write(kb.Bytes())
	sql.WriteString(" values ")
	sql.Write(vb.Bytes())

	return sql.String()
}

// turn map into update sql
// update tbl set k1=v1, k2=v2, ... where id1=idv1 and id2=idv2
func BuildUpdateSql(fld map[string]interface{}, key map[string]interface{}, tbl string) string {

	fld[updatetime] = time.Now().Format("2006-01-02 15:04:05")

	vb := buildKVPair(fld, ",")
	kb := buildKVPair(key, " and ")

	var sql bytes.Buffer
	sql.WriteString("update ")
	sql.WriteString(tbl)

	sql.WriteString(" set ")
	sql.Write(vb)
	sql.WriteString(" where ")
	sql.Write(kb)

	return sql.String()
}

func printValue(v interface{}) string {
	if v == nil {
		return "null"
	}

	var sv string

	switch v.(type) {
	case string:
		sv = "'" + v.(string) + "'"
	case int, int64, float64, float32, bool:
		sv = fmt.Sprintf("%v", v)
	default:
		sv = fmt.Sprintf("%T", v)
	}

	return sv
}

func buildKVPair(fld map[string]interface{}, sep string) []byte {
	var kb bytes.Buffer

	cnt := 0
	for k, v := range fld {
		if cnt > 0 {
			kb.WriteString(sep)
		}
		kb.WriteString(k)
		kb.WriteString("=")
		kb.WriteString(printValue(v))
		cnt++
	}
	return kb.Bytes()
}

func BuildSelectSql(tbl string, cond map[string]interface{}) string {
	kb := buildKVPair(cond, " and ")

	var sql bytes.Buffer
	sql.WriteString("select * from ")
	sql.WriteString(tbl)
	if len(kb) > 0 {
		sql.WriteString(" where ")
		sql.Write(kb)
	}
	sql.WriteString(" limit ?,?")

	return sql.String()
}

func BuildCountSql(tbl string, cond map[string]interface{}) string {
	kb := buildKVPair(cond, " and ")

	var sql bytes.Buffer
	sql.WriteString("select count(id) as cnt from ")
	sql.WriteString(tbl)

	if len(kb) > 0 {
		sql.WriteString(" where ")
		sql.Write(kb)
	}
	return sql.String()

}

func BuildDeleteSql(tbl string, id int) string {
	fld := make(map[string]interface{})
	fld["deleted"] = 1
	fld[updatetime] = time.Now().Format("2006-01-02 15:04:05")

	key := make(map[string]interface{})
	key["id"] = id

	return BuildUpdateSql(fld, key, tbl)
}

func (this *Database) QuerySingleValueString(sql string, field string, params []interface{}) (string, error) {

	if rslt, err := this.QueryDBReturnAll(sql, params...); err != nil {
		return "", err
	} else {
		if len(rslt) != 1 {
			return "", errors.New("return more than one row")
		} else {
			return toolkit.GetMapValueString(rslt[0], field, "")
		}
	}
}

// turn map into insert sql
// insert into tbl (k1,k2,k3...) values (v1,v2,v3...)
func BuildInsertSql1(m map[string]interface{}, tbl string) string {

	ts := time.Now().Format("2006-01-02 15:04:05")
	m[updatetime] = ts
	//	m[createtime] = ts

	var kb bytes.Buffer
	var vb bytes.Buffer

	cnt := 0
	for k, v := range m {
		if cnt > 0 {
			kb.WriteString(",")
			vb.WriteString(",")
		} else {
			kb.WriteString("(")
			vb.WriteString("(")
		}
		kb.WriteString(k)
		vb.WriteString(printValue1(v))
		cnt++
	}
	kb.WriteString(")")
	vb.WriteString(")")

	var sql bytes.Buffer
	sql.WriteString("insert into ")
	sql.WriteString(tbl)
	sql.Write(kb.Bytes())
	sql.WriteString(" values ")
	sql.Write(vb.Bytes())

	return sql.String()
}

func printValue1(v interface{}) string {
	var sv string

	switch v.(type) {
	case string:
		if v.(string) == "" {
			sv = "null"
		} else {
			sv = "'" + v.(string) + "'"
		}

	case int, int64, float64, float32, bool:
		sv = fmt.Sprintf("%v", v)
	default:
		sv = "unknown type"
	}

	return sv
}

func (this *Database) GetSingleValueString(sql string, field string, params ...interface{}) (string, error) {

	if rslt, err := this.QueryDBReturnAllNew(sql, params...); err != nil {
		return "", err
	} else {
		if len(rslt) != 1 {
			return "", errors.New("return more than one row")
		} else {
			return toolkit.GetMapValueString(rslt[0], field, "")
		}
	}
}

func InitSqlBuilder(ut string, ct string) {
	updatetime = ut
	//	createtime = ct
}

func (this *Database) QueryDB(w io.Writer, sqlString string, args ...interface{}) error {

	rows, err := this.DB.Query(sqlString, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	count := len(columns)

//	_, err := rows.ColumnTypes()
//	if err != nil {
//		return err
//	}

	if count == 0 {
		// no records
		return nil
	}

	for j := 0; j < count; j++ {
		if _, err := fmt.Fprintf(w, "%s\t", columns[j]); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(w, "\n"); err != nil {
		return err
	}

	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
//		entry := make(map[string]interface{})
		for i, _ := range columns {

			var v interface{}
			val := values[i]

			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}

			if _, err := fmt.Fprintf(w, "%v\t", v); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(w, "\n"); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}
