# Guia de Hardening - IceWarp Mail Server

## Configurações de Segurança Essenciais

### 1. Configuração de Sistema Base

#### Atualizações de Segurança
```bash
# Verificar versão atual
/opt/icewarp/tool.sh version

# Aplicar atualizações (após backup)
/opt/icewarp/tool.sh update --security-only

# Verificar integridade pós-atualização
/opt/icewarp/tool.sh verify --integrity
```

#### Configuração de Usuários e Permissões
```bash
# Criar usuário dedicado para IceWarp (se não existir)
useradd -r -s /bin/false icewarp

# Configurar permissões restritivas
chmod 750 /opt/icewarp/
chown -R icewarp:icewarp /opt/icewarp/
chmod 600 /opt/icewarp/config/*.conf
```

### 2. Configuração de Rede e Firewall

#### Firewall Restritivo
```bash
# Permitir apenas portas necessárias
ufw default deny incoming
ufw default allow outgoing

# Portas essenciais do IceWarp
ufw allow 25/tcp    # SMTP
ufw allow 110/tcp   # POP3
ufw allow 143/tcp   # IMAP
ufw allow 993/tcp   # IMAPS
ufw allow 995/tcp   # POP3S
ufw allow 80/tcp    # HTTP (redirecionar para HTTPS)
ufw allow 443/tcp   # HTTPS

# Acesso administrativo restrito
ufw allow from 192.168.1.0/24 to any port 32000  # WebAdmin

ufw enable
```

#### Configuração de SSL/TLS
```ini
# /opt/icewarp/config/tls.conf
[TLS]
MinProtocol=TLSv1.2
MaxProtocol=TLSv1.3
CipherSuites=ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-GCM-SHA256
DHParameters=/opt/icewarp/config/dhparam.pem
HSTS=true
HSTSMaxAge=31536000
```

### 3. Configuração de Autenticação

#### Políticas de Senha
```ini
# /opt/icewarp/config/security.conf
[Authentication]
MinPasswordLength=12
RequireComplexPasswords=true
PasswordExpiration=90
MaxLoginAttempts=3
LockoutDuration=1800
RequireTwoFactor=true

[Session]
SessionTimeout=1800
MaxConcurrentSessions=3
RequireSecureCookies=true
```

#### Configuração de 2FA
```ini
# /opt/icewarp/config/2fa.conf
[TwoFactor]
Enabled=true
Method=TOTP
BackupCodes=true
RequireForAdmin=true
GracePeriod=7
```

### 4. Configuração de Logging e Monitoramento

#### Logging Detalhado
```ini
# /opt/icewarp/config/logging.conf
[Logging]
Level=INFO
AuditLog=true
SecurityEvents=true
LoginAttempts=true
FileOperations=true
AdminActions=true

[LogRotation]
MaxSize=100MB
MaxFiles=30
Compress=true
```

#### Monitoramento de Integridade
```bash
#!/bin/bash
# /opt/icewarp/scripts/integrity_check.sh

# Verificar integridade de arquivos críticos
find /opt/icewarp/config -name "*.conf" -exec sha256sum {} \; > /var/log/icewarp/config_hashes.log

# Verificar processos suspeitos
ps aux | grep -E "(nc|netcat|bash|sh)" | grep -v grep > /var/log/icewarp/suspicious_processes.log

# Verificar conexões de rede
netstat -tulpn | grep ESTABLISHED > /var/log/icewarp/network_connections.log
```

### 5. Configuração de Email Security

#### Anti-Spam e Anti-Malware
```ini
# /opt/icewarp/config/antispam.conf
[AntiSpam]
Enabled=true
SpamAssassinEnabled=true
BayesianFilter=true
RBLChecks=true
GreyListing=true
SPFCheck=strict
DKIMVerification=true
DMARCPolicy=quarantine

[AntiVirus]
Enabled=true
ScanIncoming=true
ScanOutgoing=true
QuarantineSuspicious=true
```

#### Configuração de SPF, DKIM e DMARC
```dns
; SPF Record
example.com. IN TXT "v=spf1 ip4:192.168.1.100 include:icewarp.armazemdc.inf.br ~all"

; DKIM Record
default._domainkey.example.com. IN TXT "v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC..."

; DMARC Record
_dmarc.example.com. IN TXT "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.com; ruf=mailto:dmarc@example.com; fo=1"
```

### 6. Configuração de Backup e Recovery

#### Backup Automatizado
```bash
#!/bin/bash
# /opt/icewarp/scripts/backup.sh

BACKUP_DIR="/backup/icewarp"
DATE=$(date +%Y%m%d_%H%M%S)

# Criar diretório de backup
mkdir -p $BACKUP_DIR/$DATE

# Backup de configurações
tar -czf $BACKUP_DIR/$DATE/config.tar.gz /opt/icewarp/config/

# Backup de dados de email
tar -czf $BACKUP_DIR/$DATE/mail.tar.gz /opt/icewarp/mail/

# Backup de logs
tar -czf $BACKUP_DIR/$DATE/logs.tar.gz /var/log/icewarp/

# Remover backups antigos (manter 30 dias)
find $BACKUP_DIR -type d -mtime +30 -exec rm -rf {} \;
```

### 7. Configuração de Headers de Segurança

#### Headers HTTP Seguros
```apache
# /opt/icewarp/config/httpd.conf
Header always set Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
Header always set X-Content-Type-Options "nosniff"
Header always set X-Frame-Options "SAMEORIGIN"
Header always set X-XSS-Protection "1; mode=block"
Header always set Referrer-Policy "strict-origin-when-cross-origin"
Header always set Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'"
```

