package authorizer

import "github.com/golang-jwt/jwt"

const OlivePublicKey = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAqbE/7LnSaigDC7NQ2NUZ
YtNulnrdrE4u/jC4V2wM6Rod5Q6ndW2W9gMa3ZVJT3dya/8jJY6jLL8lHMVkRRNL
HOIrzKogJlLt9KnyjasHJRTUUvTF93+aj/E1ekEW025PyI4huf4TrbSVwckHxsB3
fW3Mc2Yos3VY+/bwtISnWDU3eRmNi1pREyO4s0GL4ai3THTBK26vAW8fdBUu1OcP
7ptM2J45y43CzyQ/yUFvIuHZPSZaYSVUv/dts3j7PLnn6Sr//D95u3zudSTmbQ+9
c6jclCgOXEXGcxWr8dg1MyjyXOruxllmRymWeyXjJSayfym0r5JfRbOAN0emEYA0
XUdmQXDucrqSRCrQToEajt/6KG9ZsiCOUGsx9+OH6jSZA+8ooeq4wQxqHTrmZQpA
kNvIX+xGbk7Hsx76Urb6OPW4Qr5pHlR0HAVd1czpOZED7ziuSnYw1JZFQlK9k3FP
ponWliHrJe9HGKof3LwAyIxpLLWN3DVlo7ESrupYnCydOKUqNO5jAePwz27vfuMN
bhrHBxaWwnDuDUVKO9fD0dOsa1y8ACRz72t7cJ4rIdlKVFGWoTXnvGNe+tLXxJ6N
pua52Ip1qLWslSs9UjuIh0n1XXBuFASpd6qvCJuSCec/dfIUT4/G9Ua1KUkj+7wF
mIlx8N8x4QLWH3C8zKA5w30CAwEAAQ==
-----END PUBLIC KEY-----
`

type OliveClaims struct {
	AuthorizedParty string `json:"azp,omitempty"`
	Email           string `json:"email,omitempty"`
	jwt.StandardClaims
}
