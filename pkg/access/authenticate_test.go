package access

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
	"log"
	"testing"
	"time"
)

var auth *Authenticator
var testCookie string

const (
	userPassword string = "TestTest1234567890"
	secret       string = "TestTest12345678909843434"
	testJWT      string = `{"role":"admin","email":"test@test.com"}`
)

func init() {
	auth = &Authenticator{}
	secret32 := sha256.Sum256([]byte(secret))
	secretKey := secret32[:]

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		log.Fatalf("Failed to create AES-256 cipher")
	}

	aesGCM, err := cipher.NewGCM(block)
	auth.aesGCM = aesGCM

	auth.keyStore = redis.NewClient(&redis.Options{
		Addr:       ":6379",
		Password:   "testtesttest1234567890", // no password
		DB:         0,                        // use default DB 0
		MaxRetries: 3,
	})

	_, err = auth.keyStore.Ping().Result()
	if err != nil {
		log.Printf("*** redis not connected ***\n")
	}

	testCookie, err = auth.encryptJWT([]byte(testJWT))
	if err != nil {
		log.Fatalf("*** failed to create encrypted test jwt***")
	}
}

// TestAuthenticator_GenerateHash compares whether the GenerateHash function
// returns the proper base64 encoded hash
func TestGenerateBase64Hash(t *testing.T) {

	h64, err := generateBase64Hash([]byte(userPassword))
	if err != nil {
		log.Fatal("GenerateHash method failed to generate a salted hash")
	}

	h, err := base64.StdEncoding.DecodeString(h64)
	if err != nil {
		log.Fatal("Failed to decode base64 encoded salted hash")
	}

	err = bcrypt.CompareHashAndPassword(h, []byte(userPassword))
	if err != nil {
		log.Fatalf("Test failed. Expected: nil Received: %v", err)
	}
}

func TestValidatePassword(t *testing.T) {

	h64, err := generateBase64Hash([]byte(userPassword))
	if err != nil {
		log.Fatal("Failed to generate base64 encoded salted hash")
	}

	// test good password
	result, err := validatePassword(h64, string([]byte(userPassword)))

	if !result {
		log.Fatalf("Expected: True Received: %t", result)
	}

	result, err = validatePassword(h64, "Teadsfaf")

	// test bad password
	if result {
		log.Fatalf("Expected: False, Received: %t", result)
	}
}

func TestAuthenticator_encryptJWT(t *testing.T) {

	json := []byte(testJWT)

	// get the base64 encoded, encrypted JWT
	result, err := auth.encryptJWT(json)
	if err != nil {
		log.Fatalf("Failed to encrypt JWT with AES-256 cipher (%v)", err)
	}

	resultBytes, err := base64.StdEncoding.DecodeString(result)
	if err != nil {
		log.Fatalf("Failed to decode Base64 encoded result (%v)", err)
	}

	nonce := resultBytes[len(resultBytes)-NONCEBYTES:]

	e := auth.aesGCM.Seal(nil, nonce, json, nil)
	eB64 := base64.StdEncoding.EncodeToString(e)

	rParsed := base64.StdEncoding.EncodeToString(
		resultBytes[:len(resultBytes)-NONCEBYTES])

	if eB64 != rParsed {
		log.Fatalf("Expected: %s Received: %s", eB64, rParsed)
	}
}

func TestAuthenticator_decryptJWT(t *testing.T) {

	expected := testJWT

	encryptedJWT, err := auth.encryptJWT([]byte(testJWT))
	if err != nil {
		log.Fatalf("Failed to generate encrypted JWT")
	}

	jwt, err := base64.StdEncoding.DecodeString(encryptedJWT)
	if err != nil {
		log.Fatalf("Failed to decode Base64 encoded, encrypyed JWT (%v)", err)
	}

	nonce := jwt[len(jwt)-NONCEBYTES:]
	cipherText := jwt[:len(jwt)-NONCEBYTES]

	plain, err := auth.aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		log.Fatalf("Failed to decrypt AES-256 message (%v)", err)
	}

	if string(plain) != expected {
		log.Fatalf("Test failed to properly decrypt reference string")
	}

	result, err := auth.decryptJWT(encryptedJWT)
	if err != nil {
		log.Fatalf("Failed to receive decrypted JWT (%v", err)
	}

	resultText := string(result)

	if expected != resultText {
		log.Fatalf("Expected: %s Received: %s", expected, resultText)
	}

}

func TestAuthenticator_getNonce(t *testing.T) {
	expected := NONCEBYTES
	result := len(getNonce(NONCEBYTES))

	if expected != result {
		log.Fatalf("Expected: %d Received: %d", expected, result)
	}
}

func TestIsValidPassword(t *testing.T) {
	// TODO: need to check for dictionary words
	goodPassword := "TTT323dsa432!!"
	badPasswords := [...]string{"Tddfdafadfasdf", "31341421241", "TTTddafadfadfa!"}

	// test valid passwords
	if result, _ := isValidPassword(goodPassword); !result {
		log.Fatalf("Expected: True Received: %t", result)
	}

	// test invalid passwords
	for i, password := range badPasswords {
		if result, _ := isValidPassword(password); result {
			log.Fatalf("Expected: False (%s) and Received: %t",
				badPasswords[i], !result)
		}
	}
}

func TestIsValidEmail(t *testing.T) {
	// TODO: add a file of good and bad email examples to run against
	if got, _ := isValidEmail("test"); got {
		log.Fatalf("Expected: False Received: %t", got)
	}

	if got, _ := isValidEmail("test@test.com"); !got {
		log.Fatalf("Expected: True (%s) Received: %t", "test@test.com", got)
	}
}

func TestAllow(t *testing.T) {

	result, err := auth.allow(testCookie)
	if err != nil && !result {
		log.Fatalf("Failed to determine whether authorized access (%v)", err)
	}
	if !result {
		log.Fatalf("Failed authorization. Expected: true Received: %t", result)
	}
}

func TestRedis(t *testing.T) {
	// test pass if unable to connect to Redis
	_, err := auth.keyStore.Ping().Result()
	if err != nil {
		t.SkipNow()
	}

	jwt, err := auth.encryptJWT([]byte(testJWT))
	if err != nil {
		log.Fatalf("Unable to generate encrypted JWT for test (%v)", err)
	}

	t.Run("Add", func(t *testing.T) {
		// check add
		_, err = auth.keyStore.Set(jwt, "admin", 0).Result()
		if err != nil {
			log.Fatalf("Failed to set key-value \"%s\":\"admin\"", jwt)
		}
	})

	t.Run("Read", func(t *testing.T) {
		val, err := auth.keyStore.Get(jwt).Result()
		if (err != nil) || (val != "admin") {
			log.Fatalf("Failed to read key-value \"%s\":\"admin\"", jwt)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		num, err := auth.keyStore.Del(jwt).Result()
		if (err != nil) || (num != 1) {
			log.Fatalf("Failed to delete key-value \"%s\":\"admin", jwt)
		}
	})

	t.Run("Timeout", func(t *testing.T) {
		// check expiration
		_, err = auth.keyStore.Set(jwt, "admin", time.Second*1).Result()
		if err != nil {
			log.Fatalf("Failed to set key-value \"%s\":\"admin\"", jwt)
		}

		time.Sleep(time.Second * 2)

		_, err = auth.keyStore.Get(jwt).Result()
		if err != redis.Nil {
			log.Fatalf("Failed to timeout key-value \"%s\":\"admin\"", jwt)
		}
	})

}
