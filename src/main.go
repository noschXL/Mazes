package main

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type cell struct {
	ID int32
	Walls [4]bool
	Dir int32
}

type CellList []cell

var WallCol = color.RGBA{R: 0, G: 255, B: 0, A: 255}
var CellSizePx int32 = 20
var FPS int32 = 120

const (
	NORTH = 0
	EAST = 1
	SOUTH = 2
	WEST = 3
)

const (
	GRID = 0
	WALL = 1
)

func (Cell *cell) SetWall (dir int32, value bool, Mazesize int32) {
	Cell.Walls[dir] = value
	if dir == NORTH && Cell.ID - Mazesize >= 0 {
		Cell.Walls[SOUTH] = value
	}else if dir == EAST && Cell.ID % Mazesize < Mazesize - 1 {
		Cell.Walls[WEST] = value
	}else if dir == SOUTH && Cell.ID + Mazesize < Mazesize * Mazesize {
		Cell.Walls[NORTH] = value
	}else if Cell.ID % Mazesize >= 0 {
		Cell.Walls[EAST] = value
	}
}

func DrawMaze(maze CellList,size int32 ,mode int, MazeRect rl.Rectangle) {
	if mode == 1 {
		for i, c := range maze {
			x := (int32(i) % size) * CellSizePx + int32(MazeRect.X)
			y := (int32(i) / size) * CellSizePx + int32(MazeRect.Y)
		
			if c.Walls[NORTH] {
				rl.DrawLine(x, y, x + CellSizePx, y, WallCol)
			}
			if c.Walls[EAST] {
				rl.DrawLine(x + CellSizePx, y, x+CellSizePx, y+CellSizePx, WallCol)
			}
			if c.Walls[SOUTH] {
				rl.DrawLine(x, y+CellSizePx, x+CellSizePx, y+CellSizePx, WallCol)
			}
			if c.Walls[WEST] {
				rl.DrawLine(x, y, x, y+CellSizePx, WallCol)
			}
		}
	}
}

const (
	CLOSED = 0
	OPEN = 1
	ORIGIN = 2
) 

func StringReplaceAt(s string, i int, c rune) string {
	r := []rune(s)
	r[i] = c
	return string(r)
}

func printwalls (walls [4]bool) {
	chars := "   \n   "
	if walls[0]{
		chars = StringReplaceAt(chars, 1, '_')
	}
	if walls[1]{
		chars = StringReplaceAt(chars, 6, '|')
	}
	if walls[2]{
		chars = StringReplaceAt(chars, 5, '_')
	}
	if walls[3]{
		chars = StringReplaceAt(chars, 4, '|')  
	}

	println(chars)

}

func CreateMaze (size int32) (CellList) {
	var maze CellList
	for i := int32(0); i < size; i++ {
		maze = append(maze, cell{ID: i, Walls: [4]bool{false, false, false, false}, Dir: -1})
	}

	return maze
}
//Updates the walls based on the direction of a cell
func (maze CellList) UpdateWalls (size int32, originpoint *int32) (CellList){

	if maze[*originpoint].Dir != -1 {
		println("The Dev is an idiot, cant even keep track of a var. Ill go look for it...")
		for i, c := range maze {
			if c.Dir == -1 {
				*originpoint = int32(i)
				println("found it and updated it, its ", i)
				break
			}
		}
	}
	
	for i := range maze {
		var tocheck [4]bool = [4]bool{true,true,true,true}
		maze[i].Walls = [4]bool{true,true,true,true}

		if (int32(i) < size) || (maze[i].Dir == NORTH) {
			tocheck[NORTH] = false
		}

		if (int32(i) > int32(len(maze)) - size - 1) || (maze[i].Dir == SOUTH) {
			tocheck[SOUTH] = false
		}
		if (int32(i) % size == 0) || (maze[i].Dir == WEST) {
			tocheck[WEST] = false
		}
		if (int32(i) % size == size - 1) || (maze[i].Dir == EAST) {
			tocheck[EAST] = false
		}

		if !(maze[i].Dir == -1) {
			maze[i].Walls[maze[i].Dir] = false
		}

		if tocheck[NORTH] {
			maze[i].Walls[NORTH] = !(maze[int32(i) - size].Dir == SOUTH) 
		}
		if tocheck[EAST] {
			maze[i].Walls[EAST] = !(maze[int32(i)+ 1].Dir == WEST)
		}
		if tocheck[SOUTH] {
			maze[i].Walls[SOUTH] = !(maze[int32(i) + size].Dir == NORTH) 
		}
		if tocheck[WEST] {
			maze[i].Walls[WEST] = !(maze[int32(i) - 1].Dir == EAST) 
		}
	}

	return maze

}

