package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/kiku99/morama/internal/config"
	"github.com/kiku99/morama/internal/storage"
	"github.com/kiku99/morama/internal/utils"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show statistics about your movie and drama collection",
	Long: `Display comprehensive statistics about your movie and drama collection.
Shows total counts, average ratings, yearly breakdowns, and more.

Examples:
  morama stats`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()
		defer func() {
			utils.LogCommandExecution("stats", args, time.Since(startTime))
		}()

		utils.LogUserAction("stats_requested", "statistics display")

		store, err := storage.NewStorage()
		if err != nil {
			utils.HandleError(
				utils.DatabaseError("Failed to initialize storage", err),
				"Storage initialization error",
			)
		}
		defer store.Close()

		stats, err := store.GetStats()
		if err != nil {
			utils.HandleError(
				utils.DatabaseError("Failed to retrieve statistics", err),
				"Statistics retrieval error",
			)
		}

		config := config.GetConfig()
		fmt.Println("📊 Collection Statistics")
		fmt.Println("=" + strings.Repeat("=", 50))

		// 전체 영화/드라마 개수 출력
		totalMovies := stats["total_movies"].(int)
		totalDramas := stats["total_dramas"].(int)
		totalEntries := totalMovies + totalDramas

		fmt.Printf("📽️  Total Movies: %d\n", totalMovies)
		fmt.Printf("📺  Total Dramas: %d\n", totalDramas)
		fmt.Printf("📚  Total Entries: %d\n\n", totalEntries)

		// 평균 평점 출력
		if totalMovies > 0 {
			avgMovieRating := stats["avg_movie_rating"].(float64)
			fmt.Printf("⭐ Average Movie Rating: %.2f/%.1f\n", avgMovieRating, config.Display.RatingScale)
		}

		if totalDramas > 0 {
			avgDramaRating := stats["avg_drama_rating"].(float64)
			fmt.Printf("⭐ Average Drama Rating: %.2f/%.1f\n", avgDramaRating, config.Display.RatingScale)
		}

		if totalEntries > 0 {
			avgOverallRating := stats["avg_overall_rating"].(float64)
			fmt.Printf("⭐ Overall Average Rating: %.2f/%.1f\n\n", avgOverallRating, config.Display.RatingScale)
		}

		// 별점 분포도 출력
		if ratingDistribution, ok := stats["rating_distribution"].(map[string]int); ok {
			fmt.Println("📈 Rating Distribution:")
			for rating, count := range ratingDistribution {
				if count > 0 {
					percentage := float64(count) / float64(totalEntries) * 100
					fmt.Printf("   %.1f stars: %d entries (%.1f%%)\n",
						parseRating(rating), count, percentage)
				}
			}
			fmt.Println()
		}

		// 연도별 통계 출력
		if yearlyStats, ok := stats["yearly_stats"].(map[string]interface{}); ok {
			fmt.Println("📅 Yearly Breakdown:")
			for yearStr, yearData := range yearlyStats {
				if yearDataMap, ok := yearData.(map[string]interface{}); ok {
					movies := yearDataMap["movies"].(int)
					dramas := yearDataMap["dramas"].(int)
					avgRating := yearDataMap["avg_rating"].(float64)

					fmt.Printf("   %s: %d movies, %d dramas (avg: %.2f)\n",
						yearStr, movies, dramas, avgRating)
				}
			}
		}

		// 마지막 시청일 출력
		if lastWatched, ok := stats["last_watched"].(string); ok && lastWatched != "" {
			fmt.Printf("\n🕒 Last Watched: %s\n", lastWatched)
		}

		utils.LogUserAction("stats_completed", fmt.Sprintf("displayed stats for %d entries", totalEntries))
	},
}

func parseRating(ratingStr string) float64 {
	var rating float64
	fmt.Sscanf(ratingStr, "%f", &rating)
	return rating
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
