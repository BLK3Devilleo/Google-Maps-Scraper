package web

import "errors"

type Location struct {
	ID   string `json:"id"`
	Path string `json:"path"` // e.g. "Valencia/Capital"
	Lat  string `json:"lat"`
	Lon  string `json:"lon"`
}

func (l *Location) Validate() error {
	if l.Path == "" {
		return errors.New("missing path")
	}
	if l.Lat == "" || l.Lon == "" {
		return errors.New("missing coordinates")
	}
	return nil
}
