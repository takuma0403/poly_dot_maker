IMAGE_NAME  := poly-dot-maker
RELEASE_TAG := $(IMAGE_NAME):latest

# .env が存在すれば読み込む
-include .env
export

GCP_PROJECT_ID    ?= your-project-id
GCP_REGION        ?= asia-northeast1
CLOUD_RUN_SERVICE ?= poly-dot-maker

IMAGE_REPO := gcr.io/$(GCP_PROJECT_ID)/$(IMAGE_NAME)

.PHONY: init-env dev dev-stop build run deploy check-env clean tidy help

## init-env: .env.example を .env にコピー (初回セットアップ)
init-env:
	@if [ -f .env ]; then \
		echo ".env already exists. Skipping."; \
	else \
		cp .env.example .env; \
		echo ".env created from .env.example"; \
	fi

## dev: ローカル開発 (air ホットリロード)
dev:
	docker compose up --build

## dev-stop: 開発コンテナを停止
dev-stop:
	docker compose down

## build: リリース用イメージをビルド
build:
	docker build -f Dockerfile.release -t $(RELEASE_TAG) .

## run: ビルド済みリリースイメージをローカル実行
run:
	docker run --rm -p 8080:8080 -e PORT=8080 $(RELEASE_TAG)

## deploy: Cloud Run にデプロイ
deploy: check-env
	docker build -f Dockerfile.release -t $(IMAGE_REPO):latest .
	docker push $(IMAGE_REPO):latest
	gcloud run deploy $(CLOUD_RUN_SERVICE) \
		--image $(IMAGE_REPO):latest \
		--platform managed \
		--region $(GCP_REGION) \
		--allow-unauthenticated \
		--project $(GCP_PROJECT_ID)

## check-env: デプロイに必要な環境変数を確認
check-env:
	@[ "$(GCP_PROJECT_ID)" != "your-project-id" ] || \
		(echo "Error: GCP_PROJECT_ID が未設定です。.env を確認してください。" && exit 1)

## tidy: go mod tidy
tidy:
	go mod tidy

## clean: ビルド成果物・コンテナ・イメージを割り当てて削除
clean:
	docker compose down --rmi local --volumes --remove-orphans
	docker rmi -f $(RELEASE_TAG) 2>/dev/null || true
	rm -rf tmp/

## help: ターゲット一覧を表示
help:
	@grep -E '^## ' Makefile | sed 's/## //'
