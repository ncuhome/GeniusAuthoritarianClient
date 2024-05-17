package jwt

import (
	"crypto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ncuhome/GeniusAuthoritarianClient/jwt/jwtClaims"
	"time"
)

const (
	Refresh = "Refresh"
	Access  = "Access"
)

type Parser struct {
	PublicKey crypto.PublicKey
}

func ParseToken[C jwtClaims.Claims](publicKey crypto.PublicKey, Type, token string, target C) (claims C, valid bool, err error) {
	var t *jwt.Token
	t, err = jwt.ParseWithClaims(
		token, target, func(t *jwt.Token) (interface{}, error) {
			return publicKey, nil
		},
		jwt.WithLeeway(time.Second*3),
	)
	if err != nil {
		return
	}

	claims, _ = t.Claims.(C)
	valid = t.Valid && claims.GetType() == Type
	return
}

func (p Parser) ParseRefreshToken(token string) (*jwtClaims.RefreshToken, bool, error) {
	return ParseToken(p.PublicKey, Refresh, token, &jwtClaims.RefreshToken{})
}

func (p Parser) ParseAccessToken(token string) (*jwtClaims.AccessToken, bool, error) {
	return ParseToken(p.PublicKey, Access, token, &jwtClaims.AccessToken{})
}
