package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	Username string
	Score    int
}

type BlockDropMessage struct {
	Action string `json:"action"` // e.g., "drop"
	Block  int    `json:"block"`  // The block number that was dropped
	Column int    `json:"column"` // The column where the block was dropped
}

// GameState struct for storing important information about the game for single player or versus
type GameState struct {
	Grid *Grid
	Time time.Time
}

type GameRoom struct {
	ID      string
	Grid    *Grid
	Players map[*websocket.Conn]*Player
}

type Cell struct {
	Number  int
	Changed bool
}

type Grid struct {
	Grid  [GridHeight][GridWidth]Cell
	Drops int
}

//rooms := make(GameRoom, 0)

const GridHeight = 8
const GridWidth = 7

func initGrid(g *Grid) {
	g.Drops = 0
	for i := 0; i < len(g.Grid); i++ {
		for j := 0; j < len(g.Grid[i]); j++ {
			g.Grid[i][j].Number = 0
			g.Grid[i][j].Changed = false
		}
	}
}

func dropBlock(block int, col int, g *Grid) int {
	fmt.Println("in drop block")

	var row int
	for i := GridHeight - 1; i >= 0; i-- {
		if g.Grid[i][col].Number == 0 {
			g.Grid[i][col].Number = block
			row = i
			break // Block is placed, exit the loop
		}
	}
	//fmt.Println("at end of dropblock")
	return row
}

func printGrid(g *Grid) {
	for i := 0; i < len(g.Grid); i++ {
		for j := 0; j < len(g.Grid[i]); j++ {
			fmt.Printf("%d ", g.Grid[i][j].Number)
		}
		fmt.Println()
	}
	fmt.Println()
}

func gravity(g *Grid) {
	fmt.Println("in gravity")

	for i := 0; i < GridHeight; i++ {
		for j := 0; j < GridWidth; j++ {
			if g.Grid[i][j].Number != 0 {
				gravityHelp(g, j)
			}
		}
	}
	//fmt.Println("at end of gravity")
}

func checkForChange(g *Grid) bool {
	fmt.Println("in check for change")
	changed := false
	for i := 0; i < GridHeight; i++ {
		for j := 0; j < GridWidth; j++ {
			if g.Grid[i][j].Changed {
				changed = true
				break
			}
		}
	}
	//fmt.Println("at end of check for change")
	return changed
}

func setCellsToUnchanged(g *Grid) {
	for i := 0; i < GridHeight; i++ {
		for j := 0; j < GridWidth; j++ {
			g.Grid[i][j].Changed = false
		}
	}
}

func gravityHelp(g *Grid, col int) {
	fmt.Println("in gravity help")

	var newColumn = make([]Cell, 0)
	var oldColumn = make([]Cell, 0)

	for i := 0; i < GridHeight; i++ {
		oldColumn = append(oldColumn, g.Grid[i][col])
	}
	fmt.Printf("oldcolumn: %v\n", oldColumn)
	fixEmptyCells(g)
	for j := GridHeight - 1; j >= 0; j-- {
		if g.Grid[j][col].Number != 0 {
			newColumn = append(newColumn, g.Grid[j][col])
		}
	}
	fmt.Printf("newcolumn: %v\n", newColumn)

	slices.Reverse(newColumn)
	fmt.Printf("correct newColumn: %v\n", newColumn)

	var zerosSlice = make([]Cell, len(oldColumn)-len(newColumn))

	newColumn = append(zerosSlice, newColumn...)
	fmt.Printf("correcter newColumn: %v\n", newColumn)
	//fmt.Println("before replace column call in gravity help")
	replaceColumn(g, newColumn, col)
	//fmt.Println("after replace column in gravity help")

	if checkForChange(g) {
		//fmt.Println("in check for change if before check break")
		//setCellsToUnchanged(g)
		checkBreak(g)
		//fmt.Println("in check for change if after check break")
	}
	//fmt.Println("at end of gravity help")

}

func fixEmptyCells(g *Grid) bool {
	fmt.Println("in fix empty cells")
	didChange := false
	for i := 0; i < GridHeight; i++ {
		for j := 0; j < GridWidth; j++ {
			if g.Grid[i][j].Number == -3 {
				g.Grid[i][j].Number = 0
				g.Grid[i][j].Changed = false
				didChange = true
			}
		}
	}
	return didChange
	//fmt.Println("at end of fix empty cells")
}

