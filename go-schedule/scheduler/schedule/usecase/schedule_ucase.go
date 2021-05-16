package usecase

import "github.com/kenshin579/tutorials-go/go-schedule/scheduler/domain"

type scheduleUsecase struct {
	scheduleStore domain.ScheduleStore
}

func NewScheduleUsecase(ss domain.ScheduleStore) domain.ScheduleUsecase {
	return &scheduleUsecase{
		scheduleStore: ss,
	}
}

func (s *scheduleUsecase) CreateJob(request domain.ScheduleRequest) error {
	return s.scheduleStore.Create(request)
}

func (s *scheduleUsecase) ListJob() ([]domain.ScheduleInfo, error) {
	return s.scheduleStore.List()
}

func (s *scheduleUsecase) Start() error {
	return s.scheduleStore.Start()
}

func (s *scheduleUsecase) Stop() error {
	return s.scheduleStore.Stop()
}
