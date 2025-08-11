# WhatsApp Go SDK

SDK idiomático em Go para integração com a **WhatsApp Cloud API (Meta Graph)**, projetado com **Arquitetura Hexagonal**, **SOLID** e foco em **observabilidade**, **testabilidade** e **portabilidade** de segredos/tokens.

> **Status**: Documento inicial de projeto. Sem código ainda — esta versão descreve organização, responsabilidades e pontos de extensão.

---

## Objetivos

* Fornecer uma **fachada única** (Client) para operações da WhatsApp Cloud API.
* Manter **modelos de domínio puros** (sem acoplamento a HTTP/Graph).
* Separar **transporte** (Graph API), **serviços** (casos de uso) e **armazenamento de segredos/tokens** via **ports & adapters**.
* Permitir múltiplas fontes de configuração (arquivo, env, Redis, Vault) sem alterar o domínio.
* Entregar **erros tipados**, **retry/backoff**, **timeouts**, métricas e tracing.

---

## Estrutura do Projeto

```
whatsapp-sdk-go/
├── cmd/
│   └── cli/                          # Utilitários opcionais (depuração, healthcheck)
├── docs/                             # Notas de design, decisões e guias
├── examples/                         # Exemplos executáveis mínimos (futuros)
├── internal/
│   ├── httpx/                        # Cliente HTTP interno (retry, backoff, circuit breaker)
│   └── telemetry/                    # Métricas, logs e tracing (adapters)
├── pkg/
│   ├── whatsapp/                     # API pública do SDK
│   │   ├── client.go                 # Facade: expõe serviços de alto nível
│   │   ├── options.go                # Configuração do Client (Version, IDs, timeouts, providers)
│   │   ├── errors.go                 # Erros tipados (HTTPError, GraphError, ValidationError)
│   │   ├── domain/                   # Modelos puros
│   │   │   ├── common.go             # Tipos base (IDs, timestamps, paginação)
│   │   │   ├── messages.go           # Modelos de envio/recebimento de mensagens
│   │   │   ├── media.go              # Metadados de mídia
│   │   │   ├── phone_numbers.go      # Phone Number, quality, display name
│   │   │   ├── registration.go       # Registro/deregistro/2FA
│   │   │   └── webhooks.go           # Payloads e eventos de webhook
│   │   ├── ports/                    # Portas (interfaces) hexagonais
│   │   │   ├── HTTPDoer.go           # Abstração de transporte HTTP
│   │   │   ├── TokenProvider.go      # Fornecedor de access token (memória, arquivo, Redis, Vault)
│   │   │   └── SecretsProvider.go    # Fornecedor de segredos (verify token, app secret)
│   │   ├── services/                 # Casos de uso orientados ao domínio
│   │   │   ├── messages_service.go   # Envio de mensagens
│   │   │   ├── phone_service.go      # Consulta de phone numbers
│   │   │   ├── registration_service.go   # Registro/verificação/2FA
│   │   │   └── webhook_service.go    # Validação/parsing de webhooks
│   │   └── transport/graph/          # Adaptadores para Meta Graph API
│   │       ├── endpoints.go          # Templating de paths por versão (vXX.Y)
│   │       ├── messages_api.go       # /{Phone-Number-ID}/messages
│   │       ├── phone_numbers_api.go  # /{WABA-ID}/phone_numbers, /{Phone-Number-ID}
│   │       └── registration_api.go   # register/deregister/request_code/verify_code/2FA
│   └── storage/                      # Implementações concretas dos providers
│       ├── token/
│       │   ├── memory.go             # Token em memória
│       │   ├── file.go               # Token em arquivo seguro
│       │   ├── redis.go              # Token em Redis (TTL)
│       │   └── oci_vault.go          # Token via OCI Vault
│       └── secrets/
│           ├── env.go                # Segredos via variáveis de ambiente
│           ├── file.go               # Segredos em arquivo
│           ├── redis.go              # Segredos em Redis
│           └── oci_vault.go          # Segredos via OCI Vault
├── testdata/                         # Fixtures JSON (respostas e payloads de webhook)
├── go.mod
└── README.md
```

### Motivações da Estrutura

