LDFLAGS := -s -w
VERSION := $(shell bash version.sh)

all: windows linux-64 linux-86 darwin

windows:
	@echo "Building for windows"
	@GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS) -X 'sooprox/cmd.Version=$(VERSION)'" -o build/api.exe .
linux-64:
	@echo "Building for linux-64"
	@GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS) -X 'sooprox/cmd.Version=$(VERSION)'" -o build/api-amd64-linux .
linux-86:
	@echo "Building for linux-86"
	@GOOS=linux GOARCH=386 go build -ldflags="$(LDFLAGS) -X 'sooprox/cmd.Version=$(VERSION)'" -o build/api-386-linux .
darwin:
	@echo "Building for darwin"
	@GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS) -X 'sooprox/cmd.Version=$(VERSION)'" -o build/api-amd64-darwin .

upx:
	@upx build/*

clean:
	@rm -rf build

test:
	@echo $(VERSION)

run:
	go run main.go


