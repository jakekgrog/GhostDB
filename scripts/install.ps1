New-Item -Path "C:\Program Files (x86)\" -Name "GhostDB" -ItemType "directory"
Move-Item -Path ..\ghostdb.exe -Destination "C:\Program Files (x86)\GhostDB\ghostdb.exe"