* **Hexagonal / Ports & Adapters**: `ports` define contratos; `transport/graph` e `storage/*` são adapters.
* **Domínio limpo**: `domain` contém exclusivamente tipos e regras do negócio (sem HTTP/JSON externos).
* **Coesão por contexto**: `services` agrupam casos de uso (messages, phone, registration, webhook).
* **Encapsulamento de transporte**: `internal/httpx` concentra resiliência (retry, timeouts, CB) e telemetria.

---

## Mapeamento: Endpoints → Serviços → Modelos

* **Mensagens**

    * **Endpoint**: `/{Version}/{Phone-Number-ID}/messages`
    * **Serviço**: `MessagesService`
    * **Modelos**: `domain.Message`, `domain.MessageResponse`, `domain.ErrorItem`

* **Phone Numbers**

    * **Endpoints**: `/{Version}/{WABA-ID}/phone_numbers`, `/{Version}/{Phone-Number-ID}`
    * **Serviço**: `PhoneService`
    * **Modelos**: `domain.PhoneNumber`, `domain.PhoneNumberList`

* **Registro / Verificação / 2FA**

    * **Endpoints**: `register`, `deregister`, `request_code`, `verify_code`, `POST /{Phone-Number-ID}` (2FA)
    * **Serviço**: `RegistrationService`
    * **Modelos**: `domain.RegistrationRequest`, `domain.VerificationCodeRequest`, `domain.OperationStatus`

* **Webhooks**

    * **Entrega**: payloads com `entry`, `changes`, `messages`, `statuses`, `errors`
    * **Serviço**: `WebhookService` (validação de verify token, parsing e utilitários)
    * **Adapter**: `transport/webhook` handler HTTP para verificação e dispatch
    * **Modelos**: `domain.WebhookEvent`, `domain.MessageEvent`, `domain.StatusEvent`

> Observação: nomes de modelos são ilustrativos e serão finalizados ao implementar a serialização oficial.

---

## Configuração e Provedores

### Options do Client

* **Version** (ex.: `v20.0`)
* **WABA ID** e **Phone Number ID**
* **HTTPDoer** (transporte) com timeout padrão
* **TokenProvider** (Access Token)
* **SecretsProvider** (Verify Token, App Secret)

### TokenProvider (ports)

* **memory**: útil em dev/testes
* **file**: arquivo seguro (legível pelo processo)
* **redis**: cache com TTL e invalidação
* **oci\_vault**: busca baseada em secret names/OCIDs

### SecretsProvider (ports)

* **env**: variáveis de ambiente
* **file**: JSON/YAML protegido
* **redis**: chaves rotativas e verify tokens
* **oci\_vault**: segredos críticos (produção)

---

## Resiliência e Observabilidade

* **Retry/Backoff**: para 429/5xx com jitter
* **Timeouts**: por requisição e total
* **Circuit Breaker**: proteção contra cascatas
* **Métricas**: latência por endpoint, taxa de erro, retries
* **Tracing**: spans por chamada (OpenTelemetry-friendly)
* **User-Agent**: identificável (`ampere-whatsapp-sdk/<versão>`)

---

## Modelo de Erros

* **HTTPError**: status, headers, fb-trace-id, corpo bruto
* **GraphError**: `code`, `type`, `message` e dados adicionais
* **ValidationError**: entradas inválidas antes do transporte

Padrão de propagação: serviço → facade, sem vazar detalhes de `internal/httpx`.

---

## Testes

* **Unitários**: mocks de `HTTPDoer`, `TokenProvider`, `SecretsProvider`.
* **Fixtures**: `testdata/` com respostas reais da collection e payloads de webhook.
* **Integração (opcional)**: smoke em ambiente sandbox, controlado por flags/vars de ambiente.

---

## Roadmap (alta prioridade)

1. Definição de **interfaces** em `ports/` (HTTPDoer, TokenProvider, SecretsProvider).
2. **endpoints.go** com strategy de versionamento.
3. `MessagesService` e `messages_api` (primeiro caso de uso completo).
4. `PhoneService` (list/get) e `RegistrationService` (register/verify/2FA).
5. Modelo de **erros tipados** e **telemetria** mínima.

---

## Contribuição

* Pull requests devem manter o isolamento entre **domínio**, **transporte** e **armazenamento**.
* Evitar dependências desnecessárias; preferir standard library quando possível.
* Cobertura de testes para serviços e parsing.

---

## Licença

Definir antes do primeiro release público.
