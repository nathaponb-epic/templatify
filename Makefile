BINARY=templatify
BINARY_WINDOWS=templatify.exe

build_linux:
	@echo Building binary...
	set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${BINARY} .
	@echo Done!

build_mac:
	@echo Building binary...
	set GOOS=darwin&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${BINARY} .
	@echo Done!

build_windows:
	@echo Building binary...
	set GOOS=windows&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${BINARY_WINDOWS} .
	@echo Done!