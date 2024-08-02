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

#all: test target release-all install
all: target release-all install

core-stratumd:
	@echo "Building nosogo to cmd/nosogo/nosogod"
	@go build $(BUILD_FLAGS) -o cmd/nosogo/nosogod cmd/nosogo/main.go

install:
	@echo "Installing nossogo to $(GOPATH)/bin"
	@go install ./cmd/nosogo

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
	@rm -rf cmd/nosogo/nosogod
	@rm -rf target
	@rm -rf $(GOPATH)/bin/nosogod
	@echo "Cleaning temp test data..."
	@echo "Done."

target/$(NOSOGOD_BINARY64):
	CGO_ENABLED=0 GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ cmd/nosogo/main.go

#test:
#	@echo "====> Running go test"
#	@go test -tags "network" $(PACKAGES)

#benchmark:
#	@go test -bench $(PACKAGES)

#functional-tests:
#	@go test -timeout=5m -tags="functional" ./test
#
#ci: test functional-tests

.PHONY: all target release-all clean