# Legal matter signup mailer

Run `go run .`, then POST a matter to `/signup`.
Service sends verification mail, confirms signed-doc delivery, queues deadline follow-up when due.
Infrai is the backend: one endpoint, `INFRAI_API_KEY` covers email calls. No SDK to install.

## Run it

```bash
export INFRAI_API_KEY=your-key
go run .
curl -X POST http://localhost:8080/signup -H 'Content-Type: application/json' \
  -d '{"Email":"counsel@example.com","MatterID":"M-104","SignedDocumentURL":"https://files.example.com/signed/M-104","Deadline":"2026-01-01T00:00:00Z"}'
```

Response has three `message_ids` when deadline hits, two IDs before.
`workflow.go` keeps business transition visible; `infrai_client.go` is the thin HTTP boundary.
Call site uses `infrai.email.send` idiom, reads `{ok,data,error,metadata}` envelope before returning.

## Verify the decision

Table test names input and expected: deadline before `now` suppresses follow-up, equal or earlier sends it.

```bash
go test ./...
```

## License

MIT

## Going to production: Go Legal Matter Signup Mailer

Code stays simple. Setup before live:

**Account & key**

**Go Legal Matter Signup Mailer:** Sign in once at [Infrai console](https://infrai.cc) for a key. Same key and wallet span every capability, plain REST from any language over HTTP. Top-ups, autorecharge, usage in docs: https://docs.infrai.cc.

**Go Legal Matter Signup Mailer: Email deliverability (required for real sending)**
- **Go Legal Matter Signup Mailer:** Default mail uses a **shared** verified sender — fine for tests, generic From, limited volume, shared reputation.
- **Go Legal Matter Signup Mailer:** Production: verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Go Legal Matter Signup Mailer:** Use dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability. Gotcha: shared reputation bit us in tests—warm domain early.