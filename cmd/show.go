package cmd

import (
	"fmt"
	"strings"

	"github.com/kiku99/morama/internal/models"
	"github.com/kiku99/morama/internal/storage"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [title]",
	Short: "Show detailed info of a movie or drama",
	Long: `Show the detailed record of a movie or drama entry.
Example:
  morama show "슬기로운 전공의 생활" --drama`,
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
			fmt.Printf("❌ Error opening database: %v\n", err)
			return
		}
		defer store.Close()

		entry, err := store.FindByTitleAndType(title, mediaType)
		if err != nil {
			fmt.Printf("❌ %v\n", err)
			return
		}

		fmt.Println(strings.Repeat("━", 60))
		fmt.Printf("📌 제목        : %s\n", entry.Title)
		fmt.Printf("🎞️ 유형        : %s\n", strings.Title(string(entry.Type)))
		fmt.Printf("⭐ 평점        : %.1f / 5.0\n", entry.Rating)
		fmt.Printf("🗓️ 시청일      : %s\n", entry.DateWatched.Format("2006-01-02"))
		fmt.Printf("💬 한줄평      : %s\n", entry.Comment)
		fmt.Println(strings.Repeat("━", 60))
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().Bool("movie", false, "Show movie entry")
	showCmd.Flags().Bool("drama", false, "Show drama entry")
}
