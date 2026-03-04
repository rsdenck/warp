# Análise SSL/TLS - IceWarp Server

## Configuração Atual

### Certificado SSL
**URL Analisada**: https://icewarp.armazemdc.inf.br/

#### Status do Certificado
- **Emissor**: [Análise limitada - servidor não forneceu detalhes completos]
- **Validade**: [Requer verificação manual]
- **Algoritmo**: [Requer verificação manual]
- **Tamanho da Chave**: [Requer verificação manual]

### Protocolos TLS Suportados
**Recomendação**: Verificar com ferramentas especializadas

```bash
# Verificação manual recomendada
nmap --script ssl-enum-ciphers -p 443 icewarp.armazemdc.inf.br
openssl s_client -connect icewarp.armazemdc.inf.br:443 -tls1_2
openssl s_client -connect icewarp.armazemdc.inf.br:443 -tls1_3
```

## Vulnerabilidades SSL/TLS Conhecidas

### Protocolos Inseguros
- **SSLv2**: DEVE estar desabilitado (vulnerável)
- **SSLv3**: DEVE estar desabilitado (POODLE)
- **TLSv1.0**: DEVE estar desabilitado (vulnerável)
- **TLSv1.1**: DEVE estar desabilitado (vulnerável)
- **TLSv1.2**: Mínimo aceitável
- **TLSv1.3**: Recomendado

### Cipher Suites Inseguros
**Evitar**:
- RC4 (todas as variantes)
- DES e 3DES
- MD5
- SHA1 para assinatura
- NULL ciphers
- EXPORT ciphers

**Recomendados**:
```
ECDHE-RSA-AES256-GCM-SHA384
ECDHE-RSA-AES128-GCM-SHA256
ECDHE-RSA-CHACHA20-POLY1305
DHE-RSA-AES256-GCM-SHA384
DHE-RSA-AES128-GCM-SHA256
```

## Configuração Segura Recomendada

### Apache/Nginx Configuration
```apache
# Apache SSL Configuration
SSLEngine on
SSLProtocol -all +TLSv1.2 +TLSv1.3
SSLCipherSuite ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256
SSLHonorCipherOrder on
SSLCompression off
SSLSessionTickets off

# HSTS Header
Header always set Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"

# OCSP Stapling
SSLUseStapling on
SSLStaplingCache shmcb:/var/run/ocsp(128000)
```

```nginx
# Nginx SSL Configuration
ssl_protocols TLSv1.2 TLSv1.3;
ssl_ciphers ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256;
ssl_prefer_server_ciphers on;
ssl_session_cache shared:SSL:10m;
ssl_session_timeout 10m;
ssl_stapling on;
ssl_stapling_verify on;

# Security Headers
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains; preload" always;
```

### IceWarp Specific Configuration
```ini
# /opt/icewarp/config/webmail.conf
[SSL]
MinProtocol=TLSv1.2
MaxProtocol=TLSv1.3
CipherList=ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-GCM-SHA256
DHParameters=/opt/icewarp/config/dhparam.pem
HSTS=true
HSTSMaxAge=31536000
HSTSIncludeSubdomains=true
```

## Testes de Segurança SSL/TLS

### Ferramentas Recomendadas

#### 1. SSL Labs Test
```bash
# Online: https://www.ssllabs.com/ssltest/
# Objetivo: Grade A+ desejável
```

#### 2. testssl.sh
```bash
# Download e execução
git clone https://github.com/drwetter/testssl.sh.git
cd testssl.sh
./testssl.sh https://icewarp.armazemdc.inf.br/
```

#### 3. OpenSSL Manual Testing
```bash
# Testar TLS 1.2
openssl s_client -connect icewarp.armazemdc.inf.br:443 -tls1_2 -cipher 'ECDHE-RSA-AES256-GCM-SHA384'

# Testar TLS 1.3
openssl s_client -connect icewarp.armazemdc.inf.br:443 -tls1_3

# Verificar certificado
openssl s_client -connect icewarp.armazemdc.inf.br:443 -showcerts
```

#### 4. Nmap SSL Scripts
```bash
# Enumerar ciphers
nmap --script ssl-enum-ciphers -p 443 icewarp.armazemdc.inf.br

# Verificar vulnerabilidades conhecidas
nmap --script ssl-* -p 443 icewarp.armazemdc.inf.br
```

## Vulnerabilidades SSL/TLS Específicas

### Heartbleed (CVE-2014-0160)
```bash
# Teste
nmap -p 443 --script ssl-heartbleed icewarp.armazemdc.inf.br
```

### POODLE (CVE-2014-3566)
```bash
# Teste SSLv3
openssl s_client -connect icewarp.armazemdc.inf.br:443 -ssl3
```

### BEAST (CVE-2011-3389)
```bash
# Verificar CBC ciphers em TLS 1.0
openssl s_client -connect icewarp.armazemdc.inf.br:443 -tls1 -cipher 'AES128-SHA'
```

