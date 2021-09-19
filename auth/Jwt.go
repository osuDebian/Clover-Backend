package auth

import (
	"Bancho/userDB"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"strings"
	"time"
)

type JwtResponseStruct struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type Claims struct {
	Iss string  `json:"iss"` //iss: 토큰 발급자 (issuer) account.thftgr
	Sub string  `json:"sub"` //sub: 토큰 제목 (subject) client_id
	Aud string  `json:"aud"` //aud: 토큰 대상자 (audience) user_id
	Exp float64 `json:"exp"` //exp: 토큰의 만료시간 (expiraton), millisecond 언제나 현재 시간보다 이후로 설정되어있어야합니다.
	Nbf float64 `json:"nbf"` //nbf: Not Before 토큰 활성시간
	Iat float64 `json:"iat"` //iat: 토큰이 발급된 시간 (issued at), 이 값을 사용하여 토큰의 age 가 얼마나 되었는지 판단 할 수 있습니다.
	Jti string  `json:"jti"` //jti: JWT의 고유 식별자로서, 주로 중복적인 처리를 방지하기 위하여 사용됩니다. 일회용 토큰에 사용하면 유용합니다.
	MapClaim jwt.MapClaims
}

const rediskey = `config.jwt.thftgr` // 15

func GetClaim(token *string) (claim *Claims) {
	s := strings.Split(*token, ".")
	if len(s) != 3 {
		return
	}
	b, err := base64.RawStdEncoding.DecodeString(s[1])
	if err != nil {
		b, err = base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			b, err = base64.RawURLEncoding.DecodeString(s[1])
			if err != nil {
				b, err = base64.URLEncoding.DecodeString(s[1])
				if err != nil {
					return
				}
			}
		}
	}
	var claims Claims
	_ = json.Unmarshal(b, &claims)
	_ = json.Unmarshal(b, &claims.MapClaim)
	return &claims
}

func GenerateJwt(userid string) (token JwtResponseStruct, err error) {

	mySigningKey := []byte(userDB.Redis[15].HGet(context.TODO(), rediskey, "key").Val())

	now := time.Now()

	jwtClaims := jwt.MapClaims{
		"aud":    "thftgr",
		"jti":    hex.EncodeToString(sha256.New().Sum([]byte(uuid.New().String()))),
		"iat":    now.Unix(),
		"nbf":    now.Unix(),
		"exp":    now.Add(time.Hour * 24).Unix(),
		"sub":    userid,
		"scopes": []string{"*"},
	}

	t, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims).SignedString(mySigningKey)

	token = JwtResponseStruct{
		TokenType:    "Bearer",
		ExpiresIn:    86400,
		AccessToken:  t,
		RefreshToken: hex.EncodeToString(sha256.New().Sum([]byte(uuid.New().String()))),
	}
	err = userDB.Redis[1].HSet(context.TODO(), "account.token:"+jwtClaims["jti"].(string),
		"accessToken", token.AccessToken,
		"refreshToken", token.RefreshToken,
		"iat", now.Unix(),
		"nbf", now.Unix(),
		"exp", now.Add(time.Hour*24).Unix(),
		"user_login_id", userid,
	).Err()

	// 리프레시 타임아웃 지정 2w
	err = userDB.Redis[1].Expire(context.TODO(), "account.token:"+jwtClaims["jti"].(string), time.Hour*24*14).Err()
	if err != nil {
		return
	}
	return

}
