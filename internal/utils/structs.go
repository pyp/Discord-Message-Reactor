package utils

type Config struct {
	Proxyless bool `json:"proxyless"`
	Threads   int  `json:"threads"`
	Retry     struct {
		RetryInterval int `json:"retry_interval"`
		MaxRetries    int `json:"max_retries"`
	} `json:"retry"`
	React struct {
		Emoji     string `json:"emoji"`
		ChannelID string `json:"channel_id"`
		MessageID string `json:"message_id"`
	} `json:"react"`
}
