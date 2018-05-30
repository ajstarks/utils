package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

type JSONFeed struct {
	Version     string `json:"version"`
	Title       string `json:"title"`
	HomePageURL string `json:"home_page_url"`
	FeedURL     string `json:"feed_url"`
	Description string `json:"description"`
	UserComment string `json:"user_comment"`
	NextURL     string `json:"next_url"`
	Icon        string `json:"icon"`
	Favicon     string `json:"favicon"`
	Author      author `json:"author"`
	Items       []item `json:"items"`
	Expired     bool   `json:"expired"`
	Hubs        []hub  `json:"hubs"`
}

type author struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Avatar string `json:"avatar"`
}

type hub struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type item struct {
	Id            string       `json:"id"`
	ContentText   string       `json:"content_text"`
	ContentHTML   string       `json:"content_html"`
	URL           string       `json:"url"`
	ExternalURL   string       `json:"external_url"`
	Title         string       `json:"title"`
	Summary       string       `json:"summary"`
	Image         string       `json:"image"`
	BannerImage   string       `json:"banner_image"`
	DatePublished string       `json:"date_published"`
	DateModified  string       `json:"date_modified"`
	Author        author       `json:"author"`
	Tags          []string     `json:"tags"`
	Attachments   []attachment `json:"attachments"`
}

type attachment struct {
	URL      string `json:"url"`
	MIMEType string `json:"mime_type"`
	Title    string `json:"title"`
	ByteSize int64  `json:"size_in_bytes"`
	Duration int64  `json:"duration_in_seconds"`
}

func ParseFeed(r io.Reader) (JSONFeed, error) {
	var feed JSONFeed
	err := json.NewDecoder(r).Decode(&feed)
	return feed, err
}

func Titles(w io.Writer, feed JSONFeed) {
	for _, items := range feed.Items {
		fmt.Fprintf(w, "Title: %s\n", items.Title)
	}
}

func Content(w io.Writer, feed JSONFeed) {
	for _, items := range feed.Items {
		if len(items.ContentText) > 0 {
			fmt.Fprintf(w, "Content: %s\n", items.ContentText)
		}
		if len(items.ContentHTML) > 0 {
			fmt.Fprintf(w, "HTML: %s\n", items.ContentHTML)
		}
	}
}

func Top(w io.Writer, feed JSONFeed) {
	fmt.Fprintln(w, "Top Level")
	fmt.Fprintf(w, "\tVersion: %s\n", feed.Version)
	fmt.Fprintf(w, "\tTitle: %s\n", feed.Title)
	fmt.Fprintf(w, "\tFeed URL: %s\n", feed.FeedURL)
	fmt.Fprintf(w, "\tDescription:%s\n", feed.Description)
	fmt.Fprintf(w, "\tComment: %s\n", feed.UserComment)
	fmt.Fprintf(w, "\tNext URL: %s\n", feed.NextURL)
	fmt.Fprintf(w, "\tIcon: %s\n", feed.Icon)
	fmt.Fprintf(w, "\tFavicon: %s\n", feed.Favicon)
	fmt.Fprintf(w, "\tExpired: %v\n", feed.Expired)
	fmt.Fprintf(w, "\tAuthor: %s\n", feed.Author.Name)
	fmt.Fprintf(w, "\tAuthor URL: %s\n", feed.Author.URL)
	fmt.Fprintf(w, "\tAvatar: %s\n", feed.Author.Avatar)
}

func Items(w io.Writer, feed JSONFeed) {
	for i, items := range feed.Items {
		fmt.Fprintf(w, "\nItem [%d]\n", i+1)
		fmt.Fprintf(w, "\tID: %s\n", items.Id)
		fmt.Fprintf(w, "\tURL: %s\n", items.URL)
		fmt.Fprintf(w, "\tExternal URL: %s\n", items.ExternalURL)
		fmt.Fprintf(w, "\tTitle: %s\n", items.Title)
		fmt.Fprintf(w, "\tContent: %s\n", items.ContentText)
		fmt.Fprintf(w, "\tHTML: %s\n", items.ContentHTML)
		fmt.Fprintf(w, "\tSummary: %s\n", items.Summary)
		fmt.Fprintf(w, "\tImage: %s\n", items.Image)
		fmt.Fprintf(w, "\tBanner Image: %s\n", items.BannerImage)
		fmt.Fprintf(w, "\tPublished on: %s\n", items.DatePublished)
		fmt.Fprintf(w, "\tModified on: %s\n", items.DateModified)
		fmt.Fprintf(w, "\tAuthor: %s\n", items.Author.Name)
		fmt.Fprintf(w, "\tAuthor URL: %s\n", items.Author.URL)
		fmt.Fprintf(w, "\tAvatar: %s\n", items.Author.Avatar)
		fmt.Fprintf(w, "\tTags:")
		for _, t := range items.Tags {
			fmt.Fprintf(w, " `%s`", t)
		}
		fmt.Fprintln(w)
		fmt.Fprintln(w, "\tAttachments:")
		for _, att := range items.Attachments {
			fmt.Fprintf(w, "\t\tTitle: %s\n", att.Title)
			fmt.Fprintf(w, "\t\tURL: %s\n", att.URL)
			fmt.Fprintf(w, "\t\tMIME Type: %s\n", att.MIMEType)
			if att.ByteSize > 0 {
				fmt.Fprintf(w, "\t\tSize: %d bytes\n", att.ByteSize)
			}
			if att.Duration > 0 {
				fmt.Fprintf(w, "\t\tDuration: %d seconds\n", att.Duration)
			}
		}
	}
}

func DumpFeed(w io.Writer, feed JSONFeed) {
	Top(w, feed)
	Items(w, feed)
}

func main() {
	var showtitle = flag.Bool("title", false, "show titles")
	var showcontent = flag.Bool("content", false, "show content")
	var showtop = flag.Bool("top", false, "show top level")
	var showall = flag.Bool("all", false, "show all attributes")
	var showitem = flag.Bool("item", false, "show items")
	flag.Parse()
	f, err := ParseFeed(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	if *showall {
		DumpFeed(os.Stdout, f)
	}
	if *showtop {
		Top(os.Stdout, f)
	}
	if *showtitle {
		Titles(os.Stdout, f)
	}
	if *showcontent {
		Content(os.Stdout, f)
	}
	if *showitem {
		Items(os.Stdout, f)
	}
}
