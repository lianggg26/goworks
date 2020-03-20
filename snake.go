package main

/*
#include <windows.h>
int KeyDown(int key) {
    // 数据兼容：因为 GetKeyState() 不接受小写字母
    if (key > 96 && key < 123)  key -= 32;

    // 获取按下的键的状态，返回 0 则表示没按，其他情况表示按了
    return (GetKeyState(key) < 0) ? 1 : 0;
}
*/
import "C"
import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type D int

const (
	Up = iota
	Down
	Left
	Right
)

//蛇身坐标数组，存的是各自+1之后的值，取出来之后-1才是真正的坐标值
type position_plus struct {
	body_x_plus int
	body_y_plus int
}

var debug bool = false

type W int

var w bool = true
var start bool = false
var conflict bool = false

const (
	Win = iota
	Lose
)

var food_x int
var food_y int
var snake_x int
var snake_y int
var dir int
var move_times int

//当前蛇身长度
var body_len int = 0

//蛇身坐标+1之后所形成的数组
var body_plus [999]position_plus

const map_width int = 30
const map_height int = 10

var a [map_width][map_height]string

//获取当前时间
func currentTime() string {
	t := time.Now()

	return t.Local().Format(time.UnixDate)
}

//写入日志
func writeLog(content string) {
	f, err := os.OpenFile("log.txt", os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("文件处理失败")
	} else {
		n, _ := f.Seek(0, 2)
		_, err = f.WriteAt([]byte(content), n)
	}

	defer f.Close()
}

//初始化地图坐标
func initMap() {
	for x := 0; x < map_width; x++ {
		for y := 0; y < map_height; y++ {
			a[x][y] = "*"
		}
	}
}

//命令行绘制地图
func printMap() {
	refreshMap()
	for y := 0; y < map_height; y++ {
		fmt.Println()
		for x := 0; x < map_width; x++ {
			//蛇头绘制
			if x == snake_x && y == snake_y {
				a[x][y] = "X"
			}

			//食物绘制
			if x == food_x && y == food_y {
				a[x][y] = "@"
			}

			//蛇身绘制
			for i := 0; i < body_len; i++ {
				if x == body_plus[i].body_x_plus-1 && y == body_plus[i].body_y_plus-1 {
					a[x][y] = "O"
				}
			}

			fmt.Print(a[x][y])
		}
	}
	fmt.Println("\n\n")
}

func randomFood() {
	for {
		rand.Seed(time.Now().UnixNano())
		food_x = rand.Intn(map_width)
		food_y = rand.Intn(map_height)

		if foodOK(food_x, food_y) {
			break
		}
	}

}

//食物坐标不与蛇头和蛇身重合
func foodOK(x int, y int) bool {

	//与蛇头重合
	if x == snake_x && y == snake_y {
		return false
	}

	//与蛇身重合
	for i := 0; i < body_len; i++ {
		if x == body_plus[i].body_x_plus-1 && y == body_plus[i].body_y_plus-1 {
			return false
		}
	}

	return true
}

//随机蛇头随机
func initSnake() {
	for {
		rand.Seed(time.Now().UnixNano())
		snake_x = rand.Intn(map_width)
		snake_y = rand.Intn(map_height)

		if snake_x != food_x || snake_y != food_y {
			break
		}
	}
	a[snake_x][snake_y] = "X"
}

//调试产生方向随机数
func initDir() {
	rand.Seed(time.Now().UnixNano())
	dir = rand.Intn(4)
}

//输赢检查
func checkContinue() bool {
	if conflict == true {
		return false
	}

	return true
}

//碰撞：失败条件
func conflictCheck(x int, y int) {
	if x < 0 || x > map_width-1 || y < 0 || y > map_height-1 {
		conflict = true
	}

	for i := 0; i < body_len; i++ {
		if x == body_plus[i].body_x_plus-1 && y == body_plus[i].body_y_plus {
			conflict = true
		}
	}
}

//移动
func move(D int) {
	if debug == true {
		writeLog("\n" + currentTime() + "移动前蛇头坐标-->" + strconv.Itoa(snake_x) + " , " + strconv.Itoa(snake_y))
	}

	var origin_x int = snake_x
	var origin_y int = snake_y

	switch D {
	case Up:
		if debug {
			writeLog("\n" + currentTime() + "蛇头上移：")

		}
		snake_y = snake_y - 1

	case Down:
		if debug {
			writeLog("\n" + currentTime() + "蛇头下移：")
		}
		snake_y = snake_y + 1
	case Left:
		if debug {
			writeLog("\n" + currentTime() + "蛇头左移：")
		}
		snake_x = snake_x - 1
	case Right:
		if debug {
			writeLog("\n" + currentTime() + "蛇头右移：")
		}
		snake_x = snake_x + 1
	}

	//移动蛇身
	moveBody(origin_x, origin_y)

	if eatFood() == true {
		if debug == true {
			writeLog("\n" + currentTime() + "要进食前蛇头坐标-->" + strconv.Itoa(snake_x) + " , " + strconv.Itoa(snake_y))
		}

		//蛇头坐标更新为原食物坐标
		snake_x = food_x
		snake_y = food_y
		if debug == true {
			writeLog("\n" + currentTime() + "新蛇头坐标：" + strconv.Itoa(snake_x) + " , " + strconv.Itoa(snake_y))
		}

		//原蛇头坐标变成蛇身， 蛇身长度+1
		var p position_plus
		if debug == true {
			writeLog("\n" + currentTime() + "初始化蛇身坐标元素(含+1)：" + strconv.Itoa(p.body_x_plus) + " , " + strconv.Itoa(p.body_y_plus))
		}

		p.body_x_plus = snake_x + 1
		p.body_y_plus = snake_y + 1

		if debug == true {
			writeLog("\n" + currentTime() + "蛇头坐标+1作为新蛇身坐标+1：" + strconv.Itoa(p.body_x_plus) + " , " + strconv.Itoa(p.body_y_plus))
		}

		if debug == true {
			writeLog("\n" + currentTime() + "原本蛇身长度：" + strconv.Itoa(body_len))
		}
		body_plus[body_len] = p
		if debug == true {
			writeLog("\n" + currentTime() + "新蛇身元素：" + strconv.Itoa(body_len) + " , " + strconv.Itoa(body_plus[body_len].body_x_plus) + " , " + strconv.Itoa(body_plus[body_len].body_y_plus))
		}
		body_len = body_len + 1
		if debug == true {
			writeLog("\n" + currentTime() + "新蛇身长度：" + strconv.Itoa(body_len))
		}

		conflictCheck(snake_x, snake_y)

		//继续重新生成食物坐标
		randomFood()

		if debug == true {
			writeLog("\n" + currentTime() + "新食物坐标：" + strconv.Itoa(food_x) + " , " + strconv.Itoa(food_y))
		}
	}
}

