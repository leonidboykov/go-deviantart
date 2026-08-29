package deviantart

import "uuid"

func uuidsToStrings(uuids ...uuid.UUID) []string {
	result := make([]string, len(uuids))
	for i := range uuids {
		result[i] = uuids[i].String()
	}
	return result
}
