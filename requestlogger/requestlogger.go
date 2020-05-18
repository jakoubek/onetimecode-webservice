package requestlogger

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type counterstruct struct {
	Requests  int64 `json:"requests"`
	Timestamp int64 `json:"timestamp"`
}

func ReadCounterfile(filename string) int64 {

	file, _ := ioutil.ReadFile(filename)

	currentCount := counterstruct{}

	_ = json.Unmarshal([]byte(file), &currentCount)

	return currentCount.Requests

}

func SaveCounterfile(filename string, requests int64) {

	currentCount := counterstruct{
		Requests:  requests,
		Timestamp: time.Now().Unix(),
	}

	file, _ := json.Marshal(currentCount)

	_ = ioutil.WriteFile(filename, file, 0644)

}
