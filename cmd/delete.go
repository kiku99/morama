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
	Short: "Delete a record",
	Long: `Deletes a record. You can delete a single entry using --id, or all entries using --all.

Examples:
  morama delete --id=3     # Delete entry with ID 3
  morama delete --all      # Delete all entries`,
	Run: func(cmd *cobra.Command, args []string) {
		if deleteID == 0 && !deleteAll {
			fmt.Println("❌ Please specify either --id or --all to delete entries.")
			return
		}
		if deleteID > 0 && deleteAll {
			fmt.Println("❌ --id and --all cannot be used together.")
			return
		}

		store, err := storage.NewStorage()
		if err != nil {
			fmt.Printf("❌ Failed to open the database: %v\n", err)
			return
		}
		defer store.Close()

		if deleteAll {
			count, err := store.DeleteAll()
			if err != nil {
				fmt.Printf("❌ Failed to delete all entries: %v\n", err)
				return
			}
			fmt.Printf("🗑️ Deleted all %d entries.\n", count)
			return
		}

		deleted, err := store.DeleteByID(deleteID)
		if err != nil {
			fmt.Printf("❌ Failed to delete entry: %v\n", err)
			return
		}
		if deleted == 0 {
			fmt.Printf("⚠️ No entry found with ID %d.\n", deleteID)
		} else {
			fmt.Printf("🗑️ Deleted entry with ID %d.\n", deleteID)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().IntVar(&deleteID, "id", 0, "ID of the entry to delete")
	deleteCmd.Flags().BoolVar(&deleteAll, "all", false, "Delete all entries")
}
