package models_test

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/google/go-cmp/cmp"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

func TestNewValidatedPoll_Invalid(t *testing.T) {
	tests := []struct {
		errMsg  string
		prompt  string
		choices []string
	}{
		{"prompt cannot be empty", "", []string{"yuzu", "clementine"}},
		{
			errMsg:  "prompt cannot be empty",
			choices: []string{"lettuce", "carrot", "green beans"},
		},
		{
			errMsg:  "prompt cannot be empty",
			choices: []string{},
		},
		{"there must be at least two choices", "What is the best fruit?", []string{"yuzu"}},
		{"there must be at least two choices", "What is the best vegetable?", []string{}},
		{
			errMsg: "there must be at least two choices",
			prompt: "What is the best color?",
		},
		{"none of the choices can be empty", "What is the best fruit?", []string{"yuzu", ""}},
		{"none of the choices can be empty", "What is the best vegetable?", []string{"", "", ""}},
		{"none of the choices can be empty", "What is the best color?",
			[]string{"red", "blue", "", "yellow", "orange"}},
		{"choices must be unique", "What is the best fruit?", []string{"yuzu", "yuzu"}},
		{"choices must be unique", "What is the best vegetable?",
			[]string{"lettuce", "carrot", "green beans", "carrot"}},
		{"choices must be unique", "What is the best color?",
			[]string{"red", "blue", "blue", "green", "yellow", "orange"}},
	}
	for _, test := range tests {
		_, _, err := models.NewValidatedPoll(test.prompt, test.choices)
		if err == nil || err.Error() != test.errMsg {
			t.Errorf("expected error with message %q, got %v", test.errMsg, err)
		}
	}
}

func TestNewValidatedPoll_Valid(t *testing.T) {
	tests := []struct {
		prompt  string
		choices []string
	}{
		{"What is the best fruit?", []string{"yuzu", "clementine"}},
		{"What is the best vegetable?", []string{"lettuce", "carrot", "green beans"}},
		{"What is the best color?", []string{"red", "blue", "green", "yellow", "orange"}},
	}
	for _, test := range tests {
		if _, _, err := models.NewValidatedPoll(test.prompt, test.choices); err != nil {
			t.Errorf("expected success, got %v", err)
		}
	}
}

func TestValidatedPollFromJSON_Invalid(t *testing.T) {
	tests := []struct {
		errMsg     string
		jsonString string
	}{
		{"invalid JSON", `{"prompt":"What is the best fruit?,"choices":["yuzu","clementine"]}`},
		{"invalid JSON",
			`{"prompt":"What is the best vegetable?","choices"["lettuce","carrot","green beans"]}`},
		{"invalid JSON",
			`{"prompt":"What is the best color?","choices":[red","blue","green","yellow","orange"]}`},

		{"prompt cannot be empty", `{"prompt":"","choices":["yuzu","clementine"]}`},
		{"prompt cannot be empty", `{"choices":["lettuce","carrot","green beans"]}`},
		{"prompt cannot be empty", `{"choices":[]}`},

		{"there must be at least two choices",
			`{"prompt":"What is the best fruit?","choices":["yuzu"]}`},
		{"there must be at least two choices",
			`{"prompt":"What is the best vegetable?","choices":[""]}`},
		{"there must be at least two choices", `{"prompt":"What is the best color?"}`},

		{"none of the choices can be empty",
			`{"prompt":"What is the best fruit?","choices":["yuzu",""]}`},
		{"none of the choices can be empty",
			`{"prompt":"What is the best vegetable?","choices":["","",""]}`},
		{"none of the choices can be empty",
			`{"prompt":"What is the best color?","choices":["red","blue","","yellow","orange"]}`},

		{"choices must be unique",
			`{"prompt":"What is the best fruit?","choices":["yuzu","yuzu"]}`},
		{"choices must be unique",
			`{"prompt":"What is the best vegetable?","choices":["lettuce","carrot","green beans","carrot"]}`},
		{"choices must be unique",
			`{"prompt":"What is the best color?","choices":["red","blue","blue","green","yellow","orange"]}`},
	}
	for _, test := range tests {
		_, _, err := models.ValidatedPollFromJSON(test.jsonString)
		if err == nil || err.Error() != test.errMsg {
			t.Errorf("expected error with message %q, got %v", test.errMsg, err)
		}
	}
}

func TestValidatedPollFromJSON_Valid(t *testing.T) {
	tests := []string{
		`{"prompt":"What is the best fruit?","choices":["yuzu","clementine"]}`,
		`{"prompt":"What is the best vegetable?","choices":["lettuce","carrot","green beans"]}`,
		`{"prompt":"What is the best color?","choices":["red","blue","green","yellow","orange"]}`,
	}
	for _, test := range tests {
		if _, _, err := models.ValidatedPollFromJSON(test); err != nil {
			t.Errorf("expected success, got %v", err)
		}
	}
}

func TestPollMarshalJSON(t *testing.T) {
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
		poll, _, err := models.NewValidatedPoll(test.prompt, test.choices)
		if err != nil {
			t.Error("unexpected error creating poll:", err)
		}
		body, err := json.Marshal(poll)
		if err != nil {
			t.Errorf("failed to marshal JSON: %v", err)
		}
		if string(body) != test.jsonString {
			t.Errorf("unexpected JSON: %s", string(body))
		}
	}
}

func TestPollMarshalUnmarshalDynamoDBAttributeValue(t *testing.T) {
	tests := []struct {
		prompt  string
		choices []string
	}{
		{"What is the best fruit?", []string{"yuzu", "clementine"}},
		{"What is the best vegetable?", []string{"lettuce", "carrot", "green beans"}},
		{"What is the best color?", []string{"red", "blue", "green", "yellow", "orange"}},
	}
	for _, test := range tests {
		inputPoll, _, err := models.NewValidatedPoll(test.prompt, test.choices)
		if err != nil {
			t.Error("unexpected error creating poll:", err)
		}
		av, err := attributevalue.MarshalMap(inputPoll)
		if err != nil {
			t.Errorf("failed to marshal map: %v", err)
		}
		var p models.Poll
		if err = attributevalue.UnmarshalMap(av, &p); err != nil {
			t.Errorf("failed to unmarshal map: %v", err)
		}
		if !cmp.Equal(&p, inputPoll, cmp.AllowUnexported(models.Poll{})) {
			t.Errorf("unexpected unmarshaled result: %+v", p)
		}
	}
}
