package discord

import "net/http"

type Discord struct {
	Client          *http.Client
	Proxy           string
	UserAgent       string
	SuperProperties string
	CapSolverKey    string
	Fingerprint     string
	Headers         map[string]string
}

type DiscordResponse struct {
	Username string `json:"username"`
}
