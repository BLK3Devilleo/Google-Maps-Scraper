package web

import (
	"context"
	"errors"
	"time"
)

var jobs []Job

const (
	StatusPending = "pending"
	StatusWorking = "working"
	StatusOK      = "ok"
	StatusFailed  = "failed"
)

const (
	JobTypeGmaps = "maps"
	JobTypeWeb   = "web"
)

type SelectParams struct {
	Status string
	Limit  int
}

type JobRepository interface {
	Get(context.Context, string) (Job, error)
	Create(context.Context, *Job) error
	Delete(context.Context, string) error
	Select(context.Context, SelectParams) ([]Job, error)
	Update(context.Context, *Job) error

	CreateLocation(context.Context, *Location) error
	GetLocations(context.Context) ([]Location, error)
	DeleteLocation(context.Context, string) error
}

type Job struct {
	ID             string
	Name           string
	Type           string
	Date           time.Time
	DateUnix       int64
	Status         string
	Data           JobData
	MaxTimeSeconds int
}

func (j *Job) Validate() error {
	if j.ID == "" {
		return errors.New("missing id")
	}

	if j.Name == "" {
		return errors.New("missing name")
	}

	if j.Status == "" {
		return errors.New("missing status")
	}

	if j.Date.IsZero() {
		return errors.New("missing date")
	}

	if err := j.Data.Validate(); err != nil {
		return err
	}

	return nil
}

type JobData struct {
	Keywords []string      `json:"keywords"`
	Lang     string        `json:"lang"`
	Zoom     int           `json:"zoom"`
	Lat      string        `json:"lat"`
	Lon      string        `json:"lon"`
	FastMode bool          `json:"fast_mode"`
	Radius   int           `json:"radius"`
	Depth    int           `json:"depth"`
	Email    bool          `json:"email"`
	Social   bool          `json:"social"`
	MaxTime  time.Duration `json:"max_time"`
	Proxies  []string      `json:"proxies"`
}

func (d *JobData) Validate() error {
	if len(d.Keywords) == 0 {
		return errors.New("missing keywords")
	}

	if d.Lang == "" {
		// Allow empty lang for web jobs as we are using URLs
		return nil
	}

	if len(d.Lang) != 2 {
		return errors.New("invalid lang")
	}

	return nil
}
