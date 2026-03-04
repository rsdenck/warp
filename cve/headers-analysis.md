# Análise de Headers de Segurança - IceWarp Server

## Headers de Segurança Essenciais

### 1. Strict-Transport-Security (HSTS)
**Status**: Requer verificação  
**Importância**: CRÍTICA

```http
Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
```

**Função**: Força conexões HTTPS, previne downgrade attacks  
**Risco se ausente**: Man-in-the-middle, protocol downgrade  
**Configuração recomendada**:
```apache
Header always set Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
```

### 2. Content-Security-Policy (CSP)
**Status**: Requer verificação  
**Importância**: ALTA

```http
Content-Security-Policy: default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'
```

**Função**: Previne XSS, injection attacks  
**Risco se ausente**: Cross-site scripting, data injection  
**Configuração recomendada**:
```apache
Header always set Content-Security-Policy "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'"
```

### 3. X-Frame-Options
**Status**: Requer verificação  
**Importância**: ALTA

```http
X-Frame-Options: DENY
```

**Função**: Previne clickjacking  
**Risco se ausente**: Clickjacking attacks  
**Configuração recomendada**:
```apache
Header always set X-Frame-Options "DENY"
```

### 4. X-Content-Type-Options
**Status**: Requer verificação  
**Importância**: MÉDIA

```http
X-Content-Type-Options: nosniff
```

**Função**: Previne MIME type sniffing  
**Risco se ausente**: MIME confusion attacks  
**Configuração recomendada**:
```apache
Header always set X-Content-Type-Options "nosniff"
```

### 5. Referrer-Policy
**Status**: Requer verificação  
**Importância**: MÉDIA

```http
Referrer-Policy: strict-origin-when-cross-origin
```

**Função**: Controla informações de referrer  
**Risco se ausente**: Information leakage  
**Configuração recomendada**:
```apache
Header always set Referrer-Policy "strict-origin-when-cross-origin"
```

### 6. X-XSS-Protection
**Status**: Requer verificação  
**Importância**: BAIXA (deprecated, mas ainda útil)

```http
X-XSS-Protection: 1; mode=block
```

**Função**: Ativa proteção XSS do browser  
**Risco se ausente**: XSS em browsers antigos  
**Configuração recomendada**:
```apache
Header always set X-XSS-Protection "1; mode=block"
```

### 7. Permissions-Policy
**Status**: Requer verificação  
**Importância**: MÉDIA

```http
Permissions-Policy: geolocation=(), microphone=(), camera=(), payment=(), usb=()
```

**Função**: Controla APIs do browser  
**Risco se ausente**: Uso não autorizado de APIs  
**Configuração recomendada**:
```apache
Header always set Permissions-Policy "geolocation=(), microphone=(), camera=(), payment=(), usb=(), magnetometer=(), gyroscope=(), speaker=()"
```

## Headers Informativos (Remover/Ocultar)

### Headers que Expõem Informações
```http
Server: IceWarp/14.x.x
X-Powered-By: PHP/7.x
X-AspNet-Version: x.x.x
```

**Risco**: Information disclosure, fingerprinting  
**Configuração para ocultar**:
```apache
# Ocultar versão do servidor
ServerTokens Prod
Header unset Server
Header unset X-Powered-By
Header unset X-AspNet-Version
Header set Server "WebServer"
```

## Configuração Completa de Headers

### Apache Configuration
```apache
# /etc/apache2/conf-available/security-headers.conf

# HSTS
Header always set Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"

# CSP
Header always set Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'"

# Clickjacking Protection
Header always set X-Frame-Options "DENY"

# MIME Type Sniffing Protection
Header always set X-Content-Type-Options "nosniff"

# XSS Protection (legacy)
Header always set X-XSS-Protection "1; mode=block"

# Referrer Policy
Header always set Referrer-Policy "strict-origin-when-cross-origin"

# Permissions Policy
Header always set Permissions-Policy "geolocation=(), microphone=(), camera=(), payment=(), usb=(), magnetometer=(), gyroscope=(), speaker=()"

# Remove server information
ServerTokens Prod
Header unset Server
Header unset X-Powered-By
Header set Server "WebServer"

# Cache Control for sensitive pages
<LocationMatch "\.(php|cgi|pl|py)$">
    Header always set Cache-Control "no-store, no-cache, must-revalidate, proxy-revalidate"
    Header always set Pragma "no-cache"
    Header always set Expires "0"
</LocationMatch>
```

