package local

import (
	"errors"
	"strconv"

	"github.com/robfig/cron/v3"

	"github.com/mitchellh/mapstructure"

	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/common/cronner"
	scheduleError "github.com/kenshin579/tutorials-go/go-schedule/scheduler/common/errors"
	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/domain"
	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/domain/job"
	"github.com/labstack/gommon/log"
)

var (
	ErrJobTypeUnknown = errors.New("err job type unknown")
	ErrJobStart       = errors.New("err no job to start")
	ErrJobStop        = errors.New("err no job to stop")
)

type localScheduleStore struct {
	cronner      *cronner.Cronner
	scheduleList []*domain.ScheduleInfo
}

func NewLocalScheduleStore(c *cronner.Cronner) domain.ScheduleStore {
	return &localScheduleStore{
		cronner:      c,
		scheduleList: make([]*domain.ScheduleInfo, 0),
	}
}

func (l *localScheduleStore) Create(request domain.ScheduleRequest) error {
	currentStatus := domain.JobStatusInit
	if len(l.scheduleList) != 0 {
		currentStatus = l.scheduleList[0].Status
	}

	if request.JobType == domain.JobTypeHttpPost {
		var fileRequest job.HttpMoveFile
		err := mapstructure.Decode(request.JobRequest, &fileRequest)
		if err != nil {
			log.Error(err)
			return err
		}

		entryID, err := l.cronner.Cron.AddJob(request.Schedule, &job.HttpMoveFile{
			Url:  fileRequest.Url,
			Body: fileRequest.Body,
		})

		if err != nil {
			log.Error(err)
			return err
		}
		l.scheduleList = append(l.scheduleList, &domain.ScheduleInfo{
			EntryID:        entryID,
			JobDescription: request.JobDescription,
			JobType:        request.JobType,
			Schedule:       request.Schedule,
			Status:         currentStatus,
		})

		return nil
	} else if request.JobType == domain.JobTypePrint {
		var printRequest job.PrintRequest

		err := mapstructure.Decode(request.JobRequest, &printRequest)
		if err != nil {
			log.Error(err)
			return err
		}

		entryID, err := l.cronner.Cron.AddJob(request.Schedule, &job.Print{Message: printRequest.Message})

		if err != nil {
			log.Error(err)
			return err
		}
		l.scheduleList = append(l.scheduleList, &domain.ScheduleInfo{
			EntryID:        entryID,
			JobDescription: request.JobDescription,
			JobType:        request.JobType,
			Schedule:       request.Schedule,
			Status:         currentStatus,
		})
		return nil
	}
	return ErrJobTypeUnknown
}

func (l *localScheduleStore) Update(request domain.ScheduleRequest) error {
	err := l.Delete(request.JobID)
	if err != nil {
		log.Error(err)
		return err
	}

	err = l.Create(request)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (l *localScheduleStore) List() ([]*domain.ScheduleInfo, error) {
	return l.scheduleList, nil
}

func (l *localScheduleStore) Get(jobID string) (domain.ScheduleInfo, error) {
	jobNum, _ := strconv.Atoi(jobID)
	entryID := cron.EntryID(jobNum)

	for _, schedule := range l.scheduleList {

		if schedule.EntryID == entryID {
			return domain.ScheduleInfo{
				EntryID:        entryID,
				JobDescription: schedule.JobDescription,
				JobType:        schedule.JobType,
				Schedule:       schedule.Schedule,
				Status:         schedule.Status,
			}, nil
		}
	}
	return domain.ScheduleInfo{}, scheduleError.ErrNotFoundJob.WithParams(jobID)
}

func (l *localScheduleStore) Delete(jobID string) error {
	jobNum, _ := strconv.Atoi(jobID)
	entryID := cron.EntryID(jobNum)

	temp := l.scheduleList[:0]

	for _, schedule := range l.scheduleList {
		if schedule.EntryID != entryID {
			temp = append(temp, schedule)
		}
	}

	checkError := func(oldSize, newSize int) error {
		if oldSize != newSize {
			return nil
		}
		return scheduleError.ErrNotFoundJob.WithParams(jobID)
	}
	err := checkError(len(l.scheduleList), len(temp))

	l.scheduleList = temp
	return err
}

func (l *localScheduleStore) Start() error {
	if len(l.scheduleList) == 0 {
		return ErrJobStart
	}
	l.setAllStatus(domain.JobStatusStart)
	l.cronner.Cron.Start()
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
