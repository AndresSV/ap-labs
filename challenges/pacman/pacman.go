package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var configInfo config
var player sprite
var ghosts []*ghost
var maze []string
var totalpoints int
var numDots int
var lives = 3

func ClearScreen() {
	fmt.Print("\x1b[2J")
	MoveCursor(0, 0)
}

func MoveCursor(row, col int) {
	fmt.Printf("\x1b[%d;%df", row+1, col+1)
}

func BlueFont(text string) string {
	return "\x1b[44m" + text + "\x1b[0m"
}

func ghostGen(number string) string {
	switch number {
	case "1":
		return "######.## # G    # ##.######"
	case "2":
		return "######.## # GG   # ##.######"
	case "3":
		return "######.## # GGG  # ##.######"
	case "4":
		return "######.## # GGGG # ##.######"
	case "0":
		return "######.## #      # ##.######"
	default:
		return ""
	}
}

type sprite struct {
	row      int
	col      int
	startRow int
	startCol int
}

type ghost struct {
	position sprite
	status   GhostStatus
}

type GhostStatus string

const (
	GhostStatusNormal GhostStatus = "Normal"
	GhostStatusBlue   GhostStatus = "Blue"
)

var ghostsStatusMx sync.RWMutex
var pillMx sync.Mutex

type config struct {
	Player           string
	Ghost            string
	GhostBlue        string
	Wall             string
	Dot              string
	Pill             string
	Death            string
	Space            string
	PillDurationSecs time.Duration
}

func getIcons() error {
	configInfo.Player = "😃"
	configInfo.Ghost = "ᗣ"
	configInfo.GhostBlue = "👾"
	configInfo.Wall = "  "
	configInfo.Dot = "· "
	configInfo.Pill = "🔆"
	configInfo.Death = "💀"
	configInfo.Space = "  "
	configInfo.PillDurationSecs = 12
	return nil
}

func createMap(number string) error {
	upperMaze, err := os.Open("./maze_elements/upperMaze.txt")
	if err != nil {
		return err
	}
	defer upperMaze.Close()

	scanner1 := bufio.NewScanner(upperMaze)
	for scanner1.Scan() {
		line := scanner1.Text()
		maze = append(maze, line)
	}

	ghostsLine := ghostGen(number)
	maze = append(maze, ghostsLine)

	println(ghostsLine)

	lowerMaze, err := os.Open("./maze_elements/lowerMaze.txt")
	if err != nil {
		return err
	}
	defer lowerMaze.Close()

	scanner2 := bufio.NewScanner(lowerMaze)
	for scanner2.Scan() {
		line := scanner2.Text()
		maze = append(maze, line)
	}

	for row, line := range maze {
		for col, char := range line {
			switch char {
			case 'P':
				player = sprite{row, col, row, col}
			case 'G':
				ghosts = append(ghosts, &ghost{sprite{row, col, row, col}, GhostStatusNormal})
			case '.':
				numDots++
			}
		}
	}

	return nil
}

func moveCursor(row, col int) {
	MoveCursor(row, col*2)
}

func printScreen() {
	ClearScreen()
	for _, line := range maze {
		for _, chr := range line {
			switch chr {
			case '#':
				fmt.Print(BlueFont(configInfo.Wall))
			case '.':
				fmt.Print(configInfo.Dot)
			case 'X':
				fmt.Print(configInfo.Pill)
			default:
				fmt.Print(configInfo.Space)
			}
		}
		fmt.Println()
	}

	moveCursor(player.row, player.col)
	fmt.Print(configInfo.Player)

	ghostsStatusMx.RLock()
	for _, g := range ghosts {
		moveCursor(g.position.row, g.position.col)
		if g.status == GhostStatusNormal {
			fmt.Printf(configInfo.Ghost)
		} else if g.status == GhostStatusBlue {
			fmt.Printf(configInfo.GhostBlue)
		}
	}
	ghostsStatusMx.RUnlock()

	moveCursor(len(maze)+1, 0)

	livesRemaining := strconv.Itoa(lives) //converts lives int to a string
	livesRemaining = getHPImage()

	fmt.Println("Score:", totalpoints, "\tLives:", livesRemaining)
}

func getHPImage() string {
	buf := bytes.Buffer{}
	for i := lives; i > 0; i-- {
		buf.WriteString(configInfo.Player)
	}
	return buf.String()
}

func inputCont() (string, error) {
	buffer := make([]byte, 100)

	cnt, err := os.Stdin.Read(buffer)
	if err != nil {
		return "", err
	}

	if cnt == 1 && buffer[0] == 0x1b {
		return "ESC", nil
	} else if cnt >= 3 {
		if buffer[0] == 0x1b && buffer[1] == '[' {
			switch buffer[2] {
			case 'A':
				return "UP", nil
			case 'B':
				return "DOWN", nil
			case 'C':
				return "RIGHT", nil
			case 'D':
				return "LEFT", nil
			}
		}
	}

	return "", nil
}

