package mdi3

// HCT (hue, chroma, tone) color space: CAM16 hue/chroma + CIE L* tone.
// Port of the parts of Google's material-color-utilities (Apache 2.0) that
// M3 scheme generation needs: CAM16 under default viewing conditions,
// HctSolver, Blend.harmonize, Contrast, TemperatureCache.complement and
// DislikeAnalyzer.

import (
	"image/color"
	"math"
)

type hct struct{ h, c, t float64 }

type vec [3]float64
type mat [3]vec

func (m mat) mul(v vec) vec {
	return vec{
		m[0][0]*v[0] + m[0][1]*v[1] + m[0][2]*v[2],
		m[1][0]*v[0] + m[1][1]*v[1] + m[1][2]*v[2],
		m[2][0]*v[0] + m[2][1]*v[1] + m[2][2]*v[2],
	}
}

func (m mat) inverse() mat {
	a, b, c := m[0], m[1], m[2]
	det := a[0]*(b[1]*c[2]-b[2]*c[1]) - a[1]*(b[0]*c[2]-b[2]*c[0]) + a[2]*(b[0]*c[1]-b[1]*c[0])
	return mat{
		{(b[1]*c[2] - b[2]*c[1]) / det, (a[2]*c[1] - a[1]*c[2]) / det, (a[1]*b[2] - a[2]*b[1]) / det},
		{(b[2]*c[0] - b[0]*c[2]) / det, (a[0]*c[2] - a[2]*c[0]) / det, (a[2]*b[0] - a[0]*b[2]) / det},
		{(b[0]*c[1] - b[1]*c[0]) / det, (a[1]*c[0] - a[0]*c[1]) / det, (a[0]*b[1] - a[1]*b[0]) / det},
	}
}

var (
	xyzFromLinrgb = mat{
		{0.41233895, 0.35762064, 0.18051042},
		{0.2126, 0.7152, 0.0722},
		{0.01932141, 0.11916382, 0.95034478},
	}
	whiteD65 = vec{95.047, 100, 108.883}
	m16      = mat{
		{0.401288, 0.650173, -0.051461},
		{-0.250268, 1.204414, 0.045854},
		{-0.002079, 0.048952, 0.953127},
	}
)

// CAM16 default viewing conditions: D65, L_A = 200/π·Y(L*50)/100, gray
// background, average surround, no discounting.
var vc = func() (v struct {
	n, aw, nbb, ncb, c, nc, fl, z float64
	rgbD                          vec
}) {
	la := 200 / math.Pi * yFromLstar(50) / 100
	w := m16.mul(whiteD65)
	f := 1.0                                                           // 0.8 + surround/10
	v.c = 0.69                                                         // lerp(0.59, 0.69, (f-0.9)*10)
	d := math.Max(0, math.Min(1, f*(1-(1/3.6)*math.Exp((-la-42)/92)))) // degree of adaptation
	for i := range 3 {
		v.rgbD[i] = d*(100/w[i]) + 1 - d
	}
	k := 1 / (5*la + 1)
	k4 := k * k * k * k
	v.fl = k4*la + 0.1*(1-k4)*(1-k4)*math.Cbrt(5*la)
	v.n = yFromLstar(50) / whiteD65[1]
	v.z = 1.48 + math.Sqrt(v.n)
	v.nbb = 0.725 / math.Pow(v.n, 0.2)
	v.ncb = v.nbb
	v.nc = f
	var rgbA vec
	for i := range 3 {
		af := math.Pow(v.fl*v.rgbD[i]*w[i]/100, 0.42)
		rgbA[i] = 400 * af / (af + 27.13)
	}
	v.aw = (2*rgbA[0] + rgbA[1] + 0.05*rgbA[2]) * v.nbb
	return v
}()

// scaledDiscount maps linear RGB (0..100) straight to the CAM16 input of the
// chromatic adaptation step; linrgbFromScaled is its inverse.
var scaledDiscount, linrgbFromScaled = func() (mat, mat) {
	m := m16
	for i := range 3 {
		m[i] = xyzRow(m16[i])
		for j := range 3 {
			m[i][j] *= vc.fl * vc.rgbD[i] / 100
		}
	}
	return m, m.inverse()
}()

func xyzRow(r vec) vec {
	var out vec
	for j := range 3 {
		out[j] = r[0]*xyzFromLinrgb[0][j] + r[1]*xyzFromLinrgb[1][j] + r[2]*xyzFromLinrgb[2][j]
	}
	return out
}

// criticalPlanes[i] is the linear value halfway between sRGB i and i+1.
var criticalPlanes = func() (p [255]float64) {
	for i := range p {
		p[i] = linearized(float64(i) + 0.5)
	}
	return p
}()