### Nginx Configuration
```nginx
# /etc/nginx/conf.d/security-headers.conf

# HSTS
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains; preload" always;

# CSP
add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'" always;

# Clickjacking Protection
add_header X-Frame-Options "DENY" always;

# MIME Type Sniffing Protection
add_header X-Content-Type-Options "nosniff" always;

# XSS Protection
add_header X-XSS-Protection "1; mode=block" always;

# Referrer Policy
add_header Referrer-Policy "strict-origin-when-cross-origin" always;

# Permissions Policy
add_header Permissions-Policy "geolocation=(), microphone=(), camera=(), payment=(), usb=(), magnetometer=(), gyroscope=(), speaker=()" always;

# Hide server information
server_tokens off;
more_set_headers "Server: WebServer";
```

## Testes de Headers de Segurança

### Ferramentas de Teste

#### 1. curl
```bash
# Verificar headers de resposta
curl -I https://icewarp.armazemdc.inf.br/

# Verificar headers específicos
curl -H "User-Agent: SecurityTest" -I https://icewarp.armazemdc.inf.br/
```

#### 2. Security Headers Online Test
```bash
# Usar: https://securityheaders.com/
# Objetivo: Grade A desejável
```

#### 3. Mozilla Observatory
```bash
# Usar: https://observatory.mozilla.org/
# Objetivo: Grade A+ desejável
```

#### 4. Script de Verificação Automatizada
```bash
#!/bin/bash
# check_headers.sh

URL="https://icewarp.armazemdc.inf.br/"
LOG_FILE="/var/log/header_check.log"

echo "=== Security Headers Check - $(date) ===" >> $LOG_FILE

# Verificar HSTS
HSTS=$(curl -s -I $URL | grep -i "strict-transport-security")
if [ -z "$HSTS" ]; then
    echo "MISSING: Strict-Transport-Security header" >> $LOG_FILE
else
    echo "FOUND: $HSTS" >> $LOG_FILE
fi

# Verificar CSP
CSP=$(curl -s -I $URL | grep -i "content-security-policy")
if [ -z "$CSP" ]; then
    echo "MISSING: Content-Security-Policy header" >> $LOG_FILE
else
    echo "FOUND: $CSP" >> $LOG_FILE
fi

# Verificar X-Frame-Options
XFRAME=$(curl -s -I $URL | grep -i "x-frame-options")
if [ -z "$XFRAME" ]; then
    echo "MISSING: X-Frame-Options header" >> $LOG_FILE
else
    echo "FOUND: $XFRAME" >> $LOG_FILE
fi

# Verificar X-Content-Type-Options
XCONTENT=$(curl -s -I $URL | grep -i "x-content-type-options")
if [ -z "$XCONTENT" ]; then
    echo "MISSING: X-Content-Type-Options header" >> $LOG_FILE
else
    echo "FOUND: $XCONTENT" >> $LOG_FILE
fi

# Verificar Server header (deve estar oculto)
SERVER=$(curl -s -I $URL | grep -i "^server:")
if [ ! -z "$SERVER" ]; then
    echo "WARNING: Server header exposed: $SERVER" >> $LOG_FILE
fi

echo "=== End Check ===" >> $LOG_FILE
```

## Headers Específicos para IceWarp

### Webmail Security Headers
```apache
# Para interface webmail
<Location "/webmail">
    Header always set X-Frame-Options "DENY"
    Header always set Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'"
    Header always set Cache-Control "no-store, no-cache, must-revalidate"
</Location>
```

