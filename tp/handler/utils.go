package handler

import (
	"encoding/hex"
	"crypto/sha512"
	"strings"
)

var Namespace = HexDigest([]byte(FAMILY_NAME))[:FAMILY_NAMESPACE_ADDRESS_LENGTH]

var (
	FAMILY_NAME                     string = "audit"
	FAMILY_VERSION                  string = "1.0"
	FAMILY_NAMESPACE_ADDRESS_LENGTH uint   = 6
	FAMILY_VERB_ADDRESS_LENGTH      uint   = 64
)

func HexDigest(value []byte) string {
	hashHandler := sha512.New()
	hashHandler.Write(value)
	return strings.ToLower(hex.EncodeToString(hashHandler.Sum(nil)))
}

// getPrefix - generates the namespace prefix from constants
func getPrefix() string {
	return HexDigest([]byte(FAMILY_NAME))[:FAMILY_NAMESPACE_ADDRESS_LENGTH]
}

// getAddress - Return the namespaced address
func getAddress(name string) string {
	return getPrefix() + HexDigest([]byte(name))[FAMILY_VERB_ADDRESS_LENGTH:]
}