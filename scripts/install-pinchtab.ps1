# Script para instalar e configurar Pinchtab no Windows
param(
    [string]$Port = "8080",
    [string]$Token = ""
)

Write-Host "🚀 Instalando Pinchtab..." -ForegroundColor Green

# Criar diretório para Pinchtab
$pinchtabDir = "$env:USERPROFILE\.pinchtab"
if (!(Test-Path $pinchtabDir)) {
    New-Item -ItemType Directory -Path $pinchtabDir -Force
    Write-Host "📁 Diretório criado: $pinchtabDir" -ForegroundColor Yellow
}

# Download do Pinchtab (simulado - você precisa do binário real)
$pinchtabExe = "$pinchtabDir\pinchtab.exe"

if (!(Test-Path $pinchtabExe)) {
    Write-Host "⬇️  Baixando Pinchtab..." -ForegroundColor Yellow
    
    # URL do release mais recente (ajustar conforme necessário)
    $downloadUrl = "https://github.com/pinchtab/pinchtab/releases/latest/download/pinchtab-windows-amd64.exe"
    
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $pinchtabExe -ErrorAction Stop
        Write-Host "✅ Pinchtab baixado com sucesso!" -ForegroundColor Green
    } catch {
        Write-Host "❌ Erro ao baixar Pinchtab: $($_.Exception.Message)" -ForegroundColor Red
        Write-Host "💡 Baixe manualmente de: https://github.com/pinchtab/pinchtab/releases" -ForegroundColor Yellow
        exit 1
    }
}

# Configurar comando de inicialização
$startCommand = "$pinchtabExe --port $Port"
if ($Token -ne "") {
    $startCommand += " --token `"$Token`""
}

Write-Host "🔧 Comando de inicialização: $startCommand" -ForegroundColor Cyan

# Criar script de inicialização
$startScript = @"
@echo off
echo 🚀 Iniciando Pinchtab na porta $Port...
echo 💡 Acesse: http://localhost:$Port/health
echo 🛑 Pressione Ctrl+C para parar
$startCommand
"@

$startScriptPath = "$pinchtabDir\start-pinchtab.bat"
$startScript | Out-File -FilePath $startScriptPath -Encoding ASCII

Write-Host "📝 Script criado: $startScriptPath" -ForegroundColor Yellow

# Perguntar se deve iniciar agora
$start = Read-Host "🚀 Iniciar Pinchtab agora? (y/N)"
if ($start -eq "y" -or $start -eq "Y") {
    Write-Host "🌐 Iniciando Pinchtab..." -ForegroundColor Green
    Write-Host "💡 Uma janela do Chrome será aberta automaticamente" -ForegroundColor Yellow
    Write-Host "🔗 Acesse TeamChat em: https://icewarp.armazemdc.inf.br/teamchat" -ForegroundColor Cyan
    
    # Iniciar Pinchtab em processo separado
    Start-Process -FilePath $startScriptPath -WorkingDirectory $pinchtabDir
    
    # Aguardar um pouco e testar conexão
    Start-Sleep -Seconds 3
    
    try {
        $response = Invoke-RestMethod -Uri "http://localhost:$Port/health" -Method GET -TimeoutSec 5
        Write-Host "✅ Pinchtab iniciado com sucesso!" -ForegroundColor Green
        Write-Host "📊 Status: $($response.status)" -ForegroundColor Green
    } catch {
        Write-Host "⚠️  Pinchtab pode estar iniciando... Aguarde alguns segundos" -ForegroundColor Yellow
    }
    
    Write-Host ""
    Write-Host "🎯 Próximos passos:" -ForegroundColor Cyan
    Write-Host "1. Teste a conexão: warpctl pinchtab health" -ForegroundColor White
    Write-Host "2. Configure grupos: warpctl zabbix configure" -ForegroundColor White
    Write-Host "3. Teste notificação web: warpctl zabbix test-web-notify" -ForegroundColor White
    Write-Host "4. Monitoramento integrado: warpctl zabbix integrated-monitoring" -ForegroundColor White
} else {
    Write-Host ""
    Write-Host "💡 Para iniciar Pinchtab manualmente:" -ForegroundColor Yellow
    Write-Host "   $startScriptPath" -ForegroundColor White
    Write-Host ""
    Write-Host "🔗 Ou execute diretamente:" -ForegroundColor Yellow
    Write-Host "   $startCommand" -ForegroundColor White
}

Write-Host ""
Write-Host "📋 Configuração completa!" -ForegroundColor Green
Write-Host "📁 Arquivos em: $pinchtabDir" -ForegroundColor Yellow