package field

import "encoding/json"

type SingleOrSlice[T any] []T

func (f *SingleOrSlice[T]) UnmarshalJSON(data []byte) error {
	var slice []T
	if err := json.Unmarshal(data, &slice); err != nil {
		var single T
		if err := json.Unmarshal(data, &single); err != nil {
			return err
		}
		*f = []T{single}
	}
	*f = slice
	return nil
}
