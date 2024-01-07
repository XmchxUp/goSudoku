package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/xmchxup/goSudoku/sat"
)

type VicotryCheckResposne struct {
	Victory bool `json:"victory"`
}

func isValidUnit(board []int) bool {
	visited := make(map[int]bool, 9)
	for _, v := range board {
		if visited[v] {
			return false
		}
		visited[v] = true
	}
	return true
}

func isValidSudoku(board [][]int) bool {
	for i := 0; i < 9; i++ {
		row := make([]int, 9)
		col := make([]int, 9)
		for j := 0; j < 9; j++ {
			row[i] = board[i][j]
			col[j] = board[j][i]
		}
		if !isValidUnit(row) || !isValidUnit(col) {
			return false
		}
	}

	for i := 0; i < 9; i += 3 {
		for j := 0; j < 9; j += 3 {
			nums := make([]int, 0, 9)
			for k := 0; k < 3; k++ {
				for l := 0; l < 3; l++ {
					nums = append(nums, board[i+k][j+l])
				}
			}
			if !isValidUnit(nums) {
				return false
			}
		}
	}
	return true
}

func solve(board [][]int) ([][]int, error) {
	formula := sat.SudokuBoardToSatFormula(board)
	assignment := sat.Assignment{}
	assignment = sat.SatisfyingAssignment(formula, assignment)
	if assignment == nil {
		return nil, errors.New("no solution")
	}
	return sat.AssignmentsToSudokuBoard(assignment, 9), nil
}

func main() {
	log.Println("Server is starting in http://localhost:9999")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if r.URL.Path != "/" {
				http.StripPrefix("/", http.FileServer(http.Dir("./ui"))).ServeHTTP(w, r)
			} else {
				http.ServeFile(w, r, "./ui/sudoku.html")
			}
		case "POST":
			log.Println(r.URL.Path)
			if r.URL.Path != "/victory_check" && r.URL.Path != "/solve" {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}

			bs, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}

			var board [][]int
			err = json.Unmarshal(bs, &board)
			if err != nil {
				http.Error(w, "Error parsing request json", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")

			if r.URL.Path == "/victory_check" {
				response, err := json.Marshal(VicotryCheckResposne{Victory: isValidSudoku(board)})
				if err != nil {
					http.Error(w, "Error generating JSON response for victory_check", http.StatusInternalServerError)
					return
				}
				w.Write(response)
				return
			}

			solution, err := solve(board)
			if err != nil {
				json.NewEncoder(w).Encode(nil)
				return
			}

			response, err := json.Marshal(solution)
			if err != nil {
				http.Error(w, "Error generating JSON response for solve", http.StatusInternalServerError)
				return
			}
			w.Write(response)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(":9999", nil))
}
