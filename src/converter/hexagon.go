package converter

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// Hexagon は正六角形の6頂点をfloat64で保持する。
type Hexagon struct {
	pts [6][2]float64
}

// ConvertHexagon は画像を正六角形ドット絵に変換して返す。
func ConvertHexagon(src image.Image, dots, colors, rotateDeg int) (*image.RGBA, error) {
	bounds := src.Bounds()
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y

	// --- パレット抽出 ---
	palette := ExtractPalette(src, colors)

	// --- 辺長・グリッドパラメータ ---
	side, wStep, hStep := computeHexParams(w, h, dots)

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

	riMax := int(math.Ceil((yMax-yMin)/hStep)) + 1
	ciMax := int(math.Ceil((xMax-xMin)/wStep)) + 1

	for ri := -1; ri <= riMax; ri++ {
		for ci := -1; ci <= ciMax; ci++ {
			hex := buildHexagon(ri, ci, side, wStep, hStep, xMin, yMin)
			// 回転適用
			if rotateDeg != 0 {
				hex = rotateHex(hex, cx, cy, rotateDeg)
			}
			renderHexagon(src, dst, hex, palette, w, h)
		}
	}

	return dst, nil
}

// computeHexParams は dots 数から辺長・水平step・垂直stepを算出する。
func computeHexParams(w, h, dots int) (side, wStep, hStep float64) {
	areaPer := float64(w*h) / float64(dots)
	// 正六角形の面積 = 3 * sqrt(3) / 2 * side^2
	side = math.Sqrt(2.0 * areaPer / (3.0 * math.Sqrt(3)))
	// pointy-topped (上が尖っている)
	wStep = math.Sqrt(3) * side
	hStep = 1.5 * side
	return
}

// buildHexagon は (ri, ci) に対応する六角形の頂点を構築する。
func buildHexagon(ri, ci int, side, wStep, hStep, xMin, yMin float64) Hexagon {
	centerX := xMin + float64(ci)*wStep
	if ri%2 != 0 {
		centerX += wStep / 2.0
	}
	centerY := yMin + float64(ri)*hStep

	var pts [6][2]float64
	for i := 0; i < 6; i++ {
		// pointy-topped の頂点は 30度回転した位置から始まる (90, 30, -30, -90, -150, 150)
		angleDeg := 60.0*float64(i) - 30.0
		angleRad := math.Pi / 180.0 * angleDeg
		pts[i][0] = centerX + side*math.Cos(angleRad)
		pts[i][1] = centerY + side*math.Sin(angleRad)
	}
	return Hexagon{pts: pts}
}

// rotateHex は六角形の全頂点を (cx, cy) 周りに deg 度回転させる。
func rotateHex(h Hexagon, cx, cy float64, deg int) Hexagon {
	var rot Hexagon
	for i := 0; i < 6; i++ {
		nx, ny := rotatePoint(h.pts[i][0], h.pts[i][1], cx, cy, deg)
		rot.pts[i][0] = nx
		rot.pts[i][1] = ny
	}
	return rot
}

// renderHexagon は1つの六角形を平均色→パレットスナップで描画する。
func renderHexagon(src image.Image, dst *image.RGBA, hex Hexagon, palette []color.RGBA, imgW, imgH int) {
	// バウンディングボックス
	minX, minY := hex.pts[0][0], hex.pts[0][1]
	maxX, maxY := minX, minY
	for i := 1; i < 6; i++ {
		minX = math.Min(minX, hex.pts[i][0])
		minY = math.Min(minY, hex.pts[i][1])
		maxX = math.Max(maxX, hex.pts[i][0])
		maxY = math.Max(maxY, hex.pts[i][1])
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
			if pointInPolygon(float64(px)+0.5, float64(py)+0.5, hex.pts[:]) {
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
