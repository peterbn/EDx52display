go clean

Remove-Item  .\EDx52Display -Force -Recurse -ErrorAction SilentlyContinue
Remove-Item .\Release.zip -ErrorAction SilentlyContinue

mkdir EDx52Display

go build

Copy-Item -Path EDx52display.exe,conf.yaml,LICENSE,README.md,names,DepInclude -Destination .\EDx52Display -Recurse

7z.exe a Release.zip .\EDx52Display
