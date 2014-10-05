package main;

import "math/rand"
import "bytes"
import "fmt"
import "time"
import "os"
import "strconv"

const QSIZE = 5

//////////////////////////////////////////////////////////////////////////////
//
//  Conwya's Game of Life in Go Channels
//

//  Communication type definitions
type Link chan int
type Status struct {
  state bool
  count int
}

//  Intiializer type
type InitFn func(x, y int) bool

//////////////////////////////////////////////////////////////////////////////
//
//  The game board
//

type Nodes struct {
  width, height int
  nodes []Node
}

func NewNodes(width, height int, initFn InitFn) *Nodes {
  self := Nodes { width, height, make([]Node, width * height) }
  for index, _ := range(self.nodes) {
    (&self.nodes[index]).Init()
  }
  for index, _ := range(self.nodes) {
    x := index % self.width
    y := index / self.width
    (&self.nodes[index]).Start(&self, x, y, initFn(x, y))
  }
  return &self
}

func (self *Nodes) GetNode(x, y int) *Node {
  for x < 0 {
    x += self.width
  }
  for y < 0 {
    y += self.height
  }
  return &self.nodes[(y % self.height) * self.width + (x % self.width)]
}

func (self *Nodes)String() string {
  var buf bytes.Buffer

  for index, node := range(self.nodes) {
    if index > 0 && index % self.width == 0 {
      buf.WriteRune('\n')
    }
    stat := <-node.Out
    if stat.state {
      buf.WriteRune('*')
    } else {
      buf.WriteRune(' ')
    }
  }

  return buf.String()
}

//////////////////////////////////////////////////////////////////////////////
//
//  Each Node contains the channels used to communicate with its peers.
//

type Node struct {
  NW, N, NE, E, SE, S, SW, W Link
  Out chan Status
}

func (self *Node) Init() {
  self.NW  = make(Link, QSIZE)
  self.N   = make(Link, QSIZE)
  self.NE  = make(Link, QSIZE)
  self.E   = make(Link, QSIZE)
  self.SE  = make(Link, QSIZE)
  self.S   = make(Link, QSIZE)
  self.SW  = make(Link, QSIZE)
  self.W   = make(Link, QSIZE)
  self.Out = make(chan Status, QSIZE)
}

func (self *Node) Start(nodes *Nodes, x, y int, initState bool) {
  // Convert initState to state
  state := 0
  if initState {
    state = 1
  }

  go Run(
    // Initial state
    state,
    // Link to neighbors
    nodes.GetNode(x - 1, y - 1).SE,
    nodes.GetNode(x + 0, y - 1).S,
    nodes.GetNode(x + 1, y - 1).SW,
    nodes.GetNode(x + 1, y + 0).W,
    nodes.GetNode(x + 1, y + 1).NW,
    nodes.GetNode(x + 0, y + 1).N,
    nodes.GetNode(x - 1, y + 1).NE,
    nodes.GetNode(x - 1, y + 0).E,
    // Outputs that go to neighbors 
    self.NW,
    self.N,
    self.NE,
    self.E,
    self.SE,
    self.S,
    self.SW,
    self.W,
    // The output to the main controller.
    self.Out)
}

//////////////////////////////////////////////////////////////////////////////
//
//  Actual cmoputation engine
//

func Run (state int,
          nw, n, ne, e, se, s, sw, w <-chan int,
          NW, N, NE, E, SE, S, SW, W chan<- int,
          Out chan<- Status) {

  // Start computation.
  for {
    // First, write our status
    NW <- state
    N  <- state
    NE <- state
    E  <- state
    SE <- state
    S  <- state
    SW <- state
    W  <- state

    // Count the state values coming in
    count := <-nw + <-n + <-ne + <-e + <-se + <-s + <-sw + <-w

    // Report to output
    Out <- Status { state > 0, count }

    // Determine our new state
    if state == 0 {
      if count == 3 {
        state = 1
      }
    } else {
      if count != 2 && count != 3 {
        state = 0
      }
    }
  }
}

//////////////////////////////////////////////////////////////////////////////
//
//  Game Initializers
//

func Always(x, y int) bool {
  return true
}

func Never(x, y int) bool {
  return false
}

func Random(threshold float64) InitFn {
  return func(x, y int) bool {
    return rand.Float64() < threshold
  }
}

//////////////////////////////////////////////////////////////////////////////

const WAIT = 100000000    // Frame every 100 ms
const NANOS_PER_SEC = 1000000000

func main() {
  args := os.Args
  width := 78
  height := 20

  if len(args) >= 2 {
    if n, err := strconv.ParseInt(args[1], 10, 32); err == nil{
      width = (int)(n)
    }
  }

  if len(args) >= 3 {
    if n, err := strconv.ParseInt(args[2], 10, 32); err == nil {
      height = (int)(n)
    }
  }

  now := time.Now().UnixNano()
  start := now
  rand.Seed(now)
  game := NewNodes(width, height, Random(0.5))
  frame := (int64)(0)
  fmt.Print("\033[2J")    // Clear screen
  fmt.Print("\033[;H")  // Home cursor
  fmt.Println("Nodes:", len(game.nodes))
  until := now + WAIT
  for {
    frame++
    now = time.Now().UnixNano()
    data := game.String()
    if now >= until {
      fmt.Print("\033[;H")  // Home cursor
      fps := frame * NANOS_PER_SEC / (now - start)
      fmt.Println("\nFrame", frame, ",", fps, "FPS")
      fmt.Println(data)
      until = now + WAIT
    }
  }
}

