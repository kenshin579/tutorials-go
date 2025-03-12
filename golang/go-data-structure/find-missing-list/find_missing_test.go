package find_missing_list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Post struct {
	PostId string
}

func TestFindMissingIdFromListOfEntity(t *testing.T) {
	parcelList := []Post{
		{
			PostId: "1",
		},
		{
			PostId: "2",
		},
		{
			PostId: "3",
		},
	}

	result := findMissingIDList(extractIDAsList(parcelList), []string{"1", "2", "4"})
	assert.Equal(t, []string{"4"}, result)
}

func findMissingIDList(lookupList []string, requestList []string) []string {
	result := make([]string, 0)
	for _, requestID := range requestList {
		if !stringInSlice(requestID, lookupList) {
			result = append(result, requestID)
		}
	}
	return result
}

// todo : stream 스타일의 library를 사용하던지 아니면 가능하면 직접 stream 형식으로 작성하면 좋을 듯하다
func extractIDAsList(postList []Post) []string {
	var list []string
	for _, element := range postList {
		list = append(list, element.PostId)
	}
	return list
}

func stringInSlice(parcelID string, parcelIDList []string) bool {
	for _, parcel := range parcelIDList {
		if parcel == parcelID {
			return true
		}
	}
	return false
}
