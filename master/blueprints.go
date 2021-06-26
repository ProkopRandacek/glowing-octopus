package main

type Blueprint struct {
	Dims      Position
	Buildings []Building
}

var noFluidBp =
Blueprint{
		Position{3, 3},
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
