package robfig

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kenshin579/tutorials-go/go-schedule/job"
	"github.com/robfig/cron/v3"
)

//https://en.wikipedia.org/wiki/Cron
//https://programmer.ink/think/cron-of-go-daily.html
func Test_Cron_Different_Format(t *testing.T) {
	c := cron.New()
	defer c.Stop()

	c.AddFunc("30 * * * *", func() {
		fmt.Println("Every hour on the half hour")
	})

	c.AddFunc("30 3-6,20-23 * * *", func() {
		fmt.Println("On the half hour of 3-6am, 8-11pm")
	})

	c.AddFunc("0 0 1 1 *", func() {
		fmt.Println("Jun 1 every year")
	})

	c.Start()

	//for {
	//	time.Sleep(time.Second)
	//}

	time.Sleep(time.Second * 5)
}

func Test_Cron_Predefined_time(t *testing.T) {
	c := cron.New()

	defer c.Stop()

	c.AddFunc("@every 1s", func() {
		fmt.Println("every 1 second")
	})

	c.AddFunc("@daily", func() {
		fmt.Println("every day on midnight")
	})

	c.AddFunc("@weekly", func() {
		fmt.Println("every week")
	})

	c.AddFunc("@hourly", func() {
		fmt.Println("every hour")
	})

	c.Start()
	time.Sleep(time.Second * 3)
}

func Test_Custom_Parse_Format(t *testing.T) {
	//parser := cron.NewParser(
	//	cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	//)

	//c := cron.New(cron.WithParser(parser))
	c := cron.New(cron.WithSeconds()) //이미 option.go에 설정이 있음

	c.AddFunc("1 * * * * *", func() { //todo: 실행이 안됨
		fmt.Println("every 1 second")
	})

	c.Start()

	time.Sleep(time.Second * 5)
}

func Test_Cron_WithLogger(t *testing.T) {
	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))

	c.AddFunc("@every 1s", func() {
		fmt.Println("hello world")
	})
	c.Start()

	time.Sleep(5 * time.Second)
}

func Test_AddJob(t *testing.T) {
	c := cron.New()
	c.AddJob("@every 1s", job.Hello{Name: "Frank"})
	c.Start()

	time.Sleep(time.Second * 5)
}

func Test_Validate_Cron_Expression_Invalid(t *testing.T) {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)

	var tests = []struct{ expr, err string }{
		{"* 5 j * * *", "failed to parse int from"},
		{"@every Xm", "failed to parse duration"},
		{"@unrecognized", "unrecognized descriptor"},
		{"* * * *", "expected 5 to 6 fields"},
		{"", "empty spec string"},
	}
	for _, tt := range tests {
		actual, err := parser.Parse(tt.expr)
		if err == nil || !strings.Contains(err.Error(), tt.err) {
			t.Errorf("%s => expected %v, got %v", tt.expr, tt.err, err)
		}
		if actual != nil {
			t.Errorf("expected nil schedule on error, got %v", actual)
		}
	}
}

//job안에 panic이 발생해도 recover하고 job의 나머지 작업이 실행된다
func Test_Cron_WithChain_Recover(t *testing.T) {
	c := cron.New()

	c.AddJob("@every 1s",
		cron.NewChain(cron.Recover(cron.DefaultLogger)).
			Then(&job.PanicJob{}))
	c.Start()

	time.Sleep(5 * time.Second)
}

//1초 실행 중인 job을 2초로 delay 시킴
func Test_Cron_WithChain_DelayIfStillRunning(t *testing.T) {
	c := cron.New()
	c.AddJob("@every 1s",
		cron.NewChain(cron.DelayIfStillRunning(cron.DefaultLogger)).
			Then(&job.DelayJob{}))
	c.Start()

	time.Sleep(10 * time.Second)
}

//실행을 skip 함
func Test_Cron_WithChain_SkipIfStillRunning(t *testing.T) {
	c := cron.New()
	c.AddJob("@every 1s",
		cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger)).
			Then(&job.SkipJob{}))
	c.Start()

	time.Sleep(10 * time.Second)
}

func Test_Cron_TimeZone(t *testing.T) {
	nyc, _ := time.LoadLocation("America/New_York")
	c := cron.New(cron.WithLocation(nyc))
	c.AddFunc("0 6 * * ?", func() {
		fmt.Println("Every 6 o'clock at New York")
	})

	c.AddFunc("CRON_TZ=Asia/Tokyo 0 6 * * ?", func() {
		fmt.Println("Every 6 o'clock at Tokyo")
	})

	c.Start()

	time.Sleep(time.Second * 5)
}

func Test_Cron_Remove_Job(t *testing.T) {
	c := cron.New()
	jobList := make([]cron.EntryID, 0)

	entryID, _ := c.AddFunc("@every 1s", func() {
		fmt.Println("job1")
	})
	jobList = append(jobList, entryID)
	entryID, _ = c.AddFunc("@every 1s", func() {
		fmt.Println("job2")
	})
	jobList = append(jobList, entryID)

	c.Start()

	entries := c.Entries()
	assert.Equal(t, 2, len(entries))
	fmt.Println("entryID", jobList[0])
	c.Remove(1)

	assert.Equal(t, 1, len(c.Entries()))
	time.Sleep(10 * time.Second)
}

//todo: job 생성이후 pause, stop 시킬 수 있나?
