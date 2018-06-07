package log_test

import (
	"testing"
	"strings"
	"time"
	"encoding/json"
	"fmt"
	"github.com/DusanKasan/log"
	"context"
)

func TestInfo(t *testing.T) {
	logged := &strings.Builder{}
	log.SetOutput(logged)

	log.Info(context.Background(), "data")

	logs, err := parseLogs(logged.String())
	if err != nil {
		t.Error(err)
		return
	}

	logs.expect(
		t,
		[]messageExpectation{
			level("info"),
			emptyRequestID(),
			createdBetween(time.Now().Add(-1*time.Hour), time.Now()),
			message("data"),
			source("log_test.go", "github.com/DusanKasan/log_test.TestInfo", 17),
			data([]dataItem{{"string",  "data", "\"data\"",nil}}),
		},
	)
}

func TestError(t *testing.T) {
	logged := &strings.Builder{}
	log.SetOutput(logged)

	log.Error(context.Background(), "data")

	logs, err := parseLogs(logged.String())
	if err != nil {
		t.Error(err)
		return
	}

	logs.expect(
		t,
		[]messageExpectation{
			level("error"),
			emptyRequestID(),
			createdBetween(time.Now().Add(-1*time.Hour), time.Now()),
			message("data"),
			source("log_test.go", "github.com/DusanKasan/log_test.TestError", 42),
			data([]dataItem{{"string",  "data", "\"data\"",nil}}),
		},
	)
}

func TestMultipleMessages(t *testing.T) {
	logged := &strings.Builder{}
	log.SetOutput(logged)

	log.Info(context.Background(), "info data")
	log.Error(context.Background(), "error data")

	logs, err := parseLogs(logged.String())
	if err != nil {
		t.Error(err)
		return
	}

	logs.expect(
		t,
		[]messageExpectation{
			level("info"),
			emptyRequestID(),
			createdBetween(time.Now().Add(-1*time.Hour), time.Now()),
			message("info data"),
			source("log_test.go", "github.com/DusanKasan/log_test.TestMultipleMessages", 67),
			data([]dataItem{{"string",  "info data", "\"info data\"",nil}}),
		},
		[]messageExpectation{
			level("error"),
			emptyRequestID(),
			createdBetween(time.Now().Add(-1*time.Hour), time.Now()),
			message("error data"),
			source("log_test.go", "github.com/DusanKasan/log_test.TestMultipleMessages", 68),
			data([]dataItem{{"string",  "error data", "\"error data\"",nil}}),
		},
	)
}

func TestGetWriter(t *testing.T) {
	logged := &strings.Builder{}
	log.SetOutput(logged)

	w := log.GetWriter(log.LevelUnknown)
	w.Write([]byte("data"))

	logs, err := parseLogs(logged.String())
	if err != nil {
		t.Error(err)
		return
	}

	logs.expect(
		t,
		[]messageExpectation{
			level("unknown"),
			emptyRequestID(),
			createdBetween(time.Now().Add(-1*time.Hour), time.Now()),
			message("data"),
			source("log_test.go", "github.com/DusanKasan/log_test.TestGetWriter", 102),
			data([]dataItem{{"string",  "data", "\"data\"",nil}}),
		},
	)
}

func TestNew(t *testing.T) {
	logged := &strings.Builder{}
	logger := log.New(logged)

	logger.Info(context.Background(), "data")

	logs, err := parseLogs(logged.String())
	if err != nil {
		t.Error(err)
		return
	}

	logs.expect(
		t,
		[]messageExpectation{
			level("info"),
			emptyRequestID(),
			createdBetween(time.Now().Add(-1*time.Hour), time.Now()),
			message("data"),
			source("log_test.go", "github.com/DusanKasan/log_test.TestNew", 127),
			data([]dataItem{{"string",  "data", "\"data\"",nil}}),
		},
	)
}

type logSource struct {
	FilePath     string
	LineNumber   int
	FunctionName string
}

type dataItem struct {
	Type      string
	Marshaled interface{}
	Printed   string
	Converted map[string]interface{}
}

type logMessage struct {
	Time      time.Time
	Level     string
	RequestID string
	Message   string
	Source    logSource
	Data      []dataItem
}

type logMessages []logMessage

type messageExpectation func(logMessage) []error

func (l logMessages) expect(t *testing.T, expectations ...[]messageExpectation) {
	if len(l) != len(expectations) {
		t.Errorf("bad number of messages logged. expected: %v, got :%v", len(expectations), len(l))
		return
	}

	for i, msgExpectations := range expectations {
		l[i].expect(t, msgExpectations...)
	}
}

func parseLogs(s string) (lm logMessages, err error) {
	jsonLogs := fmt.Sprintf("[%v]", strings.Join(strings.Split(strings.Trim(s, "\n"), "\n"), ","))
	err = json.Unmarshal([]byte(jsonLogs), &lm)
	return
}

func (m logMessage) expect(t *testing.T, expectations ...messageExpectation) {
	for _, expectation := range expectations {
		if errs := expectation(m); errs != nil {
			for _, err := range errs {
				t.Error(err.Error())
			}
		}
	}
}

func level(level string) messageExpectation {
	return func (m logMessage) []error {
		if m.Level != level {
			return []error{fmt.Errorf("bad level. expected: %q, got: %q", level, m.Level)}
		}

		return nil
	}
}

func emptyRequestID() messageExpectation {
	return func (m logMessage) []error {
		if m.RequestID != "" {
			return []error{fmt.Errorf("request ID not expected. got: %q", m.RequestID)}
		}

		return nil
	}
}

func createdBetween(from, to time.Time) messageExpectation {
	return func (m logMessage) []error {
		if m.Time.After(to) || m.Time.Before(from) {
			return []error{fmt.Errorf("bad log time. expected between %v and %v, got %v", from, to, m.Time)}
		}

		return nil
	}
}

func message(msg string) messageExpectation {
	return func (m logMessage) []error {
		if m.Message != msg {
			return []error{fmt.Errorf("bad message. expected: %q, got: %q", msg, m.Message)}
		}

		return nil
	}
}

func source(fileName string, funcName string, lineNumber int) messageExpectation {
	return func (m logMessage) []error {
		var errs []error

		if m.Source.FunctionName != funcName {
			errs = append(errs, fmt.Errorf("bad message source function. exepected: %q, got: %q", funcName, m.Source.FunctionName))
		}

		if m.Source.LineNumber != lineNumber {
			errs = append(errs, fmt.Errorf("bad message source line number. exepected: %v, got: %v", lineNumber, m.Source.LineNumber))
		}

		if !strings.HasSuffix(m.Source.FilePath, fileName) {
			errs = append(errs, fmt.Errorf("bad message source file path. exepected to end with: %q, got: %q", fileName, m.Source.FilePath))
		}

		return errs
	}
}

func data(dd []dataItem) messageExpectation {
	return func (m logMessage) []error {
		if len(dd) != len(m.Data) {
			return []error{fmt.Errorf("bad number of data items. exepected %v, got %v", len(dd), len(m.Data))}
		}

		var errs []error
		for i, d := range dd {
			if m.Data[i].Type != d.Type || m.Data[i].Printed != d.Printed || m.Data[i].Marshaled != d.Marshaled || len(m.Data[i].Converted) != 0 {
				errs = append(errs, fmt.Errorf("bad message data. exepected %v, got %v", d, m.Data[i]))
			}
		}

		return errs
	}
}

