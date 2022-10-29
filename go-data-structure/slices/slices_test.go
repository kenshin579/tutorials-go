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

func TestSlice(t *testing.T) {
	people := createSamplePerson(5)
	fmt.Println(people)

	fmt.Println(people[:2])
	fmt.Println(people[2:])
	fmt.Println(people[1:3])
}

func Test_Delete_Item_Index_From_Slice_Without_Temp(t *testing.T) {
	people := createSamplePerson(5)
	deletedPersonID := 2

	for i, person := range people {
		if person.ID == deletedPersonID { //i:2
			people = append(people[:i], people[i+1:]...) //([:2], [3:]...) -> (0, 1, 3, 4, 5)
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

// https://yourbasic.org/golang/delete-element-slice/
func Example_Delete_Item_From_Slice_Fast_Version_Changes_Order() {
	a := []string{"A", "B", "C", "D", "E"}
	i := 2

	a[i] = a[len(a)-1] // 마지막 요소(E) -> i로 복사
	a[len(a)-1] = ""   // 마지막 요소 삭제
	a = a[:len(a)-1]   // slice 크기 줄임

	fmt.Println(a)

	//Output:
	//[A B E D]
}

func Example_Delete_Item_From_Slice_Slow_Version_Maintains_Order() {
	a := []string{"A", "B", "C", "D", "E"}
	i := 2

	copy(a[i:], a[i+1:]) // a[2:] <- a[3:] 복사 (A, B, D, E)
	fmt.Println(len(a))
	a[len(a)-1] = ""
	a = a[:len(a)-1] // slice 크기 줄임
	fmt.Println(len(a))

	fmt.Println(a)

	//Output:
	//5
	//4
	//[A B D E]
}

// [inclusive:exclusive]
func Example_Slice_Index() {
	a := []string{"A", "B", "C", "D", "E"}
	fmt.Println(a[:])  //전체 [A B C D E]
	fmt.Println(a[2:]) //[C D E]
	fmt.Println(a[:2]) //[A B]
	fmt.Println()

	//Output:
	//[A B C D E]
	//[C D E]
	//[A B]
}
