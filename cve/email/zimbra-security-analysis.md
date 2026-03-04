# Análise de Segurança - Zimbra Admin Console
**Data**: 2026-03-04  
**Sistema**: Zimbra Admin Console  
**URL**: https://console.armazemdc.inf.br:9071/  
**Classificação**: CRÍTICO

## Sumário Executivo

### Situação Crítica Identificada
A exposição do Zimbra Admin Console representa **RISCO CRÍTICO** para toda a infraestrutura de email corporativo, permitindo controle administrativo total sobre contas, dados e configurações de email.

### Sistema Analisado
- **Zimbra Admin Console**: https://console.armazemdc.inf.br:9071/

### Classificação de Risco: **CRÍTICO**
- **Impacto**: Comprometimento total do sistema de email
- **Probabilidade**: Alta (console administrativo exposto)
- **Urgência**: EMERGENCIAL

## Vulnerabilidades Críticas do Zimbra

### CVEs Críticas Conhecidas
1. **CVE-2022-27925** (CVSS 9.8)
   - Zimbra Collaboration Suite RCE
   - Execução remota de código sem autenticação
   - Impacto total no servidor de email

2. **CVE-2022-37042** (CVSS 9.8)
   - Zimbra Collaboration Suite
   - Authentication bypass
   - Acesso administrativo sem credenciais

3. **CVE-2023-37580** (CVSS 8.8)
   - Zimbra Admin Console
   - Cross-site scripting (XSS)
   - Roubo de sessões administrativas

4. **CVE-2022-41352** (CVSS 8.1)
   - Zimbra Collaboration Suite
   - Directory traversal
   - Acesso a arquivos do sistema

### Riscos Específicos do Email

#### 1. Email Infrastructure Compromise
- **Acesso total** a todas as contas de email
- **Interceptação** de comunicações corporativas
- **Modificação** de emails e configurações
- **Exfiltração** de dados confidenciais

#### 2. Business Email Compromise (BEC)
- **Impersonation** de executivos
- **Fraudes financeiras** via email
- **Phishing interno** usando contas legítimas
- **Lateral movement** via credenciais

#### 3. Compliance Violations
- **LGPD/GDPR** - Exposição de dados pessoais
- **SOX** - Comprometimento de comunicações financeiras
- **HIPAA** - Vazamento de informações de saúde
- **Attorney-client privilege** - Exposição de comunicações legais

## Análise Técnica Detalhada

### Configuração Exposta
```
Internet → Zimbra Admin Console (console.armazemdc.inf.br:9071)
    ↓
Zimbra Mailbox Servers
    ↓
User Mailboxes (Todas as contas corporativas)
```

### Porta Não Padrão
- **Porta 9071** - Console administrativo
- **Possível tentativa** de "security by obscurity"
- **Facilmente descoberta** por scanners
- **Não oferece proteção** real

### Superfície de Ataque
- **Interface web** administrativa completa
- **API REST** do Zimbra
- **SOAP API** para administração
- **Possível acesso** a logs e configurações

## Testes de Segurança Defensivos

### Verificação de Versão
```bash
# Verificar versão via headers
curl -I -s https://console.armazemdc.inf.br:9071/ | grep -i "server\|zimbra"

# Verificar informações na página de login
curl -s https://console.armazemdc.inf.br:9071/ | grep -i "version\|build\|zimbra"
```

### Análise de SSL/TLS
```bash
# Verificar configuração SSL
nmap --script ssl-enum-ciphers -p 9071 console.armazemdc.inf.br
openssl s_client -connect console.armazemdc.inf.br:9071
```

### Verificação de API
```bash
# Verificar endpoints de API expostos
curl -s https://console.armazemdc.inf.br:9071/service/admin/soap
curl -s https://console.armazemdc.inf.br:9071/service/soap
```

### Headers de Segurança
```bash
# Verificar headers de segurança
curl -I https://console.armazemdc.inf.br:9071/ | \
    grep -E "(X-Frame-Options|Content-Security-Policy|Strict-Transport-Security)"
```

## Configuração Segura Recomendada

### 1. Isolamento Completo
```
[Internet] → [VPN Gateway] → [Admin Network] → [Zimbra Admin Console]
```

### 2. Network Segmentation
```
Admin Network (VLAN 10)
├── Zimbra Admin Console
├── Admin Workstations
└── Management Tools

Mail Network (VLAN 20)
├── Zimbra Mailbox Servers
├── LDAP Servers
└── Database Servers

User Network (VLAN 30)
├── User Workstations
├── Mobile Devices
└── Webmail Access
```

