package security

import (
	"encoding/json"
	"fmt"
	"os"
	"pakawai_service/cmd/auth/model"
	"pakawai_service/configs"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/nu7hatch/gouuid"
)

func getKey() []byte {
	var key []byte
	if os.Getenv("APP_SECRET") == "" {
		key = []byte(os.Getenv("APP_SECRET"))
	} else {
		key = []byte("secret")
	}

	return key
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	AuthId string `json:"auth_id"`
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

// Create the Signin handler
func MakeJWT(creds *model.User) (string, error) {
	key := getKey()
	expirationTime := time.Now().Add(10 * time.Minute)
	id := creds.Id.Hex()
	v4, _ := uuid.NewV4()
	u, _ := uuid.NewV5(v4, []byte(id))
	claims := &Claims{
		UserId: id,
		AuthId: u.String(),
		StandardClaims: jwt.StandardClaims{
			Issuer:    id,
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	user, _ := json.Marshal(creds)
	configs.SetKey(u.String(), user, 10*time.Minute)

	return tokenString, err
}

func Verify(tokenString string) (bool, error, *Claims) {
	tknStr := tokenString
	key := getKey()
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, err, nil
		}
		return false, err, nil
	}
	if !tkn.Valid {
		return false, err, nil
	}

	return true, nil, claims
}
