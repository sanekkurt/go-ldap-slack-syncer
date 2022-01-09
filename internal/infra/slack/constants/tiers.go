package constants

import (
	"time"
)

// rate limits of api requests per minute

const (
	Tier1 = 1
	Tier2 = 20
	Tier3 = 50
	Tier4 = 100

	TimeoutBetweenPostingMessages  = 1 * time.Second
	TimeoutBetweenIncomingWebhooks = 1 * time.Second
)
