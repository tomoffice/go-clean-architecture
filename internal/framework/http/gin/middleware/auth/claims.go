package auth

import "github.com/golang-jwt/jwt/v5"

type Claims[T any] struct {
	jwt.RegisteredClaims // ★ 匿名嵌入：iss/exp/… 會自動寫進來
	Private              T
}
