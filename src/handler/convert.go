package handler

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/t_takumaru/poly_dot_maker/src/converter"
)

const (
	defaultDots   = 3000
	defaultColors = 16
	defaultRotate = 0

	minColors = 5
	maxColors = 30
)

// Convert は POST /convert のハンドラー。
func Convert(c echo.Context) error {
	// --- 画像ファイル取得 ---
	file, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "imageフィールドが必要です")
	}
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ファイルのオープンに失敗しました")
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "画像のデコードに失敗しました。JPEG/PNG のみ対応しています")
	}

	// --- パラメータ取得 ---
	shape := c.FormValue("shape")
	if shape == "" {
		shape = "triangle"
	}
	dots := parseIntParam(c, "dots", defaultDots)
	colors := parseIntParam(c, "colors", defaultColors)
	rotateDeg := parseIntParam(c, "rotate", defaultRotate)

	// --- バリデーション ---
	if colors < minColors || colors > maxColors {
		return echo.NewHTTPError(http.StatusBadRequest, "colors は5〜30の範囲で指定してください")
	}
	if dots < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "dots は1以上で指定してください")
	}
	if rotateDeg%15 != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "rotate は15の倍数で指定してください")
	}

	// --- 変換 ---
	var result image.Image

	switch shape {
	case "hexagon":
		result, err = converter.ConvertHexagon(img, dots, colors, rotateDeg)
	case "square":
		result, err = converter.ConvertSquare(img, dots, colors, rotateDeg)
	case "triangle":
		fallthrough
	default:
		result, err = converter.ConvertTriangle(img, dots, colors, rotateDeg)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "変換に失敗しました: "+err.Error())
	}

	// --- PNG で返却 ---
	var buf bytes.Buffer
	if err := png.Encode(&buf, result); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "エンコードに失敗しました")
	}

	return c.Blob(http.StatusOK, "image/png", buf.Bytes())
}

// parseIntParam はフォームパラメータを int で取得し、失敗時はデフォルト値を返す。
func parseIntParam(c echo.Context, key string, defaultVal int) int {
	v := c.FormValue(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return n
}

// init で image/jpeg を登録（インポートされないと Decode が機能しない）
func init() {
	_ = jpeg.Decode
}
