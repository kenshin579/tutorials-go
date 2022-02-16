package repo

import "time"

var savedReports = []string{}

func SaveReport(report string) {
	savedReports = append(savedReports, report)

	time.Sleep(10 * time.Microsecond)
	if len(savedReports) > 1000000 {
		panic("OOM")
	}
}
