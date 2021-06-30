package main

type Blueprint struct {
	Dims      Position
	Buildings []Building
}

var noFluidBp = Blueprint{
	Position{8, 3},
	[]Building{
		//BBIASMIB
		//BBLASM B
		//BB ASM B
		Building{
			"asm",
			up,
			"",
			Position{3, 0},
		},
		Building{
			"electric-pole",
			up,
			"",
			Position{2, 2},
		},
		Building{
			"electric-pole",
			up,
			"",
			Position{6, 2},
		},
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
			right,
			"",
			Position{2, 1},
		},
		Building{
			"inserter",
			right,
			"",
			Position{6, 0},
		},
		Building{
			"belt",
			down,
			"",
			Position{7, 0},
		},
		Building{
			"belt",
			down,
			"",
			Position{7, 1},
		},
		Building{
			"belt",
			down,
			"",
			Position{7, 2},
		},
	},
}

var fluidBp = Blueprint{
	Position{10, 3},
	/*
	     BBIASMIB
	   PUBBUASM B
	   M BBLASM B
	*/
	[]Building{
		Building{
			"asm",
			left,
			"",
			Position{5, 0},
		},
		Building{
			"electric-pole",
			up,
			"",
			Position{1, 0},
		},
		Building{
			"electric-pole",
			up,
			"",
			Position{8, 1},
		},
		Building{
			"pipe",
			down,
			"",
			Position{0, 1},
		},
		Building{
			"underground-pipe",
			left,
			"",
			Position{1, 1},
		},
		Building{
			"underground-pipe",
			right,
			"",
			Position{4, 1},
		},
		Building{
			"pump",
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
			Position{2, 0},
		},
		Building{
			"belt",
			down,
			"",
			Position{9, 0},
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
			Position{2, 1},
		},
		Building{
			"belt",
			down,
			"",
			Position{9, 1},
		},
		Building{
			"belt",
			down,
			"",
			Position{1, 2},
		},
		Building{
			"belt",
			down,
			"",
			Position{2, 2},
		},
		Building{
			"belt",
			down,
			"",
			Position{9, 2},
		},
		Building{
			"inserter",
			right,
			"",
			Position{3, 0},
		},
		Building{
			"long-inserter",
			right,
			"",
			Position{3, 2},
		},
		Building{
			"inserter",
			right,
			"",
			Position{8, 0},
		},
	},
}
