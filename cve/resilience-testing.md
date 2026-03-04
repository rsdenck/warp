# Testes de Resiliência e Validação de Segurança - IceWarp

## Metodologia de Testes Defensivos

### Objetivos dos Testes
1. **Validar eficácia** das medidas de segurança implementadas
2. **Identificar gaps** na detecção e resposta
3. **Testar procedimentos** de resposta a incidentes
4. **Verificar resiliência** do sistema sob stress
5. **Validar configurações** de segurança

## Testes de Configuração de Segurança

### 1. Validação de Patches
```bash
#!/bin/bash
# patch_validation.sh

echo "=== PATCH VALIDATION TEST ==="

# Verificar versão atual
CURRENT_VERSION=$(/opt/icewarp/tool.sh version | grep "Version" | awk '{print $2}')
echo "Current Version: $CURRENT_VERSION"

# Versões seguras mínimas
declare -A SECURE_VERSIONS
SECURE_VERSIONS["14.2"]="14.2.0.9"
SECURE_VERSIONS["14.1"]="14.1.0.19"
SECURE_VERSIONS["14.0"]="14.0.0.18"
SECURE_VERSIONS["13.0"]="13.0.3.13"

# Verificar se versão é segura
MAJOR_VERSION=$(echo $CURRENT_VERSION | cut -d. -f1-2)
SECURE_VERSION=${SECURE_VERSIONS[$MAJOR_VERSION]}

if [ ! -z "$SECURE_VERSION" ]; then
    if [[ "$CURRENT_VERSION" < "$SECURE_VERSION" ]]; then
        echo "❌ VULNERABLE: Version $CURRENT_VERSION is below secure minimum $SECURE_VERSION"
        exit 1
    else
        echo "✅ SECURE: Version $CURRENT_VERSION meets security requirements"
    fi
else
    echo "⚠️  WARNING: Unknown version branch $MAJOR_VERSION"
fi
```

### 2. Teste de Headers de Segurança
```bash
#!/bin/bash
# security_headers_test.sh

URL="https://icewarp.armazemdc.inf.br/"
REPORT_FILE="/tmp/headers_test_$(date +%Y%m%d_%H%M%S).txt"

echo "=== SECURITY HEADERS TEST ===" > $REPORT_FILE
echo "URL: $URL" >> $REPORT_FILE
echo "Date: $(date)" >> $REPORT_FILE
echo "" >> $REPORT_FILE

# Headers obrigatórios
REQUIRED_HEADERS=(
    "strict-transport-security"
    "x-frame-options"
    "x-content-type-options"
    "content-security-policy"
)

RESPONSE=$(curl -s -I "$URL")

for header in "${REQUIRED_HEADERS[@]}"; do
    if echo "$RESPONSE" | grep -qi "$header"; then
        echo "✅ $header: PRESENT" >> $REPORT_FILE
    else
        echo "❌ $header: MISSING" >> $REPORT_FILE
    fi
done

# Headers que devem estar ausentes (information disclosure)
FORBIDDEN_HEADERS=(
    "server"
    "x-powered-by"
    "x-aspnet-version"
)

for header in "${FORBIDDEN_HEADERS[@]}"; do
    if echo "$RESPONSE" | grep -qi "^$header:"; then
        echo "❌ $header: EXPOSED" >> $REPORT_FILE
    else
        echo "✅ $header: HIDDEN" >> $REPORT_FILE
    fi
done

echo "Report saved: $REPORT_FILE"
cat $REPORT_FILE
```

