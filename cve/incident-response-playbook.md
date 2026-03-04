# Playbook de Resposta a Incidentes - IceWarp

## Classificação de Incidentes

### Nível 1 - CRÍTICO
- Comprometimento confirmado do servidor
- Execução de código malicioso detectada
- Exfiltração de dados em andamento
- **Tempo de Resposta**: < 15 minutos

### Nível 2 - ALTO
- Tentativas de exploração detectadas
- Anomalias de tráfego significativas
- Falhas de autenticação em massa
- **Tempo de Resposta**: < 1 hora

### Nível 3 - MÉDIO
- Atividade suspeita detectada
- Violações de política de segurança
- Alertas de monitoramento
- **Tempo de Resposta**: < 4 horas

## Procedimentos de Resposta Imediata

### FASE 1 - CONTENÇÃO (0-30 minutos)

#### 1.1 Isolamento Imediato
```bash
#!/bin/bash
# emergency_isolation.sh

echo "=== EMERGENCY ISOLATION INITIATED ==="
date >> /var/log/incident_response.log

# Bloquear tráfego suspeito
iptables -I INPUT 1 -s SUSPICIOUS_IP -j DROP
iptables -I OUTPUT 1 -d SUSPICIOUS_IP -j DROP

# Limitar conexões
iptables -A INPUT -p tcp --dport 443 -m connlimit --connlimit-above 10 -j DROP
iptables -A INPUT -p tcp --dport 80 -m connlimit --connlimit-above 10 -j DROP

# Backup de logs críticos
cp /var/log/icewarp/*.log /backup/incident_$(date +%Y%m%d_%H%M%S)/

echo "Isolation completed at $(date)" >> /var/log/incident_response.log
```

#### 1.2 Preservação de Evidências
```bash
#!/bin/bash
# evidence_collection.sh

INCIDENT_DIR="/forensics/incident_$(date +%Y%m%d_%H%M%S)"
mkdir -p $INCIDENT_DIR

# Capturar estado do sistema
ps aux > $INCIDENT_DIR/processes.txt
netstat -tulpn > $INCIDENT_DIR/network_connections.txt
lsof > $INCIDENT_DIR/open_files.txt
df -h > $INCIDENT_DIR/disk_usage.txt

# Capturar logs
cp -r /var/log/icewarp/ $INCIDENT_DIR/logs/
cp /var/log/auth.log $INCIDENT_DIR/
cp /var/log/syslog $INCIDENT_DIR/

# Hash de arquivos críticos
find /opt/icewarp -type f -exec sha256sum {} \; > $INCIDENT_DIR/file_hashes.txt

# Capturar tráfego de rede
tcpdump -i any -w $INCIDENT_DIR/network_capture.pcap &
TCPDUMP_PID=$!
sleep 300  # Capturar por 5 minutos
kill $TCPDUMP_PID

echo "Evidence collection completed: $INCIDENT_DIR"
```

### FASE 2 - ANÁLISE (30 minutos - 2 horas)

#### 2.1 Análise de Logs
```bash
#!/bin/bash
# log_analysis.sh

INCIDENT_DIR=$1
REPORT_FILE="$INCIDENT_DIR/analysis_report.txt"

echo "=== INCIDENT ANALYSIS REPORT ===" > $REPORT_FILE
echo "Generated: $(date)" >> $REPORT_FILE
echo "" >> $REPORT_FILE

# Analisar tentativas de exploração CVE-2025-14500
echo "=== CVE-2025-14500 Analysis ===" >> $REPORT_FILE
grep -i "x-file-operation" /var/log/icewarp/access.log | tail -50 >> $REPORT_FILE

# Analisar uploads suspeitos CVE-2025-52691
echo "=== CVE-2025-52691 Analysis ===" >> $REPORT_FILE
grep -E "\.(php|jsp|asp|exe)" /var/log/icewarp/access.log | tail -50 >> $REPORT_FILE

# Analisar injeção de headers CVE-2026-22907
echo "=== CVE-2026-22907 Analysis ===" >> $REPORT_FILE
grep -E "[;&|`\$\(\)]" /var/log/icewarp/access.log | tail -50 >> $REPORT_FILE

