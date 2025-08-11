package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	//"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now().UTC()
	expireBy := now.Add(expiresIn)
	jwtNow := jwt.NewNumericDate(now)
	jwtExpire := jwt.NewNumericDate(expireBy)
	claims := jwt.RegisteredClaims{
		Issuer: "chirpy", IssuedAt: jwtNow,
		ExpiresAt: jwtExpire, Subject: userID.String(),
	}
	unsigned_tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed_tok, err := unsigned_tok.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signed_tok, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	jwt, err := jwt.ParseWithClaims(
		tokenString, &claims,
		func(token *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		},
		jwt.WithIssuer("chirpy"), jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithLeeway(5*time.Second),
	)
	if err != nil {
		return uuid.Nil, err
	}
	//jwtExp, err := jwt.Claims.GetExpirationTime()
	//fmt.Printf("exp time %v \n", jwtExp.Time)
	jwtSubj, err := jwt.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	retUuid, err := uuid.Parse(jwtSubj)
	if err != nil {
		return uuid.Nil, err
	}
	return retUuid, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authRaw := headers.Get("Authorization")
	//fmt.Printf("GetBearerToken raw <%s> \n", authRaw)
	if len(authRaw) < 20 {
		return "", errors.New("authorization not found")
	}
	authS0, _ := strings.CutPrefix(authRaw, "Bearer")
	//fmt.Printf("GetBearerToken s0 <%s> \n", authS0)
	return strings.Trim(authS0, " \t"), nil
}

func MakeRefreshToken() (string, error) {
	buf := [32]byte{}
	n, err := rand.Read(buf[:])
	if err != nil || n < 32 {
		return "", errors.New("bad read rand")
	}
	//fmt.Printf(" MakeRefreshToken n %d err %v", n, err)
	ret := hex.EncodeToString(buf[:])
	return ret, nil
}
