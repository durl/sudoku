// Copyright (c) 2016, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func solveFile(filename string) (Sudoku, error) {
	f, _ := os.Open(filename)
	sudoku, _ := readSudoku(f)
	return solve(sudoku)
}

func readFile(filename string) string {
	text, _ := ioutil.ReadFile(filename)
	return strings.TrimSpace(string(text))
}

func TestSolveSimple(t *testing.T) {
	solution, err := solveFile("test/simple.in")
	if err != nil {
		t.Errorf("Could not solve: %s", err.Error())
	}
	if readFile("test/simple.solution") != solution.String() {
		t.Errorf("erroneus solution")
	}
}

func TestSolveWithGuessing(t *testing.T) {
	solution, err := solveFile("test/guessing.in")
	if err != nil {
		t.Errorf("Could not solve: %s", err.Error())
	}
	if readFile("test/guessing.solution") != solution.String() {
		t.Errorf("erroneus solution")
	}
}

func TestHard(t *testing.T) {
	solution, err := solveFile("test/hard.in")
	if err != nil {
		t.Errorf("Could not solve: %s", err.Error())
	}
	if readFile("test/hard.solution") != solution.String() {
		t.Errorf("erroneus solution")
	}
}

func TestConflictError(t *testing.T) {
	f, _ := os.Open("test/conflict.in")
	_, err := readSudoku(f)
	if err == nil {
		t.Errorf("conflict not detected")
	}
}

func TestUnsolveableError(t *testing.T) {
	_, err := solveFile("test/unsolveable.in")
	if err == nil {
		t.Errorf("not recognized as unsolveable")
	}
}
