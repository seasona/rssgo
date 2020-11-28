# Problem

This documentation is a record of problems occurred when I develop this program and how I solved them

## RSS feed can't fetch

when use http to request the RSS feed, it may occur two problems:

```text
2020/11/17 14:26:40 Error occur when fetch url: http://feeds.gawker.com/kotaku/vip, err: Failed to detect feed type
2020/11/17 14:26:41 Error occur when fetch url: http://www.gamespot.com/rss/game_updates.php, err: Failed to detect feed type
```

- some subscribe is not free like kotaku, can't get through http. 
- connection will be refused by host, it's caused by GFW in China, so this is a problem can't be fixed

## Load slow

At first, rssgo will fetch all rss subscribe and store them in the database, it allocates the same number of feeds' goroutines to get data and wait all goroutines return then the program go on. This procedure will cost about 10 seconds. It's really slow and you have to wait every time initialize rssgo.