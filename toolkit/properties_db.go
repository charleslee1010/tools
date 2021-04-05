package toolkit

import (
//	"fmt"
//	"strconv"
//	"strings"
	"database/sql"
	 _ "github.com/go-sql-driver/mysql"		
)



func (this *Properties) ReloadFromDB(module string, db *sql.DB) (error) {
	this.L.Lock()
	defer this.L.Unlock()
	
	this.Pmap = make(map[string][]string)
			
	rows, err := db.Query("select k,v from globals where deleted=0 and module=?", module)
	if err != nil {
		return  err
	}
	defer rows.Close()
	
	for rows.Next() {
		k := ""
		v := ""
		err = rows.Scan(&k, &v)
		if err != nil {
			return err
		}

		// save k, v
		this.Add(k, v)
//		
//		k = strings.TrimSpace(k)
//		v = strings.TrimSpace(v)
//		
//		if k == "" {
//			continue
//		}
//		
//		if s, ok := prop.Pmap[k]; ok {
//			// append the string to
//			prop.Pmap[k] = append(s, v)
//		} else {
//			 prop.Pmap[k] = []string {v}
//		}
	}
	
	return nil
}

func (this *Properties) ReloadFromDBFromTable(module string, db *sql.DB, table string) (error) {
	this.L.Lock()
	defer this.L.Unlock()
	
	this.Pmap = make(map[string][]string)
			
	rows, err := db.Query("select k,v from " + table + " where deleted=0 and module=?", module)
	if err != nil {
		return  err
	}
	defer rows.Close()
	
	for rows.Next() {
		k := ""
		v := ""
		err = rows.Scan(&k, &v)
		if err != nil {
			return err
		}

		// save k, v
		this.Add(k, v)
	}
	
	return nil
}
