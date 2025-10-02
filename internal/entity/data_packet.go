package entity

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

var EmptyPayloadError = errors.New("empty payload")

type DataPacket struct {
	ID        uuid.UUID
	Timestamp time.Time
	MaxValue  int
}

func NewDataPacket(id uuid.UUID, timestamp time.Time, payload []int) (DataPacket, error) {

	maxVal, err := getMax(payload)

	if err != nil {
		return DataPacket{}, fmt.Errorf("%w", err)
	}

	packet := DataPacket{
		ID:        id,
		Timestamp: timestamp,
		MaxValue:  maxVal,
	}

	return packet, nil
}

func getMax(payload []int) (int, error) {
	if len(payload) == 0 {
		return 0, EmptyPayloadError
	}
	maxVal := payload[0]
	for _, v := range payload[1:] {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal, nil
}
