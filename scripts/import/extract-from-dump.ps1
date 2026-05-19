# Пересобирает local_data.sql из data_only_dump.sql (лежит в этой же папке).
param(
    [string]$DumpPath = (Join-Path $PSScriptRoot "data_only_dump.sql"),
    [string]$OutPath = (Join-Path $PSScriptRoot "local_data.sql")
)

$ErrorActionPreference = "Stop"

if (-not (Test-Path $DumpPath)) {
    Write-Error "Не найден дамп: $DumpPath"
}

$tables = @(
    "auth.users",
    "public.profiles",
    "public.courses",
    "public.lessons",
    "public.course_assignments",
    "public.course_materials"
)

$lines = Get-Content -Path $DumpPath -Encoding UTF8
$blocks = [ordered]@{}

foreach ($table in $tables) {
    $pattern = "COPY $([regex]::Escape($table)) "
    $start = ($lines | Select-String -SimpleMatch $pattern | Select-Object -First 1).LineNumber
    if (-not $start) {
        Write-Warning "Секция COPY не найдена: $table"
        continue
    }
    $chunk = @()
    for ($i = $start - 1; $i -lt $lines.Count; $i++) {
        $chunk += $lines[$i]
        if ($lines[$i] -eq '\.') {
            break
        }
    }
    $blocks[$table] = $chunk
}

$header = @"
-- Автогенерация: backend/scripts/import/extract-from-dump.ps1
-- Источник: data_only_dump.sql

BEGIN;

SET session_replication_role = replica;

TRUNCATE TABLE
  public.course_assignments,
  public.course_materials,
  public.lessons,
  public.courses,
  public.profiles
RESTART IDENTITY CASCADE;

DELETE FROM auth.users;

"@

$body = ($blocks.Values | ForEach-Object { $_ -join "`n" }) -join "`n`n"
$footer = @"

SET session_replication_role = DEFAULT;

COMMIT;
"@

$content = $header + $body + $footer
Set-Content -Path $OutPath -Value $content -Encoding UTF8
Write-Host "Записано: $OutPath ($($blocks.Count) таблиц)"
