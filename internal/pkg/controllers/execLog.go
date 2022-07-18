package controllers

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/breathbath/go_utils/v2/pkg/fs"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	overwriteExistingLogFileMsg    = "Overwrite existing execution log file (y/n): "
	WillBeExecutedWithClientIDsMsg = "Your task will be executed on the following clients:"
	proceedWithClientIDsMsg        = "Proceed with the above client IDs (y/n): "
	statusFailed                   = "failed"
)

var (
	ErrOverwriteNotConfirmed = errors.New("write execution log requested but overwrite not confirmed")
	ErrClientIDsNotConfirmed = errors.New("client IDs from previous execution log not confirmed")
	ErrNoClientIDsToUse      = errors.New("nothing to do. Execution log file does not contain any failed client IDs")
)

type ExecutionLog struct {
	params       *options.ParameterBag
	logFilename  string
	promptReader config.PromptReader
	hostInfo     *config.HostInfo

	exists             bool
	overwriteConfirmed bool

	logInfo *ExecutionLogInfo
}

type ExecutionLogInfo struct {
	ExecutedAt time.Time     `yaml:"executed_at"`
	ExecutedBy string        `yaml:"executed_by"`
	ExecutedOn string        `yaml:"executed_on"`
	APIUser    string        `yaml:"api_user"`
	APIURL     string        `yaml:"api_url"`
	APIAuth    string        `yaml:"api_auth"`
	NumClients int           `yaml:"num_clients"`
	Failed     int           `yaml:"failed"`
	Jobs       []*models.Job `yaml:"jobs"`
}

func NewExecLog(params *options.ParameterBag,
	logFilename string,
	promptReader config.PromptReader,
	hostInfo *config.HostInfo) (el *ExecutionLog) {
	el = &ExecutionLog{
		params:       params,
		logFilename:  logFilename,
		promptReader: promptReader,
		hostInfo:     hostInfo,

		// assume overwrite unless modified to be otherwise
		overwriteConfirmed: true,

		logInfo: &ExecutionLogInfo{},
	}

	el.exists = fs.FileExists(el.logFilename)

	return el
}

func (el *ExecutionLog) ExistingLog() (exists bool) {
	return el.exists
}

func (el *ExecutionLog) ConfirmOverwrite() (confirmed bool, err error) {
	if el.promptReader == nil || config.ReadNoPrompt(el.params) {
		return true, nil
	}

	el.overwriteConfirmed, err = el.promptReader.ReadConfirmation(overwriteExistingLogFileMsg)
	if err != nil {
		return false, err
	}
	if !el.overwriteConfirmed {
		return false, ErrOverwriteNotConfirmed
	}

	return true, nil
}

func (el *ExecutionLog) ShouldWriteLog() (should bool) {
	return !el.exists || (el.exists && el.overwriteConfirmed)
}

func (el *ExecutionLog) getAuthInfo(params *options.ParameterBag) string {
	if params.ReadString(config.APIToken, "") == "" {
		return "bearer+jwt"
	}
	return "basic+apitoken"
}

func (el *ExecutionLog) SetJobs(executionResults []*models.Job) {
	el.logInfo.Jobs = executionResults
}

func (el *ExecutionLog) ReadFile() (logInfo *ExecutionLogInfo, err error) {
	fileContents, err := os.ReadFile(el.logFilename)
	if err != nil {
		return nil, err
	}

	logInfo = &ExecutionLogInfo{}
	err = yaml.Unmarshal(fileContents, &logInfo)
	if err != nil {
		return nil, err
	}
	el.logInfo = logInfo

	return el.logInfo, nil
}

func (el *ExecutionLog) WriteExecLog(executedAt time.Time, jobs []*models.Job) (err error) {
	el.SetJobs(jobs)
	if len(el.logInfo.Jobs) > 0 {
		err = el.MakeExecutionInfoHeader(executedAt)
		if err != nil {
			return err
		}

		contents, err := yaml.Marshal(el.logInfo)
		if err != nil {
			return err
		}

		err = os.WriteFile(el.logFilename, contents, 0600)
		if err != nil {
			return err
		}
	} else {
		logrus.Warn("no execution results to write")
	}
	return nil
}

func (el *ExecutionLog) MakeExecutionInfoHeader(executedAt time.Time) (err error) {
	if el.hostInfo == nil {
		el.hostInfo, err = config.GetHostInfo()
		if err != nil {
			return err
		}
	} else {
		// override the time supplied by the exec helper. useful for testing.
		executedAt = el.hostInfo.FetchedAt
	}

	el.logInfo.APIUser = el.getAPIUser()
	el.logInfo.APIURL = config.ReadAPIURL(el.params)
	el.logInfo.APIAuth = el.getAuthInfo(el.params)
	el.logInfo.ExecutedAt = executedAt
	el.logInfo.ExecutedBy = el.hostInfo.Username
	el.logInfo.ExecutedOn = el.hostInfo.Machine

	el.logInfo.NumClients = el.getNumClients()
	el.logInfo.Failed = el.getNumClientsWithFailedJobs()

	return nil
}

func (el *ExecutionLog) getAPIUser() (name string) {
	name = config.ReadAPIUser(el.params)
	return name
}

func (el *ExecutionLog) getNumClients() (numClients int) {
	cids := make(map[string]bool, 8)
	for _, job := range el.logInfo.Jobs {
		cids[job.ClientID] = true
	}
	return len(cids)
}

func (el *ExecutionLog) getNumClientsWithFailedJobs() (failedJobCount int) {
	cids := make(map[string]bool, 8)
	for _, job := range el.logInfo.Jobs {
		if job.Status == statusFailed {
			cids[job.ClientID] = true
		}
	}
	return len(cids)
}

func (el *ExecutionLog) GetAndConfirmFailedClientIDs() (clientIDs string, err error) {
	_, err = el.ReadFile()
	if err != nil {
		return "", err
	}

	ids := make(map[string]string, 0)
	for _, job := range el.logInfo.Jobs {
		if job.Status == statusFailed {
			ids[job.ClientID] = job.ClientName
		}
	}

	if len(ids) == 0 {
		return "", ErrNoClientIDsToUse
	}

	displayClientIDs(ids)

	if el.promptReader != nil && !config.ReadNoPrompt(el.params) {
		proceed, err := el.promptReader.ReadConfirmation(proceedWithClientIDsMsg)
		if err != nil {
			return "", err
		}
		if !proceed {
			return "", ErrClientIDsNotConfirmed
		}
	}

	clientIDs = strings.Join(getKeysFromMap(ids), ",")
	return clientIDs, nil
}

func getKeysFromMap(ids map[string]string) (keys []string) {
	keys = make([]string, 0, len(ids))
	for id := range ids {
		keys = append(keys, id)
	}
	return keys
}

func displayClientIDs(ids map[string]string) {
	fmt.Println(WillBeExecutedWithClientIDsMsg)
	for id, name := range ids {
		if name != "" {
			fmt.Printf("%s: %s\n", id, name)
		} else {
			fmt.Printf("%s\n", id)
		}
	}
}
