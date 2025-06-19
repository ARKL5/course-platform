# Course Platform Data Initialization Script

Write-Host "Starting course platform data initialization..." -ForegroundColor Green

# Go to project root directory
Set-Location ..

# Build and run initialization script
Write-Host "Compiling data initialization program..." -ForegroundColor Yellow
go run scripts/init-data.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "Data initialization completed successfully!" -ForegroundColor Green
    Write-Host "You can now start the server to see the course platform with real data" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Yellow
    Write-Host "   1. Run .\quick-start.ps1 to start all services" -ForegroundColor White
    Write-Host "   2. Visit http://localhost:8083 to view the course platform" -ForegroundColor White
    Write-Host "   3. Login to see the real course data" -ForegroundColor White
} else {
    Write-Host "Data initialization failed. Please check database connection and configuration" -ForegroundColor Red
}

Write-Host ""
Write-Host "Press any key to continue..."
Read-Host 