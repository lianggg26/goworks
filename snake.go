package main

/*
#include <windows.h>
int KeyDown(int key) {
    if (key > 96 && key < 123)  key -= 32;
    return (GetKeyState(key) < 0) ? 1 : 0;
}
*/
import "C"
import (
	"fmt"
	"math/ran"
	"os"
	"os/eec"
	"strconv"
	"tim"
)

var w it
vr h int
var x int
var y int
var foodX int
var foodY int
var bodyX [200]int
var bodyY [200]int
var nBody int
var gameOver bool
var score int
var dir int
var previousX, previousY, nextX, nextY int

type D int

const (
	STOP = iota
	UP
	DOWN
	LEFT
	RIGHT
)

func defaultSetup() {
	w = 50
	h = 20
}

func foodOK(foodX int, foodY int) bool {
	if foodX == x && foodY == y {
		return false
	}
	return true
}

func generateFood() {
	for {
		rand.Seed(time.Now().UnixNano())
		foodX = rand.Intn(w)
		foodY = rand.Intn(h)

		if foodOK(foodX, foodY) {
			break
		}
	}

}

func setup() {
	fmt.Println("set up")
	if len(os.Args) == 1 {
		defaultSetup()
		fmt.Println("width: " + strconv.Itoa(w) + "  height: " + strconv.Itoa(h))
	}

	if len(os.Args) == 3 {
		w, _ = strconv.Atoi(os.Args[1])
		h, _ = strconv.Atoi(os.Args[2])
		fmt.Println("width: " + strconv.Itoa(w) + "  height: " + strconv.Itoa(h))
	}

	//initiate snake and food position
	rand.Seed(time.Now().UnixNano())
	x = rand.Intn(w)
	y = rand.Intn(h)
	generateFood()

	//initiate direction and body length
	dir = STOP
	nBody = 0
}

func isBody(i int, j int) bool {
	var x = i - 1
	var y = j

	for p := 0; p < nBody; p++ {
		if x == bodyX[p] && y == bodyY[p] {
			return true
		}
	}

	return false
}

func draw() {
	for i := 0; i < w+2; i++ {
		fmt.Print("#")
	}

	fmt.Println()

	for j := 0; j < h; j++ {
		for i := 0; i < w+2; i++ {
			if i == 0 || i == w+1 {
				fmt.Print("#") // wall
			} else {
				if i-1 == x && j == y { //snake or food or body or blank
					fmt.Print("X")
				} else if i-1 == foodX && j == foodY {
					fmt.Print("@")
				} else if isBody(i, j) == true {
					fmt.Print("O")
				} else {
					fmt.Print(" ")
				}
			}
		}

		fmt.Println()
	}

	for i := 0; i < w+2; i++ {
		fmt.Print("#")
	}

	fmt.Println()
	fmt.Println("Score: " + strconv.Itoa(score))
	fmt.Println("Bodylength: " + strconv.Itoa(nBody))
}

func inputWatcher() {
	for {
		if int(C.KeyDown('w')) == 1 {
			dir = UP
		} else if int(C.KeyDown('s')) == 1 {
			dir = DOWN
		} else if int(C.KeyDown('a')) == 1 {
			dir = LEFT
		} else if int(C.KeyDown('d')) == 1 {
			dir = RIGHT
		}
	}
}

func rules(D int) {
	if nBody > 0 {
		previousX = x
		previousY = y
		bodyX[0] = previousX
		bodyY[0] = previousY

		for i := 1; i < nBody; i++ {
			nextX = bodyX[i]
			nextY = bodyY[i]
			bodyX[i] = previousX
			bodyY[i] = previousY
			previousX = nextX
			previousY = nextY
		}
	}

	switch D {
	case UP:
		y--
	case DOWN:
		y++
	case LEFT:
		x--
	case RIGHT:
		x++
	}

	if x == foodX && y == foodY { //eatFood
		generateFood()
		nBody++
		score = score + 10
	}

	if x < 0 || x > w-1 || y < 0 || y > h-1 { //set Gameover
		gameOver = true
	}

}

func main() {
	setup()
	go inputWatcher()
	for {
		time.Sleep(500000000)
		command := exec.Command("cmd", "/c", "cls")
		command.Stdout = os.Stdout
		command.Run()
		rules(dir)
		draw()

		if gameOver == true {
			fmt.Println("Game Over！")
			break
		}
	}

}
