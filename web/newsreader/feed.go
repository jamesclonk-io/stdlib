package newsreader

import (
	"regexp"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	rss "github.com/jteeuwen/go-pkg-rss"
)

var redditRx = regexp.MustCompile(`<br\/>\s+<a href="(.*)">\[link\]`)

type Feed struct {
	Title string
	URL   string
	Items []FeedItem
}

type FeedItem struct {
	Title    string
	URL      string
	Comments string
}

func (n *NewsReader) getFeeds() (feeds []Feed, err error) {
	// check if enough time has passed and a feed update is needed
	if n.timestamp.Add(time.Duration(n.Config.CacheDuration)).Before(time.Now()) {
		n.timestamp = time.Now()

		n.mutex.Lock()
		defer n.mutex.Unlock()

		go n.updateFeeds()
	}

	return n.Feeds, nil
}

func (n *NewsReader) updateFeeds() {
	var feeds []Feed
	var workChan = make(chan string, 1000)
	var feedChan = make(chan Feed, 1000)
	var doneChan = make(chan bool, 3)
	defer close(doneChan)

	channelHandler := func(feedChan chan Feed) func(*rss.Feed, []*rss.Channel) {
		return func(feed *rss.Feed, channels []*rss.Channel) {
			for _, channel := range channels {
				reddit := strings.Contains(getFeedTitleUrl(channel), "www.reddit.com") // flag feed as reddit

				var feedItems []FeedItem
				for i, item := range channel.Items {
					if i < 25 {
						feedItem := FeedItem{
							Title:    item.Title,
							URL:      getFeedItemUrl(item),
							Comments: item.Comments,
						}
						// special handling for reddit (links & comments)
						if reddit {
							// comments url for reddit rss is actually the feed item link itself
							feedItem.Comments = getFeedItemUrl(item)
							feedItem.URL = getRedditFeedItemUrl(item)
						}

						feedItems = append(feedItems, feedItem)
					}
				}

				feed := Feed{
					Title: channel.Title,
					URL:   getFeedTitleUrl(channel),
					Items: feedItems,
				}
				feedChan <- feed
			}
		}
	}

	for w := 1; w <= 3; w++ {
		go func(workChan chan string, feedChan chan Feed, doneChan chan bool) {
			for url := range workChan {
				rssFeed := rss.New(5, true, channelHandler(feedChan), nil)
				if err := rssFeed.Fetch(url, nil); err != nil &&
					// ignore these encoding errors
					!strings.Contains(err.Error(), `encoding "ISO-8859-1" declared but Decoder.CharsetReader is nil`) {
					n.log.WithFields(logrus.Fields{
						"error": err,
						"feed":  url,
					}).Error("Could not fetch feed data")
				}
			}
			doneChan <- true
		}(workChan, feedChan, doneChan)
	}

	for _, url := range n.Config.Feeds {
		workChan <- url
	}
	close(workChan)

	for w := 1; w <= 3; w++ {
		<-doneChan
	}
	close(feedChan)

	for feed := range feedChan {
		feeds = append(feeds, feed)
	}
	n.Feeds = feeds
}

func getFeedTitleUrl(feed *rss.Channel) string {
	if len(feed.Links) > 0 {
		return feed.Links[0].Href
	}
	return "#"
}

func getFeedItemUrl(item *rss.Item) string {
	if len(item.Links) > 0 {
		return item.Links[0].Href
	}
	return "#"
}

func getRedditFeedItemUrl(item *rss.Item) string {
	matched := redditRx.FindStringSubmatch(item.Description)
	if len(matched) > 1 {
		return matched[1]
	}
	return "#"
}

func isFeedCommentUrl(comments string) bool {
	return strings.HasPrefix(comments, "http")
}
