package models

const (
	MovieType   = 1
	EpisodeType = 2
)

type Episode struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Duration string   `json:"duration"`
	FileID   string   `json:"file_id"`
	Sources  []Source `json:"sources"`
}

type EpisodeResponse struct {
	Episodes []Episode `json:"episodes"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"token"`
}

type Source struct {
	MasterFile  string `json:"master_file"`
	Quality     string `json:"quality"`
	Type        string `json:"type"`
	Main        bool   `json:"main"`
	DownloadURL string `json:"download_url"`
}

type Movie struct {
	ID      int
	Name    string
	Sources []Source
	Type    int
}

type MovieResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Film    Film   `json:"film"`
}

type Film struct {
	ID              int             `json:"id"`
	Year            int             `json:"year"`
	TypeID          int             `json:"type_id"`
	CategoryID      int             `json:"category_id"`
	Name            string          `json:"name"`
	Age             string          `json:"age"`
	Duration        string          `json:"duration"`
	Language        string          `json:"language"`
	Description     string          `json:"description"`
	Genres          []string        `json:"genres"`
	Countries       []Country       `json:"countries"`
	Actors          []Person        `json:"actors"`
	Directors       []Person        `json:"directors"`
	Seasons         []Season        `json:"seasons"`
	Trailers        any             `json:"trailers"`
	Studios         []Studio        `json:"studios"`
	Thumbnails      ImageSet        `json:"thumbnails"`
	Images          Images          `json:"images"`
	LastEpisodeInfo LastEpisodeInfo `json:"last_episode_info"`
	MediaInfo       MediaInfo       `json:"media_info"`
	ParentID        int             `json:"parent_id"`
	RatingKP        float64         `json:"rating_kp"`
	RatingIMDB      float64         `json:"rating_imdb"`
	WatchTime       float64         `json:"watch_time"`
	Like            bool            `json:"like"`
	Dislike         bool            `json:"dislike"`
	Favorites       bool            `json:"favorites"`
	ForKids         bool            `json:"for_kids"`
}

type Country struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Person struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type ImageSet struct {
	Default  ImageSize `json:"default"`
	High     ImageSize `json:"high"`
	Medium   ImageSize `json:"medium"`
	Standard ImageSize `json:"standard"`
}

type ImageSize struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Images struct {
	Vertical              ImageSet  `json:"vertical"`
	VerticalWithoutName   ImageSet  `json:"vertical_without_name"`
	HorizontalWithoutName ImageSet  `json:"horizontal_without_name"`
	HorizontalWithName    ImageSet  `json:"horizontal_with_name"`
	Name                  ImageSize `json:"name"`
}

type LastEpisodeInfo struct {
	EpisodeWatchTime float64 `json:"episode_watch_time"`
	EpisodeID        int     `json:"episode_id"`
	SeasonID         int     `json:"season_id"`
	EpisodeName      string  `json:"episode_name"`
	EpisodeDuration  string  `json:"episode_duration"`
	EpisodePosition  *int    `json:"episode_position"`
	SeasonPosition   *int    `json:"season_position"`
}

type MediaInfo struct {
	AudioInfo    []AudioInfo    `json:"audio_info"`
	SubtitleInfo []SubtitleInfo `json:"subtitle_info"`
}

type Season struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	BBID int    `json:"bb_id"`
}

type Studio struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AudioInfo struct {
	StudioName  string `json:"studio_name"`
	Language    string `json:"language"`
	Abbreviated string `json:"abbreviated"`
	VoiceType   int    `json:"voice_type"`
}

type SubtitleInfo struct {
	Language string `json:"language"`
}
