ifndef GOOS
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	GOOS := darwin
else ifeq ($(UNAME_S),Linux)
	GOOS := linux
else
#$(error "$$GOOS is not defined. If you are using Windows, try to re-make using 'GOOS=windows make ...' ")
	GOOS := windows
endif
endif

BUILD_FLAGS := -ldflags "-X version.GitCommit=`git rev-parse HEAD`"

NOSOGOD_BINARY64 := nosogod-$(GOOS)_amd64
NOSOGOCLI_BINARY64 := nosogocli-$(GOOS)_amd64

VERSION := $(shell grep -E "Version[ ]*=" version/version.go | tr -d "\" " | cut -d "=" -f 2)

NOSOGOD_RELEASE64 := nosogod-$(VERSION)-$(GOOS)_amd64
NOSOGOCLI_RELEASE64 := nosogocli-$(VERSION)-$(GOOS)_amd64

NOSOGO_RELEASE64 := nosogo-$(VERSION)-$(GOOS)_amd64

all: test target release-all install
#all: target release-all install

nosogod:
	@echo "Building nosogod to bin/nosogod"
	@mkdir -p bin
	@go build $(BUILD_FLAGS) -o bin/nosogod cmd/nosogod/main.go

nosogocli:
	@echo "Building nosogocli to bin/nosogocli"
	@mkdir -p bin
	@go build $(BUILD_FLAGS) -o bin/nosogocli cmd/nosogocli/main.go

install:
	@echo "Installing nosogod to $(GOPATH)/bin"
	@go install ./bin/nosogod
	@go install ./bin/nosogocli

target:
	mkdir -p $@

binary: target/$(NOSOGOD_BINARY64) target/$(NOSOGOCLI_BINARY64)

ifeq ($(GOOS),windows)
release: binary
	cd target && cp -f $(NOSOGOD_BINARY64) $(NOSOGOD_BINARY64).exe
	cd target && cp -f $(NOSOGOCLI_BINARY64) $(NOSOGOCLI_BINARY64).exe
	cd target && md5sum $(NOSOGOD_BINARY64).exe $(NOSOGOCLI_BINARY64).exe  >$(NOSOGO_RELEASE64).md5
	cd target && zip $(NOSOGO_RELEASE64).zip $(NOSOGOD_BINARY64).exe $(NOSOGOCLI_BINARY64).exe $(NOSOGO_RELEASE64).md5
	cd target && rm -f $(NOSOGOD_BINARY64) $(NOSOGOD_BINARY64).exe $(NOSOGOCLI_BINARY64) $(NOSOGOCLI_BINARY64).exe $(NOSOGO_RELEASE64).md5
else
release: binary
	cd target && md5sum $(NOSOGOD_BINARY64) $(NOSOGOCLI_BINARY64) >$(NOSOGO_RELEASE64).md5
	cd target && tar -czf $(NOSOGO_RELEASE64).tgz $(NOSOGOD_BINARY64) $(NOSOGOCLI_BINARY64) $(NOSOGO_RELEASE64).md5
	cd target && rm -f $(NOSOGOD_BINARY64) $(NOSOGOCLI_BINARY64) $(NOSOGO_RELEASE64).md5
endif

release-all: clean
	GOOS=darwin  make release
	GOOS=linux   make release
	GOOS=windows make release

clean:
	@echo "Cleaning binaries built..."
	@rm -f bin/nosogod bin/nosogocli
	@rm -rf target
	@rm -f $(GOPATH)/bin/nosogod
	@rm -f $(GOPATH)/bin/nosogocli
	@echo "Done."

target/$(NOSOGOD_BINARY64):
	CGO_ENABLED=0 GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ cmd/nosogod/main.go

target/$(NOSOGOCLI_BINARY64):
	CGO_ENABLED=0 GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ cmd/nosogocli/main.go

test:
	@echo "====> Running tests"
	@go test -v -run "." "github.com/Friends-Of-Noso/NosoGo/tests"

benchmark:
	@echo "====> Running benchmarks"
	@go test -benchmem -run="^$$" -bench "." "github.com/Friends-Of-Noso/NosoGo/tests"

.PHONY: all target release-all clean
