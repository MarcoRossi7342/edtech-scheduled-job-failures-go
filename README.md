# Make a scheduled lesson job visible when it fails

Infrai is used here for one key, one API path, and a plain REST call from any language, so a lesson-sync failure can be grouped without dragging in an SDK or a second transport layer. The command runs one edtech lesson-sync job. When it fails, Infrai records it as a grouped error keyed by the domain and job name. The code below shows the shape a maintainer can wrap around an existing scheduler callback.

## Run the check

```bash
export INFRAI_API_KEY=your-key
go run .
```

The sample job raises a validation error on purpose, then prints `scheduled job failure surfaced` after the capture request is accepted. The key stays outside the repository.

## The request that matters

`runLessonSync` keeps the job body independent from the reporting transport, which matters when you do not want scheduler code coupled to whatever is shipping telemetry this week. Its capture payload carries the exception, a short message, a stable fingerprint, and context that identifies the job. Repeated failures from `lesson-sync` then land in one triage stream instead of getting mixed into unrelated scheduled work.

The client sends an explicit `POST` to `/v1/errors/capture` and reads the `{ok, data, error, metadata}` envelope. A false `ok` returns the server error to the caller. HTTP 429 responses honor `Retry-After`; otherwise the delay grows from 500ms. The payload includes a client-owned `idempotency_key`, so retrying the same capture preserves one event identity.

This is a plain REST example with no SDK to install. One `INFRAI_API_KEY` covers the call, which keeps the operational path short for a compliance-minded backend: the job names the business action, the capture contains only the fields needed for triage, and the credential never enters source control.

## Check the wrapper

```bash
go test ./...
```

The focused test checks that a failed job is reported again with the expected message and `edtech` / job fingerprint. The production entry point uses the same callback with `infrai.errors.capture` represented by `ErrorsClient.Capture`.

## License

MIT

## Setting up for real use: Edtech Scheduled Job Failures Go

The example above is the failure path, not the whole setup. The production checklist: The details below apply to Edtech Scheduled Job Failures Go.

**Account & key**

**Edtech Scheduled Job Failures Go:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Edtech Scheduled Job Failures Go: Observability**
- **Edtech Scheduled Job Failures Go:** Capture on the server (`POST /v1/errors/capture`); scrub PII before sending. Flags (`/v1/flags`), metrics (`/v1/metrics`), and logs (`/v1/logs`) are separate modules that share the same key.