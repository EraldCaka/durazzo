package durazzo_test

import (
	"github.com/EraldCaka/durazzo/pkg/durazzo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDurazzo_Delete(t *testing.T) {
	newDurazzo := durazzo.NewDurazzo(durazzo.Config{
		Driver: durazzo.Postgres,
		DSN:    "postgresql://postgres:postgres@localhost:5432/testdb?sslmode=disable",
	})

	type User struct {
		ID    int    `durazzo:"primary_key"`
		Name  string `durazzo:"size:100"`
		Email string `durazzo:"unique"`
	}

	err := newDurazzo.AutoMigrate(&User{})
	assert.Nil(t, err)

	insertQuery := `INSERT INTO "user" (name, email) VALUES ($1, $2)`
	_, err = newDurazzo.Db.Exec(insertQuery, "edgar", "edgar@gmail.com")
	assert.Nil(t, err)

	deleteQuery := `DELETE FROM "user" WHERE name = $1`
	_, err = newDurazzo.Db.Exec(deleteQuery, "edgar")
	assert.Nil(t, err)

	var users []User
	err = newDurazzo.Select(&users).Where("name", "edgar").Run()
	assert.Nil(t, err)

	assert.Equal(t, 0, len(users))

	dropTableQuery := `DROP TABLE IF EXISTS "user"`
	_, err = newDurazzo.Db.Exec(dropTableQuery)
	assert.Nil(t, err)
}
