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
	LastInput             string
}

type CleaningMethod int

const (
	Neutral CleaningMethod = iota
	LastInputPriority
)

var kb keybd_event.KeyBonding

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
	dir := Direction{}
	if getKeyState(VK_A) {
		dir.Left = true
		dir.LastInput = "Left"
	}
	if getKeyState(VK_D) {
		dir.Right = true
		dir.LastInput = "Right"
	}
	if getKeyState(VK_W) {
		dir.Up = true
		dir.LastInput = "Up"
	}
	if getKeyState(VK_S) {
		dir.Down = true
		dir.LastInput = "Down"
	}
	return dir
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
			cleaned.Up = true
			cleaned.Down = false
		}
	case LastInputPriority:
		if dir.Left && dir.Right {
			if dir.LastInput == "Left" {
				cleaned.Right = false
			} else {
				cleaned.Left = false
			}
		}
		if dir.Up && dir.Down {
			if dir.LastInput == "Up" {
				cleaned.Down = false
			} else {
				cleaned.Up = false
			}
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
	fmt.Println("Press 'N' for Neutral cleaning, 'L' for Last Input Priority cleaning")

	cleaningMethod := Neutral

	for {
		if getKeyState('N') {
			cleaningMethod = Neutral
			fmt.Println("\nSwitched to Neutral cleaning")
		}
		if getKeyState('L') {
			cleaningMethod = LastInputPriority
			fmt.Println("\nSwitched to Last Input Priority cleaning")
		}

		rawDir := getDirection()
		cleanedDir := cleanSOCD(rawDir, cleaningMethod)

		fmt.Printf("\rRaw: %-20s | Cleaned: %-20s | Method: %-20v",
			directionToString(rawDir), directionToString(cleanedDir),
			map[CleaningMethod]string{Neutral: "Neutral", LastInputPriority: "Last Input Priority"}[cleaningMethod])

		simulateKeyPress(cleanedDir)

		time.Sleep(16 * time.Millisecond) // ~60 fps
	}
}
