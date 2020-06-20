package persistence

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
	// "time"

	"github.com/ghostdb/ghostdb-cache-node/store/lru"
	"github.com/ghostdb/ghostdb-cache-node/store/cache"
	"github.com/ghostdb/ghostdb-cache-node/config"
)

const (
	SNAPSHOT_FILENAME = "snapshot.gz"
)

func CreateSnapshot(cache *cache.Cache, config *config.Configuration) (bool, error) {
	switch (*cache).(type) {
	case *lru.LRUCache:
		return createLruSnapshot((*cache).(*lru.LRUCache), config.EnableEncryption, config.Passphrase)
	default:
		return false, nil
	}
}

func createLruSnapshot(cache *lru.LRUCache, encryption bool, passphrase ...string) (bool, error) {
	serialized, _ := json.MarshalIndent(cache, "", " ")

	configPath, _ := os.UserConfigDir()
	snapshotPath := configPath + SNAPSHOT_FILENAME

	if _, err := os.Stat(snapshotPath); err == nil {
		os.Remove(snapshotPath)
	}

	f, err := os.OpenFile(snapshotPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		f, err = os.Create(snapshotPath)
	}

	w, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		log.Printf("failed to create new snapshot writer: %s", err.Error())
		return false, err
	}

	if encryption {
		encryptedData, err := EncryptData(serialized, passphrase[0])
		if err != nil {
			w.Close()
			log.Printf("failed to encrypt snapshot: %s", err.Error())
			return false, err
		}
		w.Write(encryptedData)
	} else {
		w.Write(serialized)
	}
	w.Close()
	log.Println("successfully created snapshot")
	return true, nil
}

// GetSnapshotFilename builds the filename for the snapshot being taken
func GetSnapshotFilename() string {
	return SNAPSHOT_FILENAME
}


// BuildCache rebuilds the cache from the byte stream of the snapshot
func BuildCacheFromSnapshot(bs *[]byte) (lru.LRUCache, error) {
	// Create a new cache instance.
	var cache lru.LRUCache

	// Create a new configuration object
	var config config.Configuration = config.InitializeConfiguration()
	
	// Unmarshal the byte stream and update the new cache object with the result.
	err := json.Unmarshal(*bs, &cache)
	
	if err != nil {
		log.Fatalf("failed to rebuild cache from snapshot: %s", err.Error())
	}

	cache.Config = config

	// Create a new doubly linked list object
	ll := lru.InitList()

	// Populate the caches hashtable and doubly linked list with the values 
	// from the unmarshalled byte stream
	for _, v := range cache.Hashtable {
		n, err := lru.Insert(ll, v.Key, v.Value, v.TTL)
		if err != nil {
			return lru.LRUCache{}, err
		}
		cache.Hashtable[v.Key] = n
	}

	// Reset the watchdog
	// wdMetricInterval := time.Duration(config.WatchdogMetricInterval)
	// TODO: Add proper handling for this.
	// cache.Watchdog = monitor.Boot(wdMetricInterval, config.EntryTimestamp)

	cache.DLL = ll

	return cache, nil
}

// ReadSnapshot reads the compressed snapshot file into
// buffer and returns a reference to the buffer
func ReadSnapshot(encryption bool, passphrase ...string) *[]byte {

	configPath, _ := os.UserConfigDir()
	snap, err := os.Open(configPath + SNAPSHOT_FILENAME)
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