### 3. Firewall Configuration
```bash
# Bloquear acesso direto da internet
iptables -A INPUT -p tcp --dport 9071 -s 0.0.0.0/0 -j DROP

# Permitir apenas via VPN
iptables -A INPUT -p tcp --dport 9071 -s VPN_NETWORK -j ACCEPT

# Permitir apenas IPs administrativos
iptables -A INPUT -p tcp --dport 9071 -s ADMIN_IP_RANGE -j ACCEPT
```

### 4. Zimbra Hardening
```bash
# Configurações de segurança via CLI
zmprov mcf zimbraAdminConsoleLoginURL https://internal-admin.company.com:9071
zmprov mcf zimbraAdminConsoleLogoutURL https://internal-admin.company.com:9071/logout
zmprov mcf zimbraWebClientLoginURL https://mail.company.com
zmprov mcf zimbraWebClientLogoutURL https://mail.company.com/logout

# Configurar timeout de sessão
zmprov mcf zimbraAdminConsoleUISessionTimeout 1800000

# Habilitar HTTPS apenas
zmprov ms `zmhostname` zimbraMailMode https
zmprov ms `zmhostname` zimbraMailSSLPort 443
```

## Monitoramento e Detecção

### Logs Críticos
- **Admin login attempts** - /opt/zimbra/log/mailbox.log
- **Configuration changes** - /opt/zimbra/log/audit.log
- **Account modifications** - /opt/zimbra/log/sync.log
- **Failed authentications** - /opt/zimbra/log/zimbra.log

### Script de Monitoramento
```bash
#!/bin/bash
# zimbra_security_monitor.sh

LOG_DIR="/opt/zimbra/log"
ALERT_FILE="/var/log/zimbra_security.log"
ALERT_EMAIL="security@company.com"

# Monitorar tentativas de login no admin console
ADMIN_FAILED_LOGINS=$(grep "authentication failed.*admin" $LOG_DIR/mailbox.log | \
    grep "$(date +%Y-%m-%d)" | wc -l)

if [ $ADMIN_FAILED_LOGINS -gt 3 ]; then
    ALERT="CRITICAL: $ADMIN_FAILED_LOGINS failed admin login attempts on Zimbra"
    echo "$(date): $ALERT" >> $ALERT_FILE
    echo "$ALERT" | mail -s "CRITICAL: Zimbra Admin Security Alert" $ALERT_EMAIL
fi

# Monitorar alterações de configuração
CONFIG_CHANGES=$(grep "modify.*config" $LOG_DIR/audit.log | \
    grep "$(date +%Y-%m-%d)" | wc -l)

if [ $CONFIG_CHANGES -gt 5 ]; then
    echo "$(date): WARNING - Multiple config changes: $CONFIG_CHANGES" >> $ALERT_FILE
fi

# Monitorar criação/modificação de contas
ACCOUNT_CHANGES=$(grep -E "(CreateAccount|ModifyAccount|DeleteAccount)" $LOG_DIR/audit.log | \
    grep "$(date +%Y-%m-%d)" | wc -l)

if [ $ACCOUNT_CHANGES -gt 10 ]; then
    echo "$(date): WARNING - Unusual account activity: $ACCOUNT_CHANGES changes" >> $ALERT_FILE
fi
```

### SIEM Integration
```json
{
  "rule_name": "Zimbra_Admin_Unauthorized_Access",
  "description": "Detect unauthorized access to Zimbra Admin Console",
  "conditions": [
    {
      "field": "destination_host",
      "operator": "equals",
      "value": "console.armazemdc.inf.br"
    },
    {
      "field": "destination_port",
      "operator": "equals",
      "value": "9071"
    },
    {
      "field": "source_ip",
      "operator": "not_in",
      "value": ["AUTHORIZED_ADMIN_IPS"]
    }
  ],
  "severity": "CRITICAL",
  "action": "IMMEDIATE_ALERT_AND_BLOCK"
}
```

## Email Security Best Practices

### SPF Configuration
```dns
; SPF Record para armazemdc.inf.br
armazemdc.inf.br. IN TXT "v=spf1 ip4:SERVER_IP include:_spf.google.com ~all"
```

### DKIM Configuration
```bash
# Gerar chave DKIM
/opt/zimbra/libexec/zmdkimkeyutil -a -d armazemdc.inf.br -s default

# Configurar DKIM
zmprov md armazemdc.inf.br zimbraDomainDKIMSelector default
zmprov md armazemdc.inf.br zimbraDKIMEnabled TRUE
```

### DMARC Configuration
```dns
; DMARC Record
_dmarc.armazemdc.inf.br. IN TXT "v=DMARC1; p=quarantine; rua=mailto:dmarc@armazemdc.inf.br; ruf=mailto:dmarc@armazemdc.inf.br; fo=1"
```

