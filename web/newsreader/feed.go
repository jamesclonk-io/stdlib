package newsreader

import (
	"fmt"
	"math"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Sirupsen/logrus"
	"github.com/mmcdole/gofeed"
	"github.com/mmcdole/gofeed/rss"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	newsreaderUpdates = promauto.NewCounter(prometheus.CounterOpts{
		Name: "jcio_stdlib_newsreader_feeds_updated",
		Help: "Total number of JCIO stdlib newsreader feeds updated.",
	})
	newsreaderFailures = promauto.NewCounter(prometheus.CounterOpts{
		Name: "jcio_stdlib_newsreader_feed_failures",
		Help: "Total number of JCIO stdlib newsreader feed failures.",
	})
)

type Feeds []Feed
type Feed struct {
	Id    int
	Title string
	URL   string
	Items []FeedItem
}

type FeedItem struct {
	Title    string
	URL      string
	Comments string
}

type Work struct {
	Index int
	URL   string
}

type CommentTranslator struct {
	defaultTranslator *gofeed.DefaultRSSTranslator
}

func NewCommentTranslator() *CommentTranslator {
	t := &CommentTranslator{}

	// We create a DefaultRSSTranslator internally so we can wrap its Translate
	// call since we only want to modify the precedence for a single field.
	t.defaultTranslator = &gofeed.DefaultRSSTranslator{}
	return t
}

func (ct *CommentTranslator) Translate(feed interface{}) (*gofeed.Feed, error) {
	rss, found := feed.(*rss.Feed)
	if !found {
		return nil, fmt.Errorf("Feed did not match expected type of *rss.Feed")
	}

	for i := range rss.Items {
		if len(rss.Items[i].Comments) > 0 {
			rss.Items[i].Link = rss.Items[i].Comments
		}
	}

	f, err := ct.defaultTranslator.Translate(rss)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (n *NewsReader) InitializeFeeds() {
	n.UpdateFeeds()
}

func (n *NewsReader) GetFeeds() (feeds Feeds, err error) {
	// check if enough time has passed and a feed update is needed
	if n.timestamp.Add(time.Duration(n.Config.CacheDuration)).Before(time.Now()) {
		go n.UpdateFeeds()
	}

	return n.Feeds, nil
}

func (n *NewsReader) UpdateFeeds() {
	n.timestamp = time.Now()
	n.mutex.Lock()
	defer n.mutex.Unlock()

	// go through feeds
	var feeds Feeds
	for idx, url := range n.Config.Feeds {
		feedParser := gofeed.NewParser()
		feedParser.RSSTranslator = NewCommentTranslator()
		parsedFeed, err := feedParser.ParseURL(url)
		if err != nil {
			if strings.Contains(err.Error(), `encoding "ISO-8859-1" declared but Decoder.CharsetReader is nil`) || // ignore these encoding errors
				strings.Contains(err.Error(), `429 Too Many Requests`) { // ignore reddit ratelimit
				continue
			} else {
				n.log.WithFields(logrus.Fields{
					"error":    err,
					"feed_id":  idx,
					"feed_url": url,
				}).Error("Could not fetch feed data")
				newsreaderFailures.Inc()
				continue
			}
		}

		var feedItems []FeedItem
		for i, item := range parsedFeed.Items {
			if i < 30 { // max items per feed
				feedItem := FeedItem{
					Title: item.Title,
					URL:   item.Link,
				}
				// // special handling
				// if strings.Contains(parsedFeed.Link, "www.reddit.com") {
				// 	if len(item.Enclosures) > 0 {
				// 		feedItem.URL = item.Enclosures[0].URL
				// 	}
				// }
				feedItems = append(feedItems, feedItem)
			}
		}
		feed := Feed{
			Id:    idx,
			Title: parsedFeed.Title,
			URL:   parsedFeed.Link,
			Items: feedItems,
		}

		// cut off feed size at 30 lines
		for feed.Lines() > 30 {
			feed.Items = feed.Items[:len(feed.Items)-1]
		}
		feeds = append(feeds, feed)
		newsreaderUpdates.Inc()
	}

	// only replace feeds if we got enough of them back
	if len(feeds) >= 4 {
		feeds.sort()
		n.Feeds = feeds
	}
}

func (f Feed) Lines() int {
	// calculate number of lines a feed will occupy
	var lines float64
	for _, item := range f.Items {
		lines += math.Ceil(float64(utf8.RuneCountInString(item.Title)) / 82.0)
	}
	return len(f.Items) + ((int(lines) - len(f.Items)) / 2)
}
