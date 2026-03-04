# Configuração de Monitoramento de Segurança - IceWarp

## Arquitetura de Monitoramento

### Componentes Principais
1. **Log Aggregation** - Centralização de logs
2. **SIEM** - Correlação de eventos
3. **Network Monitoring** - Análise de tráfego
4. **File Integrity Monitoring** - Detecção de alterações
5. **Vulnerability Scanning** - Verificação contínua

## Configuração de Logs

### Rsyslog Configuration
```bash
# /etc/rsyslog.d/icewarp-security.conf

# IceWarp Security Logs
$template IceWarpSecurityFormat,"%timestamp% %hostname% %syslogtag% %msg%\n"

# Log all IceWarp events
if $programname startswith 'icewarp' then {
    /var/log/icewarp/security.log;IceWarpSecurityFormat
    stop
}

# Authentication failures
if $msg contains 'authentication failed' then {
    /var/log/icewarp/auth_failures.log;IceWarpSecurityFormat
    @@siem-server:514
}

# Suspicious activities
if $msg contains 'x-file-operation' or $msg contains 'command injection' then {
    /var/log/icewarp/suspicious.log;IceWarpSecurityFormat
    @@siem-server:514
}
```

### Logrotate Configuration
```bash
# /etc/logrotate.d/icewarp-security

/var/log/icewarp/*.log {
    daily
    rotate 90
    compress
    delaycompress
    missingok
    notifempty
    create 0644 icewarp icewarp
    postrotate
        systemctl reload rsyslog
    endscript
}
```

## Detecção de Ameaças

### Script de Detecção CVE-2025-14500
```bash
#!/bin/bash
# detect_cve_2025_14500.sh

LOG_FILE="/var/log/icewarp/access.log"
ALERT_FILE="/var/log/security_alerts.log"
THRESHOLD=5

# Detectar tentativas de command injection via X-File-Operation
ATTEMPTS=$(grep -c "X-File-Operation.*[;&|`\$\(\)]" $LOG_FILE)

if [ $ATTEMPTS -gt $THRESHOLD ]; then
    ALERT="CRITICAL: CVE-2025-14500 exploitation attempts detected - $ATTEMPTS attempts"
    echo "$(date): $ALERT" >> $ALERT_FILE
    
    # Extrair IPs suspeitos
    SUSPICIOUS_IPS=$(grep "X-File-Operation.*[;&|`\$\(\)]" $LOG_FILE | awk '{print $1}' | sort | uniq)
    
    # Bloquear IPs automaticamente
    for IP in $SUSPICIOUS_IPS; do
        iptables -I INPUT 1 -s $IP -j DROP
        echo "$(date): Blocked IP $IP for CVE-2025-14500 attempts" >> $ALERT_FILE
    done
    
    # Enviar alerta
    echo "$ALERT" | mail -s "CRITICAL: CVE-2025-14500 Attack Detected" admin@example.com
fi
```

### Script de Detecção CVE-2025-52691
```bash
#!/bin/bash
# detect_cve_2025_52691.sh

WEBROOT="/opt/icewarp/webmail"
ALERT_FILE="/var/log/security_alerts.log"

# Detectar uploads de arquivos suspeitos
SUSPICIOUS_FILES=$(find $WEBROOT -name "*.php" -o -name "*.jsp" -o -name "*.asp" -newer $(date -d "1 hour ago" +%Y-%m-%d))

if [ ! -z "$SUSPICIOUS_FILES" ]; then
    ALERT="CRITICAL: CVE-2025-52691 - Suspicious file uploads detected"
    echo "$(date): $ALERT" >> $ALERT_FILE
    echo "Files: $SUSPICIOUS_FILES" >> $ALERT_FILE
    
    # Quarentena automática
    for FILE in $SUSPICIOUS_FILES; do
        mv "$FILE" "/quarantine/$(basename $FILE).$(date +%s)"
        echo "$(date): Quarantined file: $FILE" >> $ALERT_FILE
    done
    
    # Enviar alerta
    echo "$ALERT - Files quarantined" | mail -s "CRITICAL: CVE-2025-52691 Attack Detected" admin@example.com
