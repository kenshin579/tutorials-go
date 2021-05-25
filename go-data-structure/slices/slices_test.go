package slices

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kenshin579/tutorials-go/go-data-structure/slices/model"
)

func Test_Delete_Item_Index_From_Slice(t *testing.T) {
	people := createSamplePerson(5)
	deletedPersonID := 2

	temp := people[:0]

	for _, person := range people {
		if person.ID != deletedPersonID {
			temp = append(temp, person)
		}
	}

	people = temp
	assert.Equal(t, 4, len(people))
	assert.False(t, isPersonIDFound(people, deletedPersonID))

}

func Test_Delete_Item_Index_From_Slice_Without_Temp(t *testing.T) {
	people := createSamplePerson(5)
	deletedPersonID := 2

	for i, person := range people {
		if person.ID == deletedPersonID {
			people = append(people[:i], people[i+1:]...) //todo: 이거에 대한 설명이 필요함
		}
	}

	assert.Equal(t, 4, len(people))
	assert.False(t, isPersonIDFound(people, deletedPersonID))
}

func createSamplePerson(max int) []model.Person {
	people := make([]model.Person, 0)

	for i := 0; i < max; i++ {
		people = append(people, model.Person{
			ID:   i + 1,
			Name: fmt.Sprintf("name-%d", i),
			Age:  i,
		})
	}

	return people
}

func isPersonIDFound(people []model.Person, deletedPersonID int) bool {
	for _, person := range people {
		if person.ID == deletedPersonID {
			return true
		}
	}
	return false
}

//todo: 아래부분 스터디해보기
//https://yourbasic.org/golang/delete-element-slice/
func Example_Delete_Item_From_Slice() {
	strList := []string{"A", "B", "C", "D", "E"}
	deletedIndex := 2

	// Remove the element at index deletedIndex from strList.
	strList[deletedIndex] = strList[len(strList)-1] // Copy last element to index deletedIndex.
	strList[len(strList)-1] = ""                    // Erase last element (write zero value).
	strList = strList[:len(strList)-1]              // Truncate slice.

	fmt.Println(strList)

	//Output:
	//[A B E D]
}
