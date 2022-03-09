package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file.")
)

func main() {
	log.Println("begin")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)     // 开启CPU性能分析
		defer pprof.StopCPUProfile() // 停止CPU性能分析
	}

	for i := 0; i < 30; i++ {
		nums := fibonacci(i)
		fmt.Println(nums)
	}
}

func fibonacci(num int) int {
	if num < 2 {
		return 1
	}
	return fibonacci(num-1) + fibonacci(num-2)
}
