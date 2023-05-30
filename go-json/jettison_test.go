package go_json

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wI2L/jettison"
	"gopkg.in/guregu/null.v4"
)

type Student struct {
	Name  string     `json:"name"`
	Age   *float64   `json:"age,omitnil"` // null.Float은 동작을 하지 않음
	Age2  null.Float `json:"age2"`
	Memos []string   `json:"memos"`
}

func Test_StructFieldOmitnil(t *testing.T) {
	age := 10.2

	tests := []struct {
		name    string
		Student Student
		want    string
	}{
		{
			name: "test1",
			Student: Student{
				Name: "frank",
				Age:  &age,
			},
			want: `{"name":"frank","age":10.2,"memos":[]}`,
		},
		{
			name: "test2",
			Student: Student{
				Name: "frank",
			},
			want: `{"name":"frank","memos":[]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data, err := jettison.MarshalOpts(tt.Student,
				jettison.NilSliceEmpty(),
				jettison.NilMapEmpty(),
				jettison.DenyList(extractNullFields(tt.Student)),
			)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(data))
		})
	}
}

func extractNullFields(inter interface{}) []string {
	ignoreFields := make([]string, 0)
	v := reflect.ValueOf(inter)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		if fieldValue.Type() == reflect.TypeOf(null.Float{}) && !fieldValue.Interface().(null.Float).Valid {
			ignoreFields = append(ignoreFields, field.Tag.Get("json"))
		}

	}

	return ignoreFields
}
