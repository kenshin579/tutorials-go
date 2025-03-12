package usecase

import (
	"testing"

	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"

	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/common/config"
	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/common/cronner"

	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/domain"
	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/schedule/store/local"
)

var (
	ss domain.ScheduleStore
	su domain.ScheduleUsecase
	c  *cronner.Cronner
)

func setup() {
	cfg, _ := config.New("../../config/config.yaml")
	c, _ := cronner.New()
	ss = local.NewLocalScheduleStore(c)
	su = NewScheduleUsecase(ss, cfg)
}

func TestInitializeJobs(t *testing.T) {
	//GIVEN
	setup()

	//WHEN
	su.InitializeJobs()
	jobList, err := su.ListJob()
	if err != nil {
		log.Error(err)
	}

	//THEN
	assert.Equal(t, 2, len(jobList))

}
