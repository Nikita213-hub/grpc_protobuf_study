package helpers

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ParseTime(timeStr string) (*timestamppb.Timestamp, error) {
	if timeStr == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, err
	}

	return timestamppb.New(t), nil
}
