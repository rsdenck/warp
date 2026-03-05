# 🚀 WARPCTL + PINCHTAB INTEGRATION

## Integração Completa: Zabbix + IceWarp TeamChat + Pinchtab

Esta integração combina o poder do **Pinchtab** (automação de browser) com **Zabbix API** (monitoramento) e **IceWarp TeamChat** (notificações) para criar um sistema de monitoramento avançado e automatizado.

## 🏗️ Arquitetura

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   PINCHTAB      │    │    ZABBIX API   │    │  ICEWARP        │
│ Browser Control │◄──►│   Monitoring    │◄──►│  TeamChat       │
│                 │    │                 │    │                 │
│ • Web Automation│    │ • VDC Groups    │    │ • Notifications │
│ • Screenshots   │    │ • Active Alerts │    │ • Channels      │
│ • Form Filling  │    │ • Event Data    │    │ • HTML Messages │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         ▲                       ▲                       ▲
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │    WARPCTL      │
                    │  CLI Interface  │
                    │                 │
                    │ • Configuration │
                    │ • Monitoring    │
                    │ • Integration   │
                    └─────────────────┘
```

## 📦 Instalação

### 1. Instalar Pinchtab

```bash
# Download e instalar Pinchtab
curl -L https://github.com/pinchtab/pinchtab/releases/latest/download/pinchtab-windows-amd64.exe -o pinchtab.exe

# Ou via Go
go install github.com/pinchtab/pinchtab@latest
```

### 2. Iniciar Pinchtab

```bash
# Iniciar Pinchtab na porta 8080
./pinchtab.exe --port 8080

# Ou com token de segurança
./pinchtab.exe --port 8080 --token "seu-token-seguro"
```

### 3. Configurar WARPCTL

```yaml
# warpctl.yaml
pinchtab:
  url: http://localhost:8080
  token: "seu-token-seguro"  # opcional
  auto_start: true
  integration:
    zabbix_auto_login: true
    teamchat_auto_login: true
    screenshot_on_error: true

zabbix:
  url: https://monitoramento.armazem.cloud/api_jsonrpc.php
  web_url: https://monitoramento.armazem.cloud
  username: zabbix_api
  password: kCfOtLUG3xSn2X9uyo6lXSVX1RFHGEw7
  group_mappings:
    "4296": Horus Monitoramento  # VDC_BONJA

teamchat:
  default_channel: Horus Monitoramento
```

## 🎯 Funcionalidades Implementadas

### 1. **Comandos Pinchtab Básicos**

```bash
# Verificar saúde do Pinchtab
warpctl pinchtab health

# Listar abas abertas
warpctl pinchtab tabs

# Navegar para URL
warpctl pinchtab navigate https://monitoramento.armazem.cloud

# Tirar screenshot
warpctl pinchtab screenshot zabbix.jpg --quality 90

# Ver árvore de acessibilidade
warpctl pinchtab snapshot

# Clicar em elemento
warpctl pinchtab click e15

# Digitar texto
warpctl pinchtab type "#username" "zabbix_api"
```

### 2. **Integrações Avançadas**

```bash
# Integração automática com Zabbix Web
warpctl pinchtab zabbix-integration

# Integração automática com TeamChat Web
warpctl pinchtab teamchat-integration

# Monitoramento integrado completo
warpctl zabbix integrated-monitoring --auto-login --screenshot
```

### 3. **Monitoramento Automatizado**

```bash
# Configurar mapeamentos interativamente
warpctl zabbix configure

# Ver status completo do sistema
warpctl zabbix status

# Iniciar monitoramento com automação web
warpctl zabbix start

# Teste de notificação com HTML enterprise
warpctl zabbix test-notify --show-html
```

## 🔧 Casos de Uso

### **Caso 1: Monitoramento VDC_BONJA**

```bash
# 1. Configurar grupo VDC_BONJA
warpctl zabbix configure
# Selecionar: 312 (VDC_BONJA)
# Canal: Horus Monitoramento

