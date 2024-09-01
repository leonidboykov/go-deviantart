package field

import (
	"fmt"
	"strconv"
	"time"
)

// type Timestamp time.Time

// func (t *Timestamp) UnmarshalJSON(data []byte) error {
// 	if data[0] == '"' && data[len(data)-1] == '"' {
// 		data = data[len(`"`) : len(data)-len(`"`)]
// 	}
// 	val, err := strconv.ParseInt(string(data), 10, 64)
// 	if err != nil {
// 		return fmt.Errorf("parse int: %w", err)
// 	}
// 	*t = Timestamp(time.Unix(val, 0))
// 	return nil
// }

// func (t Timestamp) String() string {
// 	return time.Time(t).String()
// }

type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	if data[0] == '"' && data[len(data)-1] == '"' {
		data = data[len(`"`) : len(data)-len(`"`)]
	}
	val, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return fmt.Errorf("parse int: %w", err)
	}
	t.Time = time.Unix(val, 0)
	return nil
}
