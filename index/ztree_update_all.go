package index

import (
	"database/sql"
	"github.com/open-falcon/graph/g"
	"log"
)

// 索引的全量更新
func UpdateZtreeAll() {
	dbConn, err := g.GetDbConn("UpdateZtreeAllTask")
	if err != nil {
		log.Println("[ERROR] get dbConn fail", err)
	}
	err = updataztree(dbConn)
	if err != nil {
		log.Println(err)
	}
}

func updataztree(conn *sql.DB) error {
	rows, err := conn.Query("SELECT endpoint from endpoint")
	if err != nil {
		log.Println(err)
		return err
	}

	var endpoint string
	for rows.Next() {
		err := rows.Scan(&endpoint)
		if err != nil {
			log.Println(err)
			return err
		}

		StarZtreeUpdateIncrTask(endpoint, conn)
	}

	return nil
}
