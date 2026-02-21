# poly_dot_maker

画像をアップロードすると、埋め尽くし可能な様々な形状のタイルで埋め尽くしたドット絵に変換する API サーバー。

## 技術スタック

| 項目 | 内容 |
|------|------|
| 言語 | Go 1.25 |
| フレームワーク | Echo |
| 開発環境 | Docker + air（ホットリロード） |
| デプロイ先 | Google Cloud Run |

## ディレクトリ構成

```
poly_dot_maker/
├── Makefile
├── Dockerfile.develop   # 開発用（air ホットリロード）
├── Dockerfile.release   # 本番用（マルチステージビルド）
├── docker-compose.yml
├── .air.toml
├── .env.example
├── docs/
│   └── architecture.md  # アーキテクチャ設計ドキュメント
└── src/
    └── main.go          # エントリポイント
```

## 環境構築

### 前提条件

- Docker / Docker Compose
- Go 1.25+（ローカル開発時）
- gcloud CLI（デプロイ時）

### セットアップ

```bash
# 1. リポジトリをクローン
git clone https://github.com/takuma0403/poly_dot_maker.git
cd poly_dot_maker

# 2. .env を作成
make init-env

# 3. .env を編集（デプロイする場合は GCP_PROJECT_ID 等を設定）
vi .env
```

### 環境変数

`.env.example` をコピーして `.env` を作成し、必要に応じて値を編集します。

| 変数名 | 説明 | デフォルト |
|--------|------|-----------|
| `PORT` | サーバーのポート番号 | `8080` |
| `GCP_PROJECT_ID` | GCP プロジェクト ID | `your-project-id` |
| `GCP_REGION` | Cloud Run のリージョン | `asia-northeast1` |
| `CLOUD_RUN_SERVICE` | Cloud Run サービス名 | `poly-dot-maker` |

## Makefile コマンド一覧

```bash
make help        # コマンド一覧を表示
make init-env    # .env.example を .env にコピー（初回セットアップ）
make dev         # ローカル開発サーバー起動（air ホットリロード）
make build       # リリース用 Docker イメージをビルド
make run         # ビルド済みリリースイメージをローカル実行
make deploy      # Google Cloud Run にデプロイ
make tidy        # go mod tidy
```

## API エンドポイント

| メソッド | パス | 説明 |
|---------|------|------|
| `GET` | `/health` | ヘルスチェック |
| `GET` | `/hello` | Hello World |

## Cloud Run デプロイ手順

```bash
# 1. gcloud にログイン
gcloud auth login
gcloud auth configure-docker

# 2. 必要な API を有効化
gcloud services enable run.googleapis.com \
  containerregistry.googleapis.com \
  artifactregistry.googleapis.com

# 3. .env に GCP_PROJECT_ID を設定した上でデプロイ
make deploy
```

デプロイ後のサービス URL 確認：

```bash
gcloud run services describe poly-dot-maker \
  --region=asia-northeast1 \
  --format="value(status.url)"
```

## ドキュメント

- [アーキテクチャ設計 (将来実装の converter 等)](docs/architecture.md)
