package main

import (
	"encoding/json"
	"github.com/ChromeTemp/Popup"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/natefinch/npipe"
	"io"
	"time"
)

type Executable struct {
	IsLauncher bool                   `json:"is_launcher"`
	Name       string                 `json:"name"`
	Os         string                 `json:"os"`
	Arguments  string                 `json:"arguments,omitempty"`
	Extra      map[string]interface{} `json:"-"`
}

type GameStore struct {
	Executables []Executable           `json:"executables,omitempty"`
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Aliases     []string               `json:"aliases,omitempty"`
	Extra       map[string]interface{} `json:"-"`
}

func title() ([]byte, error) {
	var recV []byte
	conn, err := npipe.DialTimeout(`\\.\pipe\discord-fake-activity-embedded`, time.Second*5)
	if err != nil {
		return recV, err
	}
	defer func(conn *npipe.PipeConn) {
		_ = conn.Close()
	}(conn)
	recV, err = io.ReadAll(conn)
	if err != nil {
		return recV, err
	}
	return recV, nil
}

func main() {
	t, err := title()
	if err != nil {
		Popup.Alert("Error", "Parent process has died unexpectedly")
		return
	}
	var game GameStore
	err = json.Unmarshal(t, &game)
	if err != nil {
		Popup.Alert("Error", "Unable to read game information, try again or open an issue")
		return
	}

	rl.InitWindow(500, 200, game.Name)
	rl.SetTargetFPS(2)
	txt := "To end close this window or press the ESCAPE button"
	x := rl.MeasureText(txt, 10)
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawText(txt, 500/2-x/2, 95, 10, rl.RayWhite)
		rl.EndDrawing()
	}
	rl.CloseWindow()
}
