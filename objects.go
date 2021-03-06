package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	goutils "github.com/simonski/goutils"
)

// KPDB helper struct holds the data and keys
type KPDB struct {
	data               DB
	Filename           string
	PublicKeyFilename  string
	PrivateKeyFilename string
	EncryptionEnabled  bool
}

// DB is the thing that we serialise to JSON
type DB struct {
	Entries map[string]DBEntry
}

// DBEntry represents the a single item in the DB
type DBEntry struct {
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	LastUpdated time.Time `json:"lastUpdated"`
	Created     time.Time `json:"created"`
}

// NewKPDB constructor
func NewKPDB(filename string, pubKey string, privKey string, encryptionEnabled bool) *KPDB {
	cdb := KPDB{}
	cdb.Load(filename, pubKey, privKey, encryptionEnabled)
	return &cdb
}

// Load populates the db with the file
func (cdb *KPDB) Load(filename string, pubKey string, privKey string, encryptionEnabled bool) bool {
	cdb.Filename = goutils.EvaluateFilename(filename)
	cdb.PublicKeyFilename = goutils.EvaluateFilename(pubKey)
	cdb.PrivateKeyFilename = goutils.EvaluateFilename(privKey)
	cdb.EncryptionEnabled = encryptionEnabled

	if !goutils.FileExists(cdb.Filename) {
		db := DB{}
		db.Entries = make(map[string]DBEntry)
		cdb.data = db
	} else {
		jsonFile, err := os.Open(cdb.Filename)
		if err != nil {
			fmt.Printf("ERR %v\n", err)
			db := DB{}
			db.Entries = make(map[string]DBEntry)
			cdb.data = db
			// panic(err)
		} else {
			db := DB{}
			bytes, _ := ioutil.ReadAll(jsonFile)
			var data map[string]DBEntry
			json.Unmarshal(bytes, &data)
			db.Entries = data
			cdb.data = db
		}
	}

	return true
}

// Clear empties the db (without saving it)
func (cdb *KPDB) Clear() {
	cdb.data.Entries = make(map[string]DBEntry)
}

// Save writes the DB to disk
func (cdb *KPDB) Save() bool {
	data := cdb.data.Entries
	file, _ := json.MarshalIndent(data, "", " ")
	err := ioutil.WriteFile(cdb.Filename, file, 0644)
	if err != nil {
		fmt.Printf("%v", err)
	}
	return true
}

// GetData returns the data map of all key
func (cdb *KPDB) GetData() DB {
	return cdb.data
}

// Get returns the (DBEntry, bool) indicating it exists (or not)
func (cdb *KPDB) Get(key string) (DBEntry, bool) {
	entry, exists := cdb.data.Entries[key]
	if exists {
		decValue := entry.Value
		if cdb.EncryptionEnabled {
			decValue = cdb.Decrypt(entry.Value)
		}
		entry.Value = decValue
	}
	return entry, exists
}

// Put stores (or replaces) the key/value pair
func (cdb *KPDB) Put(key string, value string, description string) {
	entry, exists := cdb.data.Entries[key]
	encValue := value
	if cdb.EncryptionEnabled {
		encValue = cdb.Encrypt(value)
	}
	if exists {
		if value != "" {
			entry.Value = encValue
		}
		entry.LastUpdated = time.Now()
		if description != "" {
			entry.Description = description
		}
		cdb.data.Entries[key] = entry
	} else {
		entry = DBEntry{Key: key, Value: encValue, Created: time.Now(), LastUpdated: time.Now()}
		if description != "" {
			entry.Description = description
		}
		cdb.data.Entries[key] = entry
	}
}

// Delete removes the key/value pair from the DB
func (cdb *KPDB) Delete(key string) {
	delete(cdb.data.Entries, key)
}

// Encrypt helper function encrypts with public key
func (cdb *KPDB) Encrypt(value string) string {
	publicKey := LoadPublicKey(cdb.PublicKeyFilename)
	bytes := []byte(value)
	encrypted := EncryptWithPublicKey(bytes, publicKey)
	s := b64.StdEncoding.EncodeToString(encrypted)
	return s
}

// Decrypt helper function decrypts with private key
func (cdb *KPDB) Decrypt(value string) string {
	uDec, _ := b64.StdEncoding.DecodeString(value)
	privateKey := LoadPrivateKey(cdb.PrivateKeyFilename)
	bytes := []byte(uDec)
	decrypted := DecryptWithPrivateKey(bytes, privateKey)
	s := string(decrypted)
	return s
}
