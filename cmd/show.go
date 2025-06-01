package cmd

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/kiku99/morama/internal/models"
	"github.com/kiku99/morama/internal/storage"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [title]",
	Short: "선택한 영화 또는 드라마의 상세 정보를 출력합니다",
	Long: `입력한 제목의 영화 또는 드라마 기록을 상세히 보여줍니다.
예시:
  morama show "슬기로울 전공의 생활" --drama
  morama show "인셉션" --movie`,

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		isMovie, _ := cmd.Flags().GetBool("movie")
		isDrama, _ := cmd.Flags().GetBool("drama")

		if (isMovie && isDrama) || (!isMovie && !isDrama) {
			fmt.Println("❌ Please specify either --movie or --drama (but not both)")
			return
		}

		var mediaType models.MediaType
		if isMovie {
			mediaType = models.Movie
		} else {
			mediaType = models.Drama
		}

		store, err := storage.NewStorage()
		if err != nil {
			fmt.Printf("❌ 데이터베이스 열기 실패: %v\n", err)
			return
		}
		defer store.Close()

		entries, err := store.FindAllByTitleAndType(title, mediaType)
		if err != nil {
			fmt.Printf("❌ 검색 중 오류 발생: %v\n", err)
			return
		}

		if len(entries) == 0 {
			fmt.Printf("‼️ \"%s\" (%s)에 해당하는 항목이 없습니다.\n", title, mediaType)
			return
		}

		for i, entry := range entries {
			if len(entries) > 1 {
				fmt.Printf("\n📄 결과 %d/%d\n", i+1, len(entries))
			}
			printEntryBox(&entry)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().Bool("movie", false, "영화로 조회")
	showCmd.Flags().Bool("drama", false, "드라마로 조회")
}

func printEntryBox(entry *models.MediaEntry) {
	line := strings.Repeat("━", 60)
	labelWidth := 6
	c := cases.Title(language.Und)

	fmt.Println(line)
	fmt.Println(formatField("📌 제목", entry.Title, labelWidth))
	fmt.Println(formatField("🎞️ 유형", c.String(string(entry.Type)), labelWidth))
	fmt.Println(formatField("⭐ 평점", fmt.Sprintf("%.1f / 5.0", entry.Rating), labelWidth))
	fmt.Println(formatField("🗓️ 시청일", entry.DateWatched.Format("2006-01-02"), labelWidth))
	fmt.Println(formatField("💬 한줄평", entry.Comment, labelWidth))
	fmt.Println(line)
}

func formatField(label string, value string, labelWidth int) string {
	labelPadded := padStringToWidth(label, labelWidth)
	return fmt.Sprintf("%s : %s", labelPadded, value)
}
