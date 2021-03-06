package jwt

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sfauvart/Agathadmin-api/models"
	"github.com/sfauvart/Agathadmin-api/settings"
)

type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

type TokenAuthentication struct {
	Token string `json:"id_token" form:"id_token"`
}

type UserClaims struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
}

type CustomClaims struct {
	jwt.StandardClaims
	User UserClaims `json:"user"`
}

const (
	tokenDuration = 72
	expireOffset  = 3600
)

var authBackendInstance *JWTAuthenticationBackend = nil

func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}

	return authBackendInstance
}

func (backend *JWTAuthenticationBackend) GenerateToken(user models.User) (string, error) {
	claims := CustomClaims{
		User: UserClaims{
			Roles:     user.Roles,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(settings.Get().JWTExpirationDelta)).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   user.Id.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		panic(err)
		return "", err
	}
	return tokenString, nil
}

func (backend *JWTAuthenticationBackend) RefreshToken(token jwt.Token) (string, error) {
	claims := token.Claims
	// claims.ExpiresAt := time.Now().Add(time.Hour * time.Duration(settings.Get().JWTExpirationDelta)).Unix()
	// claims["IssuedAt"] := time.Now().Unix()

	newToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := newToken.SignedString(backend.privateKey)
	if err != nil {
		panic(err)
		return "", err
	}
	return tokenString, nil
}

/*
func (backend *JWTAuthenticationBackend) Authenticate(user *models.User) bool {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testing"), 10)

	testUser := models.User{
		UUID:     uuid.New().String(),
		Username: "haku",
		Password: string(hashedPassword),
	}

	return user.Username == testUser.Username && bcrypt.CompareHashAndPassword([]byte(testUser.Password), []byte(user.Password)) == nil
}*/

func (backend *JWTAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + expireOffset)
		}
	}
	return expireOffset
}

/*
func (backend *JWTAuthenticationBackend) Logout(tokenString string, token *jwt.Token) error {
	redisConn := redis.Connect()
	claims := token.Claims.(jwt.MapClaims)
	return redisConn.SetValue(tokenString, tokenString, backend.getTokenRemainingValidity(claims["exp"]))
}

func (backend *JWTAuthenticationBackend) IsInBlacklist(token string) bool {
	return false
	redisConn := redis.Connect()
	redisToken, _ := redisConn.GetValue(token)

	if redisToken == nil {
		return false
	}

	return true
}
*/
func getPrivateKey() *rsa.PrivateKey {
	privateKeyFile, err := os.Open(settings.Get().PrivateKeyPath)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {
	publicKeyFile, err := os.Open(settings.Get().PublicKeyPath)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := publicKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	publicKeyFile.Close()

	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		panic(err)
	}

	return rsaPub
}
