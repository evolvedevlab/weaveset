APP_NAME=weaveset

apiserver:
	@go build -o bin/$(APP_NAME)-apiserver ./cmd/apiserver
	@hugo -s site --minify
	@bin/$(APP_NAME)-apiserver

worker:
	@go build -o bin/$(APP_NAME)-worker ./cmd/worker
	@bin/$(APP_NAME)-worker

test:
	@go test -v -count=1 ./...
