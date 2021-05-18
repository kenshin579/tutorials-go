package usecase

import (
	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/common/config"
	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/domain"
	"github.com/labstack/gommon/log"
)

type scheduleUsecase struct {
	scheduleStore domain.ScheduleStore
	cfg           *config.Config
}

func NewScheduleUsecase(ss domain.ScheduleStore, cfg *config.Config) domain.ScheduleUsecase {
	return &scheduleUsecase{
		scheduleStore: ss,
		cfg:           cfg,
	}
}

func (s *scheduleUsecase) InitializeJobs() {
	for _, cronJob := range s.cfg.CronConfig {

		err := s.CreateJob(domain.ScheduleRequest{
			JobDescription: cronJob.Description,
			JobType:        cronJob.JobType,
			Schedule:       cronJob.Schedule,
			JobRequest:     cronJob.JobRequest,
		})
		if err != nil {
			log.Error(err)
		}
	}
}

func (s *scheduleUsecase) CreateJob(request domain.ScheduleRequest) error {
	return s.scheduleStore.Create(request)
}

func (s *scheduleUsecase) UpdateJob(request domain.ScheduleRequest) error {
	return s.scheduleStore.Update(request)
}

func (s *scheduleUsecase) ListJob() ([]*domain.ScheduleInfo, error) {
	return s.scheduleStore.List()
}

func (s *scheduleUsecase) DeleteJob(jobID string) error {
	return s.scheduleStore.Delete(jobID)
}

func (s *scheduleUsecase) GetJob(jobID string) (domain.ScheduleInfo, error) {
	return s.scheduleStore.Get(jobID)
}

func (s *scheduleUsecase) Start() error {
	return s.scheduleStore.Start()
}

func (s *scheduleUsecase) Stop() error {
	return s.scheduleStore.Stop()
}
