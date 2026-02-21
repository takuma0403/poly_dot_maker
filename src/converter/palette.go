package converter

import (
	"image"
	"image/color"
	"math"
	"math/rand"
)

// color3 は RGB を float64 で保持する軽量型。
type color3 struct {
	r, g, b float64
}

// ExtractPalette は画像から n 色のパレットを抽出する。
// 色相ビン（12分割）ごとのピクセル比率に応じて色数を按分し、
// 各ビン内で k-means クラスタリングを行う。
func ExtractPalette(img image.Image, n int) []color.RGBA {
	bounds := img.Bounds()
	w, h := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y

	const numBins = 12

	// --- 全ピクセルを取得して色相ビンに振り分け ---
	bins := make([][]color3, numBins)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rf := float64(r>>8) / 255.0
			gf := float64(g>>8) / 255.0
			bf := float64(b>>8) / 255.0
			hue := rgbToHue(rf, gf, bf)
			idx := int(hue/30.0) % numBins
			bins[idx] = append(bins[idx], color3{rf, gf, bf})
		}
	}

	total := w * h

	// --- 各ビンへの色数按分 ---
	counts := countDistribution(bins, total, n, numBins)

	// --- 各ビンで k-means を実行してパレットを作成 ---
	var palette []color.RGBA
	for i, k := range counts {
		if k == 0 || len(bins[i]) == 0 {
			continue
		}
		if k > len(bins[i]) {
			k = len(bins[i])
		}
		centers := kmeans(bins[i], k)
		for _, c := range centers {
			palette = append(palette, color.RGBA{
				R: uint8(math.Round(c.r * 255)),
				G: uint8(math.Round(c.g * 255)),
				B: uint8(math.Round(c.b * 255)),
				A: 255,
			})
		}
	}

	return palette
}

// countDistribution は各ビンに割り当てる色数を按分する。
// 小数点以下の残余は比率の降順で1色ずつ追加する。
func countDistribution(bins [][]color3, total, n, numBins int) []int {
	counts := make([]int, numBins)
	fracs := make([]float64, numBins)
	sum := 0

	for i, bin := range bins {
		if len(bin) == 0 {
			continue
		}
		exact := float64(n) * float64(len(bin)) / float64(total)
		counts[i] = int(exact)
		fracs[i] = exact - float64(counts[i])
		sum += counts[i]
	}

	// 残余を比率降順で配分
	rem := n - sum
	type binFrac struct {
		idx  int
		frac float64
	}
	var bf []binFrac
	for i, f := range fracs {
		if len(bins[i]) > 0 {
			bf = append(bf, binFrac{i, f})
		}
	}
	// 降順ソート（バブル）
	for i := 0; i < len(bf)-1; i++ {
		for j := 0; j < len(bf)-1-i; j++ {
			if bf[j].frac < bf[j+1].frac {
				bf[j], bf[j+1] = bf[j+1], bf[j]
			}
		}
	}
	for i := 0; i < rem && i < len(bf); i++ {
		counts[bf[i].idx]++
	}
	return counts
}

// rgbToHue は RGB(0-1) から Hue(0-360) を返す。無彩色は 0 を返す。
func rgbToHue(r, g, b float64) float64 {
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	delta := max - min
	if delta == 0 {
		return 0
	}
	var h float64
	switch max {
	case r:
		h = (g - b) / delta
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/delta + 2
	case b:
		h = (r-g)/delta + 4
	}
	return h * 60
}

// NearestColor はパレットから最もユークリッド距離の近い色を返す。
func NearestColor(c color.RGBA, palette []color.RGBA) color.RGBA {
	best := palette[0]
	bestDist := colorDistSq(c, palette[0])
	for _, p := range palette[1:] {
		d := colorDistSq(c, p)
		if d < bestDist {
			bestDist = d
			best = p
		}
	}
	return best
}

func colorDistSq(a, b color.RGBA) int64 {
	dr := int64(a.R) - int64(b.R)
	dg := int64(a.G) - int64(b.G)
	db := int64(a.B) - int64(b.B)
	return dr*dr + dg*dg + db*db
}

// kmeans は RGB 点列に対して k-means クラスタリングを行い、重心の slice を返す。
func kmeans(points []color3, k int) []color3 {
	if len(points) <= k {
		return points
	}

	// k-means++ 初期化
	centers := kmeanspp(points, k)

	labels := make([]int, len(points))
	for iter := 0; iter < 100; iter++ {
		changed := false
		// ラベル割り当て
		for i, p := range points {
			best, bestD := 0, dist3Sq(p, centers[0])
			for j := 1; j < k; j++ {
				d := dist3Sq(p, centers[j])
				if d < bestD {
					bestD = d
					best = j
				}
			}
			if labels[i] != best {
				labels[i] = best
				changed = true
			}
		}
		if !changed {
			break
		}
		// 重心更新
		sums := make([]color3, k)
		cnts := make([]int, k)
		for i, p := range points {
			l := labels[i]
			sums[l].r += p.r
			sums[l].g += p.g
			sums[l].b += p.b
			cnts[l]++
		}
		for j := 0; j < k; j++ {
			if cnts[j] > 0 {
				centers[j] = color3{sums[j].r / float64(cnts[j]), sums[j].g / float64(cnts[j]), sums[j].b / float64(cnts[j])}
			}
		}
	}
	return centers
}

// kmeanspp は k-means++ によって初期クラスタ中心を選択する。
func kmeanspp(points []color3, k int) []color3 {
	centers := make([]color3, 0, k)
	// 最初の中心はランダム
	centers = append(centers, points[rand.Intn(len(points))])

	for len(centers) < k {
		// 各点から最近傍中心への距離の二乗
		dists := make([]float64, len(points))
		total := 0.0
		for i, p := range points {
			minD := math.Inf(1)
			for _, c := range centers {
				d := dist3Sq(p, c)
				if d < minD {
					minD = d
				}
			}
			dists[i] = minD
			total += minD
		}
		// 距離の二乗に比例した確率でサンプリング
		r := rand.Float64() * total
		cumul := 0.0
		chosen := points[len(points)-1]
		for i, d := range dists {
			cumul += d
			if cumul >= r {
				chosen = points[i]
				break
			}
		}
		centers = append(centers, chosen)
	}
	return centers
}

func dist3Sq(a, b color3) float64 {
	dr := a.r - b.r
	dg := a.g - b.g
	db := a.b - b.b
	return dr*dr + dg*dg + db*db
}
