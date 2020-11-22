# Problem

This documentation is a record of problems occurred when I develop this program and how I solved them

## RSS feed can't fetch

when use http to request the RSS feed, it may occur two problems:

```text
2020/11/17 14:26:40 Error occur when fetch url: http://feeds.gawker.com/kotaku/vip, err: Failed to detect feed type
2020/11/17 14:26:41 Error occur when fetch url: http://www.gamespot.com/rss/game_updates.php, err: Failed to detect feed type
2020/11/17 14:26:45 Error occur when fetch url: http://thepracticaldev.com/feed, err: Failed to detect feed type
2020/11/17 14:27:00 Error occur when fetch url: http://www.reddit.com/r/cpp/.rss, err: Get "http://www.reddit.com/r/cpp/.rss": dial tcp 162.125.1.8:80: connectex: A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.
2020/11/17 14:27:00 Error occur when fetch url: http://www.reddit.com/r/worldnews/.rss, err: Get "http://www.reddit.com/r/worldnews/.rss": dial tcp 162.125.1.8:80: connectex: A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.
2020/11/17 14:27:00 Error occur when fetch url: http://programming.reddit.com/.rss, err: Get "http://programming.reddit.com/.rss": dial tcp 202.160.128.14:80: connectex: A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.
2020/11/17 14:27:00 Error occur when fetch url: http://rss.time.com/web/time/rss/top/index.xml, err: Get "http://rss.time.com/web/time/rss/top/index.xml": dial tcp 66.220.149.18:80: connectex: A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.
2020/11/17 14:27:01 Error occur when fetch url: http://feeds.guardian.co.uk/theguardian/world/rss, err: Get "http://feeds.theguardian.com/theguardian/world/rss": dial tcp 108.160.162.31:80: connectex: A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond.
```

- some subscribe is not free like kotaku, can't get through http. 
- connection will be refused by host, it's caused by GFW in china, so this is a problem can't be fixed