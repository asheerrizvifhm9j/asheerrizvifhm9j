package treed

import "github.com/unixpickle/model3d/model3d"

func splitTriangle(t *model3d.Triangle, axis model3d.Coord3D, threshold float64) (lessThan,
	greaterEqual []*model3d.Triangle) {
	var signs []bool
	for i, c := range t {
		if axis.Dot(c) >= threshold {
			signs[i] = true
		}
	}
	if signs[0] == signs[1] && signs[1] == signs[2] {
		if signs[0] {
			return []*model3d.Triangle{}, []*model3d.Triangle{t}
		} else {
			return []*model3d.Triangle{t}, []*model3d.Triangle{}
		}
	}

	// Find the majority sign
	var trueCount int
	for _, s := range signs {
		if s {
			trueCount++
		}
	}
	majority := trueCount == 2

	// The lines between the majorities and the minority intersect the plane.
	majLoop := make([]model3d.Coord3D, 0, 4)
	minLoop := make([]model3d.Coord3D, 0, 3)
	for i, c := range t {
		if signs[i] == signs[(i+1)%3] {
			majLoop = append(majLoop, c)
			continue
		}

		// Find plane intersection.
		p1 := t[i]
		p2 := t[(i+1)%3]

		// x = o + tr (ray)
		// n*x = b    (plane)
		// => n*(o+tr) = b
		// => n*o + t*(n*r) = b
		// => t*(n*r) = b - n*o
		// => t = (b - n*o)/(n*r)
		o := p1
		r := p2.Sub(p1)
		alpha := (threshold - axis.Dot(o)) / (axis.Dot(r))

		// In the case of rounding error putting us outside the triangle,
		// we abort splitting and assume the triangle is really on one
		// side.
		if alpha <= 0 {
			if signs[(i+1)%3] {
				return []*model3d.Triangle{}, []*model3d.Triangle{t}
			} else {
				return []*model3d.Triangle{t}, []*model3d.Triangle{}
			}
		} else if alpha >= 1 {
			if signs[i] {
				return []*model3d.Triangle{}, []*model3d.Triangle{t}
			} else {
				return []*model3d.Triangle{t}, []*model3d.Triangle{}
			}
		}

		midPoint := o.Add(r.Scale(alpha))
		majLoop = append(majLoop, midPoint)
		minLoop = append(minLoop, midPoint)

		if signs[i] == majority {
			minLoop = append(minLoop, p2)
		} else {
			majLoop = append(majLoop, p2)
		}
	}

	majTris := []*model3d.Triangle{
		{majLoop[0], majLoop[1], majLoop[3]},
		{majLoop[1], majLoop[2], majLoop[3]},
	}
	minTris := []*model3d.Triangle{
		{minLoop[0], minLoop[1], minLoop[2]},
	}
	if majority {
		return minTris, majTris
	} else {
		return majTris, minTris
	}
}
