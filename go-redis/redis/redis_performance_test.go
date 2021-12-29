package redis

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	TestKey = "t-100"
)

func setupSingleThreadTest() {
	client = newRedisClient()

}

func teardownSingleThreadTest() {
	client.ZRemRangeByScore(context.Background(), TestKey, "-inf", "+inf")
	client.Close()
}

func Test_Read(t *testing.T) {
	setupSingleThreadTest()
	createSamples(50000)
	defer teardownSingleThreadTest()

	start := time.Now()
	client.ZRange(context.Background(), TestKey, 0, -1)
	fmt.Printf("1.%d ms\n", time.Since(start).Milliseconds())

}

/*
redis는 single thread로 동작한다. 아래 케이스에서 어떻게 동작을 하는지 검증
1.한 쓰레드에서 오랫동안 reading하는 작업을 함
2.다른 쓰레드에서 reading/writing을 하는 경우 blocking이 되는지 확인
*/
func Test_Redis_SingleThread(t *testing.T) {
	setupSingleThreadTest()
	createSamples(1000)
	defer teardownSingleThreadTest()

	messages := make(chan int)
	var wg sync.WaitGroup

	wg.Add(2)
	go read1(&wg, messages)
	go read2(&wg, messages)

	go func() {
		for i := range messages {
			fmt.Println("messages", i)
		}
	}()
	wg.Wait()

}

func read1(wg *sync.WaitGroup, messages chan int) {
	defer wg.Done()
	fmt.Println("1.START")
	start := time.Now()
	client.ZRange(context.Background(), TestKey, 0, -1)
	fmt.Printf("1.%d ms\n", time.Since(start).Milliseconds())
	messages <- 1
	fmt.Println("1.END")
}

func read2(wg *sync.WaitGroup, messages chan int) {
	defer wg.Done()
	fmt.Println("2.START")
	time.Sleep(time.Millisecond * 3)
	start := time.Now()
	client.ZRange(context.Background(), TestKey, 0, 10)
	fmt.Printf("2.%d ms\n", time.Since(start).Milliseconds())
	messages <- 2
	fmt.Println("2.END")
}

func createSampleMembers(max int, data []byte) []*redis.Z {
	mList := make([]*redis.Z, 0)

	for i := 0; i < max; i++ {
		replace := strings.Replace(string(data), "4865bc96-454b-4a78-9efa-c8365e1b62a8", "4865bc96-454b-4a78-9efa-c8365e1b62a8"+strconv.Itoa(i+1), 1)
		mList = append(mList, &redis.Z{
			Score:  float64(i + 1),
			Member: replace,
			//Member: "msg" + strconv.Itoa(i+1),
		})
	}
	fmt.Println("size:", len(mList))

	return mList
}

func readJsonFile() []byte {
	b, err := ioutil.ReadFile(filepath.Join("./sample.json"))
	if err != nil {
		fmt.Errorf("err:%v", err)
	}
	return b
}

func createSamples(maxMember int) {
	client.ZAdd(context.Background(), TestKey, createSampleMembers(maxMember, readJsonFile())...)
}
