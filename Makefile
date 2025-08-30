
build-linux:
	@echo "Started building..."
	@GOOS=linux GOARCH=amd64 go build -o ./bin/downloader .
	@echo "Building done"

build-mac:
	@echo "Started building..."
	@GOOS=darwin GOARCH=arm64 go build -o ./bin/downloader .
	@echo "Building done"

build-windows:
	@echo "Started building for Windows..."
	@GOOS=windows GOARCH=amd64 go build -o ./bin/downloader.exe .
	@echo "Building done"

build-windows7:
	@echo "Started building for Windows 7..."
	@GOOS=windows GOARCH=amd64 go build -o ./bin/downloader_win7.exe .
	@echo "Building done for Windows 7"

deploy:
	@echo "Started building..."
	@GOOS=linux GOARCH=amd64 go build -o ./bin/film_downloader_tmp ./cmd/main.go
	@echo "Building done"

	@echo "Deploying..."
	@scp ./bin/film_downloader_tmp root@95.85.126.202:/var/www/film_downloader
	@ssh root@95.85.126.202 "rm -f /var/www/film_downloader/film_downloader && mv /var/www/film_downloader/film_downloader_tmp /var/www/film_downloader/film_downloader"
	@scp ./.env root@95.85.126.202:/var/www/film_downloader
	@echo "Restarting remote service..."
	@ssh root@95.85.126.202 "sudo -S systemctl restart film_downloader.service"
	@echo "Done"


