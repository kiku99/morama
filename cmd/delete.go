package cmd

import (
	"fmt"

	"github.com/kiku99/morama/internal/models"
	"github.com/kiku99/morama/internal/storage"
	"github.com/kiku99/morama/internal/utils"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [id or title]",
	Short: "기록을 삭제합니다",
	Long: `기록을 삭제합니다. ID 기반 단일 삭제 또는 --all로 전체 삭제가 가능합니다.

예시:
  morama delete 3 --drama          # ID 3번 드라마 삭제
  morama delete --all --movie      # 모든 영화 기록 삭제
  morama delete "인셉션" --movie     # 제목 기반 삭제`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		store, err := storage.NewStorage()
		if err != nil {
			fmt.Printf("❌ 데이터베이스 열기 실패: %v\n", err)
			return
		}
		defer store.Close()

		isMovie, _ := cmd.Flags().GetBool("movie")
		isDrama, _ := cmd.Flags().GetBool("drama")
		isAll, _ := cmd.Flags().GetBool("all")

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

		// 전체 삭제
		if isAll {
			count, err := store.DeleteAllByType(mediaType)
			if err != nil {
				fmt.Printf("❌ 전체 삭제 실패: %v\n", err)
				return
			}
			fmt.Printf("🗑️ %s 기록 %d개를 모두 삭제했습니다.\n", mediaType, count)
			return
		}

		if len(args) == 0 {
			fmt.Println("❌ 삭제할 ID 또는 제목을 입력하거나 --all 옵션을 사용하세요.")
			return
		}

		// ID인지 문자열인지 판단
		if id, err := utils.ParseID(args[0]); err == nil {
			deleted, err := store.DeleteByIDAndType(id, mediaType)
			if err != nil {
				fmt.Printf("❌ 삭제 실패: %v\n", err)
				return
			}
			if deleted == 0 {
				fmt.Printf("⚠️ ID %d (%s) 항목을 찾을 수 없습니다.\n", id, mediaType)
			} else {
				fmt.Printf("🗑️ ID %d (%s) 항목을 삭제했습니다.\n", id, mediaType)
			}
		} else {
			// 제목 기반 삭제
			title := args[0]
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
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().Bool("movie", false, "영화로 삭제")
	deleteCmd.Flags().Bool("drama", false, "드라마로 삭제")
	deleteCmd.Flags().Bool("all", false, "해당 타입의 모든 기록 삭제")
}
