package utils

import (
	"fmt"
	"strings"

	"github.com/kenshin579/tutorials-go/go-validation/article/model"
)

func Index(list interface{}, t interface{}) int {
	switch t := t.(type) {
	case string:
		for i, v := range list.([]string) {
			if v == t {
				return i
			}
		}
	case model.PostStatus:
		for i, v := range list.([]model.PostStatus) {
			if v == t {
				return i
			}
		}
	}

	return -1
}

func Include(vs interface{}, t interface{}) bool {
	return Index(vs, t) >= 0
}

func ArrayToString(a []string, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}
