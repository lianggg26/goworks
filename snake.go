package main

import (
	"fmt"
	"math/rand"
  "time"
  "strconv"
  "os/exec"
  "os"
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

const (
	Win = iota
	Lose
)

var food_x int
var food_y int
var snake_x int
var snake_y int
var dir int

//当前蛇身长度
var body_len int = 0

//蛇身坐标+1之后所形成的数组
var body_plus[999] position_plus

const map_width int = 30
const map_height int = 10


var a[map_width][map_height] string

//初始化地图坐标
func initMap() {
	for x := 0; x < map_width; x++ {
	  for y := 0; y < map_height; y++ {
		  a[x][y] = "*";
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

//随机食物坐标
func randomFood() { 
  for{
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

//产生方向随机数
func initDir() {
  rand.Seed(time.Now().UnixNano())
  dir = rand.Intn(4)
}

//坐标边界检查
func checkPosition(x int, y int) bool {
  if x < 0 || x > map_width - 1 || y < 0 || y > map_height-1  {
     return false;   
  }
  return true;
}

//移动
func move(D int) {
  if debug == true {
    fmt.Println(strconv.Itoa(snake_x) + " , "+ strconv.Itoa(snake_y))
  }

  var origin_x int = snake_x
  var origin_y int = snake_y

  switch D {
    case Up:
      if debug {
        fmt.Println("蛇头上移：")
        
      }
      snake_y = snake_y - 1
    case Down:
      if debug {
        fmt.Println("蛇头下移：")
      }
      snake_y = snake_y + 1
    case Left:
      if debug {
        fmt.Println("蛇头左移：")
      }
      snake_x = snake_x - 1
    case Right:
      if debug {
        fmt.Println("蛇头右移：")
      }
      snake_x = snake_x + 1
  }
  
  //移动蛇身
  moveBody(origin_x, origin_y)

  if eatFood() == true {
      //原蛇头坐标变成蛇身， 蛇身长度+1
      var p position_plus
      p.body_x_plus = snake_x + 1
      p.body_y_plus = snake_y + 1
      body_plus[body_len] = p
      body_len = body_len + 1      
       
      //蛇头坐标更新为原食物坐标
      snake_x = food_x
      snake_y = food_y   
      
      //重新生成食物坐标
      randomFood()
  }
}

//蛇身移动时，蛇身坐标变化
func moveBody(x int, y int) {
  //蛇身变蛇头
  var p position_plus
  p.body_x_plus = x+1
  p.body_y_plus = y+1
  
  //蛇身坐标更新
  if body_len == 0 {
    return
  } else if body_len == 1 {
    body_plus[0] = p
  } else if body_len == 2 {
    body_plus[0] =body_plus[1]
    body_plus[1] = p
  } else {
    for i := 0; i < body_len-1; i++ {
      body_plus[i] = body_plus[i+1]
    }
    body_plus[body_len-1] = p
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
  initDir()
}

//移动后清洗蛇身
func clearBody(x int, y int) {
  a[x][y] = "*"
}

//判断是否吃掉食物
func eatFood() bool {
  if snake_x == food_x && snake_y == food_y {
    body_len = body_len+1
    return true
  }
  return false
}

//开始游戏
func startGame() {
  for {
    time.Sleep(500000000)
    printMap()
    clearBody(snake_x, snake_y)
    move(dir)
    if checkPosition(snake_x, snake_y) == false {
      fmt.Println("You lose!")
      break
    }
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
