# WhatsApp Go SDK — Architecture Guide

> Status: MVP-ready · Target audience: contributors and integrators

## 1) Architectural style & principles

**Hexagonal (Ports & Adapters)**

* **Domain-first**: data models and use cases live under `pkg/whatsapp/domain`.
* **Ports**: narrow interfaces in `pkg/whatsapp/ports` isolate the core from infra.
* **Adapters**: Graph API HTTP adapter in `pkg/whatsapp/transport/graph` and the internal HTTP executor in `internal/httpx` implement the ports.
* **Services**: application services in `pkg/whatsapp/services` orchestrate domain + ports, no vendor-specific logic.

**Design principles**

* **SOLID**: small interfaces (ISP), injected dependencies (DIP), single-responsibility packages and files (SRP).
* **Pure domain**: JSON and HTTP concerns are kept out of the `domain` package.
* **Fail fast**: strict validation with explicit error types before I/O.
* **Testable-by-construction**: everything behind ports; fakes/mocks live under `internal/testutils`.

## 2) High-level module map

```
cmd/cli                 # CLI for manual flows (examples/dev)
examples/               # Minimal usage samples
internal/httpx          # HTTP Doer with retry/backoff
internal/testutils/...  # In-memory fakes for unit tests
pkg/errorsx             # Error types (HTTP/Graph/Validation)
pkg/whatsapp
  ├─ client.go          # Facade that composes services & adapters
  ├─ options.go         # Client configuration (validated)
  ├─ domain/            # Pure models (messages, phone, webhooks)
  ├─ ports/             # Hexagonal ports (HTTPDoer, TokenProvider, ...)
  ├─ services/          # Application services (messages, phone, reg, webhook)
  └─ transport/graph/   # WhatsApp Graph adapter (endpoints, requests)
```

## 3) Package responsibilities

### `pkg/whatsapp/domain`

Pure value objects and payload shapes used across the SDK. No HTTP or persistence logic. Keep stable and additive.

### `pkg/whatsapp/ports`

Narrow interfaces that the core depends on:

* `HTTPDoer` — minimal HTTP executor abstraction.
* `TokenProvider` — supplies/refreshes the OAuth token.
* `SecretsProvider` — retrieves non-rotating secrets (verify token, app secret).
* `PhoneAPI`, `RegistrationAPI` — Graph operations expressed as domain I/O.

### `pkg/whatsapp/services`

Application services orchestrate domain use-cases. They accept a minimal client/core or specific APIs via ports. Examples:

* **MessagesService**: validates inputs, builds request via transport, executes using `HTTPDoer`, decodes domain.
* **PhoneService**: delegates listing and fetching numbers to `PhoneAPI`.
* **RegistrationService**: wraps register flows against `RegistrationAPI`.
* **WebhookService**: verify token & HMAC signature + parse webhook payload.
* **WebhookDispatcher**: fan-out parsed events to handler callbacks.

### `pkg/whatsapp/transport/graph`

Adapter that knows the Graph endpoints, URL layout, and request/response encoding. Builds `*http.Request` objects; never calls the network directly.

### `internal/httpx`

Shared HTTP Doer with bounded retries, jittered exponential backoff, and context deadlines. No external deps.

### `pkg/errorsx`

Error taxonomy:

* `ValidationError` — pre-flight client-side issues.
* `HTTPError` — non-2xx with HTTP context (status, headers, trace id).
* `GraphError` — decoded Graph error envelope layered on the HTTP error.

### `pkg/whatsapp/client.go`

Public facade that composes everything:

* Validates `Options`, applies sane defaults.
* Builds default `HTTPDoer` when none is provided.
* Wires `transport/graph` adapters and constructs services.
* Exposes convenience getters used by services.

## 4) Typical request flows

### 4.1 Send a text message

1. `MessagesService.SendText` validates inputs (E.164, non-empty body).
2. Builds Graph payload (`transport/graph`).
3. Acquires token via `TokenProvider` and executes HTTP via `HTTPDoer`.
4. On 2xx, decodes into `domain.MessageSendResponse`; on error, returns `GraphError`/`HTTPError`.

