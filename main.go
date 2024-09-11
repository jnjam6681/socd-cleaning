package main

import (
	"fmt"
	"time"

	"github.com/TheTitanrain/w32"
	"github.com/micmonay/keybd_event"
)

const (
	VK_A = 0x41
	VK_D = 0x44
	VK_W = 0x57
	VK_S = 0x53
)

type Direction struct {
	Left, Right, Up, Down bool
}

type CleaningMethod int

const (
	Neutral CleaningMethod = iota
	PriorityDirection
	Alternating
)

var (
	kb               keybd_event.KeyBonding
	alternatingState bool
)

func init() {
	var err error
	kb, err = keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}
}

func getKeyState(key int) bool {
	return w32.GetAsyncKeyState(key)&0x8000 != 0
}

func getDirection() Direction {
	return Direction{
		Left:  getKeyState(VK_A),
		Right: getKeyState(VK_D),
		Up:    getKeyState(VK_W),
		Down:  getKeyState(VK_S),
	}
}

func cleanSOCD(dir Direction, method CleaningMethod) Direction {
	cleaned := dir

	switch method {
	case Neutral:
		if dir.Left && dir.Right {
			cleaned.Left = false
			cleaned.Right = false
		}
		if dir.Up && dir.Down {
			cleaned.Up = false
			cleaned.Down = false
		}
	case PriorityDirection:
		if dir.Left && dir.Right {
			cleaned.Right = false // Left has priority
		}
		if dir.Up && dir.Down {
			cleaned.Down = false // Up has priority
		}
	case Alternating:
		if dir.Left && dir.Right {
			alternatingState = !alternatingState
			cleaned.Left = alternatingState
			cleaned.Right = !alternatingState
		}
		if dir.Up && dir.Down {
			alternatingState = !alternatingState
			cleaned.Up = alternatingState
			cleaned.Down = !alternatingState
		}
	}

	return cleaned
}

func simulateKeyPress(dir Direction) {
	kb.Clear()
	if dir.Left {
		kb.AddKey(keybd_event.VK_LEFT)
	}
	if dir.Right {
		kb.AddKey(keybd_event.VK_RIGHT)
	}
	if dir.Up {
		kb.AddKey(keybd_event.VK_UP)
	}
	if dir.Down {
		kb.AddKey(keybd_event.VK_DOWN)
	}
	kb.Press()
	time.Sleep(10 * time.Millisecond)
	kb.Release()
}

func directionToString(dir Direction) string {
	result := ""
	if dir.Left {
		result += "Left "
	}
	if dir.Right {
		result += "Right "
	}
	if dir.Up {
		result += "Up "
	}
	if dir.Down {
		result += "Down "
	}
	if result == "" {
		result = "Neutral"
	}
	return result
}

func main() {
	fmt.Println("SOCD Cleaner started. Use A(Left), D(Right), W(Up), S(Down). Press Ctrl+C to exit.")
	fmt.Println("Press 'N' for Neutral, 'P' for Priority Direction, 'A' for Alternating")

	cleaningMethod := Neutral

	for {
		if getKeyState('N') {
			cleaningMethod = Neutral
			fmt.Println("\nSwitched to Neutral cleaning")
		}
		if getKeyState('P') {
			cleaningMethod = PriorityDirection
			fmt.Println("\nSwitched to Priority Direction cleaning")
		}
		if getKeyState('L') {
			cleaningMethod = Alternating
			fmt.Println("\nSwitched to Alternating cleaning")
		}

		rawDir := getDirection()
		cleanedDir := cleanSOCD(rawDir, cleaningMethod)

		fmt.Printf("\rRaw: %-20s | Cleaned: %-20s | Method: %-20v",
			directionToString(rawDir), directionToString(cleanedDir),
			map[CleaningMethod]string{Neutral: "Neutral", PriorityDirection: "Priority Direction", Alternating: "Alternating"}[cleaningMethod])

		simulateKeyPress(cleanedDir)

		time.Sleep(16 * time.Millisecond) // ~60 fps
	}
}