func linearized(c255 float64) float64 {
	n := c255 / 255
	if n <= 0.040449936 {
		return n / 12.92 * 100
	}
	return math.Pow((n+0.055)/1.055, 2.4) * 100
}

func trueDelinearized(c100 float64) float64 {
	n := c100 / 100
	if n <= 0.0031308 {
		return n * 12.92 * 255
	}
	return (1.055*math.Pow(n, 1/2.4) - 0.055) * 255
}

func delinearized(c100 float64) uint8 {
	return uint8(math.Max(0, math.Min(255, math.Round(trueDelinearized(c100)))))
}

func fromLinrgb(l vec) color.NRGBA {
	return color.NRGBA{delinearized(l[0]), delinearized(l[1]), delinearized(l[2]), 0xff}
}

const labE, labK = 216.0 / 24389, 24389.0 / 27

func yFromLstar(l float64) float64 {
	ft := (l + 16) / 116
	if ft3 := ft * ft * ft; ft3 > labE {
		return 100 * ft3
	}
	return 100 * (116*ft - 16) / labK
}

func lstarFromY(y float64) float64 { return 116*labF(y/100) - 16 }

func labF(t float64) float64 {
	if t > labE {
		return math.Cbrt(t)
	}
	return (labK*t + 16) / 116
}

func sanitizeDeg(d float64) float64 {
	d = math.Mod(d, 360)
	if d < 0 {
		d += 360
	}
	return d
}

func sanitizeRad(r float64) float64 { return math.Mod(r+8*math.Pi, 2*math.Pi) }

func adapt(c float64) float64 {
	af := math.Pow(math.Abs(c), 0.42)
	return math.Copysign(400*af/(af+27.13), c)
}

func unadapt(a float64) float64 {
	abs := math.Abs(a)
	return math.Copysign(math.Pow(math.Max(0, 27.13*abs/(400-abs)), 1/0.42), a)
}

func cam16Ab(l vec) (rA, gA, bA, a, b float64) {
	s := scaledDiscount.mul(l)
	rA, gA, bA = adapt(s[0]), adapt(s[1]), adapt(s[2])
	return rA, gA, bA, (11*rA - 12*gA + bA) / 11, (rA + gA - 2*bA) / 9
}

func hctOf(c color.Color) hct {
	n := color.NRGBAModel.Convert(c).(color.NRGBA)
	l := vec{linearized(float64(n.R)), linearized(float64(n.G)), linearized(float64(n.B))}
	rA, gA, bA, a, b := cam16Ab(l)
	u := (20*rA + 20*gA + 21*bA) / 20
	p2 := (40*rA + 20*gA + bA) / 20
	hue := sanitizeDeg(math.Atan2(b, a) * 180 / math.Pi)
	j := 100 * math.Pow(p2*vc.nbb/vc.aw, vc.c*vc.z)
	eHue := 0.25 * (math.Cos(hue*math.Pi/180+2) + 3.8)
	t := 50000 / 13.0 * eHue * vc.nc * vc.ncb * math.Hypot(a, b) / (u + 0.305)
	alpha := math.Pow(t, 0.9) * math.Pow(1.64-math.Pow(0.29, vc.n), 0.73)
	return hct{hue, alpha * math.Sqrt(j/100), lstarFromY(xyzFromLinrgb[1][0]*l[0] + xyzFromLinrgb[1][1]*l[1] + xyzFromLinrgb[1][2]*l[2])}
}

// color returns the sRGB color closest to h, reducing chroma if out of gamut.
func (h hct) color() color.NRGBA {
	if h.c < 0.0001 || h.t < 0.0001 || h.t > 99.9999 {
		g := delinearized(yFromLstar(h.t))
		return color.NRGBA{g, g, g, 0xff}
	}
	rad := sanitizeDeg(h.h) * math.Pi / 180
	y := yFromLstar(h.t)
	if c, ok := findByJ(rad, h.c, y); ok {
		return c
	}
	return fromLinrgb(bisectToLimit(y, rad))
}

