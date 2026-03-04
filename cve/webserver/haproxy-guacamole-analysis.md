# Análise de Segurança - HAProxy e Guacamole
**Data**: 2026-03-04  
**Sistemas**: HAProxy Stats, Apache Guacamole  
**Classificação**: ALTO/MÉDIO

## Sumário Executivo

### Situação Identificada
Os sistemas HAProxy e Guacamole apresentam **riscos significativos** relacionados à exposição de informações de infraestrutura e gateway de acesso remoto não protegido adequadamente.

### Sistemas Analisados
- **HAProxy Stats**: http://haproxy.armazemdc.inf.br:8080/
- **Apache Guacamole**: https://guacamole.armazem.cloud/

### Classificação de Risco
- **HAProxy**: MÉDIO (Information Disclosure)
- **Guacamole**: ALTO (Remote Access Gateway)

## Análise HAProxy

### Vulnerabilidades Identificadas
- **Interface de estatísticas** exposta publicamente
- **HTTP não criptografado** (porta 8080)
- **Information disclosure** sobre backend servers
- **Possível enumeração** de serviços internos

### Riscos Específicos
1. **Reconnaissance** - Mapeamento da infraestrutura
2. **Service enumeration** - Identificação de serviços backend
3. **Load balancer bypass** - Tentativas de acesso direto
4. **DoS information** - Status de serviços para ataques

### Configuração Atual Exposta
```
http://haproxy.armazemdc.inf.br:8080/
├── /stats - Interface de estatísticas
├── /health - Health checks
└── Possíveis endpoints adicionais
```

### Informações Expostas
- **Backend servers** e seus status
- **Health check** results
- **Connection statistics**
- **Load balancing** algorithms
- **Server weights** e configurações

## Análise Apache Guacamole

### Vulnerabilidades Identificadas
- **Gateway de acesso remoto** exposto publicamente
- **Possível bypass** de controles de rede interna
- **CVEs conhecidas** em versões vulneráveis
- **Weak authentication** se não configurado adequadamente

### CVEs Críticas Conhecidas
1. **CVE-2022-29405** (CVSS 7.5)
   - Apache Guacamole
   - Information disclosure
   - Exposição de dados de sessão

2. **CVE-2021-41767** (CVSS 8.8)
   - Apache Guacamole
   - Authentication bypass
   - Acesso não autorizado

3. **CVE-2020-9497** (CVSS 6.5)
   - Apache Guacamole
   - Cross-site scripting (XSS)
   - Roubo de credenciais

### Riscos Específicos do Guacamole
1. **Remote access** to internal systems
2. **Lateral movement** via compromised sessions
3. **Credential harvesting** from sessions
4. **Protocol exploitation** (RDP, SSH, VNC)

## Testes de Segurança Defensivos

### HAProxy Testing
```bash
# Verificar interface de stats
curl -s http://haproxy.armazemdc.inf.br:8080/stats

# Verificar informações expostas
curl -s http://haproxy.armazemdc.inf.br:8080/stats | grep -E "(server|backend|status)"

# Verificar outros endpoints
curl -s http://haproxy.armazemdc.inf.br:8080/health
curl -s http://haproxy.armazemdc.inf.br:8080/info
```

### Guacamole Testing
```bash
# Verificar versão
curl -s https://guacamole.armazem.cloud/ | grep -i "guacamole\|version"

# Verificar API endpoints
curl -s https://guacamole.armazem.cloud/api/tokens
curl -s https://guacamole.armazem.cloud/api/session

# Verificar SSL/TLS
nmap --script ssl-enum-ciphers -p 443 guacamole.armazem.cloud
```

### Headers de Segurança
```bash
# HAProxy
curl -I http://haproxy.armazemdc.inf.br:8080/

# Guacamole
curl -I https://guacamole.armazem.cloud/ | \
    grep -E "(X-Frame-Options|Content-Security-Policy|Strict-Transport-Security)"
```

## Configuração Segura - HAProxy

