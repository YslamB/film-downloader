package main

const (
	baseM3U8URL = "https://downloadfilm.belet.me/117663/98678afb-8c9d-440c-b732-604d8002b529.m3u8"
	authHeader  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTQ2ODA2NDksIlVzZXJJRCI6MTMxNjU0NCwiRGV2aWNlTW9kZSI6MX0.jV73vda6eoTAkbna7plpyr6rXV3xSO2em00tsA4YT4g"
	outputDir   = "segments"
)

func main() {
	DownloadHLS()
	// DownloadMp4()
}