### 4.2 Registration flows (request/verify/register/deregister)

* Services call `RegistrationAPI` adapter methods with domain params.
* Same transport path and error semantics as messages.

### 4.3 Webhook validation & parsing

* **Verification GET**: compare provided verify token vs secret from `SecretsProvider`.
* **POST**: verify `X-Hub-Signature-256` HMAC using app secret; then `ParseEvent` to `domain.WebhookEvent`.
* Optionally pass event to `WebhookDispatcher` and your app-level handler.

## 5) Errors & resilience

* **Validation first**: reject bad inputs before network.
* **Retries**: 429 and most 5xx are retried (cap attempts & backoff); bodies must be rewindable to enable retries.
* **Graph error decoding**: attempt to parse error envelope; keep raw bytes when decoding fails.
* **Traceability**: propagate `fb-trace-id` when present.

## 6) Configuration (`Options`)

* `Version` (e.g., `v20.0`), `WABAID`, `PhoneNumberID` — required.
* `TokenProvider` — required; **never** log the token.
* `SecretsProvider` — required for webhook features.
* `HTTPDoer` — optional; defaults to `internal/httpx` with `RetryMax`.
* `BaseURL`, `Timeout`, `RetryMax`, `UserAgent` — optional tuning knobs.

## 7) Testing strategy

**Unit tests**

* Services are tested with in-memory fakes (`internal/testutils/whatsapp/ports`).
* `internal/httpx` has isolated tests for retry, backoff, timeouts, and body rewind behavior.
* Fixtures in `testdata` emulate Graph responses.

**Practices**

* Prefer fakes over mocks; assert on behaviors and decoded models.
* Use tiny backoffs/timeouts in tests for speed.
* Keep transport tests black-box (serve httptest servers).

## 8) Concurrency & safety

* `HTTPDoer` and providers **must be safe for concurrent use**.
* Services are stateless; clients can be shared across goroutines.
* Avoid global state; inject dependencies.

## 9) Telemetry & observability (hooks)

* Wrap `internal/httpx` with your tracer/metrics if needed.
* Consider logging request ids and `fb-trace-id` from `HTTPError/GraphError`.

## 10) Versioning & compatibility

* Keep `domain` models backwards-compatible (additive fields).
* Avoid breaking changes in `ports` (treat as public contracts).
* Bump major if ports change; otherwise minor/patch.

## 11) Security

* Do not print tokens/secrets. Mask identifiers in logs.
* Keep secrets out of repo; use `SecretsProvider`.
* Validate all inputs; verify webhook signatures.

## 12) MVP scope

* **Messaging**: send text.
* **Phone**: list/get numbers.
* **Registration**: request/verify/register/set PIN/deregister.
* **Webhook**: token validation, signature verification, parse + dispatcher.

## 13) Extension guidelines

* When adding a new Graph feature:

    1. **Domain**: add models.
    2. **Ports**: extend or add a new API interface.
    3. **Transport**: implement request/response builder for Graph.
    4. **Service**: orchestrate validation + call adapter.
    5. **Examples/CLI**: add a minimal usage sample.
    6. **Tests**: fakes/fixtures; service + adapter coverage.

---

### Quick checklists

**Adding a service method**

* [ ] Validate inputs → `errorsx.ValidationError` when invalid
* [ ] Build request via `transport/graph`
* [ ] Use `Client.Do` (applies token, user agent, timeout)
* [ ] Decode 2xx into domain
* [ ] Decode non-2xx → `GraphError` else `HTTPError`
* [ ] Unit tests with httptest + fixtures

**Implementing a new adapter call**

* [ ] Construct URL from `BaseURL` + `Version`
* [ ] Set headers (Authorization by `Client.Do`)
* [ ] Ensure request bodies are rewindable for retries
* [ ] Keep adapter free of business logic

**Webhook endpoint**

* [ ] GET: verify token
* [ ] POST: verify HMAC signature
* [ ] Parse event → dispatch
* [ ] Never log secret values