//蛇身移动时，蛇身坐标变化
func moveBody(x int, y int) {

	move_times = move_times + 1

	if debug == true {
		writeLog("\n" + currentTime() + "移动次数：" + strconv.Itoa(move_times))
	}

	if debug == true {
		writeLog("\n" + currentTime() + "蛇身长度：" + strconv.Itoa(body_len))
	}

	//蛇身变蛇头
	var p position_plus
	p.body_x_plus = x + 1
	p.body_y_plus = y + 1

	//蛇身坐标更新
	if body_len == 0 {
		if debug == true {
			writeLog("\n" + currentTime() + "移动次数：" + strconv.Itoa(move_times) + "  " + "蛇身长度： " + strconv.Itoa(body_len))
		}
		return
	} else if body_len == 1 {
		if debug == true {
			writeLog("\n" + currentTime() + "移动次数：" + strconv.Itoa(move_times) + "  " + "蛇身长度： " + strconv.Itoa(body_len))
		}
		clearBody(body_plus[0].body_x_plus-1, body_plus[0].body_y_plus-1)

		body_plus[0] = p
		if debug == true {
			writeLog("\n" + currentTime() + "最新蛇身坐标：" + strconv.Itoa(body_plus[body_len].body_x_plus-1) + strconv.Itoa(body_plus[body_len].body_y_plus-1))
		}
	} else if body_len == 2 {
		if debug == true {
			writeLog("\n" + currentTime() + "移动次数：" + strconv.Itoa(move_times) + "  " + "蛇身长度： " + strconv.Itoa(body_len))
		}

		for i := 0; i < body_len; i++ {
			clearBody(body_plus[i].body_x_plus-1, body_plus[i].body_y_plus-1)
		}
		body_plus[0] = body_plus[1]
		body_plus[1] = p
		if debug == true {
			writeLog("\n" + currentTime() + "最新蛇身坐标：" + strconv.Itoa(body_plus[body_len].body_x_plus-1) + strconv.Itoa(body_plus[body_len].body_y_plus-1))
		}
	} else if body_len > 2 {
		if debug == true {
			writeLog("\n" + currentTime() + "移动次数：" + strconv.Itoa(move_times) + "  " + "蛇身长度： " + strconv.Itoa(body_len))
		}

		for l := 0; l < body_len; l++ {
			clearBody(body_plus[l].body_x_plus-1, body_plus[l].body_y_plus-1)
		}

		for i := 0; i < body_len-1; i++ {
			body_plus[i] = body_plus[i+1]
		}
		body_plus[body_len-1] = p

		if debug == true {
			writeLog("\n" + currentTime() + "最新蛇身坐标：" + strconv.Itoa(body_plus[body_len-1].body_x_plus-1) + strconv.Itoa(body_plus[body_len-1].body_y_plus-1))
		}
	}
}

//刷新地图
func refreshMap() {
	command := exec.Command("cmd", "/c", "cls")
	command.Stdout = os.Stdout
	command.Run()
}

//初始化游戏
func initGame() {
	initMap()
	randomFood()
	initSnake()
	if start == false {
		initDir()
	}
}

//移动后清洗蛇身
func clearBody(x int, y int) {
	a[x][y] = "*"
}

//判断是否吃掉食物
func eatFood() bool {
	if snake_x == food_x && snake_y == food_y {
		return true
	}
	return false
}

//监听方向
func watchKey() {
	for {
		if int(C.KeyDown('w')) == 1 {
			dir = 0
		}

		if int(C.KeyDown('s')) == 1 {
			dir = 1
		}

		if int(C.KeyDown('a')) == 1 {
			dir = 2
		}

		if int(C.KeyDown('d')) == 1 {
			dir = 3
		}
	}
}

//开始游戏
func startGame() {
	go watchKey()

	for {
		time.Sleep(500000000)
		printMap()
		if checkContinue() == true {
			clearBody(snake_x, snake_y)
		}
		move(dir)
		if checkContinue() == false {
			fmt.Println("You lose!")
			break
		}

		start = true
		printMap()
	}
}

//程序入口
func main() {
	if len(os.Args) == 2 {
		if os.Args[1] == "debug" {
			debug = true
		}
	}
	initGame()
	startGame()
}
