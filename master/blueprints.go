package main

type Blueprint struct {
	Dims      Position
	Buildings []Building
}

var twoItemBp = Blueprint{
	Position{3, 3},
	[]Building{
		//BBIASM
		//BBLASM
		//BB ASM
		Building{
			"belt",
			down,
			"",
			Position{0, 0},
		},
		Building{
			"belt",
			down,
			"",
			Position{0, 1},
		},
		Building{
			"belt",
			down,
			"",
			Position{0, 2},
		},
		Building{
			"belt",
			down,
			"",
			Position{1, 0},
		},
		Building{
			"belt",
			down,
			"",
			Position{1, 1},
		},
		Building{
			"belt",
			down,
			"",
			Position{1, 2},
		},
		Building{
			"inserter",
			right,
			"",
			Position{2, 0},
		},
		Building{
			"long-hand-inserter",
			left,
			"",
			Position{2, 1},
		},
		Building{
			"asm",
			up,
			"",
			Position{3, 0},
		},
	},
}