func replaceRow(g *Grid, row []Cell, rowIndex int) {
	fmt.Println("in replaceRow")

	for i := 0; i < GridWidth; i++ {
		g.Grid[rowIndex][i] = row[i]
	}
	//fmt.Println("at end of replace row")
}

func replaceColumn(g *Grid, col []Cell, colIndex int) {
	fmt.Println("in replaceColumn")
	for i := 0; i < GridHeight; i++ {
		g.Grid[i][colIndex] = col[i]
	}
	//fmt.Println("at end of replace column")
}

func checkBreakRow(g *Grid, row []Cell, rowIndex int) {
	fmt.Println("in checkBreakRow")
	//printGrid(g)
	len := 0
	for i := 0; i < GridWidth; i++ {
		if row[i].Number != 0 {
			len++
		} else if len != 0 {
			for j := 1; j < len+1; j++ {
				if row[i-j].Number == len {
					row[i-j].Number = -3
					//len = 0
				}
			}
			len = 0
			fmt.Println("Reset len bc of len != 0")
		} else {
			len = 0
			fmt.Println("Reset len bc of else")
		}
		fmt.Printf("len: %v\n", len)
	}

	fmt.Printf("final len: %v\n", len)
	// for i := 0; i < GridWidth; i++ {
	// 	if len == 7 {
	// 		if row[i].Number == 7 {
	// 			row[i].Number = -3
	// 		}
	// 		fixEmptyCells(g)
	// 		checkBreakRow(g, row, rowIndex)
	// 	} else {
	// 		for j := 0; j >= len; j++ {
	// 			if row[i-j].Number == len {
	// 				row[i-j].Number = -3
	// 				//len = 0
	// 			}
	// 		}
	// 	}
	// }
	for j := 0; j < len; j++ {
		if row[GridWidth-1-j].Number == len {
			row[GridWidth-1-j].Number = -3
			//len = 0
		}
	}
	if fixEmptyCells(g) {
		checkBreakRow(g, row, rowIndex)
	}

	replaceRow(g, row, rowIndex)
	fmt.Println("at end of check break row")
}

func checkBreakColumn(g *Grid, col []Cell, colIndex int) {
	fmt.Println("in checkBreakColumn")

	len := 0
	for i := 0; i < GridHeight; i++ {
		if col[i].Number != 0 {
			len++
		}
	}

	flag := false
	for j := GridHeight - 1; j > 0; j-- {
		if col[j].Number == len && !flag {
			col[j].Number = -3
			flag = true
		} else if col[j].Number != 0 && flag {
			if col[j].Number == len {
				col[j].Number = -3
			}
			col[j].Changed = true
		}
	}

	replaceColumn(g, col, colIndex)
	//fmt.Println("at end of check break column")
}

func checkBreak(g *Grid) {
	fmt.Println("in checkBreak")
	for i := 0; i < GridWidth; i++ {
		for j := 0; j < GridWidth; j++ {
			if g.Grid[i][j].Changed {
				checkBreakRow(g, g.Grid[i][:], i)
				var newColumn = make([]Cell, 0)
				for k := 0; k < GridHeight; k++ {
					newColumn = append(newColumn, g.Grid[k][j])
				}
				checkBreakColumn(g, newColumn, j)
				g.Grid[i][j].Changed = false
			}
		}
	}
	fixEmptyCells(g)
	gravity(g)
	//fmt.Println("at end of check break")
}

func raiseBarriers(g *Grid) {
	for i := 0; i < GridHeight; i++ {
		for j := 0; j < GridWidth; j++ {
			if g.Grid[i][j].Number != 0 {
				g.Grid[i-1][j].Number = g.Grid[i][j].Number
				g.Grid[i-1][j].Changed = true
			}
		}
		for k := 0; k < GridWidth; k++ {
			if i == GridHeight-1 {
				g.Grid[i][k].Number = -2
			}
		}
	}
	checkBreak(g)
}

func drop(g *Grid, colDropped int, block int) {
	fmt.Println("in drop")
	if g.Drops == 6 {
		raiseBarriers(g)
		g.Drops = 0
	}
	g.Drops++

	row := dropBlock(block, colDropped, g)
	checkBreakRow(g, g.Grid[row][:], row)
	var col = make([]Cell, 0)

	for i := 0; i < GridHeight; i++ {
		col = append(col, g.Grid[i][colDropped])
	}

	checkBreakColumn(g, col, colDropped)
	gravity(g)
	fixEmptyCells(g)
	//fmt.Println("at end of drop")

}

// func placeBlock() bool {

// }
