package cmd

import (
	"fmt"

	"github.com/kiku99/morama/internal/models"
	"github.com/kiku99/morama/internal/storage"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [title]",
	Short: "기록을 삭제합니다",
	Long: `입력한 제목의 영화 또는 드라마 기록을 삭제합니다.

예시:
  morama delete "슬기로운 전공의 생활" --drama
  morama delete "인셉션" --movie`,
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

		deleted, err := store.DeleteByTitleAndType(title, mediaType)
		if err != nil {
			fmt.Printf("❌ 삭제 실패: %v\n", err)
			return
		}

		if deleted == 0 {
			fmt.Printf("⚠️ \"%s\" (%s) 항목을 찾을 수 없습니다.\n", title, mediaType)
		} else {
			fmt.Printf("🗑️ \"%s\" (%s) 항목 %d개를 삭제했습니다.\n", title, mediaType, deleted)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().Bool("movie", false, "영화로 삭제")
	deleteCmd.Flags().Bool("drama", false, "드라마로 삭제")
}
