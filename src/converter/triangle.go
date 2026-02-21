package converter

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// Triangle は正三角形の3頂点をfloat64で保持する。
type Triangle struct {
	ax, ay float64
	bx, by float64
	cx, cy float64
}

// ConvertTriangle は画像を正三角形ドット絵に変換して返す。
func ConvertTriangle(src image.Image, dots, colors, rotateDeg int) (*image.RGBA, error) {
	bounds := src.Bounds()
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y

	// --- パレット抽出 ---
	palette := ExtractPalette(src, colors)

	// --- 辺長・グリッドパラメータ ---
	side, step, triH := computeParams(w, h, dots)

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

	riMax := int(math.Ceil((yMax-yMin)/triH)) + 1
	ciMax := int(math.Ceil((xMax-xMin)/step))*2 + 2

	for ri := -1; ri <= riMax; ri++ {
		for ci := -1; ci <= ciMax; ci++ {
			tri := buildTriangle(ri, ci, side, step, triH, xMin, yMin)
			// 回転適用
			if rotateDeg != 0 {
				tri = rotateTri(tri, cx, cy, rotateDeg)
			}
			renderTriangle(src, dst, tri, palette, w, h)
		}
	}

	return dst, nil
}

// computeParams は dots 数から辺長・step・高さを算出する。
func computeParams(w, h, dots int) (side, step, triH float64) {
	areaPer := float64(w*h) / float64(dots)
	side = math.Sqrt(4.0 * areaPer / math.Sqrt(3))
	step = side / 2.0
	triH = side * math.Sqrt(3) / 2.0
	return
}

// buildTriangle は (ri, ci) に対応する三角形の頂点を構築する。
func buildTriangle(ri, ci int, side, step, triH, xMin, yMin float64) Triangle {
	x0 := xMin + float64(ci)*step
	y0 := yMin + float64(ri)*triH

	if (ri+ci)%2 == 0 {
		// 上向き ▲: 底辺が下、頂点が上
		return Triangle{
			ax: x0, ay: y0 + triH,
			bx: x0 + side, by: y0 + triH,
			cx: x0 + step, cy: y0,
		}
	}
	// 下向き ▽: 底辺が上、頂点が下
	return Triangle{
		ax: x0, ay: y0,
		bx: x0 + side, by: y0,
		cx: x0 + step, cy: y0 + triH,
	}
}

// rotateTri は三角形の全頂点を (cx, cy) 周りに deg 度回転させる。
func rotateTri(t Triangle, cx, cy float64, deg int) Triangle {
	ax, ay := rotatePoint(t.ax, t.ay, cx, cy, deg)
	bx, by := rotatePoint(t.bx, t.by, cx, cy, deg)
	cxr, cyr := rotatePoint(t.cx, t.cy, cx, cy, deg)
	return Triangle{ax, ay, bx, by, cxr, cyr}
}

// renderTriangle は1つの三角形を平均色→パレットスナップで描画する。
func renderTriangle(src image.Image, dst *image.RGBA, tri Triangle, palette []color.RGBA, imgW, imgH int) {
	// バウンディングボックス
	minX := math.Min(tri.ax, math.Min(tri.bx, tri.cx))
	minY := math.Min(tri.ay, math.Min(tri.by, tri.cy))
	maxX := math.Max(tri.ax, math.Max(tri.bx, tri.cx))
	maxY := math.Max(tri.ay, math.Max(tri.by, tri.cy))

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
			if pointInTriangle(float64(px)+0.5, float64(py)+0.5,
				tri.ax, tri.ay, tri.bx, tri.by, tri.cx, tri.cy) {
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
