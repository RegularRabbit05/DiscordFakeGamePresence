# Discord Presence Faker
Program to fool Discord into displaying games you don't own as playing.<br>
The games are fetched at runtime from Discord's serers, as such this app won't need updates for new games.<br>
<br>
## How does it work?
This app creates a directory in your Temp folder and inside of it extracts a small executable named the way discord expects your selected game to be called.<br>
The app then starts the child executable and passes it some more information about the game such as the Window title through a Windows named pipe.<br>
The child app then creates a named window that will be picked up by Discord. Depending on Discord's configuration they overlay may appear on the child executable as well.<br>
<br>
To stop the game from showing just close the window you want to terminate. Multiple games can be played at the same time (Quests). If needed the child windows can be screen-shared.<br>
<br>
### Credits - Libraries used:
- [Fyne](https://github.com/fyne-io/fyne) (Main window)
- [RayLib / Go Bindings](https://github.com/gen2brain/raylib-go) (Sub-window, more lightweight)
- [Popup](https://github.com/ChromeTemp/Popup)
- [npipe](https://github.com/natefinch/npipe)
