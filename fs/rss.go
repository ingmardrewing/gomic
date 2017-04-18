package fs

import (
	"fmt"

	"github.com/ingmardrewing/gomic/comic"
	"github.com/ingmardrewing/gomic/config"
)

type rss struct {
	comic *comic.Comic
}

func newRss(comic *comic.Comic) *rss {
	return &rss{comic}
}

func (r *rss) RssItem(p *comic.Page) string {
	title := p.Title()
	url := p.Path()
	pubDate := p.Date()
	act := p.Act()
	description := p.Title()
	content := fmt.Sprintf(`<img src="%s">`, p.ImgUrl())
	thumbnailUrl := p.ThumnailUrl()
	imageUrl := p.ImgUrl()
	imageName := p.ImageFilename()
	return fmt.Sprintf(rssItem, title, url, pubDate, act, url, description, content, thumbnailUrl, imageUrl, imageName, thumbnailUrl)
}

func (r *rss) RssItems() string {
	pgs := r.comic.Get10LastComicPagesNewestFirst()
	h := ""
	for _, p := range pgs {
		h += r.RssItem(p)
	}
	return h
}

func (r *rss) Rss() string {
	date := DateNow()
	items := r.RssItems()
	relSelf := config.Servedrootpath() + "/feed/rss.xml"
	return fmt.Sprintf(rssTemplate, relSelf, date, items)
}

var rssTemplate = `<?xml version="1.0" encoding="UTF-8"?><rss version="2.0"
	xmlns:content="http://purl.org/rss/1.0/modules/content/"
	xmlns:wfw="http://wellformedweb.org/CommentAPI/"
	xmlns:dc="http://purl.org/dc/elements/1.1/"
	xmlns:atom="http://www.w3.org/2005/Atom"
	xmlns:sy="http://purl.org/rss/1.0/modules/syndication/"
	xmlns:slash="http://purl.org/rss/1.0/modules/slash/"
	xmlns:media="http://search.yahoo.com/mrss/"
	>

<channel>
	<title>DevAbo.de</title>
    <image>
      <url>https://devabo.de/favicon-32x32.png</url>
      <title>DevAbo.de</title>
      <link>https://devabo.de</link>
      <width>32</width>
      <height>32</height>
      <description>A science-fiction webcomic about the lives of software developers in the far, funny and dystopian future</description>
    </image>
	<atom:link href="%s" rel="self" type="application/rss+xml" />
	<link>https://DevAbo.de</link>
	<description>A science-fiction webcomic about the lives of software developers in the far, funny and dystopian future</description>
	<lastBuildDate>%s</lastBuildDate>
	<language>en-US</language>
	<sy:updatePeriod>weekly</sy:updatePeriod>
	<sy:updateFrequency>1</sy:updateFrequency>
	<generator>https://github.com/ingmardrewing/gomic</generator>
%s
	</channel>
</rss>
`

var rssItem = `  <item>
    <title>%s</title>
    <link>%s</link>
    <pubDate>%s</pubDate>
    <dc:creator><![CDATA[Ingmar Drewing]]></dc:creator>
    <category><![CDATA[%s]]></category>
    <guid>%s/index.html</guid>
    <description><![CDATA[%s]]></description>
    <content:encoded><![CDATA[%s]]></content:encoded>

    <media:thumbnail url="%s" />
    <media:content url="%s" medium="image">
      <media:title type="html">%s</media:title>
      <media:thumbnail url="%s" />
    </media:content>
  </item>
`