fi
```

### Script de Detecção CVE-2026-22907
```bash
#!/bin/bash
# detect_cve_2026_22907.sh

LOG_FILE="/var/log/icewarp/access.log"
ALERT_FILE="/var/log/security_alerts.log"

# Detectar injeção de headers maliciosos
HEADER_INJECTION=$(grep -E "User-Agent.*[;&|`\$\(\)]|Referer.*[;&|`\$\(\)]|Cookie.*[;&|`\$\(\)]" $LOG_FILE | tail -10)

if [ ! -z "$HEADER_INJECTION" ]; then
    ALERT="HIGH: CVE-2026-22907 - Header injection attempts detected"
    echo "$(date): $ALERT" >> $ALERT_FILE
    echo "$HEADER_INJECTION" >> $ALERT_FILE
    
    # Extrair e bloquear IPs
    MALICIOUS_IPS=$(echo "$HEADER_INJECTION" | awk '{print $1}' | sort | uniq)
    for IP in $MALICIOUS_IPS; do
        iptables -I INPUT 1 -s $IP -j DROP
        echo "$(date): Blocked IP $IP for header injection" >> $ALERT_FILE
    done
    
    echo "$ALERT" | mail -s "HIGH: CVE-2026-22907 Attack Detected" admin@example.com
fi
```

## Monitoramento de Integridade

### AIDE Configuration
```bash
# /etc/aide/aide.conf

# IceWarp critical files
/opt/icewarp/config f+p+u+g+s+m+c+md5+sha256
/opt/icewarp/bin f+p+u+g+s+m+c+md5+sha256
/opt/icewarp/webmail f+p+u+g+s+m+c+md5+sha256

# System files
/etc f+p+u+g+s+m+c+md5+sha256
/usr/bin f+p+u+g+s+m+c+md5+sha256
/usr/sbin f+p+u+g+s+m+c+md5+sha256

# Initialize database
aide --init
mv /var/lib/aide/aide.db.new /var/lib/aide/aide.db

# Daily check
aide --check
```

### Tripwire Alternative Script
```bash
#!/bin/bash
# file_integrity_monitor.sh

BASELINE_DIR="/var/lib/security/baselines"
CURRENT_DIR="/var/lib/security/current"
ALERT_FILE="/var/log/integrity_alerts.log"

mkdir -p $BASELINE_DIR $CURRENT_DIR

# Gerar baseline se não existir
if [ ! -f "$BASELINE_DIR/icewarp_hashes.txt" ]; then
    find /opt/icewarp -type f -exec sha256sum {} \; > $BASELINE_DIR/icewarp_hashes.txt
    echo "Baseline created: $(date)" >> $ALERT_FILE
    exit 0
fi

# Gerar hashes atuais
find /opt/icewarp -type f -exec sha256sum {} \; > $CURRENT_DIR/icewarp_hashes.txt

# Comparar com baseline
CHANGES=$(diff $BASELINE_DIR/icewarp_hashes.txt $CURRENT_DIR/icewarp_hashes.txt)

if [ ! -z "$CHANGES" ]; then
    ALERT="WARNING: File integrity changes detected in IceWarp"
    echo "$(date): $ALERT" >> $ALERT_FILE
    echo "$CHANGES" >> $ALERT_FILE
    
    # Enviar alerta detalhado
    echo -e "$ALERT\n\nChanges:\n$CHANGES" | mail -s "File Integrity Alert" admin@example.com
    
    # Atualizar baseline após verificação manual
    # cp $CURRENT_DIR/icewarp_hashes.txt $BASELINE_DIR/icewarp_hashes.txt
fi
```

## Monitoramento de Rede

### Suricata Rules para IceWarp
```bash
# /etc/suricata/rules/icewarp.rules

