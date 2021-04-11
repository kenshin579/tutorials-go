package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	sq "github.com/Masterminds/squirrel"
)

func Test_Builder_Employee(t *testing.T) {
	emp1 := EmployeeBuilder().
		Name("Michael Scott1").
		Role("manager").
		Build()
	fmt.Println(emp1)
	assert.Equal(t, "Michael Scott1", emp1.Name)
	assert.Equal(t, "manager", emp1.Role)

	emp2 := EmployeeBuilder().
		Name("Michael Scott2").
		Role("manager").
		Build()
	fmt.Println(emp2)
	assert.Equal(t, "Michael Scott2", emp2.Name)
	assert.Equal(t, "manager", emp2.Role)
}

func Test_Builder_Squirrel(t *testing.T) {

	users := sq.Select("*").
		From("users").
		Join("emails USING (email_id)")

	active := users.Where(sq.Eq{"deleted_at": nil})

	sql, _, _ := active.ToSql()

	expectedSql := "SELECT * FROM users JOIN emails USING (email_id) WHERE deleted_at IS NULL"
	assert.Equal(t, expectedSql, sql)
}
