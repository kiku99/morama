package cmd

import (
	"fmt"
	"strconv"

	"github.com/kiku99/morama/internal/models"
	"github.com/kiku99/morama/internal/storage"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [title]",
	Short: "Edit an existing movie or drama entry",
	Long: `Edit an existing movie or drama entry by its ID.
Example:
  morama edit "Drama Title" --id=3 --drama
  morama edit "Movie Title" --id=5 --movie`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		idStr, _ := cmd.Flags().GetString("id")
		isMovie, _ := cmd.Flags().GetBool("movie")
		isDrama, _ := cmd.Flags().GetBool("drama")

		// ID 파싱
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("❌ Error: Invalid ID format")
			return
		}

		// 미디어 타입 확인
		if (isMovie && isDrama) || (!isMovie && !isDrama) {
			fmt.Println("❌ Error: Please specify either --movie or --drama (but not both)")
			return
		}

		var mediaType models.MediaType
		if isMovie {
			mediaType = models.Movie
		} else {
			mediaType = models.Drama
		}

		// 스토리지 초기화
		store, err := storage.NewStorage()
		if err != nil {
			fmt.Printf("❌ Error loading storage: %v\n", err)
			return
		}
		defer store.Close()

		// 기존 항목 조회
		entries, err := store.FindAllByTitleAndType(title, mediaType)
		if err != nil {
			fmt.Printf("❌ Error finding entry: %v\n", err)
			return
		}

		// ID로 항목 찾기
		var targetEntry *models.MediaEntry
		for _, entry := range entries {
			if entry.ID == id {
				targetEntry = &entry
				break
			}
		}

		if targetEntry == nil {
			fmt.Printf("❌ Error: No entry found with ID %d for \"%s\" (%s)\n", id, title, mediaType)
			return
		}

		// Interactive rating input
		ratingPrompt := promptui.Prompt{
			Label:    "Rate",
			Validate: validateRating,
			Default:  fmt.Sprintf("%.1f", targetEntry.Rating),
		}

		ratingStr, err := ratingPrompt.Run()
		if err != nil {
			fmt.Printf("❌ Error getting rating: %v\n", err)
			return
		}

		rating, _ := strconv.ParseFloat(ratingStr, 64)

		// Interactive comment input
		commentPrompt := promptui.Prompt{
			Label:   "한줄평",
			Default: targetEntry.Comment,
		}

		comment, err := commentPrompt.Run()
		if err != nil {
			fmt.Printf("❌ Error getting comment: %v\n", err)
			return
		}

		// Update entry
		updatedEntry := models.MediaEntry{
			Title:       title,
			Type:        mediaType,
			Rating:      rating,
			Comment:     comment,
			DateWatched: targetEntry.DateWatched,
		}

		if err := store.UpdateEntry(id, updatedEntry); err != nil {
			fmt.Printf("❌ Error updating entry: %v\n", err)
			return
		}

		fmt.Println("✅ Successfully updated!")
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().String("id", "", "ID of the entry to edit")
	editCmd.Flags().Bool("movie", false, "Edit as a movie")
	editCmd.Flags().Bool("drama", false, "Edit as a drama")
	editCmd.MarkFlagRequired("id")
}
