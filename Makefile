-include .env

.PHONY: build

build-mac:
	@echo " > Building [license-hwid] for mac..."
	@GOOS=darwin GOARCH=amd64 go build -o ./bin/license_mac
	@echo " > Finished building [license-hwid]"

build-linux:
	@echo " > Building [license-hwid] for linux..."
	@GOOS=linux GOARCH=amd64 go build -o ./bin/license_linux
	@echo " > Finished building [license-hwid]"

build-windows:
	@echo " > Building [license-hwid] for windows..."
	@GOOS=windows GOARCH=amd64 go build -o ./bin/license_win.exe
	@echo " > Finished building [license-hwid]"