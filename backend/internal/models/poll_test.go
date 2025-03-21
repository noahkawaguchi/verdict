package models_test

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

func TestValidatePoll_Invalid(t *testing.T) {
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
		poll := models.NewPoll(test.prompt, test.choices)
		if err := poll.Validate(); err == nil || err.Error() != test.errMsg {
			t.Errorf("expected error with message %q, got %v", test.errMsg, err)
		}
	}
}

func TestValidatePoll_Valid(t *testing.T) {
	tests := []struct {
		prompt  string
		choices []string
	}{
		{"What is the best fruit?", []string{"yuzu", "clementine"}},
		{"What is the best vegetable?", []string{"lettuce", "carrot", "green beans"}},
		{"What is the best color?", []string{"red", "blue", "green", "yellow", "orange"}},
	}
	for _, test := range tests {
		poll := models.NewPoll(test.prompt, test.choices)
		if err := poll.Validate(); err != nil {
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
		inPoll := models.NewPoll(test.prompt, test.choices)
		body, err := json.Marshal(inPoll)
		if err != nil {
			t.Error("failed to marshal JSON:", err)
		}
		if string(body) != test.jsonString {
			t.Error("unexpected JSON:", string(body))
		}
		var outPoll *models.Poll
		if err := json.Unmarshal(body, &outPoll); err != nil {
			t.Error("failed to unmarshal JSON:", err)
		}
		if !cmp.Equal(
			inPoll,
			outPoll,
			cmp.AllowUnexported(models.Poll{}),
			cmpopts.IgnoreFields(models.Poll{}, "pollID"),
		) {
			t.Error("unexpected unmarshaled poll:", outPoll)
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
		inputPoll := models.NewPoll(test.prompt, test.choices)
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
