IMAGE_NAME  := poly-dot-maker
RELEASE_TAG := $(IMAGE_NAME):latest

.PHONY: dev build run tidy

## dev: ローカル開発 (air ホットリロード)
dev:
	docker compose up --build

## build: リリース用イメージをビルド
build:
	docker build -f Dockerfile.release -t $(RELEASE_TAG) .

## run: ビルド済みリリースイメージをローカル実行
run:
	docker run --rm -p 8080:8080 -e PORT=8080 $(RELEASE_TAG)

## tidy: go mod tidy
tidy:
	go mod tidy

## help: ターゲット一覧を表示
help:
	@grep -E '^## ' Makefile | sed 's/## //'
