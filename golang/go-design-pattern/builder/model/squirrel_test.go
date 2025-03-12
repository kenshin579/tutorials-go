package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	sq "github.com/Masterminds/squirrel"
)

func Test_Builder_Squirrel(t *testing.T) {

	users := sq.Select("*").
		From("users").
		Join("emails USING (email_id)")

	active := users.Where(sq.Eq{"deleted_at": nil})

	sql, _, _ := active.ToSql()

	expectedSql := "SELECT * FROM users JOIN emails USING (email_id) WHERE deleted_at IS NULL"
	assert.Equal(t, expectedSql, sql)
}