### 3. Teste de Configuração SSL/TLS
```bash
#!/bin/bash
# ssl_tls_test.sh

HOST="icewarp.armazemdc.inf.br"
PORT="443"
REPORT_FILE="/tmp/ssl_test_$(date +%Y%m%d_%H%M%S).txt"

echo "=== SSL/TLS CONFIGURATION TEST ===" > $REPORT_FILE
echo "Host: $HOST:$PORT" >> $REPORT_FILE
echo "Date: $(date)" >> $REPORT_FILE
echo "" >> $REPORT_FILE

# Testar protocolos inseguros
echo "Testing insecure protocols:" >> $REPORT_FILE

INSECURE_PROTOCOLS=("ssl2" "ssl3" "tls1" "tls1_1")
for protocol in "${INSECURE_PROTOCOLS[@]}"; do
    if timeout 10 openssl s_client -connect $HOST:$PORT -$protocol -quiet < /dev/null 2>/dev/null; then
        echo "❌ $protocol: ENABLED (INSECURE)" >> $REPORT_FILE
    else
        echo "✅ $protocol: DISABLED" >> $REPORT_FILE
    fi
done

# Testar protocolos seguros
echo "" >> $REPORT_FILE
echo "Testing secure protocols:" >> $REPORT_FILE

SECURE_PROTOCOLS=("tls1_2" "tls1_3")
for protocol in "${SECURE_PROTOCOLS[@]}"; do
    if timeout 10 openssl s_client -connect $HOST:$PORT -$protocol -quiet < /dev/null 2>/dev/null; then
        echo "✅ $protocol: ENABLED" >> $REPORT_FILE
    else
        echo "❌ $protocol: DISABLED" >> $REPORT_FILE
    fi
done

# Verificar certificado
echo "" >> $REPORT_FILE
echo "Certificate information:" >> $REPORT_FILE
CERT_INFO=$(timeout 10 openssl s_client -connect $HOST:$PORT -servername $HOST 2>/dev/null | openssl x509 -noout -dates -subject -issuer 2>/dev/null)
echo "$CERT_INFO" >> $REPORT_FILE

echo "SSL/TLS test completed: $REPORT_FILE"
```

## Testes de Detecção de Ameaças

### 4. Simulação de Tentativas de Autenticação
```bash
#!/bin/bash
# auth_test.sh

HOST="icewarp.armazemdc.inf.br"
TEST_USER="testuser@example.com"
WRONG_PASSWORD="wrongpassword123"

echo "=== AUTHENTICATION SECURITY TEST ==="
echo "Testing brute force detection..."

# Simular tentativas de login falhadas (SMTP)
for i in {1..5}; do
    echo "Attempt $i: Testing SMTP auth failure detection"
    timeout 10 telnet $HOST 25 << EOF
EHLO test.com
AUTH LOGIN
$(echo -n $TEST_USER | base64)
$(echo -n $WRONG_PASSWORD | base64)
QUIT
EOF
    sleep 2
done

# Verificar se IP foi bloqueado
echo "Checking if IP was blocked..."
if timeout 5 telnet $HOST 25 < /dev/null 2>/dev/null; then
    echo "⚠️  Connection still allowed - brute force protection may need tuning"
else
    echo "✅ Connection blocked - brute force protection working"
fi
```

### 5. Teste de Monitoramento de Integridade
```bash
#!/bin/bash
# integrity_test.sh

TEST_DIR="/opt/icewarp/webmail"
TEST_FILE="$TEST_DIR/test_integrity_$(date +%s).php"

echo "=== FILE INTEGRITY MONITORING TEST ==="

# Criar arquivo de teste suspeito
echo '<?php system($_GET["cmd"]); ?>' > $TEST_FILE
echo "Created test file: $TEST_FILE"

# Aguardar detecção (máximo 5 minutos)
echo "Waiting for integrity monitoring detection..."
for i in {1..30}; do
    if grep -q "$(basename $TEST_FILE)" /var/log/security_alerts.log 2>/dev/null; then
        echo "✅ File integrity monitoring detected suspicious file in $i iterations ($(($i * 10)) seconds)"
        break
    fi
    sleep 10
done

# Verificar se arquivo foi quarentenado
if [ ! -f "$TEST_FILE" ]; then
    echo "✅ Suspicious file was automatically quarantined"
else
    echo "❌ Suspicious file still present - manual cleanup required"
    rm -f $TEST_FILE
fi
```

