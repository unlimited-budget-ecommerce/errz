# Common Errors

| Code | Msg | HTTP | Category | Severity | Retryable |
|:------:|:-----:|:------:|:----------:|:----------:|:-----------:|
| CM0000 | success | 200 | business | low | false |
| CM0400 | bad request | 400 | validation | medium | false |
| CM0500 | internal server error | 500 | internal | high | true |

---

## CM0000

- **Message**: success
- **Cause**: operation completed successfully
- **Solution**: no action needed
- **HTTP Status**: 200
- **Category**: business
- **Severity**: low
- **Retryable**: false
- **Tags**: `success`

## CM0400

- **Message**: bad request
- **Cause**: invalid input or malformed request
- **Solution**: check input format and required fields
- **HTTP Status**: 400
- **Category**: validation
- **Severity**: medium
- **Retryable**: false
- **Tags**: `client`, `input`

## CM0500

- **Message**: internal server error
- **Cause**: unexpected server-side error
- **Solution**: check logs and trace ID
- **HTTP Status**: 500
- **Category**: internal
- **Severity**: high
- **Retryable**: true
- **Tags**: `server`, `bug`
