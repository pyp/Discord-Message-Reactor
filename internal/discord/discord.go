package discord

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

func New(proxy ...string) *Discord {
	var proxyURL *url.URL

	if len(proxy) > 0 {
		proxyURL, _ = url.Parse("http://" + proxy[0])
	}

	discord := &Discord{
		UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
		SuperProperties: "eyJvcyI6ICJXaW5kb3dzIiwgImJyb3dzZXIiOiAiQ2hyb21lIiwgImRldmljZSI6ICIiLCAic3lzdGVtX2xvY2FsZSI6ICJlbi1VUyIsICJicm93c2VyX3VzZXJfYWdlbnQiOiAiTW96aWxsYS81LjAgKFdpbmRvd3MgTlQgMTAuMDsgV2luNjQ7IHg2NCkgQXBwbGVXZWJLaXQvNTM3LjM2IChLSFRNTCwgbGlrZSBHZWNrbykgQ2hyb21lLzExMC4wLjU0ODEuMTkyIFNhZmFyaS81MzcuMzYiLCAiYnJvd3Sf sdzZXJfYWdlbnQiOiAiTW96aWxsYS81LjAgKFdpbmRvd3MgTlQgMTAuMDsgV2luNjQ7IHg2NCkgQXBwbGVXZWJLaXQvNTM3LjM2IChLSFRNTCwgbGlrZSBHZWNrbykgQ2hyb21lLzEwMC4wLjU0ODEuMTkyIFNhZmFyaS81MzcuMzYiLCAiYnJvd3Nlcl92ZXJzaW9uIjogIjExMC4wLjU0ODEuMTkyIiwgIm9zX3ZlcnNpb24iOiAiMTAiLCAicmVmZXJyZXIiOiAiIiwgInJlZmVycmluZ19kb21haW4iOiAiIiwgInJlZmVycmluZ19kb21haW5fY3VycmVudCI6ICIiLCAicmVsZWFzZV9jaGFubmVsIjogInN0YWJsZSIsICJjbGllbnRfYnVpbGR_fb nVtYmVyIjogMjE4NjA0LCAiY2xpZW50X2V2ZW50X3NvdXJjZSI6IG51bGx9",
	}

	if proxyURL != nil {
		discord.Proxy = proxyURL.String()
	} else {
		discord.Proxy = ""
	}

	discord.Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MaxVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					0x0a0a, 0x1301, 0x1302, 0x1303, 0xc02b, 0xc02f, 0xc02c, 0xc030,
					0xcca9, 0xcca8, 0xc013, 0xc014, 0x009c, 0x009d, 0x002f, 0x0035,
				},
				InsecureSkipVerify: true,
				CurvePreferences: []tls.CurveID{
					tls.CurveID(0x0a0a),
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
				},
			},
		},
	}

	return discord
}

func (c *Discord) SetHeaders() error {
	dfc, sdc, cfr, err := c.FetchCookie()
	if err != nil {
		return fmt.Errorf("SetHeaders: %v", err)
	}

	c.Headers = map[string]string{
		"authority":          "discord.com",
		"accept":             "*/*",
		"accept-language":    "en-US,en;q=0.9",
		"content-type":       "application/json",
		"origin":             "https://discord.com/",
		"referer":            "https://discord.com/invite/",
		"sec-ch-ua":          `"Chromium";v="110", "Not A(Brand";v="24", "Google Chrome";v="110"`,
		"cookie":             fmt.Sprintf(`locale=en-US; __dcfduid=%v; __sdcfduid=%v; __cfruid=%v`, dfc, sdc, cfr),
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "Windows",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
		"user-agent":         c.UserAgent,
		"x-debug-options":    "bugReporterEnabled",
		"x-discord-locale":   "en-US",
		"x-super-properties": c.SuperProperties,
	}

	return nil
}

func (c *Discord) FetchCookie() (string, string, string, error) {
	req, err := http.NewRequest("GET", "https://discord.com/", nil)
	if err != nil {
		return "", "", "", fmt.Errorf("FetchCookie: %s", err)
	}

	headers := map[string]string{
		"Accept":          "*/*",
		"Accept-Encoding": "gzip, deflate, br",
		"Accept-Language": "en-US",
		"Alt-Used":        "discord.com",
		"Connection":      "keep-alive",
		"Content-Type":    "application/json",
		"Host":            "discord.com",
		"Origin":          "https://discord.com",
		"Referer":         "https://discord.com/",
		"Sec-Fetch-Dest":  "empty",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Site":  "same-origin",
		"TE":              "trailers",
		"User-Agent":      c.UserAgent,
		"X-Track":         c.SuperProperties,
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return "", "", "", fmt.Errorf("FetchCookie: %s", err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("FetchCookie: status code: %d", res.StatusCode)
	}

	var CombinedCookies string
	for _, cookie := range res.Header["Set-Cookie"] {
		CombinedCookies += cookie + "; "
	}

	RegEx := regexp.MustCompile(`__dcfduid=([^;]+).*__sdcfduid=([^;]+).*__cfruid=([^;]+)`)
	matches := RegEx.FindStringSubmatch(CombinedCookies)
	if len(matches) > 3 {
		return matches[1], matches[2], matches[3], nil
	}

	return "", "", "", errors.New("FetchCookie: something went wrong")
}

func (c *Discord) React(channel string, message string, emoji string, token string) (int, error) {
	err := c.SetHeaders()
	if err != nil {
		return 0, fmt.Errorf("Check: %s", err)
	}

	req, err := http.NewRequest("PUT", "https://discord.com/api/v9/channels/"+channel+"/messages/"+message+"/reactions/"+emoji+"/%40me?location=Message&type=0", nil)
	if err != nil {
		return 0, fmt.Errorf("React: %s", err)
	}

	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	req.Header.Set("authorization", token)

	res, err := c.Client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("React: %s", err)
	}

	return res.StatusCode, nil
}