// findByJ inverts CAM16 with Newton's method on J; fails when out of gamut.
func findByJ(hue, chroma, y float64) (color.NRGBA, bool) {
	j := math.Sqrt(y) * 11
	tInner := 1 / math.Pow(1.64-math.Pow(0.29, vc.n), 0.73)
	p1 := 0.25 * (math.Cos(hue+2) + 3.8) * (50000 / 13.0) * vc.nc * vc.ncb
	hSin, hCos := math.Sincos(hue)
	for i := range 5 {
		jn := j / 100
		alpha := 0.0
		if chroma != 0 && j != 0 {
			alpha = chroma / math.Sqrt(jn)
		}
		t := math.Pow(alpha*tInner, 1/0.9)
		p2 := vc.aw * math.Pow(jn, 1/vc.c/vc.z) / vc.nbb
		gamma := 23 * (p2 + 0.305) * t / (23*p1 + 11*t*hCos + 108*t*hSin)
		a, b := gamma*hCos, gamma*hSin
		l := linrgbFromScaled.mul(vec{
			unadapt((460*p2 + 451*a + 288*b) / 1403),
			unadapt((460*p2 - 891*a - 261*b) / 1403),
			unadapt((460*p2 - 220*a - 6300*b) / 1403),
		})
		if l[0] < 0 || l[1] < 0 || l[2] < 0 {
			return color.NRGBA{}, false
		}
		fnj := xyzFromLinrgb[1][0]*l[0] + xyzFromLinrgb[1][1]*l[1] + xyzFromLinrgb[1][2]*l[2]
		if fnj <= 0 {
			return color.NRGBA{}, false
		}
		if i == 4 || math.Abs(fnj-y) < 0.002 {
			if l[0] > 100.01 || l[1] > 100.01 || l[2] > 100.01 {
				return color.NRGBA{}, false
			}
			return fromLinrgb(l), true
		}
		j -= (fnj - y) * j / (2 * fnj) // Newton step, fn'(j) ≈ 2·fn(j)/j
	}
	return color.NRGBA{}, false
}

func hueOf(l vec) float64 {
	_, _, _, a, b := cam16Ab(l)
	return math.Atan2(b, a)
}

func inCyclicOrder(a, b, c float64) bool { return sanitizeRad(b-a) < sanitizeRad(c-a) }

func lerpTo(src vec, coord float64, dst vec, axis int) vec {
	t := (coord - src[axis]) / (dst[axis] - src[axis])
	return vec{src[0] + (dst[0]-src[0])*t, src[1] + (dst[1]-src[1])*t, src[2] + (dst[2]-src[2])*t}
}

// nthVertex returns the nth of 12 edges of the RGB cube intersected with the
// plane of luminance y, or ok=false if that edge misses it.
func nthVertex(y float64, n int) (vec, bool) {
	k := xyzFromLinrgb[1]
	ca, cb := 0.0, 0.0
	if n%4 > 1 {
		ca = 100
	}
	if n%2 == 1 {
		cb = 100
	}
	var v vec
	var free int
	switch {
	case n < 4:
		v, free = vec{0, ca, cb}, 0
	case n < 8:
		v, free = vec{cb, 0, ca}, 1
	default:
		v, free = vec{ca, cb, 0}, 2
	}
	v[free] = (y - k[0]*v[0] - k[1]*v[1] - k[2]*v[2]) / k[free] // v[free] is still 0 here
	return v, 0 <= v[free] && v[free] <= 100
}

func bisectToSegment(y, target float64) (left, right vec) {
	var lHue, rHue float64
	initialized, uncut := false, true
	for n := range 12 {
		mid, ok := nthVertex(y, n)
		if !ok {
			continue
		}
		mHue := hueOf(mid)
		if !initialized {
			left, right, lHue, rHue, initialized = mid, mid, mHue, mHue, true
			continue
		}
		if uncut || inCyclicOrder(lHue, mHue, rHue) {
			uncut = false
			if inCyclicOrder(lHue, target, mHue) {
				right, rHue = mid, mHue
			} else {
				left, lHue = mid, mHue
			}
		}
	}
	return left, right
}

