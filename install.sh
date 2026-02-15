#!/bin/bash

# wip

paru -S go sway sway-contrib swappy

cd ./playerctl-monitor
go build -o /usr/bin/playerctl-monitor main.go

cd ..
cp ./sway/config ~/.config/sway/config