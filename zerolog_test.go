//go:build !no_zerolog

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package exception_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/rs/zerolog"
	"github.com/thanhminhmr/go-exception"
)

func marshalZerolog(t *testing.T, err exception.Exception) map[string]any {
	t.Helper()
	obj := err.(zerolog.LogObjectMarshaler)
	var buf bytes.Buffer
	logger := zerolog.New(&buf)
	logger.Info().Object("obj", obj).Msg("test")
	var data map[string]any
	if e := json.Unmarshal(buf.Bytes(), &data); e != nil {
		t.Fatalf("failed to unmarshal: %v", e)
	}
	objData, ok := data["obj"].(map[string]any)
	if !ok {
		t.Fatalf("expected obj to be a map, got %T", data["obj"])
	}
	return objData
}

func TestZerolog_String_TypeAndMessage(t *testing.T) {
	data := marshalZerolog(t, exception.String("TestError: message"))
	if data["type"] != "TestError" {
		t.Errorf("type: got %v, want TestError", data["type"])
	}
	if data["message"] != "message" {
		t.Errorf("message: got %v, want message", data["message"])
	}
}

func TestZerolog_String_OmitsEmptyFields(t *testing.T) {
	data := marshalZerolog(t, exception.String("TestError"))
	if data["type"] != "TestError" {
		t.Errorf("type: got %v, want TestError", data["type"])
	}
	if _, ok := data["message"]; ok {
		t.Error("expected message to be omitted")
	}
	for _, key := range []string{"cause", "suppressed", "recovered", "stack_trace", "extras"} {
		if _, ok := data[key]; ok {
			t.Errorf("expected %q to be omitted", key)
		}
	}
}

func requireZerologErr(t *testing.T, actual any, want exception.String) {
	t.Helper()
	m, ok := actual.(map[string]any)
	if !ok {
		t.Fatalf("expected error object, got %T: %v", actual, actual)
	}
	if m["type"] != want.GetType() {
		t.Errorf("type: got %v, want %v", m["type"], want.GetType())
	}
	if m["message"] != want.GetMessage() {
		t.Errorf("message: got %v, want %v", m["message"], want.GetMessage())
	}
}

func TestZerolog_SingleCause(t *testing.T) {
	err := exception.String("TestError: msg").AddCause(one)
	data := marshalZerolog(t, err)
	if data["type"] != "TestError" {
		t.Errorf("type: got %v, want TestError", data["type"])
	}
	if data["message"] != "msg" {
		t.Errorf("message: got %v, want msg", data["message"])
	}
	requireZerologErr(t, data["cause"], one)
}

func TestZerolog_MultipleCauses(t *testing.T) {
	err := exception.String("TestError: msg").AddCause(one, two)
	data := marshalZerolog(t, err)
	causes, ok := data["cause"].([]any)
	if !ok {
		t.Fatalf("expected cause to be an array, got %T", data["cause"])
	}
	if len(causes) != 2 {
		t.Fatalf("cause count: got %d, want 2", len(causes))
	}
	requireZerologErr(t, causes[0], one)
	requireZerologErr(t, causes[1], two)
}

func TestZerolog_Suppressed(t *testing.T) {
	err := exception.String("TestError: msg").AddSuppressed(one)
	data := marshalZerolog(t, err)
	requireZerologErr(t, data["suppressed"], one)
}

func TestZerolog_Recovered(t *testing.T) {
	err := exception.String("TestError: msg").SetRecovered(42)
	data := marshalZerolog(t, err)
	if data["recovered"] != float64(42) {
		t.Errorf("recovered: got %v, want 42", data["recovered"])
	}
}

func TestZerolog_StackTrace(t *testing.T) {
	err := exception.String("TestError: msg").FillStackTrace(0)
	data := marshalZerolog(t, err)
	trace, ok := data["stack_trace"].([]any)
	if !ok {
		t.Fatalf("expected stack_trace to be an array, got %T", data["stack_trace"])
	}
	if len(trace) == 0 {
		t.Fatal("expected non-empty stack_trace")
	}
	firstFrame, ok := trace[0].(map[string]any)
	if !ok {
		t.Fatalf("expected first frame to be a map, got %T", trace[0])
	}
	if fn, ok := firstFrame["function"]; !ok || fn == "" {
		t.Error("expected function to be populated")
	}
	if f, ok := firstFrame["file"]; !ok || f == "" {
		t.Error("expected file to be populated")
	}
	if l, ok := firstFrame["line"]; !ok || l == nil {
		t.Error("expected line to be populated")
	}
}

func TestZerolog_Extras(t *testing.T) {
	err := exception.String("TestError: msg").SetExtras(map[string]any{"key": "value"})
	data := marshalZerolog(t, err)
	extras, ok := data["extras"].(map[string]any)
	if !ok {
		t.Fatalf("expected extras to be a map, got %T", data["extras"])
	}
	if extras["key"] != "value" {
		t.Errorf("extras[key]: got %v, want value", extras["key"])
	}
}

func TestZerolog_JoinedErrors(t *testing.T) {
	err := exception.Join(one, two)
	data := marshalZerolog(t, err)
	causes, ok := data["cause"].([]any)
	if !ok {
		t.Fatalf("expected cause to be an array, got %T", data["cause"])
	}
	if len(causes) != 2 {
		t.Fatalf("cause count: got %d, want 2", len(causes))
	}
	requireZerologErr(t, causes[0], one)
	requireZerologErr(t, causes[1], two)
}

func TestZerolog_FullException_OmitsEmptyFields(t *testing.T) {
	err := exception.String("TestError: msg").AddCause(one)
	data := marshalZerolog(t, err)
	if data["type"] != "TestError" {
		t.Errorf("type: got %v, want TestError", data["type"])
	}
	if data["message"] != "msg" {
		t.Errorf("message: got %v, want msg", data["message"])
	}
	if _, ok := data["cause"]; !ok {
		t.Error("expected cause to be present")
	}
	for _, key := range []string{"suppressed", "recovered", "stack_trace", "extras"} {
		if _, ok := data[key]; ok {
			t.Errorf("expected %q to be omitted", key)
		}
	}
}
