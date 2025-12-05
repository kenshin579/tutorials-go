package go_file

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

func Example_실행중인_파일_이름_얻기() {
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println(filepath.Base(filename))

	//Output: file_test.go
}

func Example_Filename에서_확장명_제외히고_파일이름_가져오기() {
	filename := "test.go"
	extension := filepath.Ext(filename)
	fmt.Println(strings.TrimSuffix(filename, extension))

	//Output: test
}
