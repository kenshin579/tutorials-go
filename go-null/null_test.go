package go_null

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

type Employee struct {
	Name  string      `json:"name"`
	Phone null.String `json:"phone"`
}

func TestMarshallingWithNullString(t *testing.T) {
	employee := Employee{
		Name: "frank",
	}

	marshal, err := json.Marshal(employee)
	assert.NoError(t, err)
	assert.Equal(t, `{"name":"frank","phone":null}`, string(marshal))
}

func Test_NullStringFrom(t *testing.T) {
	from := null.StringFrom("frank")
	marshal, err := json.Marshal(from)
	assert.NoError(t, err)
	assert.Equal(t, `"frank"`, string(marshal))

}
