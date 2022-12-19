package go_maps

import (
	"fmt"

	"golang.org/x/exp/maps"
)

func Example_Maps_Copy() {
	src := map[int]string{
		200: "foo",
		300: "bar",
	}
	dest := map[int]string{}

	maps.Copy(dest, src)
	fmt.Println(dest)

	dest2 := map[int]string{
		200: "will be overwritten",
	}
	maps.Copy(dest2, src)
	fmt.Println(dest2)

	//Output:
	//map[200:foo 300:bar]
	//map[200:foo 300:bar] //todo - 왜 overwrite가 안되나?
}