### Admin Interface Security Headers
```apache
# Para interface administrativa
<Location "/admin">
    Header always set X-Frame-Options "DENY"
    Header always set Content-Security-Policy "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'"
    Header always set Cache-Control "no-store, no-cache, must-revalidate"
    Header always set X-Robots-Tag "noindex, nofollow, nosnippet, noarchive"
</Location>
```

### API Endpoints Security Headers
```apache
# Para endpoints de API
<LocationMatch "^/api/">
    Header always set Content-Type "application/json"
    Header always set X-Content-Type-Options "nosniff"
    Header always set Cache-Control "no-store"
    Header always set Access-Control-Allow-Origin "https://icewarp.armazemdc.inf.br"
</LocationMatch>
```

## Monitoramento de Headers

### Script de Monitoramento Contínuo
```bash
#!/bin/bash
# monitor_headers.sh

URL="https://icewarp.armazemdc.inf.br/"
ALERT_EMAIL="admin@example.com"
LOG_FILE="/var/log/header_monitor.log"

# Headers obrigatórios
REQUIRED_HEADERS=(
    "strict-transport-security"
    "x-frame-options"
    "x-content-type-options"
    "content-security-policy"
)

MISSING_HEADERS=()

for header in "${REQUIRED_HEADERS[@]}"; do
    if ! curl -s -I $URL | grep -qi "$header"; then
        MISSING_HEADERS+=("$header")
    fi
done

if [ ${#MISSING_HEADERS[@]} -gt 0 ]; then
    echo "$(date): ALERT - Missing security headers: ${MISSING_HEADERS[*]}" >> $LOG_FILE
    echo "Missing security headers detected on $URL: ${MISSING_HEADERS[*]}" | mail -s "Security Headers Alert" $ALERT_EMAIL
fi

# Verificar headers informativos expostos
EXPOSED_HEADERS=$(curl -s -I $URL | grep -iE "^(server|x-powered-by|x-aspnet-version):")
if [ ! -z "$EXPOSED_HEADERS" ]; then
    echo "$(date): WARNING - Information disclosure headers detected" >> $LOG_FILE
    echo "Information disclosure headers detected on $URL" | mail -s "Information Disclosure Alert" $ALERT_EMAIL
fi
```

## Recomendações de Implementação

### Fase 1 - Headers Críticos (Imediato)
1. **HSTS** - Implementar com max-age mínimo de 1 ano
2. **X-Frame-Options** - DENY para prevenir clickjacking
3. **X-Content-Type-Options** - nosniff para prevenir MIME sniffing
4. **Ocultar headers informativos** - Server, X-Powered-By

### Fase 2 - Headers Avançados (1-2 semanas)
1. **Content-Security-Policy** - Implementar gradualmente
2. **Referrer-Policy** - strict-origin-when-cross-origin
3. **Permissions-Policy** - Restringir APIs desnecessárias

### Fase 3 - Otimização (1 mês)
1. **CSP refinement** - Remover 'unsafe-inline' gradualmente
2. **HSTS preload** - Submeter para preload list
3. **Expect-CT** - Implementar Certificate Transparency

## Checklist de Headers de Segurança

### Headers Obrigatórios
- [ ] Strict-Transport-Security configurado
- [ ] X-Frame-Options configurado
- [ ] X-Content-Type-Options configurado
- [ ] Content-Security-Policy básico implementado

### Headers Recomendados
- [ ] Referrer-Policy configurado
- [ ] X-XSS-Protection configurado
- [ ] Permissions-Policy configurado
- [ ] Cache-Control para páginas sensíveis

### Information Disclosure
- [ ] Server header oculto/modificado
- [ ] X-Powered-By removido
- [ ] X-AspNet-Version removido
- [ ] Versões de software ocultas

### Monitoramento
- [ ] Script de verificação automatizada
- [ ] Alertas para headers ausentes
- [ ] Testes regulares com ferramentas online
- [ ] Logs de monitoramento configurados

---
**Documento**: Análise de Headers de Segurança  
**Versão**: 1.0  
**Data**: 2026-03-04  
**Próxima Verificação**: Semanal