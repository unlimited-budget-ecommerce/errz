# Common Errors

| Code | Message | Severity | Retryable |
|:-----:|:-----------:|:-----:|:-----:|
| CM0000 | success | low | false |
| CM0400 | bad request | medium | false |
| CM0500 | internal server error | high | true |

---

## CM0000

- **Domain**: common
- **Code**: CM0000
- **Message**: success
- **Cause**: operation completed successfully
- **Severity**: low
- **Retryable**: false

## CM0400

- **Domain**: common
- **Code**: CM0400
- **Message**: bad request
- **Cause**: invalid input or malformed request
- **Severity**: medium
- **Retryable**: false

## CM0500

- **Domain**: common
- **Code**: CM0500
- **Message**: internal server error
- **Cause**: unexpected server-side error
- **Severity**: high
- **Retryable**: true
