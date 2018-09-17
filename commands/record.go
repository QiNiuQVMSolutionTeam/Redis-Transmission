package commands

type Record struct {
	DatabaseId uint64 `json:"db"`
	Key        string `json:"key"`
	Value      string `json:"value"`
	TTL        int64  `json:"ttl"`
}
