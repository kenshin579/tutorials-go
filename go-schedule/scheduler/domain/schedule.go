package domain

import "github.com/robfig/cron/v3"

type JobType string
type JobStatus string

const (
	JobTypeUnknown  JobType = "UnknownType"
	JobTypeHttpPost JobType = "http:post"
	JobTypePrint    JobType = "print"

	JobStatusInit  JobStatus = "INIT"
	JobStatusStart JobStatus = "START"
	JobStatusStop  JobStatus = "STOP"
)

type ScheduleRequest struct {
	JobID          string      `json:"jobId"`
	JobDescription string      `json:"jobDescription" validate:"required"`
	JobType        JobType     `json:"jobType" validate:"required"`  //todo: 정의된 type만 받을 수 있도록 validate 체크 필요함
	Schedule       string      `json:"schedule" validate:"required"` //todo: cron 표현식만 받을 수 있도록 validate 체크 필요함
	JobRequest     interface{} `json:"jobRequest" validate:"required"`
}

type ScheduleResponse struct {
	JobID          int       `json:"jobId"`
	JobDescription string    `json:"jobDescription"`
	JobType        JobType   `json:"jobType"`
	Schedule       string    `json:"schedule"`
	Status         JobStatus `json:"status"`
}

type ScheduleInfo struct {
	EntryID        cron.EntryID
	JobDescription string
	JobType        JobType
	Schedule       string
	Status         JobStatus
}

type ScheduleUsecase interface {
	CreateJob(request ScheduleRequest) error
	ListJob() ([]*ScheduleInfo, error)
	Start() error
	Stop() error
	GetJob(jobID string) (ScheduleInfo, error)
	DeleteJob(jobID string) error
	UpdateJob(request ScheduleRequest) error
}

type ScheduleStore interface {
	Create(request ScheduleRequest) error
	Update(request ScheduleRequest) error
	Get(jobID string) (ScheduleInfo, error)
	List() ([]*ScheduleInfo, error)
	Delete(jobID string) error
	Start() error
	Stop() error
}
