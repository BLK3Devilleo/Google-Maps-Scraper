package gmaps

import (
	"context"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/gosom/scrapemate"
)

type WebProfile struct {
	URL         string            `json:"url"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Summary     string            `json:"summary"`
	Images      []string          `json:"images"`
	Videos      []string          `json:"videos"`
	Emails      []string          `json:"emails"`
	Socials     map[string]string `json:"socials"`
}

func (w *WebProfile) CsvRow() []string {
	socials, _ := json.Marshal(w.Socials)
	images := strings.Join(w.Images, " | ")
	videos := strings.Join(w.Videos, " | ")
	emails := strings.Join(w.Emails, " | ")

	return []string{
		w.URL,
		w.Title,
		w.Description,
		w.Summary,
		images,
		videos,
		emails,
		string(socials),
	}
}

func (w *WebProfile) CsvHeaders() []string {
	return []string{
		"url",
		"title",
		"description",
		"summary",
		"images",
		"videos",
		"emails",
		"socials",
	}
}

type WebInspectorJob struct {
	scrapemate.Job
}

func NewWebInspectorJob(parentID, targetURL string) *WebInspectorJob {
	return &WebInspectorJob{
		Job: scrapemate.Job{
			ID:       uuid.New().String(),
			ParentID: parentID,
			Method:   "GET",
			URL:      targetURL,
			Priority: scrapemate.PriorityHigh,
		},
	}
}

func (j *WebInspectorJob) Process(ctx context.Context, resp *scrapemate.Response) (any, []scrapemate.IJob, error) {
	doc, ok := resp.Document.(*goquery.Document)
	if !ok {
		return nil, nil, nil
	}

	profile := WebProfile{
		URL:     j.URL,
		Title:   doc.Find("title").Text(),
		Socials: docSocialExtractor(doc),
		Emails:  docEmailExtractor(doc),
	}

	// Extract Description/Summary
	profile.Description = doc.Find("meta[name='description']").AttrOr("content", "")
	if profile.Description == "" {
		profile.Description = doc.Find("meta[property='og:description']").AttrOr("content", "")
	}

	// Try to get a summary from the first meaningful paragraph
	doc.Find("p").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		text := strings.TrimSpace(s.Text())
		if len(text) > 50 {
			profile.Summary = text
			return false
		}
		return true
	})

	// Extract Images
	baseURL, _ := url.Parse(j.URL)
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && src != "" {
			absURL := resolveURL(baseURL, src)
			if absURL != "" {
				profile.Images = append(profile.Images, absURL)
			}
		}
	})

	// Extract Videos/iFrames
	doc.Find("video source, iframe").Each(func(_ int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			src, exists = s.Attr("data-src")
		}
		if exists && src != "" {
			absURL := resolveURL(baseURL, src)
			if absURL != "" {
				profile.Videos = append(profile.Videos, absURL)
			}
		}
	})

	return &profile, nil, nil
}

func resolveURL(base *url.URL, rel string) string {
	u, err := url.Parse(rel)
	if err != nil {
		return ""
	}
	return base.ResolveReference(u).String()
}
