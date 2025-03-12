package router

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
)

func TestJsonOmitEmpty_잘되는지_확인(t *testing.T) {
	e := Error{
		//Code:    0,
		Message: "",
		//Errors:  nil,
	}

	marshal, _ := json.Marshal(e)
	fmt.Println(string(marshal))
}

func TestRegexCapturingGroup(t *testing.T) {
	str := "Key: 'ArticleRequest.Description' Error:Field validation for 'Description' failed on the 'required' tag\n"

	r := regexp.MustCompile("Error:(.*)")
	submatch := r.FindStringSubmatch(str)
	fmt.Println(submatch)

}
