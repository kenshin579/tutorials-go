package go_uuid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
)

func TestNewRandom(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(id)
	assert.Equal(t, 16, len(id))
}

//uuid가 생성이 되었는지 체크
//Must 패턴 : https://stackoverflow.com/questions/12434565/built-in-helper-for-must-pattern-in-go
func TestMust_NewRandom(t *testing.T) {
	id := uuid.Must(uuid.NewUUID())
	fmt.Println(id)

}

func TestNewUUID(t *testing.T) {
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(id)
	assert.Equal(t, 16, len(id))
}
