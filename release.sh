#!/usr/bin/env bash

go clean

rm -rf EDx52Display

mkdir EDx52Display

go build

cp -r EDx52display.exe conf.yaml LICENSE README.md X52MFDDriver.exe names DepInclude EDx52Display/
