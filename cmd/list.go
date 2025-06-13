package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/kiku99/morama/internal/storage"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// 컬럼 비율 및 최소/최대 폭 정의
const (
	minTerminalWidth = 80
	maxTerminalWidth = 200

	// 컬럼 비율 (합계 1.0)
	idRatio      = 0.06 // 6%
	titleRatio   = 0.35 // 35%
	typeRatio    = 0.08 // 8%
	ratingRatio  = 0.10 // 10%
	dateRatio    = 0.16 // 16%
	commentRatio = 0.25 // 25%

	// 최소 폭
	minIdWidth      = 4
	minTitleWidth   = 20
	minTypeWidth    = 6
	minRatingWidth  = 8
	minDateWidth    = 13
	minCommentWidth = 15
)

// 터미널 폭 감지
func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// 감지 실패 시 기본값 반환
		return 120
	}

	// 최소/최대 제한
	if width < minTerminalWidth {
		return minTerminalWidth
	}
	if width > maxTerminalWidth {
		return maxTerminalWidth
	}

	return width
}

// 동적 컬럼 폭 계산
type tableWidths struct {
	id        int
	title     int
	entryType int
	rating    int
	date      int
	comment   int
	total     int
}

func calculateTableWidths() tableWidths {
	termWidth := getTerminalWidth()

	// 테이블 경계선과 구분자를 위한 여백 (┃ 문자들)
	borderSpace := 7 // │ 문자 7개

	// 실제 내용을 위한 폭
	contentWidth := termWidth - borderSpace

	return tableWidths{
		id:        maxInt(int(float64(contentWidth)*idRatio), minIdWidth),
		title:     maxInt(int(float64(contentWidth)*titleRatio), minTitleWidth),
		entryType: maxInt(int(float64(contentWidth)*typeRatio), minTypeWidth),
		rating:    maxInt(int(float64(contentWidth)*ratingRatio), minRatingWidth),
		date:      maxInt(int(float64(contentWidth)*dateRatio), minDateWidth),
		comment:   maxInt(int(float64(contentWidth)*commentRatio), minCommentWidth),
	}
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all movies and dramas",
	Long:  "Display all recorded movies and dramas in a formatted table",
	Run: func(cmd *cobra.Command, args []string) {
		store, err := storage.NewStorage()
		if err != nil {
			fmt.Printf("❌ Error loading storage: %v\n", err)
			return
		}
		defer store.Close()

		// Get all years
		years, err := store.GetYears()
		if err != nil {
			fmt.Printf("❌ Error getting years: %v\n", err)
			return
		}

		if len(years) == 0 {
			fmt.Println("📭 No entries found. Add some movies or dramas with 'morama add'!")
			return
		}

		// 동적 폭 계산
		widths := calculateTableWidths()

		// Display each year group
		for _, year := range years {
			entries, err := store.GetEntriesByYear(year)
			if err != nil {
				fmt.Printf("❌ Error getting entries for year %d: %v\n", year, err)
				continue
			}

			if len(entries) == 0 {
				continue
			}

			fmt.Printf("\n                                                   Watched in %d\n", year)

			// Print table header with calculated widths
			fmt.Printf("┏%s┳%s┳%s┳%s┳%s┳%s┓\n",
				strings.Repeat("━", widths.id),
				strings.Repeat("━", widths.title),
				strings.Repeat("━", widths.entryType),
				strings.Repeat("━", widths.rating),
				strings.Repeat("━", widths.date),
				strings.Repeat("━", widths.comment))

			fmt.Printf("┃%s┃%s┃%s┃%s┃%s┃%s┃\n",
				padStringToWidth("ID", widths.id),
				padStringToWidth("Title", widths.title),
				padStringToWidth("Type", widths.entryType),
				padStringToWidth("Rating", widths.rating),
				padStringToWidth("Date Watched", widths.date),
				padStringToWidth("Comment", widths.comment))

			fmt.Printf("┡%s╇%s╇%s╇%s╇%s╇%s┩\n",
				strings.Repeat("━", widths.id),
				strings.Repeat("━", widths.title),
				strings.Repeat("━", widths.entryType),
				strings.Repeat("━", widths.rating),
				strings.Repeat("━", widths.date),
				strings.Repeat("━", widths.comment))

			for _, entry := range entries {
				id := fmt.Sprintf("%d", entry.ID)
				title := truncateStringWithWidth(entry.Title, widths.title)
				entryType := string(entry.Type)
				rating := fmt.Sprintf("%.1f", entry.Rating)
				dateStr := entry.DateWatched.Format("Jan 02, 2006")
				comment := truncateStringWithWidth(entry.Comment, widths.comment)

				fmt.Printf("│%s│%s│%s│%s│%s│%s│\n",
					padStringToWidth(id, widths.id),
					padStringToWidth(title, widths.title),
					padStringToWidth(entryType, widths.entryType),
					padStringToWidth(rating, widths.rating),
					padStringToWidth(dateStr, widths.date),
					padStringToWidth(comment, widths.comment))
			}

			fmt.Printf("└%s┴%s┴%s┴%s┴%s┴%s┘\n",
				strings.Repeat("─", widths.id),
				strings.Repeat("─", widths.title),
				strings.Repeat("─", widths.entryType),
				strings.Repeat("─", widths.rating),
				strings.Repeat("─", widths.date),
				strings.Repeat("─", widths.comment))
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
