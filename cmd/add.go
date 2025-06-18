package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kiku99/morama/internal/config"
	"github.com/kiku99/morama/internal/models"
	"github.com/kiku99/morama/internal/storage"
	"github.com/kiku99/morama/internal/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new movie or drama entry",
	Long:  "Add a new movie or drama entry with interactive rating and comment input",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()
		defer func() {
			utils.LogCommandExecution("add", args, time.Since(startTime))
		}()

		title := args[0]
		utils.LogUserAction("add_entry", fmt.Sprintf("title: %s", title))

		// Determine media type from flags
		isMovie, _ := cmd.Flags().GetBool("movie")
		isDrama, _ := cmd.Flags().GetBool("drama")

		var mediaType models.MediaType
		if isMovie && isDrama {
			utils.HandleError(
				utils.ValidationError("Cannot specify both --movie and --drama flags", nil),
				"Invalid media type specification",
			)
		} else if isMovie {
			mediaType = models.Movie
		} else if isDrama {
			mediaType = models.Drama
		} else {
			utils.HandleError(
				utils.ValidationError("Must specify either --movie or --drama flag", nil),
				"Missing media type specification",
			)
		}

		// Interactive rating input
		ratingPrompt := promptui.Prompt{
			Label:    "Rate",
			Validate: validateRating,
		}

		ratingStr, err := ratingPrompt.Run()
		if err != nil {
			utils.HandleError(
				utils.UserInputError("Failed to get rating input", err),
				"Rating input error",
			)
		}

		rating, _ := strconv.ParseFloat(ratingStr, 64)

		// Interactive comment input
		commentPrompt := promptui.Prompt{
			Label: "한줄평",
		}

		comment, err := commentPrompt.Run()
		if err != nil {
			utils.HandleError(
				utils.UserInputError("Failed to get comment input", err),
				"Comment input error",
			)
		}

		// Load storage
		store, err := storage.NewStorage()
		if err != nil {
			utils.HandleError(
				utils.DatabaseError("Failed to initialize storage", err),
				"Storage initialization error",
			)
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
			utils.HandleError(
				utils.DatabaseError("Failed to save entry", err),
				"Entry save error",
			)
		}

		utils.LogUserAction("entry_added", fmt.Sprintf("title: %s, type: %s, rating: %.1f", title, mediaType, rating))
		fmt.Println("✅ Successfully saved!")
	},
}

func validateRating(input string) error {
	rating, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return fmt.Errorf("invalid number")
	}

	cfg := config.GetConfig()
	if rating < 0 || rating > cfg.Display.RatingScale {
		return fmt.Errorf("rating must be between 0 and %.1f", cfg.Display.RatingScale)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().Bool("movie", false, "Add as a movie")
	addCmd.Flags().Bool("drama", false, "Add as a drama")
}