## Backup e Disaster Recovery

### Zimbra Backup Strategy
```bash
#!/bin/bash
# zimbra_backup.sh

BACKUP_DIR="/backup/zimbra/$(date +%Y%m%d)"
mkdir -p $BACKUP_DIR

# Backup completo do Zimbra
su - zimbra -c "zmbackup -f -a all -t full --dest $BACKUP_DIR"

# Backup de configurações
su - zimbra -c "zmprov gaa > $BACKUP_DIR/all_accounts.txt"
su - zimbra -c "zmprov gad > $BACKUP_DIR/all_domains.txt"
su - zimbra -c "zmprov gas > $BACKUP_DIR/all_servers.txt"

# Backup de certificados
cp -r /opt/zimbra/ssl $BACKUP_DIR/ssl_backup/
```

### Configuration Backup
```bash
# Backup de configurações críticas
#!/bin/bash
# zimbra_config_backup.sh

CONFIG_BACKUP="/backup/zimbra/config/$(date +%Y%m%d)"
mkdir -p $CONFIG_BACKUP

# Exportar configurações
su - zimbra -c "zmprov gacf > $CONFIG_BACKUP/global_config.txt"
su - zimbra -c "zmprov gasc > $CONFIG_BACKUP/server_config.txt"
su - zimbra -c "zmprov gadl > $CONFIG_BACKUP/distribution_lists.txt"
```

## Incident Response para Email

### Cenário: Comprometimento do Admin Console
1. **Isolamento imediato**
   ```bash
   # Bloquear acesso ao admin console
   iptables -A INPUT -p tcp --dport 9071 -j DROP
   ```

2. **Verificação de integridade**
   ```bash
   # Verificar contas administrativas
   su - zimbra -c "zmprov gaaa"
   
   # Verificar últimas alterações
   grep "$(date +%Y-%m-%d)" /opt/zimbra/log/audit.log
   ```

3. **Análise forense**
   ```bash
   # Coletar evidências
   cp -r /opt/zimbra/log /forensics/zimbra_logs_$(date +%s)
   su - zimbra -c "zmprov gaa" > /forensics/current_accounts.txt
   ```

4. **Recuperação**
   - Restaurar a partir de backup limpo
   - Resetar senhas administrativas
   - Aplicar patches de segurança

## Compliance e Auditoria

### Email Retention Policies
```bash
# Configurar políticas de retenção
zmprov mcf zimbraDataSourcePurgePolicy "after:365d"
zmprov mcf zimbraMailPurgePolicy "after:2555d"  # 7 anos
```

### Audit Logging
```bash
# Habilitar auditoria completa
zmprov mcf zimbraAuditLogEnabled TRUE
zmprov mcf zimbraAuditLogLevel INFO
zmprov mcf zimbraAuditLogDestination file:///opt/zimbra/log/audit.log
```

### Legal Hold
```bash
# Implementar legal hold para investigações
zmprov ma user@domain.com zimbraMailQuota 0
zmprov ma user@domain.com zimbraFeatureMailEnabled FALSE
```

## Recomendações Específicas

### Emergenciais (0-2h)
1. **Bloquear acesso** público ao admin console
2. **Implementar VPN** obrigatória
3. **Ativar logging** máximo
4. **Verificar contas** administrativas

### Críticas (2-24h)
1. **Aplicar patches** de segurança
2. **Implementar 2FA** para admins
3. **Configurar SIEM** integration
4. **Estabelecer monitoramento** 24/7

### Altas (1-7 dias)
1. **Migrar admin console** para rede interna
2. **Implementar WAF** dedicado
3. **Configurar SPF/DKIM/DMARC** completo
4. **Estabelecer backup** automatizado

### Médias (1-4 semanas)
1. **Implementar DLP** (Data Loss Prevention)
2. **Email encryption** para dados sensíveis
3. **Regular penetration** testing
4. **Compliance audit** preparation

## Arquitetura Segura de Referência

### Recommended Email Architecture
```
[Internet]
    ↓
[Email Security Gateway]
    ↓
[Load Balancer]
    ↓
[Zimbra Proxy Servers - DMZ]
    ↓
[Zimbra Mailbox Servers - Internal]
    ↓
[LDAP/Database Servers - Secure Zone]

[Admin Network - Isolated]
    ↓
[Zimbra Admin Console - Internal Only]
```

---
**Documento**: Análise de Segurança Zimbra  
**Versão**: 1.0  
**Data**: 2026-03-04  
**Próxima Revisão**: IMEDIATA após isolamento do admin console