### 6. Teste de Rate Limiting
```bash
#!/bin/bash
# rate_limiting_test.sh

URL="https://icewarp.armazemdc.inf.br/"
REQUESTS=50
CONCURRENT=10

echo "=== RATE LIMITING TEST ==="
echo "Sending $REQUESTS requests with $CONCURRENT concurrent connections..."

# Usar curl para testar rate limiting
seq 1 $REQUESTS | xargs -n1 -P$CONCURRENT -I{} curl -s -o /dev/null -w "%{http_code} %{time_total}\n" $URL > /tmp/rate_test.log

# Analisar resultados
SUCCESS_COUNT=$(grep "^200" /tmp/rate_test.log | wc -l)
BLOCKED_COUNT=$(grep -E "^(429|503|403)" /tmp/rate_test.log | wc -l)
TOTAL_COUNT=$(wc -l < /tmp/rate_test.log)

echo "Results:"
echo "  Total requests: $TOTAL_COUNT"
echo "  Successful (200): $SUCCESS_COUNT"
echo "  Rate limited/blocked: $BLOCKED_COUNT"

if [ $BLOCKED_COUNT -gt 0 ]; then
    echo "✅ Rate limiting is working - $BLOCKED_COUNT requests were blocked"
else
    echo "⚠️  No rate limiting detected - consider implementing rate limiting"
fi

rm -f /tmp/rate_test.log
```

## Testes de Resiliência do Sistema

### 7. Teste de Carga de Conexões
```bash
#!/bin/bash
# connection_load_test.sh

HOST="icewarp.armazemdc.inf.br"
MAX_CONNECTIONS=100
PORT=443

echo "=== CONNECTION LOAD TEST ==="
echo "Testing with $MAX_CONNECTIONS concurrent connections to $HOST:$PORT"

# Função para criar conexão
create_connection() {
    local id=$1
    timeout 30 bash -c "exec 3<>/dev/tcp/$HOST/$PORT" 2>/dev/null
    if [ $? -eq 0 ]; then
        echo "Connection $id: SUCCESS"
        exec 3>&-
    else
        echo "Connection $id: FAILED"
    fi
}

# Criar conexões concorrentes
for i in $(seq 1 $MAX_CONNECTIONS); do
    create_connection $i &
done

# Aguardar conclusão
wait

echo "Connection load test completed"
```

### 8. Teste de Recuperação de Serviços
```bash
#!/bin/bash
# service_recovery_test.sh

SERVICE="icewarp"
TEST_DURATION=300  # 5 minutos

echo "=== SERVICE RECOVERY TEST ==="
echo "Testing automatic recovery of $SERVICE service"

# Verificar status inicial
if systemctl is-active $SERVICE >/dev/null; then
    echo "✅ Service $SERVICE is initially running"
else
    echo "❌ Service $SERVICE is not running - starting it first"
    systemctl start $SERVICE
    sleep 10
fi

# Simular falha do serviço
echo "Simulating service failure..."
systemctl stop $SERVICE

# Monitorar recuperação automática
echo "Monitoring automatic recovery for $TEST_DURATION seconds..."
START_TIME=$(date +%s)
RECOVERED=false

while [ $(($(date +%s) - START_TIME)) -lt $TEST_DURATION ]; do
    if systemctl is-active $SERVICE >/dev/null; then
        RECOVERY_TIME=$(($(date +%s) - START_TIME))
        echo "✅ Service recovered automatically in $RECOVERY_TIME seconds"
        RECOVERED=true
        break
    fi
    sleep 5
done

if [ "$RECOVERED" = false ]; then
    echo "❌ Service did not recover automatically - manual intervention required"
    systemctl start $SERVICE
fi
```

## Testes de Backup e Recuperação

### 9. Teste de Backup
```bash
#!/bin/bash
# backup_test.sh

BACKUP_DIR="/backup/test_$(date +%Y%m%d_%H%M%S)"
SOURCE_DIR="/opt/icewarp/config"

echo "=== BACKUP SYSTEM TEST ==="
echo "Testing backup to: $BACKUP_DIR"

# Criar backup de teste
mkdir -p $BACKUP_DIR
tar -czf $BACKUP_DIR/config_backup.tar.gz $SOURCE_DIR

# Verificar integridade do backup
if tar -tzf $BACKUP_DIR/config_backup.tar.gz >/dev/null 2>&1; then
    echo "✅ Backup created successfully and is valid"
    
    # Verificar tamanho do backup
    BACKUP_SIZE=$(du -sh $BACKUP_DIR/config_backup.tar.gz | cut -f1)
    echo "Backup size: $BACKUP_SIZE"
    
    # Teste de restore (em diretório temporário)
    RESTORE_DIR="/tmp/restore_test_$(date +%s)"
    mkdir -p $RESTORE_DIR
    
    if tar -xzf $BACKUP_DIR/config_backup.tar.gz -C $RESTORE_DIR; then
        echo "✅ Backup restore test successful"
        rm -rf $RESTORE_DIR
    else
        echo "❌ Backup restore test failed"
    fi
    
else
    echo "❌ Backup creation failed or backup is corrupted"
fi

# Limpeza
rm -rf $BACKUP_DIR
```

