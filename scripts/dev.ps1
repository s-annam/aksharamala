# Kill existing processes
taskkill /F /IM go.exe 2>$null
taskkill /F /IM node.exe 2>$null

# Start the Go server
Start-Process -NoNewWindow powershell -ArgumentList "cd $PSScriptRoot\..; go run cmd/webserver/main.go"

# Start the frontend
Start-Process -NoNewWindow powershell -ArgumentList "cd $PSScriptRoot\..\web; npm start"

Write-Host "Development servers started!"
Write-Host "Frontend: http://localhost:5173"
Write-Host "Backend: http://localhost:8081"
