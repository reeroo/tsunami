BINARY_NAME=tsunami
BUILD_FOLDER=build
ifeq ($(OS),Windows_NT)
    EXE_SUFFIX=.exe
    RMDIR=rmdir /S /Q
    MKDIR=mkdir
else
	EXE_SUFFIX=
	RMDIR=rm -rf
	MKDIR=mkdir -p
endif


.PHONY: test

build: test
	go build -o $(BUILD_FOLDER)/$(BINARY_NAME)$(EXE_SUFFIX) ./...

test:
	$(RMDIR) reports
	$(MKDIR) reports
	go test -coverprofile=reports/cover.out ./...
	go tool cover -html=reports/cover.out -o reports/coverage.html

clean:
	$(RMDIR) build reports