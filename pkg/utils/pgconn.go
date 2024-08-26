package utils

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/m3owmurrr/token_server/pkg/config"
)

func NewDBConnection(maxOpenConn, maxIdleConn, connMaxIdleTime int) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbstring())
	if err != nil {
		return nil, err
	}

	// количество одновременных соединений
	db.SetMaxOpenConns(maxOpenConn)
	// размер пула соединений
	db.SetMaxIdleConns(maxIdleConn)
	// время жизни соединения после попадания в пул (обновляется каждый раз)
	db.SetConnMaxIdleTime(time.Second * time.Duration(connMaxIdleTime))

	createTable := `CREATE TABLE IF NOT EXISTS tokens (
		GUID UUID PRIMARY KEY, 
		refresh_token TEXT NOT NULL
	);`

	if _, err := db.Exec(createTable); err != nil {
		return nil, err
	}

	return db, nil
}

func dbstring() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		config.Cfg.DB.User,
		config.Cfg.DB.Pass,
		config.Cfg.DB.Host,
		config.Cfg.DB.Port,
		config.Cfg.DB.Name,
	)
}
