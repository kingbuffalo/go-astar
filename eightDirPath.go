package astar

import (
	"fmt"
)

// Kind* constants refer to EDTile kinds for input and output.
const (
	// EDKindPlain (.) is a plain EDTile with a movement cost of 1.
	EDKindPlain = iota
	// EDKindRiver (~) is a river EDTile with a movement cost of 2.
	EDKindRiver
	// EDKindMountain (M) is a mountain EDTile with a movement cost of 3.
	EDKindMountain
	// EDKindBlocker (X) is a EDTile which blocks movement.
	EDKindBlocker
	EDKindPath
)

// EDKindRunes map EDTile kinds to output runes.
var EDKindRunes = map[int]rune{
	EDKindPlain:    '.',
	EDKindRiver:    '~',
	EDKindMountain: 'M',
	EDKindBlocker:  'X',
	EDKindPath:     '●',
}

// EDRuneKinds map input runes to EDTile kinds.
var EDRuneKinds = map[rune]int{
	'.': EDKindPlain,
	'~': EDKindRiver,
	'M': EDKindMountain,
	'X': EDKindBlocker,
}

// EDKindCosts map EDTile kinds to movement costs.
var EDKindCosts = map[int]float64{
	EDKindPlain:    1.0,
	EDKindRiver:    2.0,
	EDKindMountain: 3.0,
}

// A EDTile is a EDTile in a grid which implements Pather.
type EDTile struct {
	// Kind is the kind of EDTile, potentially affecting movement.
	Kind int
	// X and Y are the coordinates of the EDTile.
	X, Y int
	// W is a reference to the EDWorld that the EDTile is a part of.
	W EDWorld
}

func (t *EDTile) Print() {
	fmt.Print("(", t.X, ",", t.Y, ")", "=", t.Kind, "  ")
}

// PathNeighbors returns the neighbors of the EDTile, excluding blockers and
// tiles off the edge of the board.
func (t *EDTile) PathNeighbors() []Pather {
	neighbors := []Pather{}
	for _, offset := range [][]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
		{1, 1},
		{-1, 1},
		{1, -1},
		{-1, 1},
	} {
		if n := t.W.EDTile(t.X+offset[0], t.Y+offset[1]); n != nil &&
			n.Kind != EDKindBlocker {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

func (t *EDTile) BSamePoint(to Pather) bool {
	toT := to.(*EDTile)
	return t.X == toT.X && t.Y == toT.Y
}

// PathNeighborCost returns the movement cost of the directly neighboring EDTile.

func (t *EDTile) PathNeighborCost(to Pather) float64 {
	toT := to.(*EDTile)
	cost := EDKindCosts[toT.Kind]
	xdis := t.X - toT.X
	if xdis < 0 {
		xdis = -xdis
	}
	ydis := t.Y - toT.Y
	if ydis < 0 {
		ydis = -ydis
	}
	if xdis == 1 && ydis == 1 {
		return cost * 1.414
	}
	return cost
}

// PathEstimatedCost uses Manhattan distance to estimate orthogonal distance
// between non-adjacent nodes.
func (t *EDTile) PathEstimatedCost(to Pather) float64 {
	toT := to.(*EDTile)
	absX := toT.X - t.X
	if absX < 0 {
		absX = -absX
	}
	absY := toT.Y - t.Y
	if absY < 0 {
		absY = -absY
	}
	return float64(absX + absY)
}

// EDWorld is a two dimensional map of Tiles.
type EDWorld map[int]map[int]*EDTile

func NewEDWorld() EDWorld {
	return map[int]map[int]*EDTile{}
}

// EDTile gets the EDTile at the given coordinates in the EDWorld.
func (w EDWorld) EDTile(x, y int) *EDTile {
	if w[x] == nil {
		return nil
	}
	return w[x][y]
}

func (w EDWorld) SetTile(bBloc bool, x, y int) {
	k := EDKindPlain
	if bBloc {
		k = EDKindBlocker
	}
	t := &EDTile{
		X:    x,
		Y:    y,
		W:    w,
		Kind: k,
	}
	if w[x] == nil {
		w[x] = map[int]*EDTile{}
	}
	w[x][y] = t
}

func (w EDWorld) Distance(fromx, fromy, tox, toy int) (float64, bool) {
	fmt.Println("Distance", fromx, fromy, tox, toy)
	from := &EDTile{
		Kind: EDKindPlain,
		X:    fromx,
		Y:    fromy,
		W:    w,
	}
	to := &EDTile{
		Kind: EDKindPlain,
		X:    tox,
		Y:    toy,
		W:    w,
	}

	_, d, f := Path(from, to)
	return d, f
}

func (w EDWorld) Print() {
	//hLen := 62
	//wLen := 38
	//for i := 0; i < hLen; i++ {
	//for j := 0; j < wLen; j++ {
	//edt := w.EDTile(j, i)
	//if edt == nil {
	//fmt.Errorf("nil=%d,%d", j, i)
	//}
	////fmt.Print(edt.Kind)
	//fmt.Print("(", j, ",", i, ")")
	//}
	//fmt.Print("\n")
	//}

	for x, v := range w {
		for y, vv := range v {
			fmt.Print("(", x, ",", y, ")", "=", vv.Kind, "  ")
		}
		fmt.Print("\n")
	}
}
