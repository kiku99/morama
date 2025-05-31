package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kiku99/morama/internal/models"
	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func getDBPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".morama", "morama.db"), nil
}

func ensureDataDir() error {
	dbPath, err := getDBPath()
	if err != nil {
		return err
	}
	return os.MkdirAll(filepath.Dir(dbPath), 0755)
}

func NewStorage() (*Storage, error) {
	if err := ensureDataDir(); err != nil {
		return nil, err
	}

	dbPath, err := getDBPath()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	storage := &Storage{db: db}
	if err := storage.initDB(); err != nil {
		db.Close()
		return nil, err
	}

	return storage, nil
}

func (s *Storage) initDB() error {
	schema := `
	CREATE TABLE IF NOT EXISTS media (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		type TEXT CHECK(type IN ('movie', 'drama')) NOT NULL,
		rating REAL CHECK(rating >= 0 AND rating <= 5),
		comment TEXT,
		date_watched DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_media_type ON media(type);
	CREATE INDEX IF NOT EXISTS idx_media_rating ON media(rating);
	CREATE INDEX IF NOT EXISTS idx_media_date_watched ON media(date_watched);
	`

	_, err := s.db.Exec(schema)
	return err
}

func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Storage) AddEntry(entry models.MediaEntry) error {
	query := `
	INSERT INTO media (title, type, rating, comment, date_watched)
	VALUES (?, ?, ?, ?, ?)
	`

	// SQLite 호환 포맷으로 시간 저장
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := s.db.Exec(query, entry.Title, string(entry.Type), entry.Rating, entry.Comment, now)
	return err
}

// parseTime tries multiple time formats
func parseTime(timeStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05-07:00",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}

func (s *Storage) GetAllEntries() ([]models.MediaEntry, error) {
	query := `
	SELECT id, title, type, rating, comment, date_watched, created_at
	FROM media
	ORDER BY date_watched DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.MediaEntry
	for rows.Next() {
		var entry models.MediaEntry
		var typeStr string
		var dateWatchedStr, createdAtStr string

		err := rows.Scan(
			&entry.ID,
			&entry.Title,
			&typeStr,
			&entry.Rating,
			&entry.Comment,
			&dateWatchedStr,
			&createdAtStr,
		)
		if err != nil {
			return nil, err
		}

		entry.Type = models.MediaType(typeStr)

		// 시간 파싱
		if entry.DateWatched, err = parseTime(dateWatchedStr); err != nil {
			return nil, err
		}
		if entry.CreatedAt, err = parseTime(createdAtStr); err != nil {
			// created_at이 파싱 실패해도 계속 진행
			entry.CreatedAt = time.Time{}
		}

		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

func (s *Storage) GetEntriesByYear(year int) ([]models.MediaEntry, error) {
	query := `
	SELECT id, title, type, rating, comment, date_watched, created_at
	FROM media
	WHERE strftime('%Y', date_watched) = ?
	ORDER BY date_watched DESC
	`

	rows, err := s.db.Query(query, fmt.Sprintf("%d", year)) // 정수를 문자열로 변환
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.MediaEntry
	for rows.Next() {
		var entry models.MediaEntry
		var typeStr string
		var dateWatchedStr, createdAtStr string

		err := rows.Scan(
			&entry.ID,
			&entry.Title,
			&typeStr,
			&entry.Rating,
			&entry.Comment,
			&dateWatchedStr,
			&createdAtStr,
		)
		if err != nil {
			return nil, err
		}

		entry.Type = models.MediaType(typeStr)

		// 시간 파싱
		if entry.DateWatched, err = parseTime(dateWatchedStr); err != nil {
			return nil, err
		}
		if entry.CreatedAt, err = parseTime(createdAtStr); err != nil {
			// created_at이 파싱 실패해도 계속 진행
			entry.CreatedAt = time.Time{}
		}

		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

func (s *Storage) GetYears() ([]int, error) {
	query := `
	SELECT DISTINCT strftime('%Y', date_watched) as year
	FROM media
	ORDER BY year DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var years []int
	for rows.Next() {
		var yearStr string // 문자열로 먼저 받아보기
		if err := rows.Scan(&yearStr); err != nil {
			return nil, err
		}

		// 문자열을 정수로 변환
		var year int
		if _, err := fmt.Sscanf(yearStr, "%d", &year); err != nil {
			continue // 파싱 실패 시 건너뛰기
		}

		years = append(years, year)
	}

	return years, rows.Err()
}

func (s *Storage) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total count
	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM media").Scan(&total)
	if err != nil {
		return nil, err
	}
	stats["total"] = total

	// Count by type
	var movieCount, dramaCount int
	err = s.db.QueryRow("SELECT COUNT(*) FROM media WHERE type = 'movie'").Scan(&movieCount)
	if err != nil {
		return nil, err
	}
	err = s.db.QueryRow("SELECT COUNT(*) FROM media WHERE type = 'drama'").Scan(&dramaCount)
	if err != nil {
		return nil, err
	}
	stats["movies"] = movieCount
	stats["dramas"] = dramaCount

	// Average rating
	var avgRating float64
	err = s.db.QueryRow("SELECT AVG(rating) FROM media WHERE rating > 0").Scan(&avgRating)
	if err != nil {
		return nil, err
	}
	stats["avg_rating"] = avgRating

	return stats, nil
}

func (s *Storage) FindByTitleAndType(title string, mediaType models.MediaType) (*models.MediaEntry, error) {
	// media 테이블에서 제목과 유형이 일치하는 첫 번째 레코드를 가져옴
	query := `
	SELECT id, title, type, rating, comment, date_watched, created_at
	FROM media
	WHERE title = ? AND type = ?
	LIMIT 1
	`

	row := s.db.QueryRow(query, title, string(mediaType))

	var entry models.MediaEntry
	var typeStr, watchedStr, createdStr string
	err := row.Scan(
		&entry.ID, &entry.Title, &typeStr, &entry.Rating, &entry.Comment, &watchedStr, &createdStr,
	)
	if err != nil {
		return nil, fmt.Errorf("entry not found for \"%s\" (%s)", title, mediaType)
	}

	entry.Type = models.MediaType(typeStr)
	entry.DateWatched, _ = parseTime(watchedStr)
	entry.CreatedAt, _ = parseTime(createdStr)

	return &entry, nil
}
