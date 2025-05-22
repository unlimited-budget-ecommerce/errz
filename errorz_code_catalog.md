# Error code

Format: `[DOMAIN_PREFIX][4_DIGIT_CODE]` → ex. `PM0001`

## Domain Prefix

| Prefix | Domain                              |
| :----: | :---------------------------------- |
|   CM   | Common (e.g., success, bad request) |
|   AU   | Auth / Authorization                |
|   PM   | Payment                             |
|   OD   | Order                               |
|   PR   | Product                             |
|   US   | User                                |
|   IV   | Inventory                           |
|   DB   | Database                            |
|   GW   | Gateway / API Layer                 |
|   EX   | External Services                   |

## All Code

### 🟢 Common (CM)

|  Code  | Description                    |
| :----: | :----------------------------- |
| CM0000 | Success                        |
| CM0400 | Bad Request                    |
| CM0500 | Internal Server Error          |
| CM0404 | Resource Not Found             |
| CM0401 | Unauthorized                   |
| CM0403 | Forbidden                      |
| CM0429 | Too Many Requests (Rate limit) |

### 🔐 Auth (AU)

|  Code  | Description         |
| :----: | :------------------ |
| AU0001 | Invalid credentials |
| AU0002 | Token expired       |
| AU0003 | Token malformed     |

### 💳 Payment (PM)

|  Code  | Description             |
| :----: | :---------------------- |
| PM0001 | Insufficient balance    |
| PM0002 | Internal error          |
| PM0003 | Payment gateway timeout |
| PM0004 | Invalid payment method  |

### 🛒 Order (OD)

|  Code  | Description           |
| :----: | :-------------------- |
| OD0001 | Product out of stock  |
| OD0002 | Invalid order ID      |
| OD0003 | Payment not confirmed |