# CVE-2025-14500 Detection
alert http any any -> any any (msg:"CVE-2025-14500 X-File-Operation Command Injection"; content:"X-File-Operation"; http_header; pcre:"/[;&|`\$\(\)]/"; sid:1001; rev:1;)

# CVE-2025-52691 Detection
alert http any any -> any any (msg:"CVE-2025-52691 Malicious File Upload"; content:"Content-Type"; http_header; content:"multipart/form-data"; http_header; pcre:"/\.(php|jsp|asp|exe)$/"; sid:1002; rev:1;)

# CVE-2026-22907 Detection
alert http any any -> any any (msg:"CVE-2026-22907 Header Injection"; content:"|3b|"; http_header; content:"|7c|"; http_header; sid:1003; rev:1;)

# Brute Force Detection
alert tcp any any -> any 25 (msg:"SMTP Brute Force"; flags:S; threshold:type both, track by_src, count 10, seconds 60; sid:1004; rev:1;)
alert tcp any any -> any 143 (msg:"IMAP Brute Force"; flags:S; threshold:type both, track by_src, count 10, seconds 60; sid:1005; rev:1;)
```

### Network Monitoring Script
```bash
#!/bin/bash
# network_monitor.sh

LOG_FILE="/var/log/network_security.log"
ALERT_THRESHOLD=100

# Monitorar conexões por IP
declare -A CONNECTION_COUNT

while true; do
    # Contar conexões ativas por IP
    netstat -tn | awk '/ESTABLISHED/ {print $5}' | cut -d: -f1 | sort | uniq -c | while read count ip; do
        if [ $count -gt $ALERT_THRESHOLD ]; then
            echo "$(date): HIGH CONNECTION COUNT - IP: $ip, Connections: $count" >> $LOG_FILE
            
            # Bloquear IP se muito suspeito
            if [ $count -gt 200 ]; then
                iptables -I INPUT 1 -s $ip -j DROP
                echo "$(date): BLOCKED IP $ip for excessive connections ($count)" >> $LOG_FILE
            fi
        fi
    done
    
    sleep 60
done
```

## Alertas e Notificações

### Configuração do Postfix para Alertas
```bash
# /etc/postfix/main.cf
relayhost = [smtp.example.com]:587
smtp_sasl_auth_enable = yes
smtp_sasl_password_maps = hash:/etc/postfix/sasl_passwd
smtp_sasl_security_options = noanonymous
smtp_tls_security_level = encrypt

# /etc/postfix/sasl_passwd
[smtp.example.com]:587 alerts@example.com:password
```

### Script de Alertas Centralizados
```bash
#!/bin/bash
# security_alerting.sh

ALERT_EMAIL="security-team@example.com"
ALERT_LOG="/var/log/security_alerts.log"
WEBHOOK_URL="https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"

send_alert() {
    local SEVERITY=$1
    local MESSAGE=$2
    local DETAILS=$3
    
    # Log local
    echo "$(date): [$SEVERITY] $MESSAGE" >> $ALERT_LOG
    
    # Email
    echo -e "Severity: $SEVERITY\nMessage: $MESSAGE\nDetails:\n$DETAILS" | \
        mail -s "[$SEVERITY] IceWarp Security Alert" $ALERT_EMAIL
    
    # Slack (se configurado)
    if [ ! -z "$WEBHOOK_URL" ]; then
        curl -X POST -H 'Content-type: application/json' \
            --data "{\"text\":\"[$SEVERITY] $MESSAGE\"}" \
            $WEBHOOK_URL
    fi
    
    # SMS para alertas críticos (usando API)
    if [ "$SEVERITY" = "CRITICAL" ]; then
        # Implementar integração com serviço de SMS
        echo "SMS alert would be sent for: $MESSAGE"
    fi
}

