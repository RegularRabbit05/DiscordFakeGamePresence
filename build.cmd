cd exe
go build -ldflags="-s -w -H=windowsgui" -o ../embedded.exe
cd ..
go build -ldflags="-s -w -H=windowsgui" -o DiscordGameFaker.exe
del embedded.exe