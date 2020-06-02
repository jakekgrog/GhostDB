package lru

import (
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
)

const (
	SnitchLogFilename = "/ghostdb/snapshot.gz"
)

// CreateSnapshot creates a JSON serialized snapshot of the cache
// and compresses it
func CreateSnapshot(cache *LRUCache, encryption bool, passphrase ...string) {
	serialized, _ := json.MarshalIndent(cache, "", " ")

	configPath, _ := os.UserConfigDir()
	if _, err := os.Stat(configPath + SnitchLogFilename); err == nil {
		os.Remove(configPath + SnitchLogFilename)
	}

	f, err := os.OpenFile(configPath+SnitchLogFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		f, err = os.Create(configPath + SnitchLogFilename)
	}

	w, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		log.Fatalf("failed to create new snapshot writer: %s", err.Error())
	}

	if encryption {
		encryptedData, err := EncryptData(serialized, passphrase[0])
		if err != nil {
			w.Close()
			log.Fatalf("failed to encrypt snapshot: %s", err.Error())
		}
		w.Write(encryptedData)
	} else {
		w.Write(serialized)
	}
	w.Close()
}

// GetSnapshotFilename builds the filename for the snapshot being taken
func GetSnapshotFilename() string {
	return "snapshot.gz"
}

// ReadSnapshot reads the compressed snapshot file into
// buffer and returns a reference to the buffer
func ReadSnapshot(encryption bool, passphrase ...string) *[]byte {

	configPath, _ := os.UserConfigDir()
	snap, err := os.Open(configPath+SnitchLogFilename)
	if err != nil {
		log.Fatalf("failed to open snapshot: %s", err.Error())
	}

	defer snap.Close()

	file, err := gzip.NewReader(snap)

	if err != nil {
		log.Fatalf("failed to create gzip reader: %s", err.Error())
	}

	byteStream, _ := ioutil.ReadAll(file)

	if encryption {
		serializedData, err := DecryptData(byteStream, passphrase[0])
		if err != nil {
			log.Fatalf("failed to decrypt snapshot: %s", err.Error())
		}
		return &serializedData
	} else {
		return &byteStream
	}
}

// EncryptData is our encryption client to encrypt the serialized
// cache object with 128 bit AES
func EncryptData(data []byte, passphrase string) ([]byte, error) {
	gcm, err := newGCMCipher(passphrase)
	if err != nil {
		return nil, err
	}

	// Create a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// Populate the nonce with a cryptographically secure
	// random sequence
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.New("Failed to populate nonce")
	}
	// Encrypt the serialized cache object using Seal
	// Seal encrypts the authenticated data. The nonce
	// must be NonceSize() bytes long and unique for all time,
	// for a given key. Seal authenticates additional data
	// and appends the result to the destination (nonce)
	cipherText := gcm.Seal(nonce, nonce, data, nil)
	return cipherText, nil
}

// DecryptData decrypts snapshots which have at-rest encryption enabled.
func DecryptData(data []byte, passphrase string) ([]byte, error) {
	gcm, err := newGCMCipher(passphrase)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt and authenticate the ciphertext. Authenticate the 
	// additional data and if successful, append the resulting data
	// to the destination. The nonce must be NonceSize() bytes long
	// and both it and the additional data must match the value passed
	// to Seal.
	serializedBytes, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("Failed to read encrypted snapshot")
	}

	return serializedBytes, nil
}

// newGCMCipher generates a new AES-GCM cipher object
func newGCMCipher(passphrase string) (cipher.AEAD, error) {
	hash := generateHash(passphrase)

	// Generate a new aes cipher using 32 byte key
	block, err := aes.NewCipher([]byte(hash))
	if err != nil {
		return nil, errors.New("Failed to generate AES cipher")
	}

	// gcm is a mode of operation for symmetric key
	// cryptographic block ciphers
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("Failed to create block cipher")
	}

	return gcm, nil
}

func generateHash(passphrase string) string {
	// Hash the password using MD5 to ensure key is always 32 bytes.
	// MD5 is not secure, but it doesn't matter, we're not storing
	// the result.
	hasher := md5.New()
	hasher.Write([]byte(passphrase))
	// Return the hash as a hexidecimal value.
	return hex.EncodeToString(hasher.Sum(nil))
}