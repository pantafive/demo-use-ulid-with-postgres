package main

import (
	"encoding/hex"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"

	"app/testhelpers"
)

func TestULIDWithPostgres(t *testing.T) {
	id, _ := ulid.New(1659647944779, io.Reader(nil))
	fmt.Println(id.String())
	// 01G9NBKT2B0000000000000000
	// INSERT INTO test (id) VALUES ('01G9NBKT2B0000000000000000'::uuid);
	// ERROR: invalid input syntax for type uuid: "01G9NBKT2B0000000000000000"

	fmt.Println(hex.EncodeToString(id.Bytes()))
	// 01826ab9e84b00000000000000000000
	// INSERT INTO test (id) VALUES ('01826ab9e84b00000000000000000000'::uuid);
	// 1 row affected in 2 ms

	// setup
	testDB := testhelpers.NewTestDatabase(t)
	defer testDB.Close(t)
	db := sqlx.MustConnect("pgx", testDB.ConnectionString(t))

	t.Run("test ulid", func(t *testing.T) {

		db.MustExec("CREATE TABLE test ( id uuid PRIMARY KEY )")

		db.MustExec("INSERT INTO test (id) VALUES ($1::uuid);", id)

		row := db.QueryRowx("SELECT id from test where id = $1;", id)

		var idFromDB []byte
		_ = row.Scan(&idFromDB)
		assert.Equal(t, "01826ab9-e84b-0000-0000-000000000000", string(idFromDB))
	})
}
