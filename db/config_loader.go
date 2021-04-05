package db
//
//import (
////	"database/sql"
////	 _ "github.com/go-sql-driver/mysql"	
//	"time"
//	"fmt"
////	"bytes"
//////	"reflect"
////	"strconv"
//	toolkit "github.com/charles/toolkit"
////	"errors"
//)
//
//type ConfigLoaderInterface interface {
//	// it is called when a record is deleted
//	// return false, if no more config load
//	OnDelete(r map[string]interface{}) bool
//	OnInsertOrUpdate(r map[string]interface{}) bool
//}
//
//type ConfigLoader struct {
//	DB		*Database	   // database handle
//	Ts		string
//	Sql		string
//	Cli		ConfigLoaderInterface
//
//	// stats
//	CntInvalidFlag	int
//	CntInsert	int
//	CntInsertOk	int
//	CntDelete	int
//	CntDeleteOk	int
//}
//
//
//
//func (this*ConfigLoader) Write(m map[string]interface{}) {
//	if deleted, err := toolkit.GetMapValueAsInt(m, "deleted"); err != nil {
//		this.CntInvalidFlag ++
//		return
//	} else if deleted == 0 {
//		if this.Cli.OnInsertOrUpdate(m) == true {
//			this.CntInsertOk ++
//		}
//		this.CntInsert ++
//		
//	} else {
//		if this.Cli.OnDelete(m) == true {
//			this.CntDeleteOk ++
//		}
//		this.CntDelete ++
//	}	
//}
//
//func (this *ConfigLoader) Load(cate string) (string, error) {
//	
////	log.Info("load, sql=" + this.Sql)
//	// load data from db
//	if err := this.DB.QueryDBCallBack(this, this.Sql, this.Ts, cate); err != nil {
////		s := fmt.Sprintf("Can not read DB, err=%v", err)
////		log.Error(s)		
//		return "", err
//	} 
//	
//	// update ts
//	this.Ts = time.Now().Format("2006-01-02 15:04:05")
//
//	// generate report
//	report := fmt.Sprintf("Ts=%s, InvalidFlag=%v,Insert=%v, InsertOk=%v, Delete=%v, DeleteOk=%v",
//			this.Ts, 
//			this.CntInvalidFlag,
//	this.CntInsert,
//	this.CntInsertOk,
//	this.CntDelete,
//	this.CntDeleteOk)
//
//	return report, nil
//}
//
