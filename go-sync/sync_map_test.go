package go_sync

import (
	"fmt"
	"sync"
)

/*
https://jacking75.github.io/go_syncmap/
*/
func Example_Sync_Map() {
	m := sync.Map{}

	//
	m.Store("hoge", "fuga")

	if v, ok := m.Load("hoge"); ok {
		fmt.Printf("hoge: %s\n", v) // hoge: fuga
	}

	actual, loaded := m.LoadOrStore("hoge", "piyo")
	fmt.Printf("hoge: %s, loaded=%t\n", actual, loaded) //hoge: fuga, loaded=true

	m.Range(func(k, v any) bool {
		fmt.Printf("key: %s, value:%s\n", k, v) //key: hoge, value: fuga
		return true
	})

	m.Delete("hoge")

	//Output:
	//hoge: fuga
	//hoge: fuga, loaded=true
	//key: hoge, value:fuga
}
