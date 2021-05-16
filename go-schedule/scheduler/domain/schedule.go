package domain

import "github.com/robfig/cron/v3"

/*
{
    "jobDescription": "move queue -> loading",
    "jobType": "http:post",
    "schedule": "* * * * *",
    "jobRequest" : {
        "url": "https://6254da0e-f2cd-444b-a51e-0f53f1c6a700.mock.pstmn.io/api/move/files",
        "body": {
            "fileName": "test.txt",
            "destination": "/test"
        }
    }
}
*/

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
	JobDescription string      `json:"jobDescription" validate:"required"`
	JobType        JobType     `json:"jobType" validate:"required"`  //todo: 정의된 type만 받을 수 있도록 validate 체크 필요함
	Schedule       string      `json:"schedule" validate:"required"` //todo: cron 표현식만 받을 수 있도록 validate 체크 필요함
	JobRequest     interface{} `json:"jobRequest" validate:"required"`
}

type ScheduleResponse struct {
	ID string `json:"jobId"`
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
	ListJob() ([]ScheduleInfo, error)
	Start() error
	Stop() error
}

type ScheduleStore interface {
	Create(request ScheduleRequest) error
	List() ([]ScheduleInfo, error)
	Start() error
	Stop() error
}
