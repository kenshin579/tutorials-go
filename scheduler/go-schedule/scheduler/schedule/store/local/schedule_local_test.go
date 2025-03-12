package local

import (
	"testing"

	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/common/cronner"
	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/domain"
	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/domain/job"
	"github.com/stretchr/testify/assert"
)

var (
	scheduleStore domain.ScheduleStore
)

func setup() {
	c, _ := cronner.New()
	scheduleStore = NewLocalScheduleStore(c)
}

func TestCreate_JotTypeHttpPost(t *testing.T) {
	//GIVEN
	setup()
	request := domain.ScheduleRequest{
		JobDescription: "test-job",
		JobType:        domain.JobTypeHttpPost,
		Schedule:       "* * * * * *",
		JobRequest: job.HttpMoveFile{
			Url: "http://test.com",
			Body: job.MoveFileRequest{
				FileName:    "test.txt",
				Destination: "/test",
			},
		},
	}

	//WHEN
	err := scheduleStore.Create(request)

	//THEN
	assert.NoError(t, err)
	jobList, _ := scheduleStore.List()
	assert.Equal(t, 1, len(jobList))

}

func TestCreate_JotTypePrint(t *testing.T) {
	//GIVEN
	setup()
	request := domain.ScheduleRequest{
		JobDescription: "test-job",
		JobType:        domain.JobTypePrint,
		Schedule:       "* * * * * *",
		JobRequest:     job.PrintRequest{Message: "this is a test"},
	}

	//WHEN
	err := scheduleStore.Create(request)

	//THEN
	assert.NoError(t, err)
	list, _ := scheduleStore.List()
	assert.Equal(t, 1, len(list))
}