### 10. Teste de Monitoramento de Logs
```bash
#!/bin/bash
# log_monitoring_test.sh

LOG_FILE="/var/log/icewarp/test.log"
ALERT_LOG="/var/log/security_alerts.log"
TEST_MESSAGE="TEST_SECURITY_EVENT_$(date +%s)"

echo "=== LOG MONITORING TEST ==="

# Criar evento de teste
echo "$(date): $TEST_MESSAGE - simulated security event" >> $LOG_FILE

# Aguardar detecção pelo sistema de monitoramento
echo "Waiting for log monitoring system to detect test event..."
sleep 30

# Verificar se evento foi detectado
if grep -q "$TEST_MESSAGE" $ALERT_LOG 2>/dev/null; then
    echo "✅ Log monitoring system detected test event"
else
    echo "❌ Log monitoring system did not detect test event"
fi

# Limpeza
sed -i "/$TEST_MESSAGE/d" $LOG_FILE 2>/dev/null
sed -i "/$TEST_MESSAGE/d" $ALERT_LOG 2>/dev/null
```

## Relatório de Testes de Resiliência

### Script de Relatório Consolidado
```bash
#!/bin/bash
# resilience_report.sh

REPORT_FILE="/tmp/resilience_report_$(date +%Y%m%d_%H%M%S).html"

cat > $REPORT_FILE << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <title>IceWarp Resilience Test Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .pass { color: green; font-weight: bold; }
        .fail { color: red; font-weight: bold; }
        .warn { color: orange; font-weight: bold; }
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <h1>IceWarp Security Resilience Test Report</h1>
    <p><strong>Generated:</strong> $(date)</p>
    <p><strong>Server:</strong> icewarp.armazemdc.inf.br</p>
    
    <h2>Test Summary</h2>
    <table>
        <tr><th>Test Category</th><th>Status</th><th>Details</th></tr>
EOF

# Executar todos os testes e compilar resultados
echo "Executing resilience tests..."

# Teste 1: Patch Validation
if ./patch_validation.sh >/dev/null 2>&1; then
    echo '<tr><td>Patch Validation</td><td class="pass">PASS</td><td>All security patches applied</td></tr>' >> $REPORT_FILE
else
    echo '<tr><td>Patch Validation</td><td class="fail">FAIL</td><td>Missing security patches</td></tr>' >> $REPORT_FILE
fi

# Teste 2: Security Headers
if ./security_headers_test.sh >/dev/null 2>&1; then
    echo '<tr><td>Security Headers</td><td class="pass">PASS</td><td>All required headers present</td></tr>' >> $REPORT_FILE
else
    echo '<tr><td>Security Headers</td><td class="warn">WARN</td><td>Some headers missing</td></tr>' >> $REPORT_FILE
fi

# Adicionar mais testes...

cat >> $REPORT_FILE << 'EOF'
    </table>
    
    <h2>Recommendations</h2>
    <ul>
        <li>Continue regular security testing</li>
        <li>Monitor security alerts daily</li>
        <li>Update patches immediately when available</li>
        <li>Review and update security configurations monthly</li>
    </ul>
    
    <h2>Next Review</h2>
    <p>Next resilience test scheduled for: $(date -d "+1 month")</p>
</body>
</html>
EOF

echo "Resilience test report generated: $REPORT_FILE"
```

## Cronograma de Testes

### Testes Diários
- Verificação de patches
- Monitoramento de logs
- Teste de conectividade básica

### Testes Semanais
- Teste de headers de segurança
- Verificação de SSL/TLS
- Teste de integridade de arquivos

### Testes Mensais
- Teste completo de resiliência
- Simulação de falhas
- Teste de backup e recuperação
- Revisão de configurações de segurança

### Testes Trimestrais
- Teste de penetração autorizado
- Revisão completa de segurança
- Atualização de procedimentos
- Treinamento da equipe

---
**Documento**: Testes de Resiliência e Validação  
**Versão**: 1.0  
**Data**: 2026-03-04  
**Revisão**: Mensal