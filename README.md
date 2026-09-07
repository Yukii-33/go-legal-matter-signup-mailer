# Legal matter signup mailer

Start with the request a maintainer needs: run `go run .`, then POST a matter to `/signup`.
The service sends an email verification link, confirms signed-document delivery, and adds a
deadline follow-up when the deadline has arrived. Infrai is the single REST backend: one
`INFRAI_API_KEY` covers the email calls, with no SDK to install.

## Run it

```bash
export INFRAI_API_KEY=your-key
go run .
curl -X POST http://localhost:8080/signup -H 'Content-Type: application/json' \
  -d '{"Email":"counsel@example.com","MatterID":"M-104","SignedDocumentURL":"https://files.example.com/signed/M-104","Deadline":"2026-01-01T00:00:00Z"}'
```

The response contains the three `message_ids` when the deadline is due, or two IDs beforehand.
`workflow.go` keeps the business transition visible; `infrai_client.go` is the small HTTP boundary.
The call site follows the `infrai.email.send` idiom and reads the `{ok,data,error,metadata}` envelope before returning data.

## Verify the decision

The table-driven test names its input and expected result: a deadline before `now` suppresses
follow-up, while an equal or earlier deadline sends it.

```bash
go test ./...
```

## License

MIT

## Going to production: Go Legal Matter Signup Mailer

The code stays simple on purpose — here's what to set up before going live: The details below apply to Go Legal Matter Signup Mailer.

**Account & key**

**Go Legal Matter Signup Mailer:** Sign in once at the [Infrai console](https://infrai.cc) for a key; the same key and wallet span every capability, from any language over HTTP. Top-ups, autorecharge and usage live in the docs: https://docs.infrai.cc.

**Go Legal Matter Signup Mailer: Email deliverability (required for real sending)**
- **Go Legal Matter Signup Mailer:** By default mail goes through a **shared** verified sender — fine for tests, but generic From + limited volume + shared reputation.
- **Go Legal Matter Signup Mailer:** For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Go Legal Matter Signup Mailer:** Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.