# 2. Verificar problemas ativos
warpctl zabbix problems --group 4296

# 3. Iniciar monitoramento integrado
warpctl zabbix integrated-monitoring --auto-login
```

### **Caso 2: Automação Web Completa**

```bash
# 1. Iniciar Pinchtab
./pinchtab.exe --port 8080

# 2. Login automático no Zabbix
warpctl pinchtab zabbix-integration

# 3. Navegar e interagir
warpctl pinchtab navigate https://monitoramento.armazem.cloud
warpctl pinchtab snapshot
warpctl pinchtab click e10  # Clicar em elemento específico

# 4. Tirar screenshot para auditoria
warpctl pinchtab screenshot audit-$(date +%s).jpg
```

### **Caso 3: Notificações Enterprise**

```bash
# 1. Testar formato HTML
warpctl zabbix test-notify --show-html

# 2. Enviar alerta real
warpctl zabbix notify --group 4296

# 3. Monitoramento contínuo
warpctl zabbix start
```

## 🎨 Formato de Notificação Enterprise

As notificações são enviadas em **HTML enterprise** com:

- ✅ **Template profissional** com tema escuro
- ✅ **Cores da marca** (verde #00b18e, navy #13222e)
- ✅ **Layout responsivo** com Tailwind CSS
- ✅ **Dados técnicos estruturados**
- ✅ **Botões de ação interativos**
- ✅ **Footer administrativo** com Event ID e timestamp
- ✅ **Extração automática** de usuário do hostname
- ✅ **Duração calculada** do incidente

## 🔐 Segurança

### **Pinchtab Security**
- Use `--token` para proteger acesso HTTP
- Trate `~/.pinchtab/` como diretório sensível
- Inicie com contas de baixo risco primeiro

### **Credenciais**
```bash
# Zabbix API
ZABBIX_URL="https://monitoramento.armazem.cloud/api_jsonrpc.php"
ZABBIX_USER="zabbix_api"
ZABBIX_PASS="kCfOtLUG3xSn2X9uyo6lXSVX1RFHGEw7"

# IceWarp
ICEWARP_USER="ranlens.denck@armazem.cloud"
ICEWARP_PASS="!@RanDenck321"
```

## 📊 Monitoramento em Tempo Real

### **Dashboard CLI**
```bash
# Status completo do sistema
warpctl zabbix status

# Problemas ativos por grupo
warpctl zabbix problems

# Saúde da integração
warpctl zabbix integrated-monitoring
```

### **Logs e Auditoria**
- Screenshots automáticos em erros
- Logs de navegação web
- Histórico de notificações
- Métricas de performance

## 🚀 Próximos Passos

1. **Instalar Pinchtab** e iniciar na porta 8080
2. **Configurar grupos** com `warpctl zabbix configure`
3. **Testar integração** com `warpctl zabbix integrated-monitoring`
4. **Iniciar monitoramento** com `warpctl zabbix start`

## 💡 Dicas Avançadas

### **Performance**
- Use `--quality 60` para screenshots menores
- Configure `interval: 2m` para monitoramento mais frequente
- Use filtros de severidade para reduzir ruído

### **Debugging**
```bash
# Debug mode
warpctl --debug zabbix integrated-monitoring

# Screenshots para debug
warpctl pinchtab screenshot debug.jpg --quality 100

# Árvore de acessibilidade detalhada
warpctl pinchtab snapshot
```

---

## ✅ **SISTEMA COMPLETO IMPLEMENTADO**

- 🎯 **Pinchtab Integration**: Browser automation via HTTP API
- 🔍 **Zabbix Monitoring**: 326 VDC groups, API integration
- 💬 **TeamChat Notifications**: HTML enterprise format
- ⚙️ **CLI Configuration**: Interactive group mapping
- 🚀 **Automated Workflows**: Web login, screenshots, monitoring
- 📊 **Real-time Dashboard**: Status, problems, health checks

**O sistema está pronto para uso em produção!**