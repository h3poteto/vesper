package main

import (
	"fmt"
	"github.com/h3poteto/vesper/github"
	"github.com/h3poteto/vesper/twitter"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type config struct {
	ConsumerKey         string `yaml:"consumer_key"`
	ConsumerSecret      string `yaml:"consumer_secret"`
	AccessToken         string `yaml:"access_token"`
	AccessTokenSecret   string `yaml:"access_token_secret"`
	PersonalAccessToken string `yaml:"personal_access_token"`
}

func main() {
	configFile := flag.StringP("config", "c", "setting.yml", "Custom configuration file")
	flag.Parse()

	c, err := initialize(configFile)
	if err != nil {
		log.Panic(err)
	}

	tw := twitter.New(c.ConsumerKey, c.ConsumerSecret, c.AccessToken, c.AccessTokenSecret)
	tweetReport, err := tw.GenerateReport()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(tweetReport)

	gh := github.New(c.PersonalAccessToken)
	githubReport, err := gh.GenerateReport()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(githubReport)
}

func initialize(configFile *string) (*config, error) {
	buf, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return nil, err
	}

	var c config
	err = yaml.Unmarshal(buf, &c)
	return &c, nil
}
