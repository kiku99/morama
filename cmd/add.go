package cmd

import (
	"fmt"
	"strconv"

	"github.com/kiku99/morama/internal/models"
	"github.com/kiku99/morama/internal/storage"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new movie or drama entry",
	Long:  "Add a new movie or drama entry with interactive rating and comment input",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]

		// Determine media type from flags
		isMovie, _ := cmd.Flags().GetBool("movie")
		isDrama, _ := cmd.Flags().GetBool("drama")

		var mediaType models.MediaType
		if isMovie && isDrama {
			fmt.Println("❌ Error: Cannot specify both --movie and --drama flags")
			return
		} else if isMovie {
			mediaType = models.Movie
		} else if isDrama {
			mediaType = models.Drama
		} else {
			fmt.Println("❌ Error: Must specify either --movie or --drama flag")
			return
		}

		// Interactive rating input
		ratingPrompt := promptui.Prompt{
			Label:    "Rate",
			Validate: validateRating,
		}

		ratingStr, err := ratingPrompt.Run()
		if err != nil {
			fmt.Printf("❌ Error getting rating: %v\n", err)
			return
		}

		rating, _ := strconv.ParseFloat(ratingStr, 64)

		// Interactive comment input
		commentPrompt := promptui.Prompt{
			Label: "한줄평",
		}

		comment, err := commentPrompt.Run()
		if err != nil {
			fmt.Printf("❌ Error getting comment: %v\n", err)
			return
		}

		// Load storage
		store, err := storage.NewStorage()
		if err != nil {
			fmt.Printf("❌ Error loading storage: %v\n", err)
			return
		}
		defer store.Close()

		// Create new entry
		entry := models.MediaEntry{
			Title:   title,
			Type:    mediaType,
			Rating:  rating,
			Comment: comment,
		}

		// Add entry
		if err := store.AddEntry(entry); err != nil {
			fmt.Printf("❌ Error saving entry: %v\n", err)
			return
		}

		fmt.Println("✅ Successfully saved!")
	},
}

func validateRating(input string) error {
	rating, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return fmt.Errorf("invalid number")
	}
	if rating < 0 || rating > 5 {
		return fmt.Errorf("rating must be between 0 and 5")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().Bool("movie", false, "Add as a movie")
	addCmd.Flags().Bool("drama", false, "Add as a drama")
}
