package repositories

import (
	"database/sql"
	"log"
)

type queryExecuter interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// queryExecuterとr.db.Exec() の引数を受け取り、実行し、RowsAffected() で更新された行数を返す
func execQueryAndReturnAffectedRows(db queryExecuter, query string, args ...interface{}) (int64, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	if rows == 0 {
		return 0, nil
	}

	return rows, nil
}
