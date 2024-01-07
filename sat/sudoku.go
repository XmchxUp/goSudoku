package sat

import (
	"fmt"
)

// literalVar generates a unique string for each cell and digit
func literalVar(r, c, d int) string {
	return fmt.Sprintf("%d_%d_%d", r, c, d)
}

// SudokuBoardToSatFormula converts a sudoku board to a SAT formula
func SudokuBoardToSatFormula(sudokuBoard [][]int) Formula {
	N := 9
	n := 3
	formula := Formula{}

	// Initial state constraints
	for r := 0; r < N; r++ {
		for c := 0; c < N; c++ {
			d := sudokuBoard[r][c]
			if d != 0 {
				literal := Literal{Variable: literalVar(r, c, d), Value: true}
				clause := Clause{literal}
				formula = append(formula, clause)
			}
		}
	}

	// Cell constraints
	for r := 0; r < N; r++ {
		for c := 0; c < N; c++ {
			// At least one number in each cell
			clause := Clause{}
			for d := 1; d <= N; d++ {
				clause = append(clause, Literal{Variable: literalVar(r, c, d), Value: true})
			}
			formula = append(formula, clause)

			// At most one number in each cell
			for i := 1; i <= N; i++ {
				clause = Clause{}
				for j := i + 1; j <= N; j++ {
					clause = append(clause, Literal{Variable: literalVar(r, c, i), Value: false},
						Literal{Variable: literalVar(r, c, j), Value: false})
				}
				if len(clause) > 0 {
					formula = append(formula, clause)
				}
			}
		}
	}

	// Row and Column constraints
	for d := 1; d <= N; d++ {
		for r := 0; r < N; r++ {
			// Row constraint
			clause := Clause{}
			for c := 0; c < N; c++ {
				clause = append(clause, Literal{Variable: literalVar(r, c, d), Value: true})
			}
			formula = append(formula, clause)

			// Column constraint
			clause = Clause{}
			for c := 0; c < N; c++ {
				clause = append(clause, Literal{Variable: literalVar(c, r, d), Value: true})
			}
			formula = append(formula, clause)

			// At most one number in each row and column
			for c1 := 0; c1 < N; c1++ {
				for c2 := c1 + 1; c2 < N; c2++ {
					formula = append(formula, Clause{
						Literal{Variable: literalVar(r, c1, d), Value: false},
						Literal{Variable: literalVar(r, c2, d), Value: false},
					}, Clause{
						Literal{Variable: literalVar(c1, r, d), Value: false},
						Literal{Variable: literalVar(c2, r, d), Value: false},
					})
				}
			}
		}
	}

	// Block constraints
	for d := 1; d <= N; d++ {
		for br := 0; br < n; br++ {
			for bc := 0; bc < n; bc++ {
				// At least one number in each block
				clause := Clause{}
				for rr := 0; rr < n; rr++ {
					for cc := 0; cc < n; cc++ {
						clause = append(clause, Literal{Variable: literalVar(br*n+rr, bc*n+cc, d), Value: true})
					}
				}
				formula = append(formula, clause)

				// At most one number in each block
				for i := 0; i < n*n; i++ {
					for j := i + 1; j < n*n; j++ {
						r1, c1 := i/n, i%n
						r2, c2 := j/n, j%n
						formula = append(formula, Clause{
							Literal{Variable: literalVar(br*n+r1, bc*n+c1, d), Value: false},
							Literal{Variable: literalVar(br*n+r2, bc*n+c2, d), Value: false},
						})
					}
				}
			}
		}
	}

	return formula
}

// AssignmentsToSudokuBoard converts SAT variable assignments to a sudoku board
func AssignmentsToSudokuBoard(assignments map[string]bool, N int) [][]int {
	board := make([][]int, N)
	for i := range board {
		board[i] = make([]int, N)
	}

	for varIndex, val := range assignments {
		if val {
			var r, c, d int
			fmt.Sscanf(varIndex, "%d_%d_%d", &r, &c, &d)
			if board[r][c] != 0 {
				// Conflicting assignment, invalid board
				return nil
			}
			board[r][c] = d
		}
	}

	return board
}
