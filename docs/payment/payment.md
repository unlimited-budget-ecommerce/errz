# Payment Errors

| Code | Msg | HTTP | Category | Severity | Retryable |
|:------:|:-----:|:------:|:----------:|:----------:|:-----------:|
| PM0001 | insufficient balance | 402 | business | medium | false |
| PM0002 | payment gateway timeout | 504 | timeout | high | true |

---

## PM0001

- **Message**: insufficient balance
- **Cause**: user has not enough balance
- **Solution**: ask user to top-up or choose another method
- **HTTP Status**: 402
- **Category**: business
- **Severity**: medium
- **Retryable**: false
- **Tags**: `payment`, `balance`

## PM0002

- **Message**: payment gateway timeout
- **Cause**: no response from payment gateway
- **Solution**: retry after a short delay
- **HTTP Status**: 504
- **Category**: timeout
- **Severity**: high
- **Retryable**: true
- **Tags**: `payment`, `timeout`, `external`
