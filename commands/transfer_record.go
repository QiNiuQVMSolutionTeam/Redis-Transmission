package commands

import (
	"time"
)

type TransferRecord struct {
	Key   string        `json:"key"`
	Value string        `json:"value"`
	TTL   time.Duration `json:"ttl"`
}