# Exemplo de uso
# send_alert "CRITICAL" "CVE-2025-14500 exploitation detected" "Details from logs..."
```

## Dashboard de Segurança

### Grafana Dashboard Configuration
```json
{
  "dashboard": {
    "title": "IceWarp Security Dashboard",
    "panels": [
      {
        "title": "Security Alerts",
        "type": "stat",
        "targets": [
          {
            "expr": "increase(security_alerts_total[1h])",
            "legendFormat": "Alerts/Hour"
          }
        ]
      },
      {
        "title": "CVE Exploitation Attempts",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(cve_2025_14500_attempts[5m])",
            "legendFormat": "CVE-2025-14500"
          },
          {
            "expr": "rate(cve_2025_52691_attempts[5m])",
            "legendFormat": "CVE-2025-52691"
          }
        ]
      },
      {
        "title": "Top Attacking IPs",
        "type": "table",
        "targets": [
          {
            "expr": "topk(10, sum by (source_ip) (rate(attack_attempts[1h])))"
          }
        ]
      }
    ]
  }
}
```

### Prometheus Metrics
```bash
# /etc/prometheus/icewarp_exporter.py
#!/usr/bin/env python3

import time
import re
from prometheus_client import start_http_server, Counter, Gauge

# Métricas
cve_2025_14500_attempts = Counter('cve_2025_14500_attempts_total', 'CVE-2025-14500 exploitation attempts')
cve_2025_52691_attempts = Counter('cve_2025_52691_attempts_total', 'CVE-2025-52691 exploitation attempts')
auth_failures = Counter('auth_failures_total', 'Authentication failures')
active_connections = Gauge('active_connections', 'Active connections')

def parse_logs():
    with open('/var/log/icewarp/access.log', 'r') as f:
        for line in f:
            if 'x-file-operation' in line.lower() and re.search(r'[;&|`$()]', line):
                cve_2025_14500_attempts.inc()
            
            if re.search(r'\.(php|jsp|asp)', line) and 'POST' in line:
                cve_2025_52691_attempts.inc()
            
            if 'authentication failed' in line.lower():
                auth_failures.inc()

if __name__ == '__main__':
    start_http_server(8000)
    while True:
        parse_logs()
        time.sleep(60)
```

## Automação de Resposta

### Fail2Ban Configuration
```ini
# /etc/fail2ban/jail.d/icewarp.conf

[icewarp-auth]
enabled = true
port = 25,110,143,993,995
protocol = tcp
filter = icewarp-auth
logpath = /var/log/icewarp/smtp.log
maxretry = 3
bantime = 3600
findtime = 600

[icewarp-cve-2025-14500]
enabled = true
port = 80,443
protocol = tcp
filter = icewarp-cve-2025-14500
logpath = /var/log/icewarp/access.log
maxretry = 1
bantime = 86400
findtime = 300
```

```bash
# /etc/fail2ban/filter.d/icewarp-cve-2025-14500.conf
[Definition]
failregex = ^<HOST> .* ".*X-File-Operation.*[;&|`$()].*"
ignoreregex =
```

### Automated Incident Response
```bash
#!/bin/bash
# automated_response.sh

INCIDENT_THRESHOLD=10
CURRENT_ALERTS=$(tail -100 /var/log/security_alerts.log | grep "$(date +%Y-%m-%d)" | wc -l)

if [ $CURRENT_ALERTS -gt $INCIDENT_THRESHOLD ]; then
    echo "AUTOMATED RESPONSE TRIGGERED - $CURRENT_ALERTS alerts today"
    
    # Ativar modo de proteção intensiva
    /opt/security/scripts/lockdown_mode.sh
    
    # Notificar equipe de segurança
    echo "Automated incident response activated due to $CURRENT_ALERTS security alerts" | \
        mail -s "AUTOMATED RESPONSE ACTIVATED" security-team@example.com
    
    # Iniciar coleta de evidências
    /opt/security/scripts/evidence_collection.sh
fi
```

---
**Documento**: Configuração de Monitoramento de Segurança  
**Versão**: 1.0  
**Data**: 2026-03-04  
**Revisão**: Mensal