// bisectToLimit finds the most chromatic in-gamut color at luminance y and
// the target hue, by bisecting along the gamut boundary on sRGB critical planes.
func bisectToLimit(y, target float64) vec {
	left, right := bisectToSegment(y, target)
	lHue := hueOf(left)
	for axis := range 3 {
		if left[axis] == right[axis] {
			continue
		}
		var lP, rP int
		if left[axis] < right[axis] {
			lP = int(math.Floor(trueDelinearized(left[axis]) - 0.5))
			rP = int(math.Ceil(trueDelinearized(right[axis]) - 0.5))
		} else {
			lP = int(math.Ceil(trueDelinearized(left[axis]) - 0.5))
			rP = int(math.Floor(trueDelinearized(right[axis]) - 0.5))
		}
		for range 8 {
			if abs(rP-lP) <= 1 {
				break
			}
			mP := int(math.Floor(float64(lP+rP) / 2))
			mid := lerpTo(left, criticalPlanes[mP], right, axis)
			mHue := hueOf(mid)
			if inCyclicOrder(lHue, target, mHue) {
				right, rP = mid, mP
			} else {
				left, lHue, lP = mid, mHue, mP
			}
		}
	}
	return vec{(left[0] + right[0]) / 2, (left[1] + right[1]) / 2, (left[2] + right[2]) / 2}
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// harmonize shifts design's hue toward source by up to 15°, the M3 way of
// making custom colors (success/warning) sit well with the seed.
func harmonize(design, source hct) hct {
	diff := 180 - math.Abs(math.Abs(design.h-source.h)-180)
	rot := math.Min(diff*0.5, 15)
	if sanitizeDeg(source.h-design.h) > 180 {
		rot = -rot
	}
	design.h = sanitizeDeg(design.h + rot)
	return design
}

// contrast is the WCAG contrast ratio between two tones.
func contrast(a, b float64) float64 {
	ya := yFromLstar(math.Max(0, math.Min(100, a)))
	yb := yFromLstar(math.Max(0, math.Min(100, b)))
	return (math.Max(ya, yb) + 5) / (math.Min(ya, yb) + 5)
}

// foregroundTone is the tone closest to bg that reaches ratio against it,
// preferring light text on dark backgrounds (DynamicColor.foregroundTone).
func foregroundTone(bg, ratio float64) float64 {
	l := lstarFromY(ratio*(yFromLstar(bg)+5)-5) + 0.4
	if l < 0 || l > 100 {
		l = 100
	}
	d := lstarFromY((yFromLstar(bg)+5)/ratio-5) - 0.4
	if d < 0 || d > 100 {
		d = 0
	}
	lr, dr := contrast(l, bg), contrast(d, bg)
	if math.Round(bg) < 60 {
		if lr >= ratio || lr >= dr || (math.Abs(lr-dr) < 0.1 && lr < ratio && dr < ratio) {
			return l
		}
		return d
	}
	if dr >= ratio || dr >= lr {
		return d
	}
	return l
}

// temperature is warmth in Lab: hue near 50° (orange) is warm, 230° cold.
func temperature(h hct) float64 {
	c := h.color()
	l := vec{linearized(float64(c.R)), linearized(float64(c.G)), linearized(float64(c.B))}
	xyz := xyzFromLinrgb.mul(l)
	fx, fy, fz := labF(xyz[0]/whiteD65[0]), labF(xyz[1]/whiteD65[1]), labF(xyz[2]/whiteD65[2])
	a, b := 500*(fx-fy), 200*(fy-fz)
	hue := sanitizeDeg(math.Atan2(b, a) * 180 / math.Pi)
	return -0.5 + 0.02*math.Pow(math.Hypot(a, b), 1.07)*math.Cos(sanitizeDeg(hue-50)*math.Pi/180)
}

// complement is the hue at s's chroma and tone whose relative warmth mirrors
// s's, walking from the warm/cold extreme on the other side of the wheel.
func complement(s hct) hct {
	var byHue [361]hct
	var temps [361]float64
	sTemp := temperature(s)
	cold, warm := sTemp, sTemp
	coldH, warmH := s.h, s.h
	for h := range byHue {
		byHue[h] = hctOf(hct{float64(h), s.c, s.t}.color())
		temps[h] = temperature(byHue[h])
		if temps[h] < cold {
			cold, coldH = temps[h], byHue[h].h
		}
		if temps[h] > warm {
			warm, warmH = temps[h], byHue[h].h
		}
	}
	rel := 0.5
	if warm != cold {
		rel = (sTemp - cold) / (warm - cold)
	}
	start, end := coldH, warmH
	if hueBetween(s.h, coldH, warmH) {
		start, end = warmH, coldH
	}
	answer, best := byHue[int(math.Round(s.h))], 1000.0
	for add := 0.0; add <= 360; add++ {
		h := sanitizeDeg(start + add)
		if !hueBetween(h, start, end) {
			continue
		}
		i := int(math.Round(h))
		if e := math.Abs(1 - rel - (temps[i]-cold)/(warm-cold)); e < best {
			answer, best = byHue[i], e
		}
	}
	return answer
}

func hueBetween(h, a, b float64) bool {
	if a < b {
		return a <= h && h <= b
	}
	return a <= h || h <= b
}

// fixIfDisliked lightens dark yellow-greens (bile), which people dislike.
func fixIfDisliked(h hct) hct {
	if hue := math.Round(h.h); hue >= 90 && hue <= 111 && math.Round(h.c) > 16 && math.Round(h.t) < 65 {
		return hctOf(hct{h.h, h.c, 70}.color())
	}
	return h
}
