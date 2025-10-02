package entity

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetMax(t *testing.T) {

	testCases := []struct {
		name     string
		payload  []int
		expected int
		err      error
	}{
		{
			name:     "empty payload",
			payload:  []int{},
			expected: 0,
			err:      EmptyPayloadError,
		},
		{
			name:     "payload with one value",
			payload:  []int{1},
			expected: 1,
			err:      nil,
		},
		{
			name:     "payload with n value",
			payload:  []int{1, 2, 1000},
			expected: 1000,
			err:      nil,
		},
		{
			name:     "payload unsorted value",
			payload:  []int{1000, 200, 1, 5},
			expected: 1000,
			err:      nil,
		},
		{
			name:     "payload negative value",
			payload:  []int{-1000, -200, -1, -5},
			expected: -1,
			err:      nil,
		},
	}

	for _, tc := range testCases {
		packet, err := NewDataPacket(uuid.New(), time.Now(), tc.payload)

		assert.ErrorIs(t, err, tc.err)
		assert.Equal(t, tc.expected, packet.MaxValue)
	}
}
