package legacy

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/decred/dcrd/crypto/ripemd160"
)

// Base58 alphabet as used by Noso
const (
	debug       = false
	b58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

// hashMD160String returns hex RIPEMD160 of string's SHA256 digest
func hashMD160String(s string) string {
	h := ripemd160.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Convert hex string to base58
func hexToBase58(hexStr string) string {
	// Convert hex to bytes
	b, _ := hex.DecodeString(hexStr)
	return encodeBase58(b)
}

// Base58 encoding (basic)
func encodeBase58(input []byte) string {
	var result []byte
	x := new(big.Int).SetBytes(input)
	base := big.NewInt(58)
	zero := big.NewInt(0)
	mod := new(big.Int)

	for x.Cmp(zero) > 0 {
		x.DivMod(x, base, mod)
		result = append([]byte{b58Alphabet[mod.Int64()]}, result...)
	}

	// Preserve leading zeros
	for _, b := range input {
		if b != 0 {
			break
		}
		result = append([]byte{b58Alphabet[0]}, result...)
	}

	return string(result)
}

// Return a Base58 Checsum of 2 bytes/chars
func checksumBase58(data string) uint16 {
	sum := 0
	for _, ch := range data {
		sum += int(ch)
	}
	return uint16(sum & 0xFFFF)
}

// Simulate BMB58resumen: sum of ASCII values of the base58 string
func b58Resumen(s string) int64 {
	sum := 0
	for _, c := range s {
		idx := strings.IndexRune(b58Alphabet, c)
		if idx < 0 {
			return -1
		}
		sum += idx
	}
	return int64(sum)
}

// Convert decimal string to base58
func decimalToBase58(decimal int64) string {
	n := new(big.Int)
	n.SetString(strconv.FormatInt(decimal, 10), 10)
	base := big.NewInt(58)
	var result []byte
	mod := new(big.Int)

	for n.Cmp(big.NewInt(0)) > 0 {
		n.DivMod(n, base, mod)
		result = append([]byte{b58Alphabet[mod.Int64()]}, result...)
	}

	return string(result)
}

func isValidBase58(s string) bool {
	for _, ch := range s {
		if !strings.ContainsRune(b58Alphabet, ch) {
			return false
		}
	}
	return true
}

// GenerateNewAddress generates a SECP256k1 keypair and returns pub, priv, and address.
func GenerateNewAddress() (pubKeyStr string, privKeyStr string, address string) {
	for {
		// Step 1: Generate SECP256k1 keypair
		privKey, err := btcec.NewPrivateKey()
		if err != nil {
			panic("Failed to generate key: " + err.Error())
		}
		pubKey := privKey.PubKey()

		// Step 2: Serialize as hex
		pubKeyBytes := pubKey.SerializeUncompressed() // or compressed, depending on chain rules
		privKeyBytes := privKey.PubKey().SerializeUncompressed()
		pubKeyStr = hex.EncodeToString(pubKeyBytes)
		privKeyStr = hex.EncodeToString(privKeyBytes)

		// Step 3: Generate address
		address = GetAddressFromPublicKey(pubKeyStr, 1)

		// Step 4: Validate length
		if len(address) >= 20 {
			return
		}
	}
}

// GetAddressFromPublicKey generates an account has from a pubkey string
func GetAddressFromPublicKey(pubKey string, addType int) string {
	// Step 1: SHA256(pubkey)
	shaHash := sha256.Sum256([]byte(pubKey))

	// Step 2: RIPEMD160 of SHA256(pubkey)
	hashHex := hashMD160String(strings.ToUpper(hex.EncodeToString(shaHash[:])))

	// Step 3: Convert hashHex to Base58
	base58Hash := hexToBase58(hashHex)

	// Step 4: Compute "summary" (sum of ASCII chars), then to base58
	sum := b58Resumen(base58Hash)
	sumBase58 := decimalToBase58(sum)

	// Step 5: Concatenate and prepend 'N' or 'M'
	full := base58Hash + sumBase58
	prefix := "N"
	if addType == 1 {
		prefix = "M"
	}

	return prefix + full
}

func NewGetAddressFromPublicKey(pubKey string) string {
	// Step 1: SHA256(pubkey)
	shaHash := sha256.Sum256([]byte(pubKey))

	// Step 2: RIPEMD160(SHA256(pubkey))
	ripemd := ripemd160.New()
	ripemd.Write(shaHash[:])
	hash160 := ripemd.Sum(nil)
	hashHex := hex.EncodeToString(hash160)

	// Step 3: Convert hex (base16) to base58
	base58Hash := hexToBase58(hashHex) // same as B16toB58

	// Step 4: Sum ASCII values, convert sum (decimal) to base58
	sumStr := b58Resumen(base58Hash)     // gives string of an int
	sumBase58 := decimalToBase58(sumStr) // same as B10toB58

	// Step 5: Concatenate and prefix
	final := base58Hash + sumBase58
	return "N" + final
}

func FutureGetAddressFromPublicKey(pubKey string) string {
	if strings.TrimSpace(pubKey) == "" {
		return ""
	}

	// Step 1: SHA256(pubkey)
	shaHash := sha256.Sum256([]byte(pubKey))
	shaHex := hex.EncodeToString(shaHash[:])

	// Step 2: RIPEMD160 of SHA256 hex (as string, encoded in ANSI)
	ripemd := ripemd160.New()
	ripemd.Write([]byte(shaHex)) // Note: using ASCII string here, like TEncoding.ANSI
	hash160 := ripemd.Sum(nil)

	// Step 3: Base58-encode the RIPEMD160 hash
	base58Hash := encodeBase58(hash160)

	// Step 4: If starts with '1', strip it
	if strings.HasPrefix(base58Hash, "1") {
		base58Hash = base58Hash[1:]
	}

	// Step 5: Generate checksum (custom)
	checksumVal := checksumBase58(base58Hash)       // Returns uint16 or uint32
	checksumHex := fmt.Sprintf("%04X", checksumVal) // Hex string, 4 chars
	checksumBytes, _ := hex.DecodeString(checksumHex)

	// Step 6: Base58-encode checksum
	base58Checksum := encodeBase58(checksumBytes)

	// Step 7: Return final address
	return "N" + base58Hash + base58Checksum
}

// IsValidHashAddress check if an arbitrary string is a valid address
func IsValidHashAddress(address string) bool {
	if len(address) <= 20 {
		return false
	}

	firstChar := address[0]
	if firstChar != 'N' && firstChar != 'M' {
		return false
	}

	// Remove the prefix and checksum
	if len(address) < 4 {
		return false
	}

	body := address[1 : len(address)-2] // Excludes prefix and 2-char checksum

	if !isValidBase58(body) {
		return false
	}

	// Recalculate checksum
	checksumValue := b58Resumen(body)
	checksum := decimalToBase58(checksumValue)

	reconstructed := string(firstChar) + body + checksum
	return reconstructed == address
}