# Analisar falhas de autenticação
echo "=== Authentication Failures ===" >> $REPORT_FILE
grep "authentication failed" /var/log/icewarp/*.log | tail -50 >> $REPORT_FILE

# IPs mais ativos
echo "=== Top Active IPs ===" >> $REPORT_FILE
awk '{print $1}' /var/log/icewarp/access.log | sort | uniq -c | sort -nr | head -20 >> $REPORT_FILE

echo "Analysis completed: $REPORT_FILE"
```

#### 2.2 Verificação de Integridade
```bash
#!/bin/bash
# integrity_check.sh

echo "=== SYSTEM INTEGRITY CHECK ==="

# Verificar arquivos modificados recentemente
echo "Files modified in last 24 hours:"
find /opt/icewarp -type f -mtime -1 -ls

# Verificar processos suspeitos
echo "Suspicious processes:"
ps aux | grep -E "(nc|netcat|bash|sh|python|perl)" | grep -v grep

# Verificar conexões de rede suspeitas
echo "Suspicious network connections:"
netstat -tulpn | grep -E "(ESTABLISHED|LISTEN)" | grep -v -E "(443|80|25|110|143|993|995)"

# Verificar arquivos web suspeitos
echo "Suspicious web files:"
find /opt/icewarp/webmail -name "*.php" -o -name "*.jsp" -newer $(date -d "24 hours ago" +%Y-%m-%d)

# Verificar modificações em configurações
echo "Configuration changes:"
find /opt/icewarp/config -name "*.conf" -mtime -1 -ls
```

### FASE 3 - ERRADICAÇÃO (2-8 horas)

#### 3.1 Remoção de Ameaças
```bash
#!/bin/bash
# threat_removal.sh

echo "=== THREAT REMOVAL INITIATED ==="

# Remover arquivos maliciosos identificados
MALICIOUS_FILES=(
    "/opt/icewarp/webmail/shell.php"
    "/opt/icewarp/webmail/backdoor.jsp"
    "/tmp/malware.exe"
)

for file in "${MALICIOUS_FILES[@]}"; do
    if [ -f "$file" ]; then
        echo "Removing malicious file: $file"
        rm -f "$file"
        echo "Removed: $file" >> /var/log/incident_response.log
    fi
done

# Terminar processos suspeitos
SUSPICIOUS_PIDS=$(ps aux | grep -E "(nc|netcat|malware)" | grep -v grep | awk '{print $2}')
for pid in $SUSPICIOUS_PIDS; do
    echo "Killing suspicious process: $pid"
    kill -9 $pid
    echo "Killed PID: $pid" >> /var/log/incident_response.log
done

# Fechar conexões suspeitas
netstat -tulpn | grep ESTABLISHED | grep -v -E "(443|80|25|110|143|993|995)" | awk '{print $7}' | cut -d'/' -f1 | xargs -r kill

echo "Threat removal completed"
```

#### 3.2 Aplicação de Patches
```bash
#!/bin/bash
# emergency_patching.sh

echo "=== EMERGENCY PATCHING ==="

# Backup antes do patch
cp -r /opt/icewarp /backup/pre_patch_$(date +%Y%m%d_%H%M%S)

# Aplicar patches críticos
/opt/icewarp/tool.sh update --security-only --force

# Verificar versão pós-patch
/opt/icewarp/tool.sh version

# Reiniciar serviços
systemctl restart icewarp

# Verificar funcionamento
sleep 30
systemctl status icewarp

echo "Emergency patching completed"
```

### FASE 4 - RECUPERAÇÃO (8-24 horas)

#### 4.1 Restauração de Serviços
```bash
#!/bin/bash
# service_recovery.sh

echo "=== SERVICE RECOVERY ==="

# Verificar integridade dos serviços
/opt/icewarp/tool.sh verify --integrity

# Restaurar configurações seguras
cp /backup/secure_configs/*.conf /opt/icewarp/config/

# Aplicar hardening
/opt/icewarp/scripts/hardening.sh

# Reiniciar com configurações seguras
systemctl restart icewarp

# Verificar funcionalidades críticas
curl -I https://icewarp.armazemdc.inf.br/
telnet localhost 25 < /dev/null
telnet localhost 143 < /dev/null

echo "Service recovery completed"
```

#### 4.2 Monitoramento Intensivo
```bash
#!/bin/bash
# intensive_monitoring.sh

echo "=== INTENSIVE MONITORING ACTIVATED ==="

# Monitoramento de logs em tempo real
tail -f /var/log/icewarp/*.log | grep -E "(error|fail|attack|exploit)" &

# Monitoramento de processos
while true; do
    ps aux | grep -E "(nc|netcat|bash|sh)" | grep -v grep
    sleep 60
done &

# Monitoramento de rede
while true; do
    netstat -tulpn | grep ESTABLISHED | grep -v -E "(443|80|25|110|143|993|995)"
    sleep 300
done &

echo "Intensive monitoring activated"
```

## Comunicação de Incidentes

### Template de Comunicação Inicial
```
ASSUNTO: [CRÍTICO] Incidente de Segurança - IceWarp Server

RESUMO:
- Data/Hora: [TIMESTAMP]
- Severidade: [CRÍTICO/ALTO/MÉDIO]
- Sistema Afetado: IceWarp Mail Server
- Status: [CONTIDO/EM ANÁLISE/RESOLVIDO]

IMPACTO:
- Serviços afetados: [LISTA]
- Usuários impactados: [NÚMERO]
- Dados comprometidos: [SIM/NÃO/INVESTIGANDO]

AÇÕES TOMADAS:
- [LISTA DE AÇÕES]

PRÓXIMOS PASSOS:
- [CRONOGRAMA]

CONTATO: [RESPONSÁVEL]
```

### Escalação de Incidentes
1. **Nível 1**: Administrador de Sistema
2. **Nível 2**: Especialista em Segurança + Gerente de TI
3. **Nível 3**: CISO + Diretoria + Jurídico
4. **Nível 4**: Autoridades + Clientes (se necessário)

## Lições Aprendidas

### Template de Relatório Pós-Incidente
```markdown
# Relatório Pós-Incidente

## Resumo Executivo
- **Incidente**: [DESCRIÇÃO]
- **Data**: [DATA]
- **Duração**: [TEMPO]
- **Impacto**: [DESCRIÇÃO]

## Cronologia
- [TIMESTAMP] - Detecção inicial
- [TIMESTAMP] - Contenção
- [TIMESTAMP] - Análise
- [TIMESTAMP] - Erradicação
- [TIMESTAMP] - Recuperação

## Causa Raiz
- [ANÁLISE DETALHADA]

## Ações Corretivas
- [ ] Ação 1 - Responsável - Prazo
- [ ] Ação 2 - Responsável - Prazo

## Melhorias Implementadas
- [LISTA DE MELHORIAS]

## Recomendações
- [RECOMENDAÇÕES FUTURAS]
```

## Contatos de Emergência

### Equipe de Resposta
- **Líder de Incidentes**: [CONTATO]
- **Administrador de Sistema**: [CONTATO]
- **Especialista em Segurança**: [CONTATO]
- **Gerente de TI**: [CONTATO]

### Fornecedores
- **IceWarp Support**: support.icewarp.com
- **ISP**: [CONTATO]
- **Fornecedor de Segurança**: [CONTATO]

### Autoridades
- **CERT.br**: cert@cert.br
- **Polícia Federal**: [CONTATO]

---
**Documento**: Playbook de Resposta a Incidentes  
**Versão**: 1.0  
**Data**: 2026-03-04  
**Revisão**: Semestral