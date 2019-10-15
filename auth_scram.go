package xmpp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/pbkdf2"
	"hash"
	"strconv"
	"strings"
)

type ScramAuth struct {
	Username       string
	Password       string
	clientNonce    string
	initialMessage string
	hashFunc       func() hash.Hash
	hashSize       int
}

func NewScramSHA1Authenticator(creds Credentials) Authenticator {
	return &ScramAuth{
		Username:    creds.Username,
		Password:    creds.Password,
		clientNonce: uuid.New().String(),
		hashFunc:    sha1.New,
		hashSize:    sha1.Size,
	}
}

func (auth *ScramAuth) Name() string {
	return "SCRAM-SHA-1"
}

func (auth *ScramAuth) Data() string {
	auth.initialMessage = fmt.Sprintf("n=%s,r=%s", auth.Username, auth.clientNonce)
	return base64.StdEncoding.EncodeToString([]byte("n,," + auth.initialMessage))
}

func (auth *ScramAuth) hash(data []byte) []byte {
	h := auth.hashFunc()
	h.Write(data)
	return h.Sum(nil)
}

func (auth *ScramAuth) hmac(key []byte, data []byte) []byte {
	h := hmac.New(auth.hashFunc, key)
	h.Write(data)
	return h.Sum(nil)
}

func (auth *ScramAuth) Challenge(challenge string) string {
	firstMessageBytes, _ := base64.StdEncoding.DecodeString(challenge)
	serverFirstMessage := string(firstMessageBytes)

	serverNonce, salt, iterations := decodeChallenge(serverFirstMessage)

	if !strings.HasPrefix(serverNonce, auth.clientNonce) {
		return ""
	}

	clientFinalMessageBare := fmt.Sprintf("c=biws,r=%s", serverNonce)
	saltedPassword := pbkdf2.Key([]byte(auth.Password), salt, iterations, auth.hashSize, auth.hashFunc)
	clientKey := auth.hmac(saltedPassword, []byte("Client Key"))
	storedKey := auth.hash(clientKey)
	authMessage := []byte(fmt.Sprintf("%s,%s,%s", auth.initialMessage, serverFirstMessage, clientFinalMessageBare))
	clientSignature := auth.hmac(storedKey, authMessage)
	clientProof := xorBytes(clientKey, clientSignature)
	clientProof64 := base64.StdEncoding.EncodeToString(clientProof)
	//serverKey := auth.hmac(saltedPassword, []byte("Server Key"))
	//serverSignature := auth.hmac(serverKey, authMessage)
	clientFinalMessage := fmt.Sprintf("%s,p=%s", clientFinalMessageBare, clientProof64)

	return base64.StdEncoding.EncodeToString([]byte(clientFinalMessage))
}

func decodeChallenge(data string) (r string, s []byte, i int) {
	parts := strings.Split(data, ",")
	for _, p := range parts {
		l := strings.SplitN(p, "=", 2)
		switch l[0] {
		case "r":
			r = l[1]
		case "s":
			if bytes, err := base64.StdEncoding.DecodeString(l[1]); err == nil {
				s = bytes
			} else {
				panic(err)
			}
		case "i":
			if v, err := strconv.Atoi(l[1]); err == nil {
				i = v
			}
		}
	}
	return
}

func xorBytes(a, b []byte) []byte {
	if len(a) != len(b) {
		panic("xorBytes: different lengths")
	}
	res := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = a[i] ^ b[i]
	}
	return res
}