func (maze CellList) Setup (mazetype int32, size int32) {
	if mazetype == CLOSED {
		for i := range maze {
			maze[i].Walls[NORTH] = true
			maze[i].Walls[EAST] = true
			maze[i].Walls[SOUTH] = true
			maze[i].Walls[WEST] = true
		}
	}else if mazetype == OPEN {
		for i := range maze {
			maze[i].Walls[NORTH] = false
			maze[i].Walls[EAST] = false
			maze[i].Walls[SOUTH] = false
			maze[i].Walls[WEST] = false	

			if int32(i) < size {
				maze[i].Walls[NORTH] = true
			}

			if int32(i) > int32(len(maze)) - size - 1 {
				maze[i].Walls[SOUTH] = true
			}
			if int32(i) % size == 0 {
				maze[i].Walls[WEST] = true
			}
			if int32(i) % size == size - 1 {
				maze[i].Walls[EAST] = true
			}
		}
	}else if mazetype == ORIGIN {
		for i := range maze{
			maze[i].Dir = EAST

			if int32(i) % size == size - 1 {
				maze[i].Dir = SOUTH
			}
			if i == len(maze) - 1 {
				maze[i].Dir = -1
			}
		}
	}
}

func (maze CellList)OriginShiftStep (size int32, originpoint *int32) {
	Dir := int32(-1)
	if maze[*originpoint].Dir != -1 {
		println("The Dev is an idiot, cant even keep track of a var. Ill go look for it...")
		for i, c := range maze {
			if c.Dir == -1 {
				*originpoint = int32(i)
				println("found it and updated it, its", i)
				break
			}
		}
	}

	var tocheck [4]bool = [4]bool{true,true,true,true}
	if int32(*originpoint) < size{
		tocheck[NORTH] = false
	}
	if int32(*originpoint) > int32(len(maze)) - size - 1 {
		tocheck[SOUTH] = false
	}
	if int32(*originpoint) % size == 0{
		tocheck[WEST] = false
	}
	if int32(*originpoint) % size == size - 1 {
		tocheck[EAST] = false
	}

	for i := 0; i < 100; i++ {
		Dir = rand.Int31n(4)

		if tocheck[Dir] {
			break
		}
	}

	if Dir == -1 {
		for i := 0; i < 3; i++ {
			if tocheck[i] {
				Dir = int32(i)
				break
			}
		}
		panic("couldnt find a spot to move to, is it a 1x1 maze?")
	}

	offset := int32(-1)

	if Dir == NORTH {
		offset = -size
	}else if Dir == EAST {
		offset = 1
	}else if Dir == SOUTH {
		offset = size
	}else {
		offset = -1
	}


	maze[*originpoint + offset].Dir = -1
	maze[*originpoint].Dir = Dir
	*originpoint = *originpoint + offset
}

func DrawHelpPoint(point int32, size int32, MazeRect rl.Rectangle, col color.RGBA) {
	x := (point % size) * CellSizePx + int32(MazeRect.X)
	y := (point / size) * CellSizePx + int32(MazeRect.Y)

	rl.DrawCircle(x + CellSizePx / 2, y + CellSizePx / 2, float32(CellSizePx) / 4, col)
	
}

func DrawHelpSquare(point int32, size int32, MazeRect rl.Rectangle, col color.RGBA) {
	x := (point % size) * CellSizePx + int32(MazeRect.X)
	y := (point / size) * CellSizePx + int32(MazeRect.Y)

	rl.DrawRectangle(x + CellSizePx / 4, y + CellSizePx / 4, CellSizePx / 2, CellSizePx / 2, col)
	
}

func DrawLinesBetweenPoints (PointList []int32, size int32, MazeRect rl.Rectangle,col color.RGBA) {

	if len(PointList) < 2 {
		return
	}

	for i := 0; i < len(PointList) - 1; i++ {
		x1 := (PointList[i] % size) * CellSizePx + int32(MazeRect.X) + CellSizePx / 2
		y1 := (PointList[i] / size) * CellSizePx + int32(MazeRect.Y) + CellSizePx / 2
		x2 := (PointList[i + 1] % size) * CellSizePx + int32(MazeRect.X) + CellSizePx / 2 
		y2 := (PointList[i + 1] / size) * CellSizePx + int32(MazeRect.Y) + CellSizePx / 2 

		rl.DrawLineEx(rl.Vector2{X: float32(x1),Y: float32(y1)}, rl.Vector2{X: float32(x2),Y: float32(y2)}, float32(CellSizePx) / 2, col)
	}
}

