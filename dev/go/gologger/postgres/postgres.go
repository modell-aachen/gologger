package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/lib/pq"
	"github.com/modell-aachen/gologger/interfaces"
	"github.com/pkg/errors"
)

const (
	host     = "localhost"
	port     = 5432
	dbname   = "foswiki_logs_9"
	logtable = "logs"
)

type postgresStore struct {
	db *sql.DB
}

func (store postgresStore) Close() {
	store.db.Close()
}

func (store postgresStore) CleanUp() (err error) {
	_, err = store.db.Exec("DELETE FROM logs WHERE time < now() - interval '3 month'")
	return err
}

func (store postgresStore) Read(startTime time.Time, endTime time.Time, source interfaces.SourceString, levels []interfaces.LevelString) (entries []interfaces.LogRow, err error) {
	if endTime.IsZero() {
		endTime = time.Now()
	}
	var rows *sql.Rows
	if source != interfaces.SourceString("") {
		rows, err = store.db.Query("SELECT time, source, level, misc, fields from logs WHERE time > $1 AND time < $2 AND source = $3 AND level = ANY($4) ORDER BY time ASC, id ASC", startTime, endTime, source, pq.Array(levels))
	} else {
		rows, err = store.db.Query("SELECT time, source, level, misc, fields from logs WHERE time > $1 AND time < $2 AND level = ANY($3) ORDER BY time ASC, id ASC", startTime, endTime, pq.Array(levels))
	}

	if err != nil {
		return entries, errors.Wrap(err, "Could not read logs from database")
	}
	defer rows.Close()

	var rowsX []interfaces.LogRow
	rowsX = make([]interfaces.LogRow, 0)
	for rows.Next() {
		var (
			row     interfaces.LogRow
			rowTime time.Time
			source  interfaces.SourceString
			level   interfaces.LevelString

			miscString   []byte
			fieldsString []byte
			misc         interfaces.LogGeneric
			fields       interfaces.LogFields
		)
		err = rows.Scan(&rowTime, &source, &level, &miscString, &fieldsString)
		if err != nil {
			fmt.Printf("Could not scan row: %s\n%s\n", err, rows)
			continue
		}
		err = json.Unmarshal(fieldsString, &fields)
		if err != nil {
			fmt.Printf("Could not unmarshal fields: %s\n%s\n", err, fieldsString)
			continue
		}
		err = json.Unmarshal(miscString, &misc)
		if err != nil {
			fmt.Printf("Could not unmarshal extra fields: %s\n%s\n", err, miscString)
			continue
		}

		row = interfaces.LogRow{
			rowTime,
			source,
			level,
			misc,
			fields,
		}
		rowsX = append(rowsX, row)
	}

	return rowsX, nil
}

func (store postgresStore) Store(logRow interfaces.LogRow) error {
	jsonMisc, err := json.Marshal(logRow.Misc)
	if err != nil {
		return errors.Wrap(err, "Could not marshal misc data")
	}

	jsonFields, err := json.Marshal(logRow.Fields)
	if err != nil {
		return errors.Wrap(err, "Could not marshal fields")
	}

	_, err = store.db.Exec("INSERT INTO logs VALUES(DEFAULT, $1, $2, $3, $4, $5)", logRow.Time, logRow.Source, logRow.Level, jsonMisc, jsonFields)
	if err != nil {
		return errors.Wrap(err, "Insertion failed")
	}
	return nil
}

/*
sudo -u postgres psql -c "CREATE DATABASE 'foswiki_logs' WITH ENCODING 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8' TEMPLATE template0"
*/

func CreateInstance() (interfaces.LogStore, error) {
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		return nil, errors.New("POSTGRES_USER not set")
	}
	password := os.Getenv("POSTGRES_PASSWORD")

	err := setupDatabase(user, password)
	if err != nil {
		return nil, err
	}

	psqlInfo := fmt.Sprintf("dbname=%s sslmode=disable user='%s' password='%s'", dbname, user, password)
	db, err := sql.Open("postgres", psqlInfo)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = setupTables(db)
	if err != nil {
		return nil, err
	}

	store := postgresStore{
		db,
	}

	return store, nil
}

func setupTables(db *sql.DB) error {
	count := 0
	err := db.QueryRow("SELECT count(*) FROM pg_type WHERE typname='debuglevels'").Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = db.Exec("CREATE TYPE debuglevels AS ENUM('notice', 'event', 'debug', 'info', 'warning', 'error', 'fatal')")
		if err != nil {
			return errors.Wrap(err, "Error creating debuglevels enum")
		}
	}

	err = db.QueryRow("SELECT count(*) FROM pg_tables WHERE tablename=$1", logtable).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		sequence := logtable + "_id_seq"
		_, err = db.Exec(
			`CREATE SEQUENCE ` + sequence + `;
			CREATE TABLE ` + logtable + ` (
				id integer NOT NULL DEFAULT nextval('` + sequence + `'),
				time timestamptz,
				source text,
				level debuglevels,
				misc jsonb,
				fields jsonb
			);
			ALTER SEQUENCE ` + sequence + ` OWNED BY ` + logtable + `.id`,
		)
		if err != nil {
			return errors.Wrap(err, "Error creating tables")
		}
	}
	return nil
}

func setupDatabase(user string, password string) error {
	psqlInfo := fmt.Sprintf("dbname=%s sslmode=disable user='%s' password='%s'", "postgres", user, password)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	count := 0
	err = db.QueryRow("SELECT count(*) FROM pg_database WHERE datname=$1", dbname).Scan(&count)
	if err != nil {
		return errors.Wrap(err, "Error checking database")
	}
	if count == 0 {
		_, err = db.Exec("CREATE DATABASE " + dbname + " WITH ENCODING 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8' TEMPLATE template0")
		if err != nil {
			return errors.Wrap(err, "Error creating database")
		}
	}
	err = db.Close()
	if err != nil {
		return errors.Wrap(err, "Could not close connection after checking database")
	}

	return nil

}
