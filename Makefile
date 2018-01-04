COVERALLS_TOKEN := "IVQwNa8dypGgtaLmBkFSoBcRcCl0tlqui"
GITHUB_API_TOKEN := ""

.PHONY: build restore test cover vet lint

all: restore test build

build:
	env CGO_ENABLED=0 go build -o submerger .
	chmod +x submerger
	env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o submerger-nas .
	chmod +x submerger-nas

restore:
	go get -u github.com/golang/dep/cmd/dep
	${GOPATH}/bin/dep ensure

test:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure
	go test -cover `go list ./... | grep -v /vendor/`

cover:
	go test -cover `go list ./... | grep -v /vendor/`

lint:
	golint `go list ./... | grep -v /vendor/`

vet:
	go vet `go list ./... | grep -v /vendor/`

cover-remote:
	go get -u github.com/golang/lint/golint
	go get -u github.com/mattn/goveralls
	go get -u github.com/golang/dep/cmd/dep
	dep ensure
	go test -covermode=count -coverprofile=cover.tmp
	goveralls -service travis-ci -coverprofile cover.tmp
