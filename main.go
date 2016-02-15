// Copyright (c) 2016, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Sudoku struct {
	fields [9][9]int

	rows   [9][9]bool
	cols   [9][9]bool
	groups [3][3][9]bool

	complete int
}

func (s *Sudoku) set(row, col, val int) {
	s.fields[row][col] = val
	if val > 0 {
		i := val - 1
		s.rows[row][i] = true
		s.cols[col][i] = true
		s.groups[row/3][col/3][i] = true

		s.complete++
	}
}

func (s *Sudoku) has(row, col, val int) bool {
	i := val - 1
	return s.rows[row][i] || s.cols[col][i] || s.groups[row/3][col/3][i]
}

func (s *Sudoku) isComplete() bool {
	return s.complete == 9*9
}

func (s *Sudoku) String() string {
	var buf bytes.Buffer
	for i, row := range s.fields {
		for i, field := range row {
			if field == 0 {
				buf.WriteString("_")
			} else {
				buf.WriteString(strconv.Itoa(field))
			}
			if i < 8 {
				buf.WriteString(" ")
			}
		}
		if i < 8 {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

func readSudoku(r io.Reader) (Sudoku, error) {
	var sudoku Sudoku

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)

	var x, y int
	for scanner.Scan() {
		if y == 9 {
			return sudoku, fmt.Errorf("too many elements")
		}

		elem := scanner.Text()
		if elem == "_" {
			sudoku.set(y, x, 0)
		} else {
			i, err := strconv.Atoi(elem)
			if err != nil || i < 0 || i > 9 {
				return sudoku, fmt.Errorf("invalid input element: '%s'", elem)
			}
			if i > 0 && sudoku.has(y, x, i) {
				return sudoku, fmt.Errorf("conflicting element in row %d, column %d: %d", y+1, x+1, i)
			}
			sudoku.set(y, x, i)
		}

		x++
		if x == 9 {
			x = 0
			y++
		}
	}
	return sudoku, nil
}

func solve(sudoku Sudoku) (Sudoku, error) {
	if sudoku.isComplete() {
		return sudoku, nil
	}

	var bestX, bestY int
	var bestMissing int

	for y, row := range sudoku.fields {
		for x, field := range row {
			if field != 0 {
				continue
			}
			missing := 9
			for i := 1; i <= 9; i++ {
				if sudoku.has(y, x, i) {
					missing--
				}
			}
			if bestMissing == 0 || bestMissing > missing {
				bestX = x
				bestY = y
				bestMissing = missing
			}
		}
	}

	for i := 1; i <= 9; i++ {
		if !sudoku.has(bestY, bestX, i) {
			old := sudoku
			sudoku.set(bestY, bestX, i)
			solution, err := solve(sudoku)
			if err != nil {
				sudoku = old
				continue
			}
			return solution, nil
		}
	}
	return sudoku, fmt.Errorf("unsolveable")
}

func main() {
	sudoku, err := readSudoku(os.Stdin)
	if err != nil {
		quit(err)
	}
	solution, err := solve(sudoku)
	if err != nil {
		quit(err)
	}
	fmt.Println(solution.String())
}

func quit(err error) {
	fmt.Println(err)
	os.Exit(1)
}
