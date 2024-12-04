package repo

import (
	"database/sql"
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

// COUNTER
func (t *TursoRepo) GetRedirectCount() (int, error) {
	var count int
	err := t.db.QueryRow("SELECT count FROM redirects_count WHERE id = 1").Scan(&count)
	return count, err
}

func (t *TursoRepo) IncrementRedirectCount() error {
	result, err := t.db.Exec("UPDATE redirects_count SET count = count + 1 WHERE id = 1")
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
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
		// v1: Create access_logs table for tracking request details
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
