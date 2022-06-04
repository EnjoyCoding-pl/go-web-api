package providers

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type jwtProvider struct {
	issuer string
	secret string
}

func NewJwtProvider(issuer string, secret string) *jwtProvider {
	return &jwtProvider{
		issuer: issuer,
		secret: secret,
	}
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId int
}

func (p *jwtProvider) Generate(userId int) (*string, error) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{

			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    p.issuer,
		},
		UserId: userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(p.secret))

	if err != nil {
		return nil, err
	}

	return &signedToken, nil
}

func (p *jwtProvider) Validate(token string) (*int, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(p.secret), nil
	})

	if err != nil {

		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, jwt.ErrTokenExpired) {
			return nil, jwt.ErrTokenExpired
		}
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*UserClaims)

	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return &claims.UserId, nil
}
