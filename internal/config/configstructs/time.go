package configstructs

import (
	"time"
)

type Time struct {
	MaximumWorking time.Duration `yaml:"maximumWorking"`
}
