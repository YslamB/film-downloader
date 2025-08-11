package main

import (
	myhash "film-downloader/hash"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	accessToken string
	filmID      string
	seasonID    string
	episodeID   string
	quality     string
	secureKey   string
	envHash     string
	expiresAt   int64
)

func init() {
	err := godotenv.Load("./.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	accessToken = os.Getenv("ACCESS_TOKEN")
	filmID = os.Getenv("FILM_ID")
	seasonID = os.Getenv("SEASON_ID")
	episodeID = os.Getenv("EPISODE_ID")
	quality = os.Getenv("QUALITY")
	envHash = os.Getenv("HASH")
	expiresAt, err = strconv.ParseInt(os.Getenv("EXPIRES_AT"), 10, 64)

	if err != nil {
		log.Fatalf("EXPIRES_AT must be integer: %e", err)
	}

	secureKey = "w3r1Sec4re_Token_"
}

type Movie struct {
	Name   string
	Source string
}

func main() {
	// generate hash
	// thirtyDays := 30 * 24 * time.Hour
	// hash, expiresAt := myhash.GenerateExpirableHash(thirtyDays, secureKey, "Assa")
	// fmt.Println(hash, expiresAt)

	// solve hash
	if !myhash.VerifyHash(envHash, secureKey, "Assa", expiresAt) {
		log.Fatalf("Subscription is expiret")
	}

	fmt.Println("🥳 👏👏 Subscription is active")

	var movies []Movie
	var err error

	if episodeID == "" && seasonID == "" && filmID != "" {
		source, err := GetFilmSourceURL(filmID)

		if err != nil {
			log.Fatal(err)
		}

		movies = append(movies, source)
	}

	if seasonID != "" {
		movies, err = GetEpisodesWithSeasonID(seasonID, episodeID)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("✅ Received Source files...", movies)

	for i := range movies {
		DownloadMp4(movies[i], time.Now().Format("2006-01-02"))
	}

	fmt.Println("👮‍♀️ everithing is ok 🎯")
}
