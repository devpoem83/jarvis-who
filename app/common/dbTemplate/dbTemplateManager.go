package dbTemplate

import (

	"context"
	"database/sql"
	"log"
	"reflect"
	"strconv"
	"strings"

	"gitlab.eduwill.net/dev_team/jarvis-who/app/base/config"
	"gitlab.eduwill.net/dev_team/jarvis-who/app/common"
)

func makeSqlAndParam(sqlId string, paramMap map[string]interface{}) (string, []interface{}) {

	var parameter []interface{}
	sql := config.GetCache(sqlId)
	//originSql := sql
	fullSql := sql

	prefixParam := "#{"
	suffixParam := "}"

	count := strings.Count(sql, prefixParam)
	if count > 0 {
		splitArray := strings.Split(sql, prefixParam)
		if splitArray != nil {
			var attrParam string
			for _, str := range splitArray {
				if strings.Contains(str, suffixParam) {
					key := strings.Split(str, suffixParam)[0]
					value := paramMap[key]
					if value != nil {
						parameter = append(parameter, value)
						attrParam = prefixParam + key + suffixParam
						sql = strings.Replace(sql, attrParam, "?", 1)

						if reflect.TypeOf(value).Name() == "string" {
							fullSql = strings.Replace(fullSql, attrParam, "'"+value.(string)+"'", 1)
						} else if reflect.TypeOf(value).Name() == "int64" {
							fullSql = strings.Replace(fullSql, attrParam, strconv.FormatInt(value.(int64), 10), 1)
						} else {
							fullSql = strings.Replace(fullSql, attrParam, strconv.Itoa(value.(int)), 1)
						}
					}
				}
			}
		}
	}
	common.Logger.Debug("DB SQL   : ", fullSql)

	return sql, parameter
}

func SelectList(db *sql.DB, sqlId string, paramMap map[string]interface{}) ([]map[string]interface{}, error) {

	var list []map[string]interface{}
	strSql, parameter := makeSqlAndParam(sqlId, paramMap)

	if db != nil && strSql != "" {
		rows, err := db.Query(strSql, parameter...)
		cols, err := rows.Columns()

		if err != nil {
			log.Panic()
			common.Logger.Error("panic 발생 후!")
		}

		for rows.Next() {
			columns := make([]interface{}, len(cols))
			columnPointers := make([]interface{}, len(cols))

			for i, _ := range columns {
				columnPointers[i] = &columns[i]
			}

			if err := rows.Scan(columnPointers...); err != nil {
				return nil, err
			}

			obj := make(map[string]interface{})
			for i, colName := range cols {
				val := columnPointers[i].(*interface{})
				obj[colName] = *val
			}

			list = append(list, obj)
		}
	} else {
		common.Logger.Warn("DB is nil or query is empty!")
		common.Logger.Warn("sqlId : " + sqlId)
	}

	common.Logger.Debug("DB Data  : ", list)
	return list, nil
}

// Select : return 1 row
func SelectOne(db *sql.DB, sqlId string, paramMap map[string]interface{}) (map[string]interface{}, error) {

	var object map[string]interface{}
	strSql, parameter := makeSqlAndParam(sqlId, paramMap)

	if db != nil && strSql != "" {
		rows, _ := db.Query(strSql, parameter...)
		cols, _ := rows.Columns()

		if rows.Next() {

			columns := make([]interface{}, len(cols))
			columnPointers := make([]interface{}, len(cols))

			for i, _ := range columns {
				columnPointers[i] = &columns[i]
			}

			if err := rows.Scan(columnPointers...); err != nil {
				return nil, err
			}

			obj := make(map[string]interface{})
			for i, colName := range cols {
				val := columnPointers[i].(*interface{})
				obj[colName] = *val
			}

			object = obj
		}
	} else {
		common.Logger.Warn("DB is nil or query is empty!")
	}
	common.Logger.Debug("DB Result : ", object)
	return object, nil
}

func SelectValue(db *sql.DB, sqlId string, paramMap map[string]interface{}) (int, error) {

	var value int
	strSql, parameter := makeSqlAndParam(sqlId, paramMap)

	if db != nil && strSql != "" {
		row := db.QueryRow(strSql, parameter...)

		err := row.Scan(&value)
		if err != nil {
			panic(err)
			return 0, err
		} else {
			common.Logger.Info("value  : ", value)
		}
		
	}
	return value, nil
}

// Insert or Update or Delete : doesnt't return rows
func Exec(db *sql.DB, sqlId string, paramMap map[string]interface{}) (int, int) {

	var lastInsertId int
	var rowsAffected int

	strSql, parameter := makeSqlAndParam(sqlId, paramMap)

	if db != nil && strSql != "" {
		result, err := db.Exec(strSql, parameter...)

		common.Logger.Debug("result : ", result)
		common.Logger.Debug("err : ", err)

		if err != nil {
			common.Logger.Error("error : ", err)
		} else {
			if result != nil {
				lastInsertId64, _ := result.LastInsertId()
				rowsAffected64, _ := result.RowsAffected()

				lastInsertId = int(lastInsertId64)
				rowsAffected = int(rowsAffected64)

			} else {
				common.Logger.Warn("DB is nil or query is empty!")
			}
		}
	}
	common.Logger.Debug("DB Result : [", lastInsertId, rowsAffected, "]")
	return lastInsertId, rowsAffected
}

func ExecTx(tx *sql.Tx, sqlId string, paramMap map[string]interface{}) (int, int, error) {

	var lastInsertId int
	var rowsAffected int

	strSql, parameter := makeSqlAndParam(sqlId, paramMap)

	if strSql != "" {
		result, err := tx.Exec(strSql, parameter...)
		if err != nil {
			common.Logger.Error("err : ", err)
			return 0, 0, err
		} else {
			lastInsertId64, _ := result.LastInsertId()
			rowsAffected64, _ := result.RowsAffected()

			lastInsertId = int(lastInsertId64)
			rowsAffected = int(rowsAffected64)
		}
	}

	return lastInsertId, rowsAffected, nil
}

func BeginTx(db *sql.DB) (*sql.Tx, error) {
	common.Logger.Debug("begin transaction!")
	ctx := context.Background()
	if tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable}); err != nil {
		return nil, err
	} else {
		return tx, nil
	}
}

func RollbackTx(tx *sql.Tx) {
	common.Logger.Debug("rollback transaction!")
	if err := tx.Rollback(); err != nil {
		common.Logger.Error("rollback transaction error : ", err)
	} else {
		common.Logger.Debug("rollback transaction success!")
	}
}

func CommitTx(tx *sql.Tx) {
	common.Logger.Debug("commit transaction!")
	if err := tx.Commit(); err != nil {
		common.Logger.Error("commit transaction error : ", err)
	} else {
		common.Logger.Debug("commit transaction success!")
	}
}
