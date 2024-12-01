package presentation

import "encoding/json"

func desserializeMessage[T any](data []byte) (T, error) {
	var output T
	err := json.Unmarshal(data, &output)

	return output, err

}
