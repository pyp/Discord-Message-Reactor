package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/its-vichy/GoCycle"
	"github.com/zenthangplus/goccm"

	"goreactor/internal/discord"
	"goreactor/internal/logging"
	"goreactor/internal/utils"
	"io/ioutil"
	"os"
	"time"
)

var (
	SuccessAmount int
	FailedAmount  int
)

/*
discord community is full of skids, if you are going to skid this, please just give me credit!
github.com/pyp | tokens.sellix.io | doozle
*/

func main() {
	utils.Clear()

	go func() {
		for {
			utils.SetTitle(fmt.Sprintf("GoReactor © doozle 2023 - success=%d failed=%d - tokens.sellix.io", SuccessAmount, FailedAmount))
		}
	}()

	data, err := ioutil.ReadFile("assets/config.json")
	if err != nil {
		logging.Logger.Error().Err(err).Msg("Failed to read assets/config.json.")
		os.Exit(1)
	}

	config, err := utils.ReadConfig(data)
	if err != nil {
		logging.Logger.Error().Msg("Failed to read assets/config.json.")
		os.Exit(1)
	}

	tokens, err := GoCycle.NewFromFile("assets/in/input.txt")
	if err != nil {
		logging.Logger.Error().Err(err).Msg("Failed to read assets/in/input.txt.")
		os.Exit(1)
	} else if len(tokens.List) == 0 {
		logging.Logger.Error().Msg("No tokens found inside assets/in/input.txt.")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}

	proxies, err := GoCycle.NewFromFile("assets/in/proxies.txt")
	if err != nil {
		logging.Logger.Error().Err(err).Msg("Failed to read assets/in/proxies.txt.")
		os.Exit(1)
	} else if len(proxies.List) == 0 {
		logging.Logger.Error().Msg("No proxies found inside assets/in/proxies.txt.")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}

	if utils.ContainsEmoji(config.React.Emoji) {
		config.React.Emoji = url.QueryEscape(config.React.Emoji)
	}

	c := goccm.New(config.Threads)
	logging.Logger.Warn().Int("threads", config.Threads).Int("tokens", len(tokens.List)).Int("proxies", len(proxies.List)).Msg("Starting discord token reacter...")
	for _, Token := range tokens.List {
		c.Wait()

		go func(Token string) {
			defer c.Done()

			if len(Token) < 25 {
				return
			}

			for retry := 0; retry < config.Retry.MaxRetries; retry++ {
				disc := discord.New()

				if !config.Proxyless {
					disc = discord.New(proxies.Next())
				}

				response, err := disc.React(config.React.ChannelID, config.React.MessageID, config.React.Emoji, strings.Split(Token, ":")[3])
				if err != nil {
					FailedAmount++
					logging.Logger.Error().Err(err).Int("status", response).Str("token", strings.Split(Token, ":")[3][:25]+"...").Msg("Something went wrong reacting.")
				}

				switch response {
				case 200, 204:
					SuccessAmount++
					logging.Logger.Info().Int("status", response).Msg(fmt.Sprintf("%v successfully reacted.", strings.Split(Token, ":")[3][:25]+"..."))
				case 429:
					FailedAmount++
					logging.Logger.Error().Int("status", response).Msg(fmt.Sprintf("%v is getting ratelimited.", strings.Split(Token, ":")[3][:25]+"..."))
					time.Sleep(time.Duration(config.Retry.RetryInterval))
					continue
				case 400, 401, 403:
					logging.Logger.Error().Int("status", response).Msg(fmt.Sprintf("%v failed to react.", strings.Split(Token, ":")[3][:25]+"..."))
					FailedAmount++
				default:
					FailedAmount++
				}
				break
			}
		}(Token)
	}
	c.WaitAllDone()
}
