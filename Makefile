
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
