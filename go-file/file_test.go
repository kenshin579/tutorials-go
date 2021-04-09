package go_file

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func Example_실행중인_파일_이름_얻기() {
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println(filepath.Base(filename))

	//Output: file_test.go

}
