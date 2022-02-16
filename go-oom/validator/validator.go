package validator

import "github.com/kenshin579/tutorials-go/go-oom/repo"

func ValidateAndSave(report string) {
	if validate(report) == true {
		repo.SaveReport(report)
	}
}

func validate(report string) bool {
	return true
}
