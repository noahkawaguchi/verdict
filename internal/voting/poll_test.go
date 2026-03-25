package voting_test

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/noahkawaguchi/verdict/internal/voting"
)

func TestValidatePoll_Invalid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, errMsg, prompt string
		choices              []string
	}{
		{
			"empty string prompt",
			"prompt cannot be empty",
			"",
			[]string{"yuzu", "clementine"},
		},
		{
			name:    "empty prompt with choices",
			errMsg:  "prompt cannot be empty",
			choices: []string{"lettuce", "carrot", "green beans"},
		},
		{
			name:    "empty prompt empty choices",
			errMsg:  "prompt cannot be empty",
			choices: []string{},
		},
		{
			"single choice",
			"there must be at least two choices",
			"What is the best fruit?",
			[]string{"yuzu"},
		},
		{
			"empty choices slice",
			"there must be at least two choices",
			"What is the best vegetable?",
			[]string{},
		},
		{
			name:   "nil choices",
			errMsg: "there must be at least two choices",
			prompt: "What is the best color?",
		},
		{
			"empty string in choices",
			"none of the choices can be empty",
			"What is the best fruit?",
			[]string{"yuzu", ""},
		},
		{
			"all empty string choices",
			"none of the choices can be empty",
			"What is the best vegetable?",
			[]string{"", "", ""},
		},
		{
			"empty string choice among valid",
			"none of the choices can be empty",
			"What is the best color?",
			[]string{"red", "blue", "", "yellow", "orange"},
		},
		{
			"duplicate two choices",
			"choices must be unique",
			"What is the best fruit?",
			[]string{"yuzu", "yuzu"},
		},
		{
			"duplicate four choices",
			"choices must be unique",
			"What is the best vegetable?",
			[]string{"lettuce", "carrot", "green beans", "carrot"},
		},
		{
			"duplicate six choices",
			"choices must be unique",
			"What is the best color?",
			[]string{"red", "blue", "blue", "green", "yellow", "orange"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			poll := voting.NewPoll(tt.prompt, tt.choices)
			if err := poll.Validate(); err == nil || err.Error() != tt.errMsg {
				t.Errorf("expected error with message %q, got %v", tt.errMsg, err)
			}
		})
	}
}

func TestValidatePoll_Valid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		prompt  string
		choices []string
	}{
		{"What is the best fruit?", []string{"yuzu", "clementine"}},
		{"What is the best vegetable?", []string{"lettuce", "carrot", "green beans"}},
		{"What is the best color?", []string{"red", "blue", "green", "yellow", "orange"}},
	}

	for _, tt := range tests {
		t.Run(tt.prompt, func(t *testing.T) {
			t.Parallel()
			poll := voting.NewPoll(tt.prompt, tt.choices)
			if err := poll.Validate(); err != nil {
				t.Errorf("expected success, got %v", err)
			}
		})
	}
}

func TestPollMarshalUnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		prompt     string
		choices    []string
		jsonString string
	}{
		{
			"What is the best fruit?",
			[]string{"yuzu", "clementine"},
			`{"prompt":"What is the best fruit?","choices":["yuzu","clementine"]}`,
		},
		{
			"What is the best vegetable?",
			[]string{"lettuce", "carrot", "green beans"},
			`{"prompt":"What is the best vegetable?","choices":["lettuce","carrot","green beans"]}`,
		},
		{
			"What is the best color?",
			[]string{"red", "blue", "green", "yellow", "orange"},
			`{"prompt":"What is the best color?","choices":["red","blue","green","yellow","orange"]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.prompt, func(t *testing.T) {
			t.Parallel()
			inPoll := voting.NewPoll(tt.prompt, tt.choices)
			body, err := json.Marshal(inPoll)
			if err != nil {
				t.Error("failed to marshal JSON:", err)
			}
			if string(body) != tt.jsonString {
				t.Error("unexpected JSON:", string(body))
			}
			var outPoll *voting.Poll
			if err := json.Unmarshal(body, &outPoll); err != nil {
				t.Error("failed to unmarshal JSON:", err)
			}
			if !cmp.Equal(
				inPoll,
				outPoll,
				cmp.AllowUnexported(voting.Poll{}),
				cmpopts.IgnoreFields(voting.Poll{}, "pollID"),
			) {
				t.Error("unexpected unmarshaled poll:", outPoll)
			}
		})
	}
}

func TestPollMarshalUnmarshalDynamoDBAttributeValue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		prompt  string
		choices []string
	}{
		{"What is the best fruit?", []string{"yuzu", "clementine"}},
		{"What is the best vegetable?", []string{"lettuce", "carrot", "green beans"}},
		{"What is the best color?", []string{"red", "blue", "green", "yellow", "orange"}},
	}

	for _, tt := range tests {
		t.Run(tt.prompt, func(t *testing.T) {
			t.Parallel()
			inputPoll := voting.NewPoll(tt.prompt, tt.choices)
			av, err := attributevalue.MarshalMap(inputPoll)
			if err != nil {
				t.Errorf("failed to marshal map: %v", err)
			}
			var p voting.Poll
			if err = attributevalue.UnmarshalMap(av, &p); err != nil {
				t.Errorf("failed to unmarshal map: %v", err)
			}
			if !cmp.Equal(&p, inputPoll, cmp.AllowUnexported(voting.Poll{})) {
				t.Errorf("unexpected unmarshaled result: %+v", p)
			}
		})
	}
}
