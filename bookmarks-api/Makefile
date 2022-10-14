build:
	go build -o ./bin/bookmarks ./cmd/bookmarks/bookmarks.go

run:
	./bin/bookmarks -config ./config/config.development.yaml

run_reload:
	CompileDaemon -build="go build -o ./bin/bookmarks ./cmd/bookmarks/bookmarks.go" -command="./bin/bookmarks -config ./config/config.development.yaml" -color
