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
	ORDER BY id DESC
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
	ORDER BY id DESC
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

func (s *Storage) FindAllByTitleAndType(title string, mediaType models.MediaType) ([]models.MediaEntry, error) {
	query := `
	SELECT id, title, type, rating, comment, date_watched, created_at
	FROM media
	WHERE title = ? AND type = ?
	ORDER BY id DESC
	`

	rows, err := s.db.Query(query, title, string(mediaType))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.MediaEntry
	for rows.Next() {
		var entry models.MediaEntry
		var typeStr, watchedStr, createdStr string

		if err := rows.Scan(
			&entry.ID, &entry.Title, &typeStr, &entry.Rating, &entry.Comment, &watchedStr, &createdStr,
		); err != nil {
			return nil, err
		}

		entry.Type = models.MediaType(typeStr)
		entry.DateWatched, _ = parseTime(watchedStr)
		entry.CreatedAt, _ = parseTime(createdStr)

		entries = append(entries, entry)
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("entry not found for \"%s\" (%s)", title, mediaType)
	}

	return entries, nil
}

// 업데이트: ID 기반
func (s *Storage) UpdateEntry(id int, entry models.MediaEntry) error {
	query := `
	UPDATE media 
	SET title = ?, type = ?, rating = ?, comment = ?, date_watched = ?
	WHERE id = ?
	`

	now := time.Now().Format("2006-01-02 15:04:05")
	result, err := s.db.Exec(query, entry.Title, string(entry.Type), entry.Rating, entry.Comment, now, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no entry found with ID %d", id)
	}

	return nil
}

func (s *Storage) DeleteByID(id int) (int64, error) {
	query := `DELETE FROM media WHERE id = ?`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Storage) DeleteAll() (int64, error) {
	query := `DELETE FROM media`
	result, err := s.db.Exec(query)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Storage) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total counts
	var totalMovies, totalDramas int
	err := s.db.QueryRow("SELECT COUNT(*) FROM media WHERE type = 'movie'").Scan(&totalMovies)
	if err != nil {
		return nil, err
	}
	err = s.db.QueryRow("SELECT COUNT(*) FROM media WHERE type = 'drama'").Scan(&totalDramas)
	if err != nil {
		return nil, err
	}
	stats["total_movies"] = totalMovies
	stats["total_dramas"] = totalDramas

	// Average ratings
	var avgMovieRating, avgDramaRating, avgOverallRating float64
	s.db.QueryRow("SELECT AVG(rating) FROM media WHERE type = 'movie' AND rating > 0").Scan(&avgMovieRating)
	s.db.QueryRow("SELECT AVG(rating) FROM media WHERE type = 'drama' AND rating > 0").Scan(&avgDramaRating)
	s.db.QueryRow("SELECT AVG(rating) FROM media WHERE rating > 0").Scan(&avgOverallRating)

	stats["avg_movie_rating"] = avgMovieRating
	stats["avg_drama_rating"] = avgDramaRating
	stats["avg_overall_rating"] = avgOverallRating

	// Rating distribution
	ratingDistribution := make(map[string]int)
	rows, err := s.db.Query(`
		SELECT 
			CASE 
				WHEN rating >= 4.5 THEN '4.5'
				WHEN rating >= 4.0 THEN '4.0'
				WHEN rating >= 3.5 THEN '3.5'
				WHEN rating >= 3.0 THEN '3.0'
				WHEN rating >= 2.5 THEN '2.5'
				WHEN rating >= 2.0 THEN '2.0'
				WHEN rating >= 1.5 THEN '1.5'
				WHEN rating >= 1.0 THEN '1.0'
				ELSE '0.5'
			END as rating_range,
			COUNT(*) as count
		FROM media
		GROUP BY rating_range
		ORDER BY rating_range DESC
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ratingRange string
			var count int
			if err := rows.Scan(&ratingRange, &count); err == nil {
				ratingDistribution[ratingRange] = count
			}
		}
	}
	stats["rating_distribution"] = ratingDistribution

	// Yearly breakdown
	yearlyStats := make(map[string]interface{})
	yearRows, err := s.db.Query(`
		SELECT 
			strftime('%Y', date_watched) as year,
			COUNT(CASE WHEN type = 'movie' THEN 1 END) as movies,
			COUNT(CASE WHEN type = 'drama' THEN 1 END) as dramas,
			AVG(rating) as avg_rating
		FROM media
		GROUP BY year
		ORDER BY year DESC
	`)
	if err == nil {
		defer yearRows.Close()
		for yearRows.Next() {
			var year string
			var movies, dramas int
			var avgRating float64
			if err := yearRows.Scan(&year, &movies, &dramas, &avgRating); err == nil {
				yearlyStats[year] = map[string]interface{}{
					"movies":     movies,
					"dramas":     dramas,
					"avg_rating": avgRating,
				}
			}
		}
	}
	stats["yearly_stats"] = yearlyStats

	// Last watched
	var lastWatched string
	s.db.QueryRow("SELECT date_watched FROM media ORDER BY date_watched DESC LIMIT 1").Scan(&lastWatched)
	stats["last_watched"] = lastWatched

	return stats, nil
}
