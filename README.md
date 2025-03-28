# スーパー支払い君.com

## Run
- go run main.go

## Test
- go test

## DB
- sqlite (in memory)

## Endpoints
- POST `/api/invoices`
- Request Body

| フィールド | コメント |
| --- | --- |
| client_id | 必須。フロントエンドは既知と想定。|
| amount | 必須。|
| due_date | 必須。YYYY-MM-DD |

- GET `/api/invoices`

- Request Query

| フィールド | コメント |
| --- | --- |
| start_date | 空許容 |
| due_date | 空許容 |

- 以下は認証のための追加です。叩くとJWTがもらえます。
    - `/api/login/company/1`
    - `/api/login/company/2`