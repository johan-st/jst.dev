package repo

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	ulid "github.com/oklog/ulid/v2"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type TursoRepo struct {
	db *sql.DB
}

func NewTursoRepo() (*TursoRepo, error) {
	dbURL := os.Getenv("DB_TURSO_URL")
	dbToken := os.Getenv("DB_TURSO_TOKEN")

	if dbURL == "" || dbToken == "" {
		log.Fatal("Missing required database environment variables")
	}

	db, err := sql.Open("libsql", dbURL+"?authToken="+dbToken)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	repo := &TursoRepo{
		db: db,
	}

	if err := repo.RunMigrations(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (t *TursoRepo) Close() error {
	return t.db.Close()
}

// Health returns an error if the database is not reachable
func (t *TursoRepo) Health() error {
	return t.db.Ping()
}

// BLOG
type BlogPost struct {
	Title  string    `db:"title"`
	Date   time.Time `db:"date"`
	Short  string    `db:"short"`
	Body   string    `db:"body"`
	Author Author    `db:"author"`
	Slug   string    `db:"slug"`
}

type Author struct {
	Name     string `db:"name"`
	ImageURL string `db:"image_url"`
}

// BlogPostsFeatured returns a list of featured blog posts
// TODO: implement
func (t *TursoRepo) BlogPostsFeatured(limit, offset int) ([]BlogPost, error) {
	if limit < 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}
	if offset < 0 {
		return nil, fmt.Errorf("offset must be greater than 0")
	}

	posts := []BlogPost{}
	for i := 0; i < limit; i++ {
		posts = append(posts, BlogPost{
			Slug:  fmt.Sprintf("%s-%d", "slug", i),
			Title: "Boost your creativity",
			Date:  time.Now().AddDate(-1, i, -37).Add(time.Hour * 12),
			Body:  "body",
			Short: "Libero neque aenean tincidunt nec consequat tempor. Viverra odio id velit adipiscing id. Nisi vestibulum orci eget bibendum dictum. Velit viverra posuere vulputate volutpat nunc. Nunc netus sit faucibus.",
			Author: Author{
				Name:     "Michael Foster",
				ImageURL: "https://images.unsplash.com/photo-1519244703995-f4e0f30006d5?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80",
			},
		})
	}
	return posts, nil
}

func (t *TursoRepo) BlogPostBySlug(slug string) (BlogPost, error) {
	return BlogPost{
		Slug:  slug,
		Title: "Boost your creativity",
		Date:  time.Now().AddDate(-1, 0, 0).Add(time.Hour * 12),
		Body:  "#### Boost your creativity\n\nbody for " + slug + "\n\n```\nbody for " + slug + "\n```",
		Short: "short for " + slug,
		Author: Author{
			Name:     "Michael Foster",
			ImageURL: "https://images.unsplash.com/photo-1519244703995-f4e0f30006d5?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80",
		},
	}, nil
}

// LOG
// AccessLog represents a single access request to the service
type AccessLog struct {
	ID            int64     `db:"id"`
	Timestamp     time.Time `db:"timestamp"`
	RemoteAddr    string    `db:"remote_addr"`
	RequestMethod string    `db:"request_method"`
	RequestURI    string    `db:"request_uri"`
	Protocol      string    `db:"protocol"`
	StatusCode    int       `db:"status_code"`
	UserAgent     string    `db:"user_agent"`
	Referer       string    `db:"referer"`
}

type EventLog struct {
	Id          ulid.ULID `db:"id"` // Unique and monotonically sortable (ref: https://github.com/oklog/ulid)
	Timestamp   time.Time `db:"timestamp"`
	ShortCode   string    `db:"short_code"`
	Description string    `db:"description"`
	Severity    Severity  `db:"severity"`
}

type Severity string

const (
	Info    Severity = "info"
	Warning Severity = "warning"
	Error   Severity = "error"
)

// AccessLogInsert inserts a new access log entry into the database
func (t *TursoRepo) AccessLogInsert(req *http.Request) error {
	remoteAddr := req.RemoteAddr
	// use X-Forwarded-For header if available (e.g. from reverse proxy)
	if forwarded := req.Header.Get("X-Forwarded-For"); forwarded != "" {
		remoteAddr = forwarded
	}
	_, err := t.db.Exec(`
		INSERT INTO access_logs (
			timestamp,
			remote_addr,
			request_method,
			request_uri, 
			protocol,
			status_code,
			user_agent,
			referer
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		time.Now(),
		remoteAddr,
		req.Method,
		req.RequestURI,
		req.Proto,
		http.StatusTemporaryRedirect,
		req.UserAgent(),
		req.Referer(),
	)
	return err
}

// AccessLogsGet retrieves a list of access logs from the database
func (t *TursoRepo) AccessLogsGet(page, pageSize int) ([]AccessLog, error) {
	rows, err := t.db.Query(`
		SELECT 
			id,
			timestamp,
			remote_addr,
			request_method,
			request_uri,5
			protocol,
			status_code,
			user_agent,
			referer
		FROM access_logs
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?`,
		pageSize, (page-1)*pageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []AccessLog
	for rows.Next() {
		var entry AccessLog
		err := rows.Scan(
			&entry.ID,
			&entry.Timestamp,
			&entry.RemoteAddr,
			&entry.RequestMethod,
			&entry.RequestURI,
			&entry.Protocol,
			&entry.StatusCode,
			&entry.UserAgent,
			&entry.Referer,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, entry)
	}

	return logs, rows.Err()
}

func (t *TursoRepo) AccessLogsCountInTimeSpan(from, to time.Time) (int, error) {
	var count int
	err := t.db.QueryRow(`
		SELECT COUNT(*) 
		FROM access_logs
		WHERE timestamp BETWEEN ? AND ?`,
		from, to,
	).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (t *TursoRepo) AccessLogsCount() (int, error) {
	var count int
	err := t.db.QueryRow(`
		SELECT COUNT(*) 
		FROM access_logs`).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// MIGRATIONS

func (t *TursoRepo) RunMigrations() error {
	// First, ensure schema_migrations table exists
	_, err := t.db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		return err
	}

	migrations := []string{
		// v1: Create initial tables
		`CREATE TABLE IF NOT EXISTS access_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			remote_addr TEXT NOT NULL,
			request_method TEXT NOT NULL,
			request_uri TEXT NOT NULL,
			protocol TEXT NOT NULL,
			status_code INTEGER NOT NULL,
			user_agent TEXT,
			referer TEXT
		);
		CREATE TABLE IF NOT EXISTS event_logs (
			id TEXT PRIMARY KEY NOT NULL,
			timestamp TIMESTAMP NOT NULL,
			short_code TEXT NOT NULL,
			description TEXT NOT NULL,
			severity TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS blog_posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			id_author INTEGER NOT NULL,
			title TEXT NOT NULL,
			date TIMESTAMP NOT NULL,
			short TEXT NOT NULL,
			body TEXT NOT NULL,
			slug TEXT NOT NULL UNIQUE
		);
		CREATE TABLE IF NOT EXISTS authors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			image_url TEXT NOT NULL
		);`,
	}

	// Run each migration in a transaction
	for version, migration := range migrations {
		tx, err := t.db.Begin()
		if err != nil {
			return err
		}

		// Check if migration was already applied
		var exists bool
		err = tx.QueryRow("SELECT 1 FROM schema_migrations WHERE version = ?", version+1).Scan(&exists)
		if err != nil && err != sql.ErrNoRows {
			tx.Rollback()
			return err
		}

		// Skip if already applied
		if exists {
			tx.Rollback()
			continue
		}

		// Apply migration
		if _, err := tx.Exec(migration); err != nil {
			tx.Rollback()
			return err
		}

		// Record migration
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version+1); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

/**
 * Schema_v1
 *
 * CREATE TABLE redirects (
 * 	id INTEGER PRIMARY KEY AUTOINCREMENT,
 * 	count INTEGER DEFAULT 0
 * );
 */
