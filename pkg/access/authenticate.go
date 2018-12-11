package access

import (
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"regexp"
	"time"
	"unicode"
)

const (
	NONCEBYTES int = 12
)

// Authenticate provides operations to create user accounts, validate user
// accounts and permissions, reset user accounts, issue bearer tokens, and
// restrict user accounts.

type AuthenticateService interface {
	encryptJWT(json []byte) (string, error)
	decryptJWT(encryptedToken string) ([]byte, error)
	allow(encryptedToken string) (bool, error)
}

type Authenticator struct {
	port     string
	userDb   *sql.DB
	keyStore *redis.Client
	secret   []byte // SHA-256 hash of secret key
	aesGCM   cipher.AEAD
}

type JWT struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// isValidPassword checks to determine a minimum password length of 12 chars.
// Password requirements include:
// * minimum of 12 characters, maximum of 50
// * minimum of 3 uppercase chars
// * minimum of 3 lowercase chars
// * minimum of 3 numbers
// * minimum of 1 special characters
func isValidPassword(password string) (bool, error) {
	// reject if password length is less than 12 characters
	var upper, lower, numbers, special = 0, 0, 0, 0

	chars := len(password)
	if !(chars >= 12 && chars <= 50) {
		return false, nil
	}

	for _, s := range password {
		switch {
		case unicode.IsLower(s):
			lower++
		case unicode.IsUpper(s):
			upper++
		case unicode.IsNumber(s):
			numbers++
		case unicode.IsPunct(s) || unicode.IsSymbol(s):
			special++
		default:
			return false, nil // illegal character detected
		}
	}

	if lower >= 3 && upper >= 3 && numbers >= 3 && special >= 1 {
		return true, nil
	} else {
		return false, nil
	}
}

// isValidEmail checks whether email is properly formatted.
// TODO: Build a more robust email authenticator
func isValidEmail(email string) (bool, error) {
	valid := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+" +
		"@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9]" +
		"(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	// TODO: return meaningful error values
	if len(email) > 254 || !valid.MatchString(email) {
		return false, nil
	} else {
		return true, nil
	}
}

//

// generateBase64Hash returns a bcrypt salted, Base64 encoded hash to store
// in the authentication database.
func generateBase64Hash(p []byte) (string, error) {
	h, err := bcrypt.GenerateFromPassword(p, 10)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(h), nil
}

// validatePassword checks the recorded Base64 encoded bcrypt, salted hash
// with the provided password
func validatePassword(hashB64 string, password string) (bool, error) {

	h, err := base64.StdEncoding.DecodeString(hashB64)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword(h, []byte(password))
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// generateNonce provides a byte array of random bytes for use in AES-256
// ciphers.
func getNonce(bytes int) []byte {
	nonce := make([]byte, bytes)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf("Failed to generate AES nonce (%v)", err)
	}

	return nonce
}

func (a *Authenticator) encryptJWT(json []byte) (string, error) {

	nonce := getNonce(NONCEBYTES)

	c := a.aesGCM.Seal(nil, nonce, json, nil)
	cNonce := append(c, nonce[:]...)
	cB64 := base64.StdEncoding.EncodeToString(cNonce)

	return cB64, nil
}

func (a *Authenticator) decryptJWT(encryptedToken string) ([]byte, error) {

	cipherJoined, err := base64.StdEncoding.DecodeString(encryptedToken)
	if err != nil {
		return nil, nil
	}

	nonce := cipherJoined[len(cipherJoined)-NONCEBYTES:]
	cipherText := cipherJoined[:len(cipherJoined)-NONCEBYTES]

	plainText, err := a.aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

func (a *Authenticator) allow(encryptedJWT string) (bool, error) {
	// check redis for key, value
	val, err := a.keyStore.Get(encryptedJWT).Result()
	if err != redis.Nil {
		log.Printf("(%v) Encrypted JWT not stored in Redis or Redis connection "+
			"invalid.", err)
		plainJwt, err := a.decryptJWT(encryptedJWT)
		if err != nil {
			log.Printf("Unable to decrypt JWT")
			return false, nil
		}

		// unmarshal plaintext JWT
		jwt := &JWT{}
		err = json.Unmarshal(plainJwt, jwt)
		if err != nil {
			log.Printf("Error (%v): Failed to unmarshal (%s) decrypyed JWT",
				err, plainJwt)
			return false, nil
		}

		// check if role properly defined in JWT
		if (jwt.Role == "admin") || (jwt.Role == "user") ||
			(jwt.Role == "facilitator") || (jwt.Role == "analyst") {

			err := a.keyStore.Set(encryptedJWT, jwt.Role, time.Hour).Err()
			if err != nil {
				log.Printf("(%v) Failed to store encrypted JWT in Redis", err)
			}
			return true, nil
		}

	}

	// return true and no errors is val in Redis key store is one of the following
	if val == "admin" || val == "user" || val == "facilitator" || val == "analyst" {
		return true, nil
	} else {
		return false, nil
	}
}
