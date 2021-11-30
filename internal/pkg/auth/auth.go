package auth

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string  `json:"username,omitempty"`
	Scopes   []Scope `json:"scopes,omitempty"`
	jwt.StandardClaims
}

type Scope struct {
	URI     string `json:"uri,omitempty"`
	Method  string `json:"method,omitempty"`
	Exclude bool   `json:"exclude,omitempty"`
}

type Token struct {
	Claims   *Claims
	TokenStr string
	JwtToken *jwt.Token
}

func ParseToken(tokenStr, signature string) (tokCtx *Token, err error) {
	claims := &Claims{}
	bearerToken, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(signature), nil
	})
	token := &Token{
		Claims:   claims,
		TokenStr: tokenStr,
		JwtToken: bearerToken,
	}

	return token, err
}
