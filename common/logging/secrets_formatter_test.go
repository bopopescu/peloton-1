package logging

import (
	"encoding/base64"
	"testing"

	mesos "code.uber.internal/infra/peloton/.gen/mesos/v1"
	"code.uber.internal/infra/peloton/.gen/peloton/api/v0/task"
	"code.uber.internal/infra/peloton/.gen/peloton/private/hostmgr/hostsvc"
	"code.uber.internal/infra/peloton/common"
	jobmgrtask "code.uber.internal/infra/peloton/jobmgr/task"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	testPath      = "/tmp/testpath"
	testSecretStr = "my-secret"
)

// TestDbStatementsFormatting tests that logs containing DB statements and UQL
// queries are filtered and those that contain secret are redacted
func TestDbStatementsFormatting(t *testing.T) {
	formatter := SecretsFormatter{&logrus.JSONFormatter{}}
	b, err := formatter.Format(
		logrus.WithField(common.DBStmtLogField, "INSERT INTO secret_info where"))
	assert.NoError(t, err)
	assert.NotContains(t, string(b), "secret_info")
	assert.Contains(t, string(b), redactedStr)

	// if formatter sees secret_info as part of the DBUqlLogField, it will
	// redact the subsequent DBArgsLogField field in the entry.
	b, err = formatter.Format(
		logrus.WithFields(logrus.Fields{
			common.DBUqlLogField:  "INSERT INTO secret_info where",
			common.DBArgsLogField: "/tmp/secret-path, bXktc2VjcmV0"}))
	assert.NoError(t, err)
	assert.NotContains(t, string(b), "secret_info")
	assert.Contains(t, string(b), redactedStr)
}

// TestLaunchableTasksFormatting tests that logs containing LaunchableTask or
// LaunchTasksRequest are filtered and if a taskconfig contains secret volume,
// the secret data is redacted in the log
func TestLaunchableTasksFormatting(t *testing.T) {
	// setup LaunchableTask, a list of LaunchableTask and LaunchTasksRequest
	// such that the task config contains secret volume
	launchableTaskWithSecret := &hostsvc.LaunchableTask{
		Config: &task.TaskConfig{
			Container: &mesos.ContainerInfo{
				Volumes: []*mesos.Volume{
					jobmgrtask.CreateSecretVolume(testPath, testSecretStr),
				},
			},
		},
	}
	launchableTasksList := []*hostsvc.LaunchableTask{
		launchableTaskWithSecret,
	}
	launchableTasksRequest := &hostsvc.LaunchTasksRequest{
		Tasks: launchableTasksList,
	}

	formatter := SecretsFormatter{&logrus.JSONFormatter{}}

	// launchableTasksRequest contains secret data, it should be redacted
	b, err := formatter.Format(logrus.WithField("req", launchableTasksRequest))
	assert.NoError(t, err)
	// make sure the secret path is kept as it is and is not redacted
	assert.Contains(t, string(b), testPath)
	// make sure the log string contains redactedStr and not the original secret
	assert.NotContains(t, string(b),
		base64.StdEncoding.EncodeToString([]byte(testSecretStr)))
	assert.Contains(t, string(b),
		base64.StdEncoding.EncodeToString([]byte(redactedStr)))

	b, err = formatter.Format(logrus.WithField("list", launchableTasksList))
	assert.NoError(t, err)
	// make sure the secret path is kept as it is and is not redacted
	assert.Contains(t, string(b), testPath)
	// make sure the log string contains redactedStr and not the original secret
	assert.NotContains(t, string(b),
		base64.StdEncoding.EncodeToString([]byte(testSecretStr)))
	assert.Contains(t, string(b),
		base64.StdEncoding.EncodeToString([]byte(redactedStr)))

	b, err = formatter.Format(logrus.WithField("task", launchableTaskWithSecret))
	assert.NoError(t, err)
	// make sure the secret path is kept as it is and is not redacted
	assert.Contains(t, string(b), testPath)
	// make sure the log string contains redactedStr and not the original secret
	assert.NotContains(t, string(b),
		base64.StdEncoding.EncodeToString([]byte(testSecretStr)))
	assert.Contains(t, string(b),
		base64.StdEncoding.EncodeToString([]byte(redactedStr)))
}