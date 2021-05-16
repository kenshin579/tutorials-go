package local

import (
	"errors"

	"github.com/mitchellh/mapstructure"

	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/domain/task"

	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/common/cronner"
	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/domain"
	"github.com/labstack/gommon/log"
)

var (
	ErrJobTypeUnknown = errors.New("err job type unknown")
	ErrJobStart       = errors.New("err no job to start")
	ErrJobStop        = errors.New("err no job to stop")
)

type localScheduleStore struct {
	cronner      *cronner.Cronner
	scheduleList []domain.ScheduleInfo
}

func NewLocalScheduleStore(c *cronner.Cronner) domain.ScheduleStore {
	return &localScheduleStore{
		cronner:      c,
		scheduleList: make([]domain.ScheduleInfo, 0),
	}
}

func (l *localScheduleStore) Create(request domain.ScheduleRequest) error {

	if request.JobType == domain.JobTypeHttpPost {
		var fileRequest task.HttpMoveFile
		err := mapstructure.Decode(request.JobRequest, &fileRequest)
		if err != nil {
			log.Error(err)
			return err
		}

		entryID, err := l.cronner.Cron.AddJob(request.Schedule, &task.HttpMoveFile{
			Url:  fileRequest.Url,
			Body: fileRequest.Body,
		})

		if err != nil {
			log.Error(err)
			return err
		}
		l.scheduleList = append(l.scheduleList, domain.ScheduleInfo{
			EntryID:        entryID,
			JobDescription: request.JobDescription,
			JobType:        request.JobType,
			Schedule:       request.Schedule,
			Status:         domain.JobStatusInit,
		})

		return nil
	} else if request.JobType == domain.JobTypePrint {
		var printRequest task.PrintRequest

		err := mapstructure.Decode(request.JobRequest, &printRequest)
		if err != nil {
			log.Error(err)
			return err
		}

		entryID, err := l.cronner.Cron.AddJob(request.Schedule, &task.Print{Message: printRequest.Message})

		if err != nil {
			log.Error(err)
			return err
		}
		l.scheduleList = append(l.scheduleList, domain.ScheduleInfo{
			EntryID:        entryID,
			JobDescription: request.JobDescription,
			JobType:        request.JobType,
			Schedule:       request.Schedule,
			Status:         domain.JobStatusInit,
		})
		return nil
	}
	return ErrJobTypeUnknown
}

func (l *localScheduleStore) List() ([]domain.ScheduleInfo, error) {
	return l.scheduleList, nil
}

func (l *localScheduleStore) Start() error {
	if len(l.scheduleList) == 0 {
		return ErrJobStart
	}
	l.cronner.Cron.Start()
	l.setAllStatus(domain.JobStatusStart)
	return nil
}

func (l *localScheduleStore) Stop() error {
	if len(l.scheduleList) == 0 {
		return ErrJobStop
	}

	stop := l.cronner.Cron.Stop()
	l.setAllStatus(domain.JobStatusStop)
	err := stop.Err()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (l *localScheduleStore) setAllStatus(status domain.JobStatus) {
	for _, schedule := range l.scheduleList {
		schedule.Status = status
	}
}
