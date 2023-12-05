package driver

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

var (
	PostgresHost   = os.Getenv("POSTGRES_HOST")
	PostgresPort   = os.Getenv("POSTGRES_PORT")
	PostgresUserNm = os.Getenv("POSTGRES_USER_NM")
	PostgresUserPw = os.Getenv("POSTGRES_USER_PW")
	PostgresDbName = os.Getenv("POSTGRES_DB_NAME")
)

// OpenDB establishes connection with psql databse
func OpenDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=localhost port=5432 user=pgadmin password=Matwyenko1_ dbname=finances sslmode=disable") //PostgresHost,
	//PostgresPort,
	//PostgresUserNm,
	//PostgresUserPw,
	//PostgresDbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
