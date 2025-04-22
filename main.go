package main

// ~/go/bin/tinygo flash -target=pico -monitor main.go
import (
	"fmt"
	"machine"
	"machine/usb/hid/joystick"
	"time"
)

type Button struct {
	JoySticdID int
	Pin        machine.Pin
	Label      string
}

type State struct {
	Buttons  []Button
	Current  []bool
	Previous []bool
}

func NewState(buttons []Button) *State {
	for _, b := range buttons {
		b.Pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}
	return &State{
		Buttons:  buttons,
		Current:  make([]bool, len(buttons)),
		Previous: make([]bool, len(buttons)),
	}
}

func (s *State) Load() {
	s.Previous, s.Current = s.Current, s.Previous
	for i := range s.Buttons {
		s.Current[i] = !s.Buttons[i].Pin.Get()
	}
}

func (s *State) Changed() bool {
	for i := range s.Buttons {
		if s.Current[i] != s.Previous[i] {
			return true
		}
	}
	return false
}

// func (s *State) Send(js *joystick.Joystick) {
// 	for i := range s.Buttons {
// 		js.SetButton(s.Buttons[i].JoySticdID, s.Current[i])
// 	}
// 	js.SendState()
// }

func main() {
	js := joystick.Port()
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// TODO hats!
	state := NewState([]Button{
		{0, machine.GP0, "UP"},
		{1, machine.GP1, "RIGHT"},
		{2, machine.GP2, "LEFT"},
		{3, machine.GP3, "RIGHT"},
		{4, machine.GP4, "A"},
		{5, machine.GP5, "B"},
		{6, machine.GP6, "X"},
		{7, machine.GP7, "Y"},
		{8, machine.GP8, "L"},
		{9, machine.GP9, "R"},
		{10, machine.GP10, "SELECT"},
		{11, machine.GP11, "START"},
	})

	for {
		time.Sleep(5 * time.Millisecond)
		state.Load()
		if state.Changed() {
			for i := range state.Buttons {
				if state.Current[i] {
					fmt.Printf("Button %s\n", state.Buttons[i].Label)
				}
				js.SetButton(state.Buttons[i].JoySticdID, state.Current[i])
			}
			js.SendState()
		}
	}
}
