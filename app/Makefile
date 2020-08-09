# This how we want to name the binary output
BINARY=../bin/PROJECTNAME

# These are the values we want to pass for VERSION and BUILD
VERSION=`git describe --abbrev=6 --always --tag`

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X main.version=${VERSION}"
GOFLAGS=-a -tags PROJECTNAME -installsuffix PROJECTNAME -mod=vendor

build: bin

bin:
	@echo "  >  Building binary \"PROJECTNAME\" $(VERSION)..."
	@go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY)

bin-linux:
	@echo "  >  Building linux-amd64 binary \"PROJECTNAME\" $(VERSION)..."
	@GOOS=linux GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY)-linux

bin-windows:
	@echo "  >  Building windows-amd64 binary \"PROJECTNAME\" $(VERSION)..."
	@GOOS=windows GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY).exe

bin-cross-platform: bin-linux bin-windows

# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	if [ -f ${BINARY} ] ; then rm ${BINARY}-linux ; fi
	if [ -f ${BINARY} ] ; then rm ${BINARY}.exe ; fi

.PHONY: build bin bin-linux bin-windows
