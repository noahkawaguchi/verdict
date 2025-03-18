package models

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

func TestPollValidateFields_EmptyFields(t *testing.T) {
	tests := []struct {
		errMsg  string
		prompt  string
		choices []string
	}{
		{"prompt cannot be empty", "", []string{"hello", "world"}},
		{"prompt cannot be empty", "", []string{}},
		{"there must be at least two choices", "What is the best fruit?", []string{"hello"}},
		{"there must be at least two choices", "What is the best fruit?", []string{}},
		{"none of the choices can be empty", "What is the best fruit?", []string{"", ""}},
		{"none of the choices can be empty", "What is the best fruit?", []string{"hello", "", "world"}},
	}
	for _, test := range tests {
		poll := NewPoll(test.prompt, test.choices)
		if err := poll.ValidateFields(); err == nil || err.Error() != test.errMsg {
			t.Errorf("expected error with message %q, got %v", test.errMsg, err)
		}
	}
}

func TestPollValidateFields_NonUniqueChoices(t *testing.T) {
	tests := [][]string{
		{"hello", "hello", "world"},
		{"one", "two", "two", "three"},
		{"ha", "ha", "ha", "ha", "ha", "ha"},
	}
	for _, test := range tests {
		poll := NewPoll("What is the best vegetable?", test)
		if err := poll.ValidateFields(); err.Error() != "choices must be unique" {
			t.Errorf("expected error with message \"choices must be unique,\" got %v", err)
		}
	}
}

func TestPollValidateFields_Valid(t *testing.T) {
	tests := []struct {
		prompt  string
		choices []string
	}{
		{"What is the best fruit?", []string{"yuzu", "clementine"}},
		{"What is the best vegetable?", []string{"lettuce", "carrot", "green beans"}},
		{"What is the best color?", []string{"red", "blue", "green", "yellow", "orange"}},
	}
	for _, test := range tests {
		poll := NewPoll(test.prompt, test.choices)
		if err := poll.ValidateFields(); err != nil {
			t.Errorf("expected success, got %v", err)
		}
	}
}

func TestPollMarshalUnmarshalJSON(t *testing.T) {
	tests := []struct {
		prompt     string
		choices    []string
		jsonString string
	}{
		{"What is the best fruit?", []string{"yuzu", "clementine"},
			`{"prompt":"What is the best fruit?","choices":["yuzu","clementine"]}`},
		{"What is the best vegetable?", []string{"lettuce", "carrot", "green beans"},
			`{"prompt":"What is the best vegetable?","choices":["lettuce","carrot","green beans"]}`},
		{"What is the best color?", []string{"red", "blue", "green", "yellow", "orange"},
			`{"prompt":"What is the best color?","choices":["red","blue","green","yellow","orange"]}`},
	}
	for _, test := range tests {
		poll := NewPoll(test.prompt, test.choices)
		body, err := json.Marshal(poll)
		if err != nil {
			t.Errorf("failed to marshal JSON: %v", err)
		}
		if string(body) != test.jsonString {
			t.Errorf("unexpected JSON: %s", string(body))
		}
		var result *Poll
		if err = json.Unmarshal(body, &result); err != nil {
			t.Errorf("failed to unmarshal JSON: %v", err)
		}
		if result.prompt != poll.prompt || !reflect.DeepEqual(result.choices, poll.choices) {
			t.Errorf("unexpected unmarshaled result: %+v", result)
		}
		if result.pollID == "" {
			t.Error("failed to automatically generate poll ID")
		}
	}
}

func TestPollMarshalUnmarshalDynamoDBAttributeValue(t *testing.T) {
	tests := []*Poll{
		{"poll1", "What is the best fruit?", []string{"yuzu", "clementine"}},
		{"poll2", "What is the best vegetable?", []string{"lettuce", "carrot", "green beans"}},
		{"poll5", "What is the best color?", []string{"red", "blue", "green", "yellow", "orange"}},
	}
	for _, test := range tests {
		av, err := attributevalue.MarshalMap(test)
		if err != nil {
			t.Errorf("failed to marshal map: %v", err)
		}
		var p *Poll
		if err = attributevalue.UnmarshalMap(av, &p); err != nil {
			t.Errorf("failed to unmarshal map: %v", err)
		}
		if p.pollID != test.pollID || p.prompt != test.prompt ||
			!reflect.DeepEqual(p.choices, test.choices) {
			t.Errorf("unexpected unmarshaled result: %+v", p)
		}
	}
}
