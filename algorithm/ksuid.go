package algorithm

import (
	"github.com/segmentio/ksuid"
)

func NewKsuid() string {
	return ksuid.New().String()
}
