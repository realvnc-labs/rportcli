package models

import (
	"strconv"
	"time"

	"github.com/breathbath/go_utils/v2/pkg/testing"
)

type JobResult struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

type Job struct {
	Jid        string    `json:"jid"`
	Status     string    `json:"status"`
	FinishedAt time.Time `json:"finished_at"`
	ClientID   string    `json:"client_id"`
	ClientName string    `json:"client_name,omitempty"`
	Command    string    `json:"command"`
	Cwd        string    `json:"cwd"`
	Shell      string    `json:"shell"`
	Pid        int       `json:"pid"`
	StartedAt  time.Time `json:"started_at"`
	CreatedBy  string    `json:"created_by"`
	MultiJobID string    `json:"multi_job_id"`
	TimeoutSec int       `json:"timeout_sec"`
	Error      string    `json:"error"`
	Result     JobResult `json:"result"`
	IsSudo     bool      `json:"sudo"`
	IsScript   bool      `json:"is_script"`
}

type WsCommand struct {
	ClientIds           []string  `json:"client_ids"`
	GroupIds            *[]string `json:"group_ids,omitempty"`
	ExecuteConcurrently bool      `json:"execute_concurrently"`
	AbortOnError        bool      `json:"abort_on_error"`
	IsSudo              bool      `json:"sudo"`
	TimeoutSec          int       `json:"timeout_sec"`
	Cwd                 string    `json:"cwd"`
	Command             string    `json:"command"`
}

func (j *Job) KeyValues() []testing.KeyValueStr {
	return []testing.KeyValueStr{
		{
			Key:   "Job ID",
			Value: j.Jid,
		},
		{
			Key:   "Status",
			Value: j.Status,
		},
		{
			Key:   "Command Output",
			Value: j.Result.Stdout,
		},
		{
			Key:   "Command Error Output",
			Value: j.Result.Stderr,
		},
		{
			Key:   "Started at",
			Value: j.StartedAt.Format(time.RFC3339),
		},
		{
			Key:   "Finished at",
			Value: j.FinishedAt.Format(time.RFC3339),
		},
		{
			Key:   "Client ID",
			Value: j.ClientID,
		},
		{
			Key:   "Command",
			Value: j.Command,
		},
		{
			Key:   "Shell",
			Value: j.Shell,
		},
		{
			Key:   "Pid",
			Value: strconv.Itoa(j.Pid),
		},
		{
			Key:   "Timeout sec",
			Value: strconv.Itoa(j.TimeoutSec),
		},
		{
			Key:   "Created By",
			Value: j.CreatedBy,
		},
		{
			Key:   "Multi Job ID",
			Value: j.MultiJobID,
		},
	}
}
