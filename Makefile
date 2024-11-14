ifndef GOOS
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	GOOS := darwin
else ifeq ($(UNAME_S),Linux)
	GOOS := linux
else
$(error "$$GOOS is not defined. If you are using Windows, try to re-make using 'GOOS=windows make ...' ")
endif
endif

BUILD_FLAGS := -ldflags "-X nosogod/version.GitCommit=`git rev-parse HEAD`"

NOSOGOD_BINARY64 := nosogod-$(GOOS)_amd64

VERSION := $(shell awk -F= '/Version =/ {print $$2}' version/version.go | tr -d "\" ")

NOSOGOD_RELEASE64 := nosogod-$(VERSION)-$(GOOS)_amd64

NOSOGO_RELEASE64 := nosogo-$(VERSION)-$(GOOS)_amd64

all: test target release-all install
#all: target release-all install

nosogod:
	@echo "Building nosogod to bin/nosogod"
	@mkdir -p bin
	@go build $(BUILD_FLAGS) -o bin/nosogod cmd/nosogod/main.go

install:
	@echo "Installing nossogod to $(GOPATH)/bin"
	@go install ./bin/nosogod

target:
	mkdir -p $@

binary: target/$(NOSOGOD_BINARY64)

ifeq ($(GOOS),windows)
release: binary
	cd target && cp -f $(NOSOGOD_BINARY64) $(NOSOGOD_BINARY64).exe
	cd target && md5sum $(NOSOGOD_BINARY64).exe  >$(NOSOGO_RELEASE64).md5
	cd target && zip $(NOSOGO_RELEASE64).zip $(NOSOGOD_BINARY64).exe $(NOSOGO_RELEASE64).md5
	cd target && rm -f $(NOSOGOD_BINARY64) $(NOSOGOD_BINARY64).exe $(NOSOGO_RELEASE64).md5
else
release: binary
	cd target && md5sum $(NOSOGOD_BINARY64) >$(NOSOGO_RELEASE64).md5
	cd target && tar -czf $(NOSOGO_RELEASE64).tgz $(NOSOGOD_BINARY64) $(NOSOGO_RELEASE64).md5
	cd target && rm -f $(NOSOGOD_BINARY64) $(NOSOGO_RELEASE64).md5
endif

release-all: clean
#	GOOS=darwin  make release
	GOOS=linux   make release
	GOOS=windows make release

clean:
	@echo "Cleaning binaries built..."
	@rm -f bin/nosogod
	@rm -rf target
	@rm -f $(GOPATH)/bin/nosogod
#	@echo "Cleaning temp test data..."
	@echo "Done."

target/$(NOSOGOD_BINARY64):
	CGO_ENABLED=0 GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ cmd/nosogod/main.go

test:
	@echo "====> Running go test"
	@go test ./tests

#benchmark:
#	@go test -bench ./tests

#functional-tests:
#	@go test -timeout=5m -tags="functional" ./test
#
#ci: test functional-tests

.PHONY: all target release-all clean
