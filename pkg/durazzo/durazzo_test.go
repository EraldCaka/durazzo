package durazzo_test

import (
	"github.com/EraldCaka/durazzo/pkg/durazzo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDurazzo_AutoMigrate(t *testing.T) {
	newDurazzo := durazzo.NewDurazzo(durazzo.Config{
		Driver: durazzo.Postgres,
		DSN:    "postgresql://postgres:1234@localhost:5432/keeper?sslmode=disable",
	})
	type User struct {
		ID    int    `durazzo:"primary_key"`
		Name  string `durazzo:"size:100"`
		Email string `durazzo:"unique"`
	}

	err := newDurazzo.AutoMigrate(&User{})
	assert.Nil(t, err)
}
