package http

import (
	"net/http"

	"github.com/labstack/gommon/log"

	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/common/errors"

	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/domain"
	"github.com/labstack/echo"
)

type scheduleHandler struct {
	scheduleUsecase domain.ScheduleUsecase
}

func NewScheduleHandler(g *echo.Group, su domain.ScheduleUsecase) *scheduleHandler {
	handler := &scheduleHandler{
		scheduleUsecase: su,
	}

	scheduler := g.Group("/scheduler")

	scheduler.POST("/jobs", handler.CreateJob)
	scheduler.GET("/jobs", handler.ListJob)
	scheduler.GET("/jobs/:jobId", handler.GetJob)
	scheduler.DELETE("/jobs/:jobId", handler.DeleteJob)

	scheduler.POST("/start", handler.StartScheduler)
	scheduler.POST("/stop", handler.StopScheduler)

	return handler
}

func (s *scheduleHandler) CreateJob(c echo.Context) error {
	request := domain.ScheduleRequest{}

	if err := c.Bind(&request); err != nil {
		return errors.ErrBinding
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	err := s.scheduleUsecase.CreateJob(request)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *scheduleHandler) ListJob(c echo.Context) error {
	jobList, err := s.scheduleUsecase.ListJob()
	if err != nil {
		log.Error(err)
		return err
	}
	return c.JSON(http.StatusOK, jobList)
}

func (s *scheduleHandler) GetJob(c echo.Context) error {
	job, err := s.scheduleUsecase.GetJob(c.Param("jobId"))
	if err != nil {
		log.Error(err)
		return err
	}
	return c.JSON(http.StatusOK, job)
}

func (s *scheduleHandler) DeleteJob(c echo.Context) error {
	return s.scheduleUsecase.DeleteJob(c.Param("jobId"))
}

func (s *scheduleHandler) StartScheduler(c echo.Context) error {
	err := s.scheduleUsecase.Start()
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (s *scheduleHandler) StopScheduler(c echo.Context) error {
	err := s.scheduleUsecase.Stop()
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
