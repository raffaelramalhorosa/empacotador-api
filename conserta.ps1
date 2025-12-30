Write-Host " Corrigindo problemas de imports..." -ForegroundColor Cyan

# Limpa caches
Write-Host " Limpando caches..."
go clean -modcache
go clean -cache

# Atualiza go.mod
Write-Host " Atualizando go.mod..."
go mod tidy

# Baixa dependências
Write-Host " Baixando dependências..."
go mod download

#  Tentar build
Write-Host " Testando build..."
$buildResult = go build 2>&1

if ($LASTEXITCODE -eq 0) {
    Write-Host " Build bem-sucedido!" -ForegroundColor Green
} else {
    Write-Host " Erro no build:" -ForegroundColor Red
    Write-Host $buildResult
}

Write-Host ""
Write-Host " Estrutura do projeto:" -ForegroundColor Cyan
tree /F