package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/pwpon500/labs-bot/pkg/notifier"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"

	"github.com/nlopes/slack"
	"github.com/spf13/viper"
)

var (
	version   = "0.1"
	configLoc = flag.String("config", "", "path to config")
)

type config struct {
	Subreddits   []string `mapstructure:"subreddits"`
	SlackToken   string   `mapstructure:"slack_token"`
	SlackChannel string   `mapstructure:"slack_channel"`
	TriggerWords []string `mapstructure:"trigger_words"`
	RedditAuth   struct {
		ClientID     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"`
		Username     string `mapstructure:"username"`
		Password     string `mapstructure:"password"`
	} `mapstructure:"reddit_auth"`
}

func main() {
	flag.Parse()
	if *configLoc == "" {
		log.Fatalln("Please specify a config location with the `-config` flag")
	}

	viper.SetConfigFile(*configLoc)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read in config: %v\n", err)
	}

	conf := &config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		log.Fatalf("Failed to unmarshal config into struct: %v\n", err)
	}

	botConf := reddit.BotConfig{
		Agent: "linux:labs-bot:" + version + "by /u/Pwpon500",
		App: reddit.App{
			ID:       conf.RedditAuth.ClientID,
			Secret:   conf.RedditAuth.ClientSecret,
			Username: conf.RedditAuth.Username,
			Password: conf.RedditAuth.Password,
		},
		Rate: 0,
	}

	bot, err := reddit.NewBot(botConf)
	if err != nil {
		log.Fatalln(err)
	}

	cfg := graw.Config{
		Subreddits:        conf.Subreddits,
		SubredditComments: conf.Subreddits,
	}

	api := slack.New(conf.SlackToken)

	notif := notifier.New(api, conf.TriggerWords, conf.SlackChannel)
	_, wait, err := graw.Run(notif, bot, cfg)
	if err != nil {
		log.Fatalf("Failed to listen for reddit updates: %v\n", err)
	}

	fmt.Println("graw run failed: ", wait())
}
