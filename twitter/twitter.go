package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"net/url"
	"time"
)

type Twitter struct {
	client *anaconda.TwitterApi
}

func New(consumerKey, consumerSecret, accessToken, accessTokenSecret string) *Twitter {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	return &Twitter{
		client: api,
	}
}

func (t *Twitter) GenerateReport() (string, error) {
	tweets, err := t.pickup()
	if err != nil {
		return "", err
	}
	report := t.dailyReport(tweets)
	return report, nil
}

// pickup pick up tweet and return tweet list
func (t *Twitter) pickup() ([]anaconda.Tweet, error) {
	params := url.Values{}
	params.Add("exclude_replies", "true")
	params.Add("include_rts", "false")
	tl, err := t.client.GetUserTimeline(params)
	if err != nil {
		return nil, err
	}

	var pickupTweet []anaconda.Tweet
	year, month, day := time.Now().Date()
	loc, _ := time.LoadLocation("UTC")
	startTime := time.Date(year, month, day, 1, 0, 0, 0, loc)
	for _, tweet := range tl {
		tweetTime, _ := time.Parse("Mon Jan 2 15:04:05 +0000 2006", tweet.CreatedAt)
		if tweetTime.After(startTime) {
			pickupTweet = append(pickupTweet, tweet)
		}
	}
	return pickupTweet, nil
}

func (t *Twitter) dailyReport(tweets []anaconda.Tweet) string {
	report := ""
	tweets = t.reverseTweets(tweets)
	for _, t := range tweets {
		tweetTime, _ := time.Parse("Mon Jan 2 15:04:05 +0000 2006", t.CreatedAt)
		tweetTime = tweetTime.Local()
		report = report + "> " + t.Text + "\n\n"
		report = report + "@" + t.User.ScreenName + " - " + tweetTime.Format("Mon Jan 2 15:04 2006") + "\n\n"
	}
	return report
}

func (t *Twitter) reverseTweets(tweets []anaconda.Tweet) []anaconda.Tweet {
	length := len(tweets)
	reverse := make([]anaconda.Tweet, length)
	for i, t := range tweets {
		reverse[length-i-1] = t
	}
	return reverse
}