### CRIME/BREACH
```bash
# Verificar compressão SSL
openssl s_client -connect icewarp.armazemdc.inf.br:443 -comp
```

### Sweet32 (CVE-2016-2183)
```bash
# Verificar 3DES
openssl s_client -connect icewarp.armazemdc.inf.br:443 -cipher '3DES'
```

## Monitoramento Contínuo

### Script de Monitoramento SSL
```bash
#!/bin/bash
# ssl_monitor.sh

DOMAIN="icewarp.armazemdc.inf.br"
LOG_FILE="/var/log/ssl_monitor.log"
ALERT_EMAIL="admin@example.com"

# Verificar expiração do certificado
EXPIRY_DATE=$(openssl s_client -connect $DOMAIN:443 -servername $DOMAIN 2>/dev/null | openssl x509 -noout -dates | grep notAfter | cut -d= -f2)
EXPIRY_EPOCH=$(date -d "$EXPIRY_DATE" +%s)
CURRENT_EPOCH=$(date +%s)
DAYS_UNTIL_EXPIRY=$(( ($EXPIRY_EPOCH - $CURRENT_EPOCH) / 86400 ))

if [ $DAYS_UNTIL_EXPIRY -lt 30 ]; then
    echo "$(date): WARNING - SSL certificate expires in $DAYS_UNTIL_EXPIRY days" >> $LOG_FILE
    echo "SSL certificate for $DOMAIN expires in $DAYS_UNTIL_EXPIRY days" | mail -s "SSL Certificate Expiry Warning" $ALERT_EMAIL
fi

# Verificar grade SSL Labs (requer API key)
# GRADE=$(curl -s "https://api.ssllabs.com/api/v3/analyze?host=$DOMAIN" | jq -r '.endpoints[0].grade')
# if [ "$GRADE" != "A+" ] && [ "$GRADE" != "A" ]; then
#     echo "$(date): WARNING - SSL grade is $GRADE" >> $LOG_FILE
# fi

# Verificar protocolos inseguros
INSECURE_PROTOCOLS=$(nmap --script ssl-enum-ciphers -p 443 $DOMAIN 2>/dev/null | grep -E "(SSLv2|SSLv3|TLSv1\.0|TLSv1\.1)")
if [ ! -z "$INSECURE_PROTOCOLS" ]; then
    echo "$(date): WARNING - Insecure SSL/TLS protocols detected" >> $LOG_FILE
    echo "Insecure SSL/TLS protocols detected on $DOMAIN" | mail -s "SSL Security Warning" $ALERT_EMAIL
fi
```

## Recomendações de Implementação

### Prioridade Alta
1. **Desabilitar protocolos inseguros** (SSLv2, SSLv3, TLS 1.0, TLS 1.1)
2. **Implementar TLS 1.3** como protocolo preferencial
3. **Configurar cipher suites seguros** (AEAD preferred)
4. **Habilitar HSTS** com preload
5. **Implementar OCSP Stapling**

### Prioridade Média
1. **Certificate Transparency** monitoring
2. **HTTP Public Key Pinning** (HPKP) - com cuidado
3. **Perfect Forward Secrecy** (PFS)
4. **Automated certificate renewal**

### Prioridade Baixa
1. **DNS-based Authentication of Named Entities** (DANE)
2. **Certificate Authority Authorization** (CAA) records
3. **Expect-CT** header

## Checklist de Verificação SSL/TLS

### Configuração Básica
- [ ] TLS 1.2 mínimo habilitado
- [ ] TLS 1.3 habilitado (se suportado)
- [ ] SSLv2/SSLv3 desabilitados
- [ ] TLS 1.0/1.1 desabilitados
- [ ] Cipher suites seguros configurados
- [ ] Cipher suites inseguros desabilitados

### Certificado
- [ ] Certificado válido e não expirado
- [ ] Cadeia de certificados completa
- [ ] Algoritmo de assinatura seguro (SHA-256+)
- [ ] Tamanho de chave adequado (2048+ RSA, 256+ ECDSA)
- [ ] SAN (Subject Alternative Names) configurado

### Headers de Segurança
- [ ] HSTS habilitado
- [ ] HSTS com includeSubDomains
- [ ] HSTS com preload (opcional)
- [ ] Expect-CT configurado (opcional)

### Funcionalidades Avançadas
- [ ] OCSP Stapling habilitado
- [ ] Perfect Forward Secrecy (PFS)
- [ ] Compressão SSL desabilitada
- [ ] Session tickets seguros
- [ ] Renegociação segura

### Monitoramento
- [ ] Monitoramento de expiração de certificado
- [ ] Testes regulares de configuração SSL
- [ ] Alertas para mudanças de configuração
- [ ] Logs de conexões SSL/TLS

---
**Documento**: Análise SSL/TLS IceWarp  
**Versão**: 1.0  
**Data**: 2026-03-04  
**Próxima Verificação**: Mensal