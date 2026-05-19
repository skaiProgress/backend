# Импорт local_data.sql в Docker Postgres (запускать из каталога backend).
param(
    [switch]$SkipMigrate
)

$ErrorActionPreference = "Stop"
$Backend = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
$SqlFile = Join-Path $PSScriptRoot "local_data.sql"

Push-Location $Backend
try {
    $pg = docker compose ps postgres --format "{{.Status}}" 2>$null
    if (-not $pg -or $pg -notmatch "Up") {
        Write-Host "Запуск Postgres..."
        docker compose up -d postgres
        $deadline = (Get-Date).AddMinutes(2)
        while ((Get-Date) -lt $deadline) {
            $health = docker compose ps postgres --format "{{.Health}}"
            if ($health -eq "healthy") { break }
            Start-Sleep -Seconds 2
        }
    }

    if (-not $SkipMigrate) {
        Write-Host "Миграции..."
        $env:DATABASE_URL = "postgres://postgres:postgres@localhost:54329/aiqadam?sslmode=disable"
        go run ./cmd/migrate -command up
    }

    Write-Host "Импорт данных..."
    Get-Content $SqlFile -Raw -Encoding UTF8 | docker compose exec -T postgres `
        psql -U postgres -d aiqadam -v ON_ERROR_STOP=1

    Write-Host "Готово. Проверка:"
    docker compose exec -T postgres psql -U postgres -d aiqadam -c `
        "SELECT COUNT(*) AS users FROM auth.users; SELECT COUNT(*) AS courses FROM public.courses;"
}
finally {
    Pop-Location
}
