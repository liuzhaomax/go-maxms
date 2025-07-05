package core

import (
	"fmt"
)

func GenAppSignature(id string, secret string, userId string, nonce string) string {
	raw := fmt.Sprintf(`{"app_id":"%s","user_id":"%s","nonce":"%s"}`, id, userId, nonce)
	return HmacSHA256Str(BASE64EncodeStr(raw), secret)
}
