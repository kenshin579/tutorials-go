package go_testing

import "log"

var show = func(v ...interface{}) {
	log.Println(v...)
}

func printSize(n int) {
	if n < 10 {
		show("SMALL")
	} else {
		show("LARGE")
	}
}
