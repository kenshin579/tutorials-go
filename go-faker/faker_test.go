package go_faker

import (
	"fmt"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/kenshin579/tutorials-go/go-faker/model"
	"github.com/labstack/gommon/log"
)

func Test_Faker_Without_Tag(t *testing.T) {
	employee := model.Employee{}
	err := faker.FakeData(&employee)
	if err != nil {
		log.Error(err)
	}

	fmt.Println(employee)
}