type stack []stackobj

type stackobj struct {
	ID int32 //corresponds to the cells ID
	rank int32 //favor of the cell beeing searched next cycle
	previos int32 //previous cell ID
}

func sortStack (s *stack) { //basic bubblesort
	changed := false
	for !changed {
		changed = false
		for i, obj := range *s { 
			if i == len(*s) - 1 {
				continue
			}

			if obj.rank > (*s)[i + 1].rank {
				tmp := obj
				(*s)[i] = (*s)[i + 1]
				(*s)[i + 1] = tmp
				changed = true
			}
		} 
	}
}

func IDInStack (s *stack, ID int32)  bool {
	for _, obj := range *s {
		if obj.ID == ID {
			return true
		}
	}
	return false
}

func FloodFillStep (s *stack, maze CellList, size int32) {
	var NewStackObjs stack

	for i, obj := range *s {
		c := maze[obj.ID]

		//already searched
		if obj.rank == 1 {
			continue
		}

		if !c.Walls[NORTH] && !IDInStack(s, c.ID - size) && !IDInStack(&NewStackObjs, c.ID - size) {
			NewStackObjs = append(NewStackObjs, stackobj{c.ID - size, -1, c.ID})
		}
		if !c.Walls[EAST] && !IDInStack(s, c.ID + 1) && !IDInStack(&NewStackObjs, c.ID + 1) {
			NewStackObjs = append(NewStackObjs, stackobj{c.ID + 1, -1, c.ID})
		}
		if !c.Walls[SOUTH] && !IDInStack(s, c.ID + size) && !IDInStack(&NewStackObjs, c.ID + size) {
			NewStackObjs = append(NewStackObjs, stackobj{c.ID + size, -1, c.ID})
		}
		if !c.Walls[WEST] && !IDInStack(s, c.ID - 1) {
			NewStackObjs = append(NewStackObjs, stackobj{c.ID - 1, -1, c.ID})
		}

		(*s)[i].rank = 1 //searched
	
	}

	*s = append(*s, NewStackObjs...)

}

func FindIDInStack (s *stack, ID int32) int32{
	for i, obj := range *s {
		if obj.ID == ID {
			return int32(i)
		}
	}

	println("ID not in stack, returning -1...")
	return -1
}

func FindSolutionStep (s *stack, path *[]int32) bool{ //its a pointer for drawing purposes
	latest := (*path)[len(*path) - 1]
	objIndex := FindIDInStack(s, latest)
	nextID := (*s)[objIndex].previos
	*path = append(*path, nextID)

	return nextID == -1
}

