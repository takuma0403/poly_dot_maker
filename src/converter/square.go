package converter

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// Square は正方形の4頂点をfloat64で保持する。
type Square struct {
	pts [4][2]float64
}

// ConvertSquare は画像を正方形ドット絵に変換して返す。
func ConvertSquare(src image.Image, dots, colors, rotateDeg int) (*image.RGBA, error) {
	bounds := src.Bounds()
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y

	// --- パレット抽出 ---
	palette := ExtractPalette(src, colors)

	// --- 辺長・グリッドパラメータ ---
	side := computeSquareParams(w, h, dots)

	// --- 出力画像（白で初期化）---
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, dst.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// --- グリッド生成・レンダリング ---
	cx := float64(w) / 2.0
	cy := float64(h) / 2.0

	// 回転がある場合は対角線長分だけグリッドを広げる
	diag := math.Sqrt(float64(w*w + h*h))
	var (
		xMin, yMin float64
		xMax, yMax float64
	)
	if rotateDeg != 0 {
		xMin = cx - diag
		yMin = cy - diag
		xMax = cx + diag
		yMax = cy + diag
	} else {
		xMin, yMin = 0, 0
		xMax, yMax = float64(w), float64(h)
	}

	riMax := int(math.Ceil((yMax-yMin)/side)) + 1
	ciMax := int(math.Ceil((xMax-xMin)/side)) + 1

	for ri := 0; ri <= riMax; ri++ {
		for ci := 0; ci <= ciMax; ci++ {
			sq := buildSquare(ri, ci, side, xMin, yMin)
			// 回転適用
			if rotateDeg != 0 {
				sq = rotateSquare(sq, cx, cy, rotateDeg)
			}
			renderSquare(src, dst, sq, palette, w, h)
		}
	}

	return dst, nil
}

// computeSquareParams は dots 数から辺長を算出する。
func computeSquareParams(w, h, dots int) (side float64) {
	areaPer := float64(w*h) / float64(dots)
	// 正方形の面積 = side^2
	side = math.Sqrt(areaPer)
	return
}

// buildSquare は (ri, ci) に対応する正方形の頂点を構築する。
func buildSquare(ri, ci int, side, xMin, yMin float64) Square {
	x0 := xMin + float64(ci)*side
	y0 := yMin + float64(ri)*side

	return Square{
		pts: [4][2]float64{
			{x0, y0},
			{x0 + side, y0},
			{x0 + side, y0 + side},
			{x0, y0 + side},
		},
	}
}

// rotateSquare は正方形の全頂点を (cx, cy) 周りに deg 度回転させる。
func rotateSquare(s Square, cx, cy float64, deg int) Square {
	var rot Square
	for i := 0; i < 4; i++ {
		nx, ny := rotatePoint(s.pts[i][0], s.pts[i][1], cx, cy, deg)
		rot.pts[i][0] = nx
		rot.pts[i][1] = ny
	}
	return rot
}

// renderSquare は1つの正方形を平均色→パレットスナップで描画する。
func renderSquare(src image.Image, dst *image.RGBA, sq Square, palette []color.RGBA, imgW, imgH int) {
	// バウンディングボックス
	minX, minY := sq.pts[0][0], sq.pts[0][1]
	maxX, maxY := minX, minY
	for i := 1; i < 4; i++ {
		minX = math.Min(minX, sq.pts[i][0])
		minY = math.Min(minY, sq.pts[i][1])
		maxX = math.Max(maxX, sq.pts[i][0])
		maxY = math.Max(maxY, sq.pts[i][1])
	}

	x0 := clamp(int(math.Floor(minX)), 0, imgW-1)
	y0 := clamp(int(math.Floor(minY)), 0, imgH-1)
	x1 := clamp(int(math.Ceil(maxX)), 0, imgW-1)
	y1 := clamp(int(math.Ceil(maxY)), 0, imgH-1)

	// 内側ピクセルの収集と平均色の算出
	var sumR, sumG, sumB int64
	count := 0
	var inside [][2]int

	for py := y0; py <= y1; py++ {
		for px := x0; px <= x1; px++ {
			if pointInPolygon(float64(px)+0.5, float64(py)+0.5, sq.pts[:]) {
				r, g, b, _ := src.At(px, py).RGBA()
				sumR += int64(r >> 8)
				sumG += int64(g >> 8)
				sumB += int64(b >> 8)
				count++
				inside = append(inside, [2]int{px, py})
			}
		}
	}

	if count == 0 {
		return
	}

	// 平均色 → 最近傍パレット色
	avg := color.RGBA{
		R: uint8(sumR / int64(count)),
		G: uint8(sumG / int64(count)),
		B: uint8(sumB / int64(count)),
		A: 255,
	}
	snapped := NearestColor(avg, palette)

	// 描画
	for _, px := range inside {
		dst.SetRGBA(px[0], px[1], snapped)
	}
}
