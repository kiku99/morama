package cmd

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/kiku99/morama/internal/models"
	"github.com/kiku99/morama/internal/storage"
	"github.com/kiku99/morama/internal/utils"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [title]",
	Short: "Display detailed information about a selected movie or drama",
	Long: `Shows detailed information about the movie or drama with the given title.
예시:
  morama show "언젠가는 슬기로울 전공의생활" --drama
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
				fmt.Printf("\n📄 Result %d/%d\n", i+1, len(entries))
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
	fmt.Println(formatField("📌 Title", entry.Title, labelWidth))
	fmt.Println(formatField("🎞️ Type", c.String(string(entry.Type)), labelWidth))
	fmt.Println(formatField("⭐ Rating", fmt.Sprintf("%.1f / 5.0", entry.Rating), labelWidth))
	fmt.Println(formatField("🗓️ Watched Date", entry.DateWatched.Format("2006-01-02"), labelWidth))
	fmt.Println(formatField("💬 Comment", entry.Comment, labelWidth))
	fmt.Println(line)
}

func formatField(label string, value string, labelWidth int) string {
	labelPadded := utils.PadStringToWidth(label, labelWidth)
	return fmt.Sprintf("%s : %s", labelPadded, value)
}
