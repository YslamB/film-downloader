package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	outputDir   string
	accessToken string
	filmID      string
	seasonID    string
	episodeID   string
)

func init() {
	err := godotenv.Load("./.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	outputDir = os.Getenv("OUTPUT_DIR")
	accessToken = os.Getenv("ACCESS_TOKEN")
	filmID = os.Getenv("FILM_ID")
	seasonID = os.Getenv("SEASON_ID")
	episodeID = os.Getenv("EPISODE_ID")
}

func main() {

	var baseM3U8URLs []string
	var err error

	if episodeID == "" && seasonID == "" && filmID != "" {
		baseM3U8URLs[0], err = GetFilmSourceURL(filmID)

		if err != nil {
			log.Fatal(err)
		}
	}

	if seasonID != "" {
		baseM3U8URLs, err = GetEpisodesWithSeasonID(seasonID, episodeID)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("basem3u8 url: ", baseM3U8URLs)

	for i := range baseM3U8URLs {
		fmt.Println("✅Received Source files...")
		DownloadMp4(baseM3U8URLs[i], fmt.Sprintf("%d_%s", i, outputDir))
	}

	fmt.Println("👮‍♀️ everithing is ok 🎯")
}
