package web

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"sort"
)

var (
	cancels sync.Map
)

type Service struct {
	repo       JobRepository
	dataFolder string
}

func NewService(repo JobRepository, dataFolder string) *Service {
	return &Service{
		repo:       repo,
		dataFolder: dataFolder,
	}
}

func (s *Service) Create(ctx context.Context, job *Job) error {
	job.DateUnix = job.Date.Unix()
	job.MaxTimeSeconds = int(job.Data.MaxTime.Seconds())
	return s.repo.Create(ctx, job)
}

func (s *Service) All(ctx context.Context) ([]Job, error) {
	return s.repo.Select(ctx, SelectParams{})
}

func (s *Service) Get(ctx context.Context, id string) (Job, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if strings.Contains(id, "/") || strings.Contains(id, "\\") || strings.Contains(id, "..") {
		return fmt.Errorf("invalid file name")
	}

	s.Stop(id)

	datapath := filepath.Join(s.dataFolder, id+".csv")

	if _, err := os.Stat(datapath); err == nil {
		if err := os.Remove(datapath); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) Update(ctx context.Context, job *Job) error {
	job.DateUnix = job.Date.Unix()
	job.MaxTimeSeconds = int(job.Data.MaxTime.Seconds())
	return s.repo.Update(ctx, job)
}

func (s *Service) RegisterCancel(id string, cancel context.CancelFunc) {
	cancels.Store(id, cancel)
}

func (s *Service) UnregisterCancel(id string) {
	cancels.Delete(id)
}

func (s *Service) Stop(id string) {
	if cancel, ok := cancels.Load(id); ok {
		if c, ok := cancel.(context.CancelFunc); ok {
			c()
		}
		cancels.Delete(id)
	}
}

func (s *Service) SelectPending(ctx context.Context) ([]Job, error) {
	return s.repo.Select(ctx, SelectParams{Status: StatusPending, Limit: 1})
}

func (s *Service) GetCSV(_ context.Context, id string) (string, error) {
	if strings.Contains(id, "/") || strings.Contains(id, "\\") || strings.Contains(id, "..") {
		return "", fmt.Errorf("invalid file name")
	}

	datapath := filepath.Join(s.dataFolder, id+".csv")

	if _, err := os.Stat(datapath); os.IsNotExist(err) {
		return "", fmt.Errorf("csv file not found for job %s", id)
	}

	return datapath, nil
}

func (s *Service) GetResults(id string) ([]map[string]string, error) {
	datapath, err := s.GetCSV(context.Background(), id)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(datapath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var results []map[string]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		row := make(map[string]string)
		for i, value := range record {
			if i < len(headers) {
				row[headers[i]] = value
			}
		}
		results = append(results, row)
		
		// Limit to prevent crashing the browser with huge files
		if len(results) >= 500 {
			break
		}
	}

	return results, nil
}

func (s *Service) DeleteResults(id string, titlesToDelete []string) error {
	datapath, err := s.GetCSV(context.Background(), id)
	if err != nil {
		return err
	}

	file, err := os.Open(datapath)
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		file.Close()
		return err
	}

	var titleIdx = -1
	for i, h := range headers {
		if h == "title" {
			titleIdx = i
			break
		}
	}

	var records [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			file.Close()
			return err
		}

		shouldDelete := false
		if titleIdx != -1 {
			rowTitle := record[titleIdx]
			for _, t := range titlesToDelete {
				if rowTitle == t {
					shouldDelete = true
					break
				}
			}
		}

		if !shouldDelete {
			records = append(records, record)
		}
	}
	file.Close()

	// Write back the remaining records
	outFile, err := os.Create(datapath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	if err := writer.Write(headers); err != nil {
		return err
	}
	if err := writer.WriteAll(records); err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateLocation(ctx context.Context, loc *Location) error {
	loc.ID = uuid.New().String()
	if err := loc.Validate(); err != nil {
		return err
	}
	return s.repo.CreateLocation(ctx, loc)
}

func (s *Service) GetLocations(ctx context.Context) ([]Location, error) {
	return s.repo.GetLocations(ctx)
}

func (s *Service) DeleteLocation(ctx context.Context, id string) error {
	return s.repo.DeleteLocation(ctx, id)
}

func (s *Service) GetAllWebsites(ctx context.Context) ([]string, error) {
	jobs, err := s.repo.Select(ctx, SelectParams{Status: StatusOK})
	if err != nil {
		return nil, err
	}

	websiteSet := make(map[string]struct{})

	for _, job := range jobs {
		datapath := filepath.Join(s.dataFolder, job.ID+".csv")
		file, err := os.Open(datapath)
		if err != nil {
			continue
		}
		reader := csv.NewReader(file)
		headers, err := reader.Read()
		if err != nil {
			file.Close()
			continue
		}

		websiteIdx := -1
		for i, h := range headers {
			if h == "website" || h == "url" {
				websiteIdx = i
				break
			}
		}

		if websiteIdx == -1 {
			file.Close()
			continue
		}

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}
			if websiteIdx < len(record) {
				ws := strings.TrimSpace(record[websiteIdx])
				if ws != "" &&
					!strings.Contains(ws, "facebook.com") &&
					!strings.Contains(ws, "instagram.com") &&
					!strings.Contains(ws, "twitter.com") &&
					!strings.Contains(ws, "linkedin.com") &&
					(strings.HasPrefix(ws, "http://") || strings.HasPrefix(ws, "https://")) {
					websiteSet[ws] = struct{}{}
				}
			}
		}
		file.Close()
	}

	var websites []string
	for ws := range websiteSet {
		websites = append(websites, ws)
	}
	sort.Strings(websites)
	return websites, nil
}
