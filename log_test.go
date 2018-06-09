package log_test

import (
	"github.com/DusanKasan/log"
	"context"
	"strings"
	"fmt"
	"encoding/json"
	"time"
	"testing"
	log2 "log"
)

type logMsg struct {
	Time      time.Time
	Level     string
	RequestID string
	Message   string
	Source    struct{
		FilePath     string
		LineNumber   int
		FunctionName string
	}
	Data      []struct{
		Type      string
		Marshaled interface{}
		Printed   string
		Converted map[string]interface{}
	}
}

func ExampleInfo() {
	logged := &strings.Builder{}
	log.SetOutput(logged)

	log.Info(context.Background(), "data")

	jsonLogs := fmt.Sprintf("[%v]", strings.Join(strings.Split(strings.Trim(logged.String(), "\n"), "\n"), ","))
	var logs []logMsg
	if err := json.Unmarshal([]byte(jsonLogs), &logs); err != nil {
		panic(err)
	}

	fmt.Printf("Log messages recorded: %v\n", len(logs))
	fmt.Printf("Message: %q\n", logs[0].Message)
	fmt.Printf("Level: %q\n", logs[0].Level)
	fmt.Printf("Created in last hour?: %v\n", logs[0].Time.After(time.Now().Add(-1*time.Hour))) // can't check dynamic value in an example
	fmt.Printf("RequestID: %q\n", logs[0].RequestID) // will be empty because no request ID in context
	fmt.Printf("Source:\n")
	segments := strings.Split(logs[0].Source.FilePath, "/")
	fmt.Printf("\tFile: %q\n", segments[len(segments)-1]) // only get the file name from path
	fmt.Printf("\tLine number: %v\n", logs[0].Source.LineNumber)
	fmt.Printf("\tFunction name: %q\n", logs[0].Source.FunctionName)
	fmt.Printf("Data count: %v\n", len(logs[0].Data))
	fmt.Printf("Data[0]:\n")
	fmt.Printf("\tType: %q\n", logs[0].Data[0].Type)
	fmt.Printf("\tMarshaled: %q\n", logs[0].Data[0].Marshaled)
	fmt.Printf("\tPrinted: %q\n", logs[0].Data[0].Printed)
	fmt.Printf("\tConverted: %q\n", logs[0].Data[0].Converted)

	// Output:
	// Log messages recorded: 1
	// Message: "data"
	// Level: "info"
	// Created in last hour?: true
	// RequestID: ""
	// Source:
	// 	File: "log_test.go"
	// 	Line number: 36
	// 	Function name: "github.com/DusanKasan/log_test.ExampleInfo"
	// Data count: 1
	// Data[0]:
	// 	Type: "string"
	// 	Marshaled: "data"
	// 	Printed: "\"data\""
	// 	Converted: map[]
}

func ExampleError() {
	logged := &strings.Builder{}
	log.SetOutput(logged)

	log.Error(context.Background(), "data")

	jsonLogs := fmt.Sprintf("[%v]", strings.Join(strings.Split(strings.Trim(logged.String(), "\n"), "\n"), ","))
	var logs []logMsg
	if err := json.Unmarshal([]byte(jsonLogs), &logs); err != nil {
		panic(err)
	}

	fmt.Printf("Log messages recorded: %v\n", len(logs))
	fmt.Printf("Message: %q\n", logs[0].Message)
	fmt.Printf("Level: %q\n", logs[0].Level)
	fmt.Printf("Created in last hour?: %v\n", logs[0].Time.After(time.Now().Add(-1*time.Hour))) // can't check dynamic value in an example
	fmt.Printf("RequestID: %q\n", logs[0].RequestID) // will be empty because no request ID in context
	fmt.Printf("Source:\n")
	segments := strings.Split(logs[0].Source.FilePath, "/")
	fmt.Printf("\tFile: %q\n", segments[len(segments)-1]) // only get the file name from path
	fmt.Printf("\tLine number: %v\n", logs[0].Source.LineNumber)
	fmt.Printf("\tFunction name: %q\n", logs[0].Source.FunctionName)
	fmt.Printf("Data count: %v\n", len(logs[0].Data))
	fmt.Printf("Data[0]:\n")
	fmt.Printf("\tType: %q\n", logs[0].Data[0].Type)
	fmt.Printf("\tMarshaled: %q\n", logs[0].Data[0].Marshaled)
	fmt.Printf("\tPrinted: %q\n", logs[0].Data[0].Printed)
	fmt.Printf("\tConverted: %q\n", logs[0].Data[0].Converted)

	// Output:
	// Log messages recorded: 1
	// Message: "data"
	// Level: "error"
	// Created in last hour?: true
	// RequestID: ""
	// Source:
	// 	File: "log_test.go"
	// 	Line number: 83
	// 	Function name: "github.com/DusanKasan/log_test.ExampleError"
	// Data count: 1
	// Data[0]:
	// 	Type: "string"
	// 	Marshaled: "data"
	// 	Printed: "\"data\""
	// 	Converted: map[]
}

