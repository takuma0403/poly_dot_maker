package converter

import "math"

// rotatePoint は点 (x, y) を中心 (cx, cy) 周りに deg 度回転させた座標を返す。
func rotatePoint(x, y, cx, cy float64, deg int) (float64, float64) {
	rad := float64(deg) * math.Pi / 180.0
	cos := math.Cos(rad)
	sin := math.Sin(rad)
	dx := x - cx
	dy := y - cy
	return cx + dx*cos - dy*sin, cy + dx*sin + dy*cos
}

// pointInTriangle は点 (px, py) が三角形 A(ax,ay) B(bx,by) C(cx,cy) の内側か判定する。
// 境界上の点も内側とみなす。
func pointInTriangle(px, py, ax, ay, bx, by, cx, cy float64) bool {
	d1 := cross(px, py, ax, ay, bx, by)
	d2 := cross(px, py, bx, by, cx, cy)
	d3 := cross(px, py, cx, cy, ax, ay)

	hasNeg := (d1 < 0) || (d2 < 0) || (d3 < 0)
	hasPos := (d1 > 0) || (d2 > 0) || (d3 > 0)
	return !(hasNeg && hasPos)
}

// cross は点 P に対して辺 AB の符号付き面積（の2倍）を返す。
func cross(px, py, ax, ay, bx, by float64) float64 {
	return (px-bx)*(ay-by) - (ax-bx)*(py-by)
}

// clamp は v を [lo, hi] の範囲に収める。
func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// pointInPolygon は ray casting 法で点 (px, py) が多角形の内側か判定する。
func pointInPolygon(px, py float64, pts [][2]float64) bool {
	inside := false
	n := len(pts)
	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		xi, yi := pts[i][0], pts[i][1]
		xj, yj := pts[j][0], pts[j][1]
		intersect := ((yi > py) != (yj > py)) && (px < (xj-xi)*(py-yi)/(yj-yi)+xi)
		if intersect {
			inside = !inside
		}
	}
	return inside
}