func makeMove(oldRow, oldCol int, dir string) (newRow, newCol int) {
	newRow, newCol = oldRow, oldCol
	switch dir {
	case "UP":
		newRow = newRow - 1
		if newRow < 0 {
			newRow = len(maze) - 1
		}
	case "DOWN":
		newRow = newRow + 1
		if newRow == len(maze)-1 {
			newRow = 0
		}
	case "RIGHT":
		newCol = newCol + 1
		if newCol == len(maze[0]) {
			newCol = 0
		}
	case "LEFT":
		newCol = newCol - 1
		if newCol < 0 {
			newCol = len(maze[0]) - 1
		}
	}
	if maze[newRow][newCol] == '#' {
		newRow = oldRow
		newCol = oldCol
	}
	return
}

func movePacman(dir string) {
	player.row, player.col = makeMove(player.row, player.col, dir)

	removeDot := func(row, col int) {
		maze[row] = maze[row][0:col] + " " + maze[row][col+1:]
	}

	switch maze[player.row][player.col] {
	case '.':
		numDots--
		totalpoints++
		removeDot(player.row, player.col)
	case 'X':
		totalpoints += 10
		removeDot(player.row, player.col)
		go processPill()
	}
}

func updateGhosts(ghosts []*ghost, ghostStatus GhostStatus) {
	ghostsStatusMx.Lock()
	defer ghostsStatusMx.Unlock()
	for _, g := range ghosts {
		g.status = ghostStatus
	}
}

var pillTimer *time.Timer

func processPill() {
	pillMx.Lock()
	updateGhosts(ghosts, GhostStatusBlue)
	if pillTimer != nil {
		pillTimer.Stop()
	}
	pillTimer = time.NewTimer(time.Second * configInfo.PillDurationSecs)
	pillMx.Unlock()
	<-pillTimer.C
	pillMx.Lock()
	pillTimer.Stop()
	updateGhosts(ghosts, GhostStatusNormal)
	pillMx.Unlock()
}

func drawDirection() string {
	dir := rand.Intn(4)
	move := map[int]string{
		0: "UP",
		1: "DOWN",
		2: "RIGHT",
		3: "LEFT",
	}
	return move[dir]
}

func moveGhosts() {
	for _, g := range ghosts {
		dir := drawDirection()
		g.position.row, g.position.col = makeMove(g.position.row, g.position.col, dir)
	}
}

func init() {
	cbTerm := exec.Command("stty", "cbreak", "-echo")
	cbTerm.Stdin = os.Stdin

	err := cbTerm.Run()
	if err != nil {
		log.Fatalln("unable to activate cbreak mode:", err)
	}
}

func main() {

	validGhosts := false

	var number string

	for !validGhosts {
		fmt.Print("Input the number ghost between [0-4]: ")
		reader := bufio.NewReader(os.Stdin)
		number, _ = reader.ReadString('\n')
		fmt.Println("number: " + number)
		value, err := strconv.Atoi(strings.TrimSuffix(number, "\n"))
		if err != nil {
			fmt.Println("Not a number")
		}
		if value < 0 || value > 4 {
			log.Println("Expected.")
			log.Println("Number between [0-4].")
		} else {
			validGhosts = true
		}
	}

	err := createMap(strings.TrimSuffix(number, "\n"))
	if err != nil {
		log.Println("failed to load maze:", err)
		return
	}

	err = getIcons()
	if err != nil {
		log.Println("failed to load configuration:", err)
		return
	}

	input := make(chan string)
	go func(ch chan<- string) {
		for {
			input, err := inputCont()
			if err != nil {
				log.Print("error reading input:", err)
				ch <- "ESC"
			}
			ch <- input
		}
	}(input)

	for {
		select {
		case inp := <-input:
			if inp == "ESC" {
				lives = 0
			}
			movePacman(inp)
		default:
		}

		moveGhosts()

		for _, g := range ghosts {
			if player.row == g.position.row && player.col == g.position.col {
				ghostsStatusMx.RLock()
				if g.status == GhostStatusNormal {
					lives = lives - 1
					if lives != 0 {
						moveCursor(player.row, player.col)
						fmt.Print(configInfo.Death)
						moveCursor(len(maze)+2, 0)
						ghostsStatusMx.RUnlock()
						updateGhosts(ghosts, GhostStatusNormal)
						time.Sleep(1000 * time.Millisecond)
						player.row, player.col = player.startRow, player.startCol
					}
				} else if g.status == GhostStatusBlue {
					ghostsStatusMx.RUnlock()
					updateGhosts([]*ghost{g}, GhostStatusNormal)
					g.position.row, g.position.col = g.position.startRow, g.position.startCol
				}
			}
		}

		printScreen()

		if numDots == 0 || lives <= 0 {
			if lives == 0 {
				moveCursor(player.row, player.col)
				fmt.Print(configInfo.Death)
				moveCursor(player.startRow, player.startCol-1)
				fmt.Print("GAME OVER")
				moveCursor(len(maze)+2, 0)
			}
			break
		}

		time.Sleep(200 * time.Millisecond)
	}
}