### 8. Configuração de Rate Limiting

#### Proteção contra Brute Force
```ini
# /opt/icewarp/config/ratelimit.conf
[RateLimit]
Enabled=true
MaxConnectionsPerIP=10
MaxRequestsPerMinute=60
BanDuration=3600
WhitelistIPs=192.168.1.0/24

[BruteForceProtection]
MaxFailedLogins=3
BanDuration=1800
MonitoringWindow=300
```

### 9. Configuração de WAF (Web Application Firewall)

#### ModSecurity Rules
```apache
# /opt/icewarp/config/modsecurity.conf
SecRuleEngine On
SecRequestBodyAccess On
SecResponseBodyAccess Off

# Bloquear injeção de comandos
SecRule ARGS "@detectSQLi" \
    "id:1001,phase:2,block,msg:'SQL Injection Attack Detected',logdata:'Matched Data: %{MATCHED_VAR} found within %{MATCHED_VAR_NAME}'"

# Bloquear headers maliciosos
SecRule REQUEST_HEADERS:X-File-Operation "@rx [;&|`$()]" \
    "id:1002,phase:1,block,msg:'Malicious X-File-Operation Header Detected'"

# Bloquear upload de arquivos suspeitos
SecRule FILES_NAMES "@rx \.(php|jsp|asp|aspx|exe|bat|cmd)$" \
    "id:1003,phase:2,block,msg:'Suspicious File Upload Detected'"
```

### 10. Monitoramento e Alertas

#### Script de Monitoramento
```bash
#!/bin/bash
# /opt/icewarp/scripts/security_monitor.sh

LOG_FILE="/var/log/icewarp/security_monitor.log"
ALERT_EMAIL="admin@example.com"

# Verificar tentativas de login suspeitas
FAILED_LOGINS=$(grep "authentication failed" /var/log/icewarp/smtp.log | wc -l)
if [ $FAILED_LOGINS -gt 10 ]; then
    echo "$(date): ALERT - $FAILED_LOGINS failed login attempts detected" >> $LOG_FILE
    echo "High number of failed login attempts detected" | mail -s "IceWarp Security Alert" $ALERT_EMAIL
fi

# Verificar arquivos suspeitos
SUSPICIOUS_FILES=$(find /opt/icewarp/webmail -name "*.php" -newer $(date -d "1 hour ago" +%Y-%m-%d) 2>/dev/null | wc -l)
if [ $SUSPICIOUS_FILES -gt 0 ]; then
    echo "$(date): ALERT - $SUSPICIOUS_FILES suspicious files detected" >> $LOG_FILE
    echo "Suspicious files detected in webmail directory" | mail -s "IceWarp Security Alert" $ALERT_EMAIL
fi

# Verificar uso de CPU/memória anômalo
CPU_USAGE=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | awk -F'%' '{print $1}')
if (( $(echo "$CPU_USAGE > 80" | bc -l) )); then
    echo "$(date): ALERT - High CPU usage: $CPU_USAGE%" >> $LOG_FILE
fi
```

### 11. Configuração de Disaster Recovery

#### Plano de Recuperação
```bash
#!/bin/bash
# /opt/icewarp/scripts/disaster_recovery.sh

BACKUP_SOURCE="/backup/icewarp/latest"
RECOVERY_LOG="/var/log/icewarp/recovery.log"

echo "$(date): Starting disaster recovery process" >> $RECOVERY_LOG

# Parar serviços
systemctl stop icewarp

# Restaurar configurações
tar -xzf $BACKUP_SOURCE/config.tar.gz -C /

# Restaurar dados de email
tar -xzf $BACKUP_SOURCE/mail.tar.gz -C /

# Verificar integridade
/opt/icewarp/tool.sh verify --integrity

# Reiniciar serviços
systemctl start icewarp

echo "$(date): Disaster recovery completed" >> $RECOVERY_LOG
```

## Checklist de Hardening

### Configuração Inicial
- [ ] Aplicar todas as atualizações de segurança
- [ ] Configurar firewall restritivo
- [ ] Implementar SSL/TLS forte
- [ ] Configurar usuários e permissões

### Autenticação e Autorização
- [ ] Implementar políticas de senha robustas
- [ ] Configurar autenticação de dois fatores
- [ ] Limitar tentativas de login
- [ ] Configurar timeout de sessão

### Monitoramento e Logging
- [ ] Habilitar logging detalhado
- [ ] Configurar rotação de logs
- [ ] Implementar monitoramento de integridade
- [ ] Configurar alertas de segurança

### Proteção de Email
- [ ] Configurar anti-spam e anti-malware
- [ ] Implementar SPF, DKIM e DMARC
- [ ] Configurar rate limiting
- [ ] Implementar greylisting

### Backup e Recovery
- [ ] Configurar backup automatizado
- [ ] Testar procedimentos de restore
- [ ] Documentar plano de disaster recovery
- [ ] Implementar backup offsite

### Segurança Web
- [ ] Configurar headers de segurança
- [ ] Implementar WAF/ModSecurity
- [ ] Bloquear uploads suspeitos
- [ ] Configurar CSP (Content Security Policy)

---
**Documento**: Guia de Hardening IceWarp  
**Versão**: 1.0  
**Data**: 2026-03-04  
**Revisão**: Trimestral