package cmd

import (
	"fmt"

	"github.com/kiku99/morama/internal/storage"
	"github.com/spf13/cobra"
)

var (
	deleteID  int
	deleteAll bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "기록을 삭제합니다",
	Long: `기록을 삭제합니다. --id로 단일 항목을 삭제하거나 --all로 전체 기록을 삭제할 수 있습니다.

예시:
  morama delete --id=3     # ID 3번 항목 삭제
  morama delete --all      # 전체 기록 삭제`,
	Run: func(cmd *cobra.Command, args []string) {
		if deleteID == 0 && !deleteAll {
			fmt.Println("❌ 삭제하려면 --id 또는 --all 중 하나를 지정하세요.")
			return
		}
		if deleteID > 0 && deleteAll {
			fmt.Println("❌ --id와 --all은 동시에 사용할 수 없습니다.")
			return
		}

		store, err := storage.NewStorage()
		if err != nil {
			fmt.Printf("❌ 데이터베이스 열기 실패: %v\n", err)
			return
		}
		defer store.Close()

		if deleteAll {
			count, err := store.DeleteAll()
			if err != nil {
				fmt.Printf("❌ 전체 삭제 실패: %v\n", err)
				return
			}
			fmt.Printf("🗑️ 모든 기록 %d개를 삭제했습니다.\n", count)
			return
		}

		deleted, err := store.DeleteByID(deleteID)
		if err != nil {
			fmt.Printf("❌ 삭제 실패: %v\n", err)
			return
		}
		if deleted == 0 {
			fmt.Printf("⚠️ ID %d 항목을 찾을 수 없습니다.\n", deleteID)
		} else {
			fmt.Printf("🗑️ ID %d 항목을 삭제했습니다.\n", deleteID)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().IntVar(&deleteID, "id", 0, "삭제할 항목 ID")
	deleteCmd.Flags().BoolVar(&deleteAll, "all", false, "전체 기록 삭제")
}