func main() {

	rl.SetConfigFlags(rl.FlagWindowHighdpi)


	var Resolution rl.Vector2
	
	Resolution.X = 1200
	Resolution.Y = 1200
	rl.InitWindow(int32(Resolution.X), int32(Resolution.Y), "a Maze ing")

	rl.SetTargetFPS(120)

	Mazesize := int32(min(Resolution.X, Resolution.Y) / float32(CellSizePx)) - 1

	MazeWidth := int32(min(Resolution.X, Resolution.Y)) - CellSizePx * 2 //buffer so it doesnt fill the edges
	MazeRect := rl.Rectangle{X: float32(CellSizePx) / 2, Y: float32(CellSizePx) / 2 + 40, Width: float32(MazeWidth), Height: float32(MazeWidth)}

	Maze := CreateMaze(Mazesize * Mazesize)
	Maze.Setup(ORIGIN, Mazesize)
	originpoint := int32(len(Maze) - 1)
	Maze.UpdateWalls(Mazesize, &originpoint)

	MainStack := stack{stackobj{0, -1, -1}}

	TopText := "Originshift - Press space to start/stop"
	BelowText := "Progress: "

	TextWidth := rl.MeasureText(TopText, 20)
	BelowTextWidth := rl.MeasureText(BelowText, 20)

	enabled := false
	help1 := true
	solve := false
	showingpath := false
	resetting := false

	waitTillReset := false
	resettimer := time.Now()

	iterations := 0
	progress := 0.0

	var Path []int32 = []int32{Mazesize * Mazesize - 1}

	speed := 0 // smol = fast

	defer rl.CloseWindow()

	start := time.Now()

	for !rl.WindowShouldClose(){

		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkGray)
		rl.DrawText(TopText, (int32(Resolution.X) - TextWidth) / 2, 0, 20, rl.RayWhite)
		rl.DrawText(BelowText, (int32(Resolution.X) - TextWidth) / 2, 20, 20, rl.RayWhite)
		DrawMaze(Maze, Mazesize, WALL, MazeRect)
		if help1 && !solve && !showingpath{
			DrawHelpPoint(originpoint, Mazesize, MazeRect, color.RGBA{255,0,0,255})
		}
		if solve && help1 && !showingpath{
			for _, obj := range MainStack {
				DrawHelpSquare(obj.ID, Mazesize, MazeRect, color.RGBA{255,0,0,255})
			}	
		}
		if showingpath {
			for _, ID := range Path[:len(Path) - 1] {
				DrawHelpPoint(ID, Mazesize, MazeRect, color.RGBA{255,255,255,255})
			}
			DrawLinesBetweenPoints(Path[:len(Path) - 1], Mazesize, MazeRect, color.RGBA{255,255,255,255})
		}

		rl.DrawRectangleLines((int32(Resolution.X) - TextWidth) / 2 + BelowTextWidth, 20, TextWidth - BelowTextWidth, 20, color.RGBA{255,255,255,255})
		rl.DrawRectangle((int32(Resolution.X) - TextWidth) / 2 + BelowTextWidth, 20, int32(float64(TextWidth - BelowTextWidth) * progress), 20, color.RGBA{255,255,255,255})
		rl.EndDrawing()


		if (time.Since(start) >= time.Duration(speed)) && enabled && !solve  && !waitTillReset{
			start = time.Now()
			for i := 0.0; i < math.Pow(float64(Mazesize), 3) / float64(FPS * 20); i++{
				Maze.OriginShiftStep(Mazesize, &originpoint)
			}
			iterations += int(math.Pow(float64(Mazesize), 3) / float64(FPS * 20)) + 1
			Maze.UpdateWalls(Mazesize, &originpoint)
			progress = float64(iterations) / float64(math.Pow(float64(Mazesize), 3))
		}

		if time.Since(start) >= time.Duration(speed) && solve && enabled{
			start = time.Now()
			if IDInStack(&MainStack, int32(len(Maze) - 1)) {
				waitTillReset = FindSolutionStep(&MainStack, &Path)
				showingpath = true
				if waitTillReset {
					resettimer = time.Now()
					solve = false

				}
			}else {
				FloodFillStep(&MainStack, Maze, Mazesize)
			}
		}

		if time.Since(resettimer) >= time.Duration(5 * math.Pow(10, 9)) && waitTillReset && enabled {
			resetting = true
			waitTillReset = false
		}

		if rl.IsKeyPressed(rl.KeySpace){
			enabled = !enabled
		}
		if rl.IsKeyPressed(rl.KeyOne){
			help1 = !help1
		}

		if rl.IsKeyDown(rl.KeyLeft) {
			speed += 1000000
		}
		if rl.IsKeyDown(rl.KeyRight) {
			speed -= 1000000
		}
		if (rl.IsKeyPressed(rl.KeyTwo) || iterations > int(Mazesize * Mazesize * Mazesize)) && !waitTillReset{
			solve = true
		}
		if rl.IsKeyPressed(rl.KeyR)  || resetting{
			enabled = true
			solve = false
			help1 = true
			showingpath = false
			resetting = false
			iterations = 0
			waitTillReset = false
			MainStack = stack{stackobj{0, -1, -1}}
			Path = []int32{Mazesize * Mazesize - 1}		
			Maze = CreateMaze(Mazesize * Mazesize)
			Maze.Setup(ORIGIN, Mazesize)
			originpoint = int32(len(Maze) - 1)
			Maze.UpdateWalls(Mazesize, &originpoint)
	
		}
		if rl.IsKeyPressed(rl.KeyThree) && !showingpath{
			for i := 0.0; i < math.Pow(float64(Mazesize), 3); i++ {
				Maze.OriginShiftStep(Mazesize, &originpoint)
			}
			Maze.UpdateWalls(Mazesize, &originpoint)
			iterations = int(Mazesize * Mazesize * Mazesize - 500)
		}

	}
}
