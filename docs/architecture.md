# poly_dot_maker アーキテクチャ設計

> このドキュメントはプロジェクト全体のアーキテクチャ設計計画です。
> Hello World フェーズ完了後、このドキュメントを参照して converter 等の実装を進める。

## 技術スタック

- **言語**: Go 1.23
- **フレームワーク**: Echo
- **開発環境**: air (ホットリロード) + Docker
- **デプロイ先**: Google Cloud Run (無料枠)

---

## プロジェクト構成（全体像）

```
poly_dot_maker/
├── Makefile
├── Dockerfile.develop        # air によるホットリロード開発用
├── Dockerfile.release        # マルチステージビルド・本番用
├── docker-compose.yml        # ローカル開発用
├── .air.toml                 # air 設定
├── .dockerignore
├── .gitignore
├── .env.example
├── docs/
│   └── architecture.md       # このファイル
└── src/
    ├── main.go               # エントリポイント・Echoサーバー起動
    ├── handler/
    │   └── convert.go        # POST /convert ハンドラー (将来実装)
    └── converter/            # テッセレーション変換ロジック (将来実装)
        ├── converter.go      # 変換インターフェース
        ├── square.go
        ├── triangle.go
        └── hexagon.go
```

---

## Makefile ターゲット

| ターゲット | 内容 |
|-----------|------|
| `make dev` | Dockerfile.develop でコンテナ起動 (air ホットリロード) |
| `make build` | Dockerfile.release でイメージビルド |
| `make run` | ビルド済みリリースイメージをローカル実行 |
| `make deploy` | `gcloud run deploy` でCloud Runにデプロイ |
| `make tidy` | go mod tidy |

---

## API 仕様（将来実装）

```
POST /convert
Content-Type: multipart/form-data

Fields:
  image  : 画像ファイル (JPEG/PNG)
  shape  : "square" | "triangle" | "hexagon"  (default: "square")
  scale  : 整数 (タイルサイズ px、default: 10)

Response:
  200 OK
  Content-Type: image/png
  Body: 変換後の PNG 画像バイナリ

GET /health
  200 OK  {"status":"ok"}
```

---

## converter 設計（将来実装）

```go
type Converter interface {
    Convert(src image.Image, scale int) (image.Image, error)
}

func New(shape string) (Converter, error)
```

| shape | アルゴリズム |
|-------|-------------|
| `square` | グリッドを `scale×scale` の正方形で分割し、各セルの平均色で塗りつぶし |
| `triangle` | 行ごとに▲▽を交互配置。各三角形の重心周辺ピクセルの平均色を採用 |
| `hexagon` | オフセット座標系（奇数行を右方向にシフト）で六角形を配置 |

描画には標準ライブラリ `image`・`image/draw` + `fogleman/gg` を使用予定。

---

## 環境変数

| 変数名 | 説明 | デフォルト |
|--------|------|-----------|
| `PORT` | サーバーが Listen するポート | `8080` |
