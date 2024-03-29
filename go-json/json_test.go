package go_json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/kenshin579/tutorials-go/go-json/model"

	"github.com/stretchr/testify/assert"
)

func ExampleJsonMarshal_Struct_To_Json_구조체_다_값이_있는_경우() {
	response := model.StudentResponse{
		Data:        3,
		Message:     "this is a message",
		MessageList: []string{"msg`", "msg1"},
		Errors: map[string]string{
			"error": "error1",
		},
		Student: model.Student{
			Name: "Frank",
		},
		StudentList: []model.Student{{Name: "Frank1"}},
	}

	jsonResponse, _ := json.MarshalIndent(response, "", "\t")
	fmt.Println("jsonResponse", string(jsonResponse))

	// Output:
	// jsonResponse {
	//	"data": 3,
	//	"message": "this is a message",
	//	"messageList": [
	//		"msg`",
	//		"msg1"
	//	],
	//	"errors": {
	//		"error": "error1"
	//	},
	//	"student": {
	//		"Name": "Frank"
	//	},
	//	"studentList": [
	//		{
	//			"Name": "Frank1"
	//		}
	//	]
	// }
}

func ExampleJsonMarshal_Struct_To_Json_구조체_값이_다_없는_경우() {
	response := model.StudentResponse{}
	jsonResponse, _ := json.MarshalIndent(response, "", "\t")
	fmt.Println("jsonResponse", string(jsonResponse))

	// Output:
	// jsonResponse {
	//	"student": {
	//		"Name": ""
	//	}
	// }
}

type Student2 struct {
	Name string
}

type Response2 struct {
	Data        interface{}       `json:"data,omitempty"`
	Message     string            `json:"message,omitempty"`
	MessageList []string          `json:"messageList,omitempty"`
	Errors      map[string]string `json:"errors,omitempty"`
	Student     *Student2         `json:"student,omitempty"`
	StudentList []Student2        `json:"studentList,omitempty"`
}

func ExampleJsonMarshal_Struct_To_Json_구조체_값이_다_없는_경우_구조쳊가_Nil인_경우에_Json에_포함되지_않는다() {
	response := Response2{}
	jsonResponse, _ := json.MarshalIndent(response, "", "\t")
	fmt.Println("jsonResponse", string(jsonResponse))

	// Output:
	// jsonResponse {}
}

func Test_Unmarshal_RawString(t *testing.T) {
	jsonStr := `{"data": 1, "studentList": [{"Name": "Frank"}]}`
	response2 := Response2{}
	err := json.Unmarshal([]byte(jsonStr), &response2)
	assert.NoError(t, err)

	assert.Equal(t, len(response2.StudentList), 1)
}

// https://stackoverflow.com/questions/58434023/how-to-parse-http-response-body-to-json-format-in-golang
func TestGet의_Response_값이_Json_Array인_경우_sliceOfMap으로_decode하는_방법(t *testing.T) {
	url := "https://skishore.github.com/inkstone/all.json"

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	keys := make([]map[string]interface{}, 0)
	err = json.NewDecoder(resp.Body).Decode(&keys)
	if err != nil {
		panic(err)
	}
	fmt.Println("first item", keys[0])
}

// https://coderwall.com/p/4c2zig/decode-top-level-json-array-into-a-slice-of-structs-in-golang
func TestJson_Array를_sliceOfMap로_decode하는_방법(t *testing.T) {
	keysBody := []byte(`[{"id": 1,"key": "-"},{"id": 2,"key": "-"},{"id": 3,"key": "-"}]`)
	keys := make([]map[string]interface{}, 0)
	err := json.Unmarshal(keysBody, &keys)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), keys[0]["id"])
}

type idKey struct {
	ID  float64
	Key string
}

func TestJson_Unmarshal_ListOfStruct(t *testing.T) {
	keysBody := []byte(`[{"id": 1,"key": "-"},{"id": 2,"key": "-"},{"id": 3,"key": "-"}]`)
	keys := make([]idKey, 0)
	err := json.Unmarshal(keysBody, &keys)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), keys[0].ID)
}

// Dto에서 필요한 속성만 정의해서 사용하면 된다
func TestJsonArray를_sliceOfStruct로_decode하는_방법(t *testing.T) {
	courseList := getScenarioNodeFromTestJsonFile("data/array.json")

	assert.Equal(t, "Lessons 1-4", courseList[0].Name)
}

func getScenarioNodeFromTestJsonFile(path string) []model.CourseResponse {
	var courseResponseList []model.CourseResponse

	jsonFile, err := ioutil.ReadFile(filepath.Join(path))
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonFile, &courseResponseList)
	if err != nil {
		panic(err)
	}

	return courseResponseList
}
