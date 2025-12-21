package ctxdata

import "github.com/golang-jwt/jwt/v4"

const Identify = "perfric"

// secretKey: 用于签名的密钥
// iat: 令牌签发时间戳
// seconds: 令牌有效期（秒）
// uid: 用户唯一标识
func GetJwtToken(secretKey string, iat, seconds int64, uid string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[Identify] = uid

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}
