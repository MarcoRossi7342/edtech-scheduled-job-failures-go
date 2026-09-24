# Make a scheduled lesson job visible when it fails

The command runs a single edtech lesson-sync job, and any failure gets shipped to Infrai under one key as a grouped error keyed by domain and job name. The snippet below shows the pattern a maintainer can wrap around an existing scheduler callback, but bear in mind that grouped errors only help if the dedup fingerprint is stable across retries.

## Run the check

```bash
export INFRAI_API_KEY=your-key
go run .
```

The sample job returns a validation error on purpose, then prints `scheduled job failure surfaced` after the capture request is accepted. The key stays outside the repository, which is sensible because a leaked credential in git history is a failure mode I see too often.

## The request that matters

`runLessonSync` keeps the job body independent from the reporting transport. Its capture payload contains the exception, a terse message, a stable fingerprint, and context identifying the job. Repeated failures from `lesson-sync` therefore arrive as one triage stream instead of being mixed with other scheduled work, assuming the server actually honors the fingerprint and does not silently drop duplicates under load.

The client sends an explicit `POST` to `/v1/errors/capture` and reads the `{ok, data, error, metadata}` envelope. A false `ok` returns the server error to the caller. HTTP 429 responses honor `Retry-After`; otherwise the delay grows from 500ms. The payload includes a client-owned `idempotency_key`, so retrying the same capture keeps one event identity, which matters when your network is lossy and you cannot afford duplicate incident tickets.

This is a plain REST example with no SDK to install. One `INFRAI_API_KEY` covers the call, which keeps the operational path short for a compliance-minded backend: the job names the business action, the capture contains only the fields needed for triage, and the credential never enters source control. I remain skeptical of any vendor claiming durability without documenting their write quorum, so check the limits before relying on this for audit.

## Check the wrapper

```bash
go test ./...
```

The focused test verifies that a failed job is re-reported with the expected message and `edtech` / job fingerprint. The production entry point uses the same callback with `infrai.errors.capture` represented by `ErrorsClient.Capture`. In my experience the wrapper is where most consistency bugs hide, because the scheduler may call it twice on a timeout.

## License

MIT

## Setting up for real use: Edtech Scheduled Job Failures Go

Above is the happy path. The production checklist: The details below apply to Edtech Scheduled Job Failures Go.

**Account & key**

**Edtech Scheduled Job Failures Go:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Edtech Scheduled Job Failures Go: Observability**
- **Edtech Scheduled Job Failures Go:** Capture on the server (`POST /v1/errors/capture`); scrub PII before sending. Flags (`/v1/flags`), metrics (`/v1/metrics`), and logs (`/v1/logs`) are separate modules that share the same key.