func ExampleGetWriter() {
	logged := &strings.Builder{}
	log.SetOutput(logged)

	w := log.GetWriter(log.LevelUnknown)
	w.Write([]byte("data"))

	jsonLogs := fmt.Sprintf("[%v]", strings.Join(strings.Split(strings.Trim(logged.String(), "\n"), "\n"), ","))
	var logs []logMsg
	if err := json.Unmarshal([]byte(jsonLogs), &logs); err != nil {
		panic(err)
	}

	fmt.Printf("Log messages recorded: %v\n", len(logs))
	fmt.Printf("Message: %q\n", logs[0].Message)
	fmt.Printf("Level: %q\n", logs[0].Level)
	fmt.Printf("Created in last hour?: %v\n", logs[0].Time.After(time.Now().Add(-1*time.Hour))) // can't check dynamic value in an example
	fmt.Printf("RequestID: %q\n", logs[0].RequestID) // will be empty because no request ID in context
	fmt.Printf("Source:\n")
	segments := strings.Split(logs[0].Source.FilePath, "/")
	fmt.Printf("\tFile: %q\n", segments[len(segments)-1]) // only get the file name from path
	fmt.Printf("\tLine number: %v\n", logs[0].Source.LineNumber)
	fmt.Printf("\tFunction name: %q\n", logs[0].Source.FunctionName)
	fmt.Printf("Data count: %v\n", len(logs[0].Data))
	fmt.Printf("Data[0]:\n")
	fmt.Printf("\tType: %q\n", logs[0].Data[0].Type)
	fmt.Printf("\tMarshaled: %q\n", logs[0].Data[0].Marshaled)
	fmt.Printf("\tPrinted: %q\n", logs[0].Data[0].Printed)
	fmt.Printf("\tConverted: %q\n", logs[0].Data[0].Converted)

	// Output:
	// Log messages recorded: 1
	// Message: "data"
	// Level: "unknown"
	// Created in last hour?: true
	// RequestID: ""
	// Source:
	// 	File: "log_test.go"
	// 	Line number: 131
	// 	Function name: "github.com/DusanKasan/log_test.ExampleGetWriter"
	// Data count: 1
	// Data[0]:
	// 	Type: "string"
	// 	Marshaled: "data"
	// 	Printed: "\"data\""
	// 	Converted: map[]
}

func ExampleNew() {
	logged := &strings.Builder{}

	logger := log.New(logged)
	logger.Info(context.Background(), "data")

	jsonLogs := fmt.Sprintf("[%v]", strings.Join(strings.Split(strings.Trim(logged.String(), "\n"), "\n"), ","))
	var logs []logMsg
	if err := json.Unmarshal([]byte(jsonLogs), &logs); err != nil {
		panic(err)
	}

	fmt.Printf("Log messages recorded: %v\n", len(logs))
	fmt.Printf("Message: %q\n", logs[0].Message)
	fmt.Printf("Level: %q\n", logs[0].Level)
	fmt.Printf("Created in last hour?: %v\n", logs[0].Time.After(time.Now().Add(-1*time.Hour))) // can't check dynamic value in an example
	fmt.Printf("RequestID: %q\n", logs[0].RequestID) // will be empty because no request ID in context
	fmt.Printf("Source:\n")
	segments := strings.Split(logs[0].Source.FilePath, "/")
	fmt.Printf("\tFile: %q\n", segments[len(segments)-1]) // only get the file name from path
	fmt.Printf("\tLine number: %v\n", logs[0].Source.LineNumber)
	fmt.Printf("\tFunction name: %q\n", logs[0].Source.FunctionName)
	fmt.Printf("Data count: %v\n", len(logs[0].Data))
	fmt.Printf("Data[0]:\n")
	fmt.Printf("\tType: %q\n", logs[0].Data[0].Type)
	fmt.Printf("\tMarshaled: %q\n", logs[0].Data[0].Marshaled)
	fmt.Printf("\tPrinted: %q\n", logs[0].Data[0].Printed)
	fmt.Printf("\tConverted: %q\n", logs[0].Data[0].Converted)

	// Output:
	// Log messages recorded: 1
	// Message: "data"
	// Level: "info"
	// Created in last hour?: true
	// RequestID: ""
	// Source:
	// 	File: "log_test.go"
	// 	Line number: 178
	// 	Function name: "github.com/DusanKasan/log_test.ExampleNew"
	// Data count: 1
	// Data[0]:
	// 	Type: "string"
	// 	Marshaled: "data"
	// 	Printed: "\"data\""
	// 	Converted: map[]
}

func ExampleAddDataSerializationConverter() {
	logged := &strings.Builder{}
	log.SetOutput(logged)

	log.AddDataSerializationConverter("myconverter", func(i interface{}) (interface{}, error) {
		if s, ok := i.(string); ok {
			return "123" + s, nil
		}

		return nil, nil
	})

	log.Info(context.Background(), "data", 1)

	jsonLogs := fmt.Sprintf("[%v]", strings.Join(strings.Split(strings.Trim(logged.String(), "\n"), "\n"), ","))
	var logs []logMsg
	if err := json.Unmarshal([]byte(jsonLogs), &logs); err != nil {
		panic(err)
	}

	fmt.Printf("Log messages recorded: %v\n", len(logs))
	fmt.Printf("Message: %q\n", logs[0].Message)
	fmt.Printf("Level: %q\n", logs[0].Level)
	fmt.Printf("Created in last hour?: %v\n", logs[0].Time.After(time.Now().Add(-1*time.Hour))) // can't check dynamic value in an example
	fmt.Printf("RequestID: %q\n", logs[0].RequestID) // will be empty because no request ID in context
	fmt.Printf("Source:\n")
	segments := strings.Split(logs[0].Source.FilePath, "/")
	fmt.Printf("\tFile: %q\n", segments[len(segments)-1]) // only get the file name from path
	fmt.Printf("\tLine number: %v\n", logs[0].Source.LineNumber)
	fmt.Printf("\tFunction name: %q\n", logs[0].Source.FunctionName)
	fmt.Printf("Data count: %v\n", len(logs[0].Data))
	fmt.Printf("Data[0]:\n")
	fmt.Printf("\tType: %q\n", logs[0].Data[0].Type)
	fmt.Printf("\tMarshaled: %q\n", logs[0].Data[0].Marshaled)
	fmt.Printf("\tPrinted: %q\n", logs[0].Data[0].Printed)
	fmt.Printf("\tConverted: %q\n", logs[0].Data[0].Converted)
	fmt.Printf("Data[1]:\n")
	fmt.Printf("\tType: %q\n", logs[0].Data[1].Type)
	fmt.Printf("\tMarshaled: %v\n", logs[0].Data[1].Marshaled)
	fmt.Printf("\tPrinted: %q\n", logs[0].Data[1].Printed)
	fmt.Printf("\tConverted: %q\n", logs[0].Data[1].Converted) // will be empty because it's not a string

	// Output:
	// Log messages recorded: 1
	// Message: "data"
	// Level: "info"
	// Created in last hour?: true
	// RequestID: ""
	// Source:
	// 	File: "log_test.go"
	// 	Line number: 233
	// 	Function name: "github.com/DusanKasan/log_test.ExampleAddDataSerializationConverter"
	// Data count: 2
	// Data[0]:
	// 	Type: "string"
	// 	Marshaled: "data"
	// 	Printed: "\"data\""
	// 	Converted: map["myconverter":"123data"]
	// Data[1]:
	// 	Type: "int"
	// 	Marshaled: 1
	// 	Printed: "1"
	// 	Converted: map[]
}

func BenchmarkWrites(b *testing.B) {
	ctx := context.Background()
	builder1 := &strings.Builder{}
	log.SetOutput(builder1)

	builder2 := &strings.Builder{}
	log2.SetOutput(builder2)

	b.Run("ThisPackage", func(bb *testing.B) {
		for n := 0; n < bb.N; n++ {
			builder1.Reset()
			log.Info(ctx, n)
		}
	})

	b.Run("NativeLog", func(bb *testing.B) {
		for n := 0; n < bb.N; n++ {
			builder1.Reset()
			log2.Print(n)
		}
	})
}

func Benchmark10Fields(b *testing.B) {
	ctx := context.Background()
	builder := &strings.Builder{}
	log.SetOutput(builder)

	b.Run("ThisPackage", func(bb *testing.B) {
		for n := 0; n < bb.N; n++ {
			log.Info(ctx, "message", 1, 2, 3, 4, 5, 6, 7, 8, 9, 0)
		}
	})

}
