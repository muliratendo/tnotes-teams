package service

import (
	"testing"

	"github.com/tendo-mulira/tnotes-teams/internal/ent"
)

func TestCalculateProgress(t *testing.T) {
	tests := []struct {
		name     string
		subtasks []*ent.Subtask
		expected int
	}{
		{
			name:     "no subtasks",
			subtasks: []*ent.Subtask{},
			expected: 0,
		},
		{
			name: "all completed",
			subtasks: []*ent.Subtask{
				{IsCompleted: true},
				{IsCompleted: true},
			},
			expected: 100,
		},
		{
			name: "half completed",
			subtasks: []*ent.Subtask{
				{IsCompleted: true},
				{IsCompleted: false},
			},
			expected: 50,
		},
		{
			name: "none completed",
			subtasks: []*ent.Subtask{
				{IsCompleted: false},
				{IsCompleted: false},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateProgress(tt.subtasks)
			if got != tt.expected {
				t.Errorf("CalculateProgress() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
