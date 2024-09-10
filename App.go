package main

import (
	_ "embed"
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ChromeTemp/Popup"
	"github.com/natefinch/npipe"
	"net/http"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
)

const AppLink = "https://github.com/RegularRabbit05/DiscordFakeGamePresence"

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

func openAppLink() *exec.Cmd {
	cmd := exec.Command("cmd", "/c", "start", AppLink)
	_ = cmd.Start()
	return cmd
}

func downloadActivities() []GameStore {
	resp, err := http.Get("https://discord.com/api/v9/applications/detectable")
	if err != nil {
		Popup.Alert("Check your connection / is Discord down?", "Unable to download game store from Discord's servers")
		os.Exit(1)
	}
	var store []GameStore
	err = json.NewDecoder(resp.Body).Decode(&store)
	if err != nil {
		Popup.Alert(
			"Unable to decode game store",
			"You might be getting rate limited or check if a new version of this software is available at "+AppLink,
		)
		_ = openAppLink().Wait()
		os.Exit(1)
	}
	sort.Slice(store, func(i, j int) bool {
		return store[i].Name < store[j].Name
	})
	return store
}

//go:embed embedded.exe
var embeddedLauncher []byte

func launchGame(game GameStore) {
	showErrNoDiscordExe := func() {
		Popup.Alert("Error", "Discord hasn't provided any executable name for this application or they are all invalid")
	}

	if len(game.Executables) == 0 {
		showErrNoDiscordExe()
		return
	}

	var selected Executable
	selected.Name = ""
	for i := range game.Executables {
		if strings.ToLower(game.Executables[i].Os) == "win32" || strings.ToLower(game.Executables[i].Os) == "win64" {
			if strings.ContainsAny(game.Executables[i].Name, "<>\"|?*") {
				continue
			}
			selected = game.Executables[i]
			break
		}
	}
	if selected.Name == "" {
		showErrNoDiscordExe()
		return
	}

	send, err := json.Marshal(game)
	if err != nil {
		Popup.Alert("Error", "An error has occurred preparing the data")
		return
	}

	newName := strings.ReplaceAll(selected.Name, "/", "\\")
	split := strings.Split(newName, "\\")
	exeName := newName
	folder := os.TempDir() + "\\discord-fake-activity-embedded\\"
	err = os.MkdirAll(folder, 777)
	if err != nil {
		Popup.Alert("Error", "Unable to create temp folder please report this! ("+err.Error()+")")
		return
	}
	if len(split) > 1 {
		exeName = split[len(split)-1]
		split = split[:len(split)-1]
		for i := range split {
			folder += split[i] + "\\"
		}
	}
	finalName := folder + exeName
	finalName = strings.ReplaceAll(finalName, "\\", "/")

	err = os.MkdirAll(path.Dir(finalName), 777)
	if err != nil {
		Popup.Alert("Error", "Unable to create temp folder please report this! ("+err.Error()+")")
		return
	}

	file, err := os.Create(finalName)
	if err != nil {
		Popup.Alert("Error", "Unable to extract executable please report this! ("+err.Error()+")")
		return
	}
	_ = file.Chmod(777)
	_, err = file.Write(embeddedLauncher)
	if err != nil {
		Popup.Alert("Error", "Unable to write executable please report this! ("+err.Error()+")")
		_ = file.Close()
		return
	}
	_ = file.Close()

	cmd := exec.Command(finalName, strings.Split(selected.Arguments, " ")...)
	err = cmd.Start()
	if err != nil {
		Popup.Alert("Error", "Unable to run executable please report this! ("+err.Error()+")")
		return
	}
	go func() {
		_ = cmd.Wait()
	}()

	ln, err := npipe.Listen(`\\.\pipe\discord-fake-activity-embedded`)
	if err != nil {
		Popup.Alert("Error", "Unable to open named pipe, already starting another game?")
		return
	}
	defer func(ln *npipe.PipeListener) {
		_ = ln.Close()
	}(ln)

	conn, err := ln.Accept()
	if err != nil {
		Popup.Alert("Error", "Child died before any connection could be completed")
		return
	}
	_, _ = conn.Write(send)
	_ = conn.Close()
}

func makeList(gameStore []GameStore) *widget.List {
	gameList := widget.NewList(
		func() int {
			return len(gameStore)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(gameStore[i].Name)
		})
	gameList.Resize(fyne.Size{Width: 390, Height: 550})
	gameList.Move(fyne.Position{Y: 42})
	gameList.OnSelected = func(id widget.ListItemID) {
		gameList.UnselectAll()
		launchGame(gameStore[id])
	}
	return gameList
}

//go:embed assets/GitHub.svg
var githubSvg []byte

//go:embed assets/Discord.png
var discordPng []byte

func overlayImage() *widget.Button {
	btn := widget.NewButton("", func() {
		openAppLink()
	})
	btn.SetIcon(fyne.NewStaticResource("githubLogo.svg", githubSvg))
	btn.Resize(fyne.Size{Width: 32, Height: 32})
	btn.Move(fyne.Position{X: -8, Y: 600 - 32})
	return btn
}

func main() {
	gameStore := downloadActivities()
	application := app.New()
	window := application.NewWindow("Select a game to imitate")
	window.Resize(fyne.Size{Width: 400, Height: 600})
	window.SetFixedSize(true)
	window.SetIcon(fyne.NewStaticResource("discordLogo.png", discordPng))

	layout := container.NewWithoutLayout()

	creditsButton := overlayImage()

	gameList := makeList(gameStore)
	layout.Add(gameList)
	layout.Add(creditsButton)

	searchBar := widget.NewEntry()
	searchBar.PlaceHolder = "Search game by title from the ones that Discord supports"
	searchBar.OnChanged = func(search string) {
		if search == "" {
			layout.Remove(gameList)
			gameList = makeList(gameStore)
			layout.Add(gameList)
		} else {
			layout.Remove(gameList)
			var filter []GameStore
			for _, game := range gameStore {
				check := strings.HasPrefix(strings.ToLower(game.Name), strings.ToLower(search))
				if !check && len(search) > 5 && !strings.Contains(search, " ") {
					check = strings.Contains(strings.ToLower(game.Name), strings.ToLower(search))
				}
				if check {
					filter = append(filter, game)
				}
			}
			gameList = makeList(filter)
			layout.Add(gameList)
		}
		layout.Remove(creditsButton)
		layout.Add(creditsButton)
		layout.Refresh()
	}
	searchBar.Resize(fyne.Size{Width: 394, Height: 36})
	layout.Add(searchBar)

	window.SetContent(layout)
	window.ShowAndRun()
}
