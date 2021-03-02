package models

import (
	"github.com/breathbath/go_utils/utils/testing"
	"strconv"
	"time"
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
	Command    string    `json:"command"`
	Shell      string    `json:"shell"`
	Pid        int       `json:"pid"`
	StartedAt  time.Time `json:"started_at"`
	CreatedBy  string    `json:"created_by"`
	MultiJobID string    `json:"multi_job_id"`
	TimeoutSec int       `json:"timeout_sec"`
	Error      string    `json:"error"`
	Result     JobResult `json:"result"`
}

type WsCommand struct {
	Command             string    `json:"command"`
	ClientIds           []string  `json:"client_ids"`
	GroupIds            *[]string `json:"group_ids,omitempty"`
	TimeoutSec          int       `json:"timeout_sec"`
	ExecuteConcurrently bool      `json:"execute_concurrently"`
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
