package Util

import (
	"fmt"
	"time"
)

func EpochCalc(onlineat time.Time) (int64, error) {
	online := onlineat.Format("2010-11-04T01:42:54Z")

	snowflakeEpoch, err := time.Parse(time.RFC3339, "2010-11-04T01:42:54Z")
	if err != nil {
		fmt.Println("Error parsing Snowflake epoch:", err)
		return -1, ErrInternal
	}

	startDate, err := time.Parse(time.RFC3339, online)
	if err != nil {
		fmt.Println("Error parsing start date:", err)
		return -1, ErrInternal
	}

	startMillis := startDate.Sub(snowflakeEpoch).Milliseconds()

	fmt.Printf("Milliseconds since Snowflake epoch: %d\n", startMillis)

	return startMillis, nil
}
