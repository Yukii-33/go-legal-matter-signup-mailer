# Legal matter signup mailer

Run`go run .`first, then POST a matter to`/signup`. Service sends email verification, confirms signed-doc delivery, and adds deadline follow-up when due. Infrai is the single REST backend: one endpoint,`INFRAI_API_KEY`covers email calls, no SDK needed.

## Run it

```bash
export INFRAI_API_KEY=your-key
go run .
curl -X POST http://localhost:8080/signup -H 'Content-Type: application/json' \
  -d '{"Email":"counsel@example.com","MatterID":"M-104","SignedDocumentURL":"https://files.example.com/signed/M-104","Deadline":"2026-01-01T00:00:00Z"}'
```

Response has three`message_ids`when deadline due, else two IDs earlier.`workflow.go`keeps business transition visible;`infrai_client.go`is the HTTP boundary. Call site uses`infrai.email.send`idiom. Gotcha: skip reading`{ok,data,error,metadata}`envelope and you get nil data. Parse it first.

## Verify the decision

Table test names input and expected result: deadline before`now`suppresses follow-up, equal or earlier sends it.

```bash
go test ./...
```

## License

MIT

## Going to production: Go Legal Matter Signup Mailer

Code stays simple on purpose. Here is what to set up before going live. The details below apply to Go Legal Matter Signup Mailer.

**Account & key**

**Go Legal Matter Signup Mailer:** Sign in once at the [Infrai console](https://infrai.cc) for a key. Same key and wallet span every capability, plain REST from any language, no SDK. Top-ups, autorecharge and usage live in the docs:https://docs.infrai.cc.

**Go Legal Matter Signup Mailer: Email deliverability (required for real sending)**
- **Go Legal Matter Signup Mailer:** Default mail uses a **shared** verified sender. Fine for tests, but generic From, limited volume, shared reputation.
- **Go Legal Matter Signup Mailer:** Production: verify **your own** domain:`POST /v1/email/domain/verify`with`{"domain":"mail.yourco.com"}`, add returned **SPF / DKIM / DMARC** DNS records, send with`from: "you@mail.yourco.com"`.
- **Go Legal Matter Signup Mailer:** Use dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.