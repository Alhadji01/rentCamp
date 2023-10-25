package helper

import (
	"rentcamp/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func generateToken(signKey string, id int, username string) string {
	var claims = jwt.MapClaims{}
	claims["id"] = id
	claims["username"] = username
	claims["role"] = "admin"
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()

	var sign = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	validToken, err := sign.SignedString([]byte(signKey))
	if err != nil {
		logrus.Error("JWT : claims isn't valid, ", err.Error())
		return ""
	}

	return validToken
}

func generateRefreshToken(signKey string, accessToken string) string {
	var claims = jwt.MapClaims{}
	claims["user"] = accessToken
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	var sign = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshToken, err := sign.SignedString([]byte(signKey))
	if err != nil {
		logrus.Error("JWT : claims isn't valid, ", err.Error())
		return ""
	}

	return refreshToken
}

func GenerateJWT(signKey string, refreshKey string, userId int, username string) map[string]any {
	var res = map[string]any{}

	var accessToken = generateToken(signKey, userId, username)
	if accessToken == "" {
		logrus.Error("JWT : Cannot generate token", nil)
		return nil
	}

	var refreshToken = generateRefreshToken(refreshKey, accessToken)
	if refreshToken == "" {
		logrus.Error("JWT : Cannot generate refresh token", nil)
		return nil
	}

	res["access_token"] = accessToken
	res["refresh_token"] = refreshToken

	return res
}

func RefreshJWT(accessToken string, refreshToken *jwt.Token, signKey string) map[string]any {
	var res = map[string]any{}

	expTime, err := refreshToken.Claims.GetExpirationTime()
	if err != nil {
		logrus.Error("JWT : get token expiration error, ", err.Error())
		return nil
	}

	if refreshToken.Valid && expTime.After(time.Now()) {
		var newClaim = jwt.MapClaims{}

		newToken, err := jwt.ParseWithClaims(accessToken, newClaim, func(t *jwt.Token) (interface{}, error) {
			return []byte(signKey), nil
		})

		if err != nil {
			logrus.Error("JWT error : ", err.Error())
			return nil
		}

		newClaim = newToken.Claims.(jwt.MapClaims)
		newClaim["iat"] = time.Now().Unix()
		newClaim["exp"] = time.Now().Add(time.Hour * 1).Unix()

		var newRefreshClaim = refreshToken.Claims.(jwt.MapClaims)
		newRefreshClaim["exp"] = time.Now().Add(time.Hour * 24).Unix()

		var newRefreshToken = jwt.NewWithClaims(refreshToken.Method, newRefreshClaim)
		newSignedRefreshToken, _ := newRefreshToken.SignedString(refreshToken.Signature)

		res["access_token"] = newToken.Raw
		res["refresh_token"] = newSignedRefreshToken

		return res
	}

	return nil
}

func ExtractToken(token *jwt.Token) any {
	if token.Valid {
		var claims = token.Claims

		expTime, _ := claims.GetExpirationTime()
		if expTime.After(time.Now()) {
			var MapClaim = claims.(jwt.MapClaims)
			newMap := map[string]any{}
			newMap["id"] = MapClaim["id"]
			newMap["username"] = MapClaim["username"]
			newMap["role"] = MapClaim["role"]
			return newMap
		}

		logrus.Error("JWT : token expired", nil)
		return nil
	}

	return nil
}

func Middleware() echo.MiddlewareFunc {
	cfg := config.Config{}
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(cfg.Secret),
		SigningMethod: "HS256",
	})
}