### 1. Proteger Interface de Stats
```haproxy
# /etc/haproxy/haproxy.cfg

# Configuração segura de stats
stats enable
stats uri /admin/stats
stats realm HAProxy\ Statistics
stats auth admin:STRONG_PASSWORD_HERE
stats refresh 30s
stats hide-version

# Restringir acesso por IP
stats http-request deny unless { src 192.168.1.0/24 }
stats http-request deny unless { src 10.0.0.0/8 }
```

### 2. Implementar HTTPS
```haproxy
# Configurar SSL/TLS para stats
frontend stats_frontend
    bind *:8443 ssl crt /etc/ssl/certs/haproxy.pem
    stats enable
    stats uri /admin/stats
    stats realm HAProxy\ Statistics
    stats auth admin:STRONG_PASSWORD
    
    # Headers de segurança
    http-response set-header Strict-Transport-Security "max-age=31536000; includeSubDomains"
    http-response set-header X-Frame-Options "DENY"
    http-response set-header X-Content-Type-Options "nosniff"
```

### 3. Logging e Monitoramento
```haproxy
# Configurar logging detalhado
global
    log 127.0.0.1:514 local0 info
    
defaults
    log global
    option httplog
    option dontlognull
    
# Log de acesso às stats
frontend stats_frontend
    capture request header User-Agent len 64
    capture request header Host len 32
```

## Configuração Segura - Guacamole

### 1. Autenticação Forte
```properties
# /etc/guacamole/guacamole.properties

# Database authentication
mysql-hostname: localhost
mysql-port: 3306
mysql-database: guacamole_db
mysql-username: guacamole_user
mysql-password: STRONG_DATABASE_PASSWORD

# LDAP authentication (recomendado)
ldap-hostname: ldap.company.com
ldap-port: 636
ldap-encryption-method: ssl
ldap-user-base-dn: ou=users,dc=company,dc=com
ldap-username-attribute: sAMAccountName
```

### 2. Configuração de Sessão
```properties
# Timeout de sessão
session-timeout: 1800000  # 30 minutos

# Configurações de segurança
enable-websocket: false
skip-if-unavailable: mysql
```

### 3. Reverse Proxy com Nginx
```nginx
# /etc/nginx/sites-available/guacamole

server {
    listen 443 ssl http2;
    server_name guacamole.armazem.cloud;
    
    # SSL Configuration
    ssl_certificate /etc/ssl/certs/guacamole.crt;
    ssl_certificate_key /etc/ssl/private/guacamole.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    
    # Security Headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    
    # Rate Limiting
    limit_req zone=guacamole burst=10 nodelay;
    
    # IP Whitelisting
    allow 192.168.1.0/24;
    allow 10.0.0.0/8;
    deny all;
    
    location / {
        proxy_pass http://127.0.0.1:8080/guacamole/;
        proxy_buffering off;
        proxy_http_version 1.1;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection $http_connection;
        proxy_cookie_path /guacamole/ /;
    }
}
```

## Monitoramento e Detecção

### HAProxy Monitoring
```bash
#!/bin/bash
# haproxy_monitor.sh

STATS_URL="http://haproxy.armazemdc.inf.br:8080/stats"
LOG_FILE="/var/log/haproxy_security.log"

# Monitorar acessos às stats
STATS_ACCESS=$(grep "GET /stats" /var/log/haproxy.log | grep "$(date +%Y-%m-%d)" | wc -l)

if [ $STATS_ACCESS -gt 50 ]; then
    echo "$(date): WARNING - High stats access: $STATS_ACCESS requests" >> $LOG_FILE
fi

# Monitorar backend failures
BACKEND_FAILURES=$(grep "backend.*DOWN" /var/log/haproxy.log | grep "$(date +%Y-%m-%d)" | wc -l)

if [ $BACKEND_FAILURES -gt 0 ]; then
    echo "$(date): ALERT - Backend failures detected: $BACKEND_FAILURES" >> $LOG_FILE
fi
```

