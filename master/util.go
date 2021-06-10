package main

func contains(s []Position, e Position) bool {
	for _, a := range s {
		if a.x == e.x && a.y == e.y {
			return true
		}
	}
	return false
}

func removeVal(s []Position, e Position) []Position {
	for i, a := range s {
		if a.x == e.x && a.y == e.y {
			return remove(s, i)
		}
	}
	return s
}

func remove(slice []Position, s int) []Position {
    return append(slice[:s], slice[s+1:]...)
}
