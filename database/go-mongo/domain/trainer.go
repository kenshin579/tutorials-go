package domain

import "fmt"

const (
	CollectionName = "trainers"
)

type Trainer struct {
	ID   string `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
	Age  int    `bson:"age" json:"age"`
	City string `bson:"city" json:"city"`
}

func CreateTrainerSample(id int) Trainer {
	return Trainer{
		ID:   fmt.Sprintf("trainer-%d", id),
		Name: "Ash",
		Age:  10,
		City: "Pallet Town",
	}
}

func CreateTrainersSample(max int) []Trainer {
	result := make([]Trainer, 0)
	for i := 0; i < max; i++ {
		result = append(result, Trainer{
			ID:   fmt.Sprintf("trainer-%d", i),
			Name: fmt.Sprintf("Frank%d", i),
			Age:  5,
			City: "city",
		})
	}
	return result
}
