package services

import (
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dgrijalva/jwt-go"
)

func GetJWTToken(jwtString string) (*jwt.Token, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(jwtString, jwt.MapClaims{})

	if err != nil {
		log.Println("Could not parse JWT token")
		return nil, err
	}

	return token, nil
}

func getJWTClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := GetJWTToken(tokenString)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("Failed to retrive claims")
	}

	return claims, nil
}

func GetTokenFromHeader(request events.APIGatewayProxyRequest) (string, error) {
	tokenHeader, ok := request.Headers["Authorization"]

	if !ok {
		return "", errors.New("Authorization header not found")
	}

	tokenHeader = strings.Replace(tokenHeader, "Bearer ", "", -1)
	return tokenHeader, nil
}

func GetJWTEmail(tokenString string) (string, error) {
	claims, err := getJWTClaims(tokenString)

	if err != nil {
		return "", err
	}

	return claims["email"].(string), nil
}

func GetJWTPublicKey(tokenString string) (string, error) {
	claims, err := getJWTClaims(tokenString)

	if err != nil {
		return "", err
	}

	return claims["publicKey"].(string), nil
}
