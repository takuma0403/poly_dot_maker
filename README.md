# poly_dot_maker

画像をアップロードすると、埋め尽くし可能な様々な形状のタイルで埋め尽くしたドット絵に変換する API サーバー。

**🌐 デモ:** <https://poly-dot-maker-962005752553.asia-northeast1.run.app/>

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
├── static/
│   └── index.html       # 変換 UI（Web フロントエンド）
├── docs/
│   └── architecture.md  # アーキテクチャ設計ドキュメント
└── src/
    ├── main.go          # エントリポイント
    └── handler/
        └── convert.go   # POST /convert ハンドラー
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
| `GET` | `/` | 変換 UI ページ（index.html）を返す |
| `GET` | `/health` | ヘルスチェック |
| `POST` | `/convert` | 画像を変換して PNG を返す |

### POST /convert

`multipart/form-data` でリクエストします。

**リクエストフィールド**

| フィールド | 型 | 必須 | デフォルト | 説明 |
|-----------|-----|------|----------|------|
| `image` | ファイル | ✅ | — | 変換元画像（JPEG / PNG） |
| `shape` | 文字列 | | `triangle` | 埋め尽くす図形 (`triangle`, `hexagon`, `square`) |
| `dots` | 整数 | | `3000` | 点の総数（1以上） |
| `colors` | 整数 | | `16` | 減色後の色数（5〜30） |
| `rotate` | 整数 | | `0` | 回転角度（15の倍数、単位: 度） |

**レスポンス**

| 条件 | ステータス | Content-Type | ボディ |
|------|----------|------------|------|
| 成功 | `200 OK` | `image/png` | 変換後の PNG 画像 |
| バリデーションエラー | `400 Bad Request` | `application/json` | エラーメッセージ |
| サーバーエラー | `500 Internal Server Error` | `application/json` | エラーメッセージ |

**curl 例**

```bash
curl -X POST https://poly-dot-maker-962005752553.asia-northeast1.run.app/convert \
  -F "image=@photo.jpg" \
  -F "shape=hexagon" \
  -F "dots=5000" \
  -F "colors=20" \
  -F "rotate=0" \
  --output result.png
```

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
