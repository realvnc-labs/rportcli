package models

import "time"

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
