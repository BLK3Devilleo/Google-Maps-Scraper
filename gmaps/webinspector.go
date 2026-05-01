package gmaps

import (
	"context"
	"encoding/json"
	"net/url"
	"regexp"
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
	Phones      []string          `json:"phones"`
	Socials     map[string]string `json:"socials"`
}

func (w *WebProfile) CsvRow() []string {
	socials, _ := json.Marshal(w.Socials)
	images := strings.Join(w.Images, " | ")
	videos := strings.Join(w.Videos, " | ")
	emails := strings.Join(w.Emails, " | ")
	phones := strings.Join(w.Phones, " | ")

	return []string{
		w.URL,
		w.Title,
		w.Description,
		w.Summary,
		images,
		videos,
		emails,
		phones,
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
		"phones",
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
		URL:         j.URL,
		Title:       strings.TrimSpace(doc.Find("title").First().Text()),
		Description: strings.TrimSpace(doc.Find("meta[name='description']").AttrOr("content", "")),
		Socials:     docSocialExtractor(doc),
		Emails:      docEmailExtractor(doc),
		Phones:      docPhoneExtractor(doc),
	}

	if profile.Description == "" {
		profile.Description = strings.TrimSpace(doc.Find("meta[property='og:description']").AttrOr("content", ""))
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

func docPhoneExtractor(doc *goquery.Document) []string {
	var phones []string
	body := doc.Find("body").Text()
	// Robust phone regex (common formats)
	re := regexp.MustCompile(`(?i)(?:(?:\+|00)\d{1,3}\s?|0)?(?:[6789]\d{2}\s?\d{3}\s?\d{3}|[6789]\d{8}|[6789]\d{2}\s?\d{2}\s?\d{2}\s?\d{2}|[6789]\d{2}\s?\d{3}\s?\d{4})`)
	matches := re.FindAllString(body, -1)
	
	seen := make(map[string]bool)
	for _, m := range matches {
		m = strings.TrimSpace(m)
		if len(m) >= 9 && !seen[m] {
			phones = append(phones, m)
			seen[m] = true
		}
	}
	return phones
}