### Guacamole Monitoring
```bash
#!/bin/bash
# guacamole_monitor.sh

LOG_FILE="/var/log/guacamole_security.log"
GUAC_LOG="/var/log/tomcat/catalina.out"

# Monitorar tentativas de login
FAILED_LOGINS=$(grep "Authentication failed" $GUAC_LOG | grep "$(date +%Y-%m-%d)" | wc -l)

if [ $FAILED_LOGINS -gt 10 ]; then
    ALERT="WARNING: $FAILED_LOGINS failed login attempts on Guacamole"
    echo "$(date): $ALERT" >> $LOG_FILE
    echo "$ALERT" | mail -s "Guacamole Security Alert" admin@company.com
fi

# Monitorar conexões suspeitas
SUSPICIOUS_CONNECTIONS=$(grep -E "(RDP|SSH|VNC).*failed" $GUAC_LOG | grep "$(date +%Y-%m-%d)" | wc -l)

if [ $SUSPICIOUS_CONNECTIONS -gt 5 ]; then
    echo "$(date): WARNING - Suspicious connection attempts: $SUSPICIOUS_CONNECTIONS" >> $LOG_FILE
fi
```

## Hardening Adicional

### HAProxy Security
```haproxy
# Configurações de segurança adicionais
global
    # Ocultar versão
    stats hide-version
    
    # Configurações de SSL
    ssl-default-bind-ciphers ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-GCM-SHA256
    ssl-default-bind-options no-sslv3 no-tlsv10 no-tlsv11
    
defaults
    # Timeouts de segurança
    timeout connect 5000ms
    timeout client 50000ms
    timeout server 50000ms
    
    # Headers de segurança
    http-response set-header Server "WebServer"
    http-response del-header X-Powered-By
```

### Guacamole Security
```xml
<!-- /etc/guacamole/logback.xml -->
<configuration>
    <appender name="GUAC" class="ch.qos.logback.core.FileAppender">
        <file>/var/log/guacamole/guacamole.log</file>
        <encoder>
            <pattern>%d{HH:mm:ss.SSS} [%thread] %-5level %logger{36} - %msg%n</pattern>
        </encoder>
    </appender>
    
    <!-- Log de segurança detalhado -->
    <logger name="org.apache.guacamole.auth" level="INFO"/>
    <logger name="org.apache.guacamole.net" level="INFO"/>
    
    <root level="INFO">
        <appender-ref ref="GUAC"/>
    </root>
</configuration>
```

## Recomendações Específicas

### HAProxy - Prioridade MÉDIA
1. **Migrar stats para HTTPS** (porta 8443)
2. **Implementar autenticação** forte
3. **Restringir acesso** por IP
4. **Ocultar informações** sensíveis

### Guacamole - Prioridade ALTA
1. **Implementar IP whitelisting** imediato
2. **Configurar 2FA** obrigatório
3. **Implementar rate limiting**
4. **Configurar SIEM** integration

### Ambos os Sistemas
1. **Aplicar patches** de segurança
2. **Implementar monitoramento** 24/7
3. **Configurar backup** automatizado
4. **Estabelecer procedimentos** de resposta

## Plano de Remediação

### Fase 1 - Imediata (0-24h)
- [ ] Restringir acesso HAProxy stats por IP
- [ ] Implementar autenticação no HAProxy
- [ ] Configurar IP whitelisting no Guacamole
- [ ] Ativar logging detalhado

### Fase 2 - Curto Prazo (1-7 dias)
- [ ] Migrar HAProxy stats para HTTPS
- [ ] Implementar 2FA no Guacamole
- [ ] Configurar rate limiting
- [ ] Estabelecer monitoramento

### Fase 3 - Médio Prazo (1-4 semanas)
- [ ] Implementar WAF para ambos
- [ ] Configurar SIEM integration
- [ ] Estabelecer backup automatizado
- [ ] Treinar equipe em procedimentos

## Arquitetura Segura Recomendada

### Network Segmentation
```
[Internet]
    ↓
[WAF/Load Balancer]
    ↓
[DMZ - Reverse Proxies]
    ↓
[Internal Network]
    ↓
[Management Network - HAProxy Stats]
    ↓
[Secure Network - Backend Services]
```

---
**Documento**: Análise HAProxy e Guacamole  
**Versão**: 1.0  
**Data**: 2026-03-04  
**Próxima Revisão**: Após implementação de controles básicos