package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	_ "modernc.org/sqlite" // sqlite driver

	"github.com/gosom/google-maps-scraper/web"
)

type repo struct {
	db *sql.DB
}

func New(path string) (web.JobRepository, error) {
	db, err := initDatabase(path)
	if err != nil {
		return nil, err
	}

	return &repo{db: db}, nil
}

func (repo *repo) Get(ctx context.Context, id string) (web.Job, error) {
	const q = `SELECT * from jobs WHERE id = ?`

	row := repo.db.QueryRowContext(ctx, q, id)

	return rowToJob(row)
}

func (repo *repo) Create(ctx context.Context, job *web.Job) error {
	item, err := jobToRow(job)
	if err != nil {
		return err
	}

	const q = `INSERT INTO jobs (id, name, status, data, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`

	_, err = repo.db.ExecContext(ctx, q, item.ID, item.Name, item.Status, item.Data, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM jobs WHERE id = ?`

	_, err := repo.db.ExecContext(ctx, q, id)

	return err
}

func (repo *repo) Select(ctx context.Context, params web.SelectParams) ([]web.Job, error) {
	q := `SELECT * from jobs`

	var args []any

	if params.Status != "" {
		q += ` WHERE status = ?`

		args = append(args, params.Status)
	}

	q += " ORDER BY created_at DESC"

	if params.Limit > 0 {
		q += " LIMIT ?"

		args = append(args, params.Limit)
	}

	rows, err := repo.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ans []web.Job

	for rows.Next() {
		job, err := rowToJob(rows)
		if err != nil {
			return nil, err
		}

		ans = append(ans, job)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ans, nil
}

func (repo *repo) Update(ctx context.Context, job *web.Job) error {
	item, err := jobToRow(job)
	if err != nil {
		return err
	}

	const q = `UPDATE jobs SET name = ?, status = ?, data = ?, updated_at = ? WHERE id = ?`

	_, err = repo.db.ExecContext(ctx, q, item.Name, item.Status, item.Data, item.UpdatedAt, item.ID)

	return err
}

type scannable interface {
	Scan(dest ...any) error
}

func rowToJob(row scannable) (web.Job, error) {
	var j job

	err := row.Scan(&j.ID, &j.Name, &j.Status, &j.Data, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		return web.Job{}, err
	}

	ans := web.Job{
		ID:     j.ID,
		Name:   j.Name,
		Status: j.Status,
		Date:   time.Unix(j.CreatedAt, 0).UTC(),
	}

	err = json.Unmarshal([]byte(j.Data), &ans.Data)
	if err != nil {
		return web.Job{}, err
	}

	ans.DateUnix = ans.Date.Unix()
	ans.MaxTimeSeconds = int(ans.Data.MaxTime.Seconds())

	return ans, nil
}

func jobToRow(item *web.Job) (job, error) {
	data, err := json.Marshal(item.Data)
	if err != nil {
		return job{}, err
	}

	return job{
		ID:        item.ID,
		Name:      item.Name,
		Status:    item.Status,
		Data:      string(data),
		CreatedAt: item.Date.Unix(),
		UpdatedAt: time.Now().UTC().Unix(),
	}, nil
}

type job struct {
	ID        string
	Name      string
	Status    string
	Data      string
	CreatedAt int64
	UpdatedAt int64
}

func initDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(30 * time.Minute)

	_, err = db.Exec("PRAGMA busy_timeout = 5000")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA synchronous=NORMAL")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA cache_size=1000")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, createSchema(db)
}

func createSchema(db *sql.DB) error {
	if err := createJobsTable(db); err != nil {
		return err
	}
	return createLocationsTable(db)
}

func createJobsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS jobs (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			status TEXT NOT NULL,
			data TEXT NOT NULL,
			created_at INT NOT NULL,
			updated_at INT NOT NULL
		)
	`)
	return err
}

func createLocationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS locations (
			id TEXT PRIMARY KEY,
			path TEXT NOT NULL,
			lat TEXT NOT NULL,
			lon TEXT NOT NULL,
			created_at INT NOT NULL
		)
	`)
	return err
}

func (repo *repo) CreateLocation(ctx context.Context, loc *web.Location) error {
	const q = `INSERT INTO locations (id, path, lat, lon, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := repo.db.ExecContext(ctx, q, loc.ID, loc.Path, loc.Lat, loc.Lon, time.Now().Unix())
	return err
}

func (repo *repo) GetLocations(ctx context.Context) ([]web.Location, error) {
	const q = `SELECT id, path, lat, lon FROM locations ORDER BY path ASC`
	rows, err := repo.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []web.Location
	for rows.Next() {
		var l web.Location
		if err := rows.Scan(&l.ID, &l.Path, &l.Lat, &l.Lon); err != nil {
			return nil, err
		}
		res = append(res, l)
	}
	return res, nil
}

func (repo *repo) DeleteLocation(ctx context.Context, id string) error {
	const q = `DELETE FROM locations WHERE id = ?`
	_, err := repo.db.ExecContext(ctx, q, id)
	return err
}
