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
	PASSWORD string = "TestTest1234567890"
	SECRET   string = "TestTest12345678909843434"
	JSON     string = `{"role":"admin","email":"test@test.com"}`
)

func init() {
	auth = &Authenticator{}
	secret32 := sha256.Sum256([]byte(SECRET))
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

	testCookie, err = auth.encryptJWT([]byte(JSON))
	if err != nil {
		log.Fatalf("*** failed to create encrypted test jwt***")
	}
}

// TestAuthenticator_GenerateHash compares whether the GenerateHash function
// returns the proper base64 encoded hash
func TestGenerateBase64Hash(t *testing.T) {

	h64, err := generateBase64Hash([]byte(PASSWORD))
	if err != nil {
		t.Errorf("failed to generate a bcrypt salted hash - %v", err)
	}

	h, err := base64.StdEncoding.DecodeString(h64)
	if err != nil {
		t.Errorf("failed to decode base64 encoded hash - %v", err)
	}

	err = bcrypt.CompareHashAndPassword(h, []byte(PASSWORD))
	if err != nil {
		t.Errorf("expected: nil received: %v", err)
	}
}

func TestValidatePassword(t *testing.T) {

	h64, err := generateBase64Hash([]byte(PASSWORD))
	if err != nil {
		t.Errorf("failed to generate base64 encoded salted hash - %v", err)
	}

	t.Run("valid", func(t *testing.T) {
		result, err := validatePassword(h64, string([]byte(PASSWORD)))
		if !result || err != nil {
			t.Errorf("expected: true received: %t", result)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		result, err := validatePassword(h64, "Teadsfaf")
		if result || err == nil {
			t.Errorf("failed to reject invalid password - %v", err)
		}
	})
}

func TestAuthenticator_encryptJWT(t *testing.T) {

	json := []byte(JSON)

	// get the base64 encoded, encrypted JWT
	result, err := auth.encryptJWT(json)
	if err != nil {
		t.Errorf("failed to encrypt jwt with aes-256 cipher - %v", err)
	}

	resultBytes, err := base64.StdEncoding.DecodeString(result)
	if err != nil {
		t.Errorf("failed to decode base64 result - %v", err)
	}

	nonce := resultBytes[len(resultBytes)-NONCEBYTES:]

	e := auth.aesGCM.Seal(nil, nonce, json, nil)
	eB64 := base64.StdEncoding.EncodeToString(e)

	rParsed := base64.StdEncoding.EncodeToString(
		resultBytes[:len(resultBytes)-NONCEBYTES])

	if eB64 != rParsed {
		t.Errorf("expected: %s received: %s", eB64, rParsed)
	}
}

func TestAuthenticator_decryptJWT(t *testing.T) {

	expected := JSON

	encryptedJWT, err := auth.encryptJWT([]byte(JSON))
	if err != nil {
		t.Errorf("failed to generate encrypted jwt - %v", err)
	}

	jwt, err := base64.StdEncoding.DecodeString(encryptedJWT)
	if err != nil {
		t.Errorf("failed to decode base64 encoded, encrypyed jwt - %v", err)
	}

	nonce := jwt[len(jwt)-NONCEBYTES:]
	cipherText := jwt[:len(jwt)-NONCEBYTES]

	plain, err := auth.aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		t.Errorf("failed to decrypt aes-256 message - %v", err)
	}

	if string(plain) != expected {
		t.Errorf("failed to properly decrypt reference string - %v",
			string(plain))
	}

	result, err := auth.decryptJWT(encryptedJWT)
	if err != nil {
		t.Errorf("failed to decrypt jwt - %v", err)
	}

	resultText := string(result)

	if expected != resultText {
		t.Errorf("expected: %s received: %s", expected, resultText)
	}

}

func TestAuthenticator_getNonce(t *testing.T) {
	expected := NONCEBYTES
	result := len(getNonce(NONCEBYTES))

	if expected != result {
		t.Errorf("expected: %d received: %d", expected, result)
	}
}

func TestIsValidPassword(t *testing.T) {

	t.Run("valid", func(t *testing.T) {
		goodPassword := "TTT323dsa432!!"
		if result, _ := isValidPassword(goodPassword); !result {
			t.Errorf("expected: true received: %t", result)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		badPasswords := [...]string{"Tddfdafadfasdf", "31341421241", "TTTddafadfadfa!"}

		for i, password := range badPasswords {
			if result, _ := isValidPassword(password); result {
				t.Errorf("expected: false (%s) and received: %t",
					badPasswords[i], !result)
			}
		}
	})

}

func TestIsValidEmail(t *testing.T) {

	t.Run("invalid", func(t *testing.T) {
		if got, _ := isValidEmail("test"); got {
			t.Errorf("accepted invalid email - %v", "test")
		}
	})

	t.Run("valid", func(t *testing.T) {
		if got, _ := isValidEmail("test@test.com"); !got {
			t.Errorf("failed to allow a valid email - %v", "test@test.com")
		}
	})

}

func TestAllow(t *testing.T) {

	t.Run("goodJWT", func(t *testing.T) {
		result, err := auth.allow(testCookie)
		if err != nil && !result {
			t.Errorf("failed to authenticate properly encrypted JWT - %v", err)
		}
	})

	t.Run("badJWT", func(t *testing.T) {
		result, err := auth.allow("reallyPoorJWT")
		if err != nil && result {
			t.Errorf("failed authorization. expected: false received: %t", result)
		}
	})
}

func TestRedis(t *testing.T) {
	// test pass if unable to connect to Redis
	_, err := auth.keyStore.Ping().Result()
	if err != nil {
		t.Skipf("failed to connect to redis - %v", err)
	}

	jwt, err := auth.encryptJWT([]byte(JSON))
	if err != nil {
		t.Fatalf("failed to generate encrypted jwt - %v", err)
	}

	t.Run("Add", func(t *testing.T) {
		// check add
		_, err = auth.keyStore.Set(jwt, "admin", 0).Result()
		if err != nil {
			t.Errorf("failed to set key-value \"%s\":\"admin\"", jwt)
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
			t.Logf("redis delete failed (%v)", err)
			t.Fail()
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
