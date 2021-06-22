package domain

const (
	CollectionName = "trainers"
)

type Trainer struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
	City string `bson:"city"`
}

func CreateSampleTrainer() []Trainer {
	var trainers []Trainer

	// Some dummy data to add to the Database
	ash := Trainer{Name: "Ash", Age: 10, City: "Pallet Town"}
	misty := Trainer{Name: "Misty", Age: 10, City: "Cerulean City"}
	brock := Trainer{Name: "Brock", Age: 15, City: "Pewter City"}
	trainers = append(trainers, ash)
	trainers = append(trainers, misty)
	trainers = append(trainers, brock)
	return trainers
}
