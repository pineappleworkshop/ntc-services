package handlers

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func parsePagination(c echo.Context) (int64, int64, error) {
	var page int
	if c.QueryParam("page") != "" {
		var err error
		page, err = strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			c.Logger().Error(http.StatusInternalServerError, "Paginating page failed")
		}
	} else {
		page = 1
	}
	if page < 1 {
		err := errors.New("pagination page cannot be less than 1")
		c.Logger().Error(http.StatusInternalServerError, err)
		return -1, -1, err
	}

	var limit int
	if c.QueryParam("limit") != "" {
		var err error
		limit, err = strconv.Atoi(c.QueryParam("limit"))
		if err != nil {
			c.Logger().Error(http.StatusInternalServerError, "Paginating limit failed")
		}
	} else {
		limit = 100
	}

	if limit > 100 {
		err := errors.New("pagination limit cannot be greater the 100")
		c.Logger().Error(http.StatusInternalServerError, err)
		return -1, -1, err
	}
	if limit%20 != 0 {
		err := errors.New("pagination limit must be in increments of 20")
		c.Logger().Error(http.StatusInternalServerError, err)
		return -1, -1, err
	}

	return int64(page), int64(limit), nil
}

func validateBTCAddress(address string) bool {
	_, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	return err == nil
}

//func extractTapRootAddrFromPubKey(pkScript []byte) (string, error) {
//	// Check that the script is the correct length for P2TR (1-byte version + 1-byte length + 32-byte key)
//	if len(pkScript) != 34 {
//		return "", fmt.Errorf("invalid pkScript length for P2TR")
//	}
//
//	// Check that the version byte is 0x51 (indicating P2TR)
//	if pkScript[0] != 0x51 {
//		return "", fmt.Errorf("invalid version byte for P2TR")
//	}
//
//	// Check that the length byte is 0x20 (indicating 32 bytes)
//	if pkScript[1] != 0x20 {
//		return "", fmt.Errorf("invalid length byte for P2TR public key")
//	}
//
//	byte32TaprootPubKey, err := convertTo32ByteArray(pkScript[2:])
//	if err != nil {
//		return "", err
//	}
//
//	var bit5Data []byte
//	var currentByte byte
//	var nextByte byte
//	var bit5Value byte
//
//	for i := 0; i < 32; i++ {
//		currentByte = byte32TaprootPubKey[i]
//		if i < 31 {
//			nextByte = byte32TaprootPubKey[i+1]
//		}
//
//		for j := 0; j < 8; j += 5 {
//			if j == 0 {
//				bit5Value = currentByte >> 3
//			} else {
//				bit5Value = (currentByte << (j - 3)) | (nextByte >> (8 - j + 3))
//				bit5Value = bit5Value & 0x1F // mask to keep only 5 bits
//			}
//			bit5Data = append(bit5Data, bit5Value)
//		}
//	}
//
//	segwit, err := bech32.EncodeM("bc", bit5Data[:64])
//	if err != nil {
//		return "", err
//	}
//
//	return segwit, nil
//}

//func convertTo32ByteArray(data []byte) ([32]byte, error) {
//	var arr [32]byte
//	if len(data) != 32 {
//		return arr, fmt.Errorf("expected a slice of length 32, got %d", len(data))
//	}
//	copy(arr[:], data)
//	return arr, nil
//}
//
//func extractTapRootAddrFromPubKey(pkScript []byte) (string, error) {
//	// Check that the script is the correct length for P2TR (1-byte version + 1-byte length + 32-byte key)
//	if len(pkScript) != 34 {
//		return "", fmt.Errorf("invalid pkScript length for P2TR")
//	}
//
//	// Check that the version byte is 0x50 (indicating P2TR)
//	if pkScript[0] != 0x50 {
//		return "", fmt.Errorf("invalid version byte for P2TR")
//	}
//
//	// Check that the length byte is 0x20 (indicating 32 bytes)
//	if pkScript[1] != 0x20 {
//		return "", fmt.Errorf("invalid length byte for P2TR public key")
//	}
//
//	byte32TaprootPubKey, err := convertTo32ByteArray(pkScript[2:])
//	if err != nil {
//		return "", err
//	}
//
//	// Convert 32-byte public key to Bech32m
//	segwit, err := bech32.EncodeM("bc", append([]byte{0}, byte32TaprootPubKey[:]...))
//	if err != nil {
//		return "", err
//	}
//
//	return segwit, nil
//}
//
//func convertTo32ByteArray(data []byte) ([32]byte, error) {
//	var arr [32]byte
//	if len(data) != 32 {
//		return arr, fmt.Errorf("expected a slice of length 32, got %d", len(data))
//	}
//	copy(arr[:], data)
//	return arr, nil
//}
//
//
//func convertTo5bit(data [32]byte) []byte {
//	var result []byte
//	var value byte
//	var remainder byte
//
//	for _, b := range data {
//		value = (remainder << 5) | (b >> 3)
//		remainder = b & 0x07
//		result = append(result, value)
//
//		value = remainder << 2
//		result = append(result, value)
//	}
//
//	return result
//}
//
//func extractTapRootAddrFromPubKey(pkScript []byte) (string, error) {
//	// Check that the script is the correct length for P2TR (1-byte version + 1-byte length + 32-byte key)
//	if len(pkScript) != 34 {
//		return "", fmt.Errorf("invalid pkScript length for P2TR")
//	}
//
//	// Check that the version byte is 0x50 (indicating P2TR)
//	if pkScript[0] != 0x50 {
//		return "", fmt.Errorf("invalid version byte for P2TR")
//	}
//
//	// Check that the length byte is 0x20 (indicating 32 bytes)
//	if pkScript[1] != 0x20 {
//		return "", fmt.Errorf("invalid length byte for P2TR public key")
//	}
//
//	byte32TaprootPubKey, err := convertTo32ByteArray(pkScript[2:])
//	if err != nil {
//		return "", err
//	}
//
//	data := append([]byte{0x01}, convertTo5bit(byte32TaprootPubKey)...)
//
//	// Convert 32-byte public key to Bech32m
//	segwit, err := bech32.EncodeM("bc", data)
//	if err != nil {
//		return "", err
//	}
//
//	return segwit, nil
//}

// Convert the 32 bytes into 5-bit groups for Bech32m encoding.

func convertTo32ByteArray(data []byte) ([32]byte, error) {
	var arr [32]byte
	if len(data) != 32 {
		return arr, fmt.Errorf("expected a slice of length 32, got %d", len(data))
	}
	copy(arr[:], data)
	return arr, nil
}

//func convertTo5bit(data [32]byte) []byte {
//	var groups []byte
//	accumulator, bits := 0, 0
//
//	for _, b := range data {
//		accumulator = (accumulator << 8) | int(b)
//		bits += 8
//
//		for bits >= 5 {
//			bits -= 5
//			groups = append(groups, byte(accumulator>>bits)&0x1F)
//		}
//	}
//
//	// Handle any remaining bits
//	if bits > 0 {
//		groups = append(groups, byte(accumulator<<(5-bits)))
//	}
//
//	return groups
//}

func convertTo5bit(data [32]byte) []byte {
	var groups []byte
	accumulator, bits := 0, 0

	for _, b := range data {
		accumulator = (accumulator << 8) | int(b)
		bits += 8

		for bits >= 5 {
			bits -= 5
			value := byte((accumulator >> bits) & 0x1F)
			groups = append(groups, value)
		}
	}

	return groups
}

func extractTapRootAddrFromPubKey(pkScript []byte) (string, error) {
	// Your checks here ...

	byte32TaprootPubKey, err := convertTo32ByteArray(pkScript[2:])
	if err != nil {
		return "", err
	}

	// SegWit version 1 (Taproot) byte + converted pubkey
	data := append([]byte{0x01}, convertTo5bit(byte32TaprootPubKey)...)

	// Convert 32-byte public key to Bech32m
	segwit, err := bech32.EncodeM("bc", data) // Ensure this uses Bech32m encoding!
	if err != nil {
		return "", err
	}

	return segwit, nil
}

//func extractTapRootAddrFromPkScript(pkScript string) (string, error) {
//	// ... [rest of your functions] ...
//
//	data, err := hex.DecodeString(pkScript[4:]) // Skipping the version and length bytes
//	if err != nil {
//		return "", err
//	}
//	var byte32TaprootPubKey [32]byte
//	copy(byte32TaprootPubKey[:], data)
//
//	dataForBech32m := append([]byte{0x01}, convertTo5bit(byte32TaprootPubKey)...)
//
//	// Convert 32-byte public key to Bech32m
//	segwit, err := bech32.EncodeM("bc", dataForBech32m) // Ensure this uses Bech32m encoding!
//	if err != nil {
//		return "", err
//	}
//
//	return segwit, nil
//}

func extractTapRootAddrFromPkScript(pkScript []byte) (string, error) {
	// Check that the script is the correct length for P2TR (1-byte version + 1-byte length + 32-byte key)
	if len(pkScript) != 34 {
		return "", fmt.Errorf("invalid pkScript length for P2TR")
	}

	// Check that the version byte is 0x50 (indicating P2TR)
	if pkScript[0] != 0x51 {
		return "", fmt.Errorf("invalid version byte for P2TR")
	}

	// Check that the length byte is 0x20 (indicating 32 bytes)
	if pkScript[1] != 0x20 {
		return "", fmt.Errorf("invalid length byte for P2TR public key")
	}

	data := pkScript[2:]
	var byte32TaprootPubKey [32]byte
	copy(byte32TaprootPubKey[:], data)

	dataForBech32m := append([]byte{0x01}, convertTo5bit(byte32TaprootPubKey)...)

	// Convert 32-byte public key to Bech32m
	segwit, err := bech32.EncodeM("bc", dataForBech32m) // Ensure this uses Bech32m encoding!
	if err != nil {
		return "", err
	}

	return segwit, nil
}

func determinePublicKeyType(pkScript []byte) string {
	switch {
	case len(pkScript) == 25 && pkScript[0] == 0x76 && pkScript[1] == 0xa9:
		return "P2PKH"
	case len(pkScript) == 23 && pkScript[0] == 0xa9 && pkScript[1] == 0x14:
		return "P2SH"
	case len(pkScript) == 22 && pkScript[0] == 0x00 && pkScript[1] == 0x14:
		return "P2WPKH"
	case len(pkScript) == 34 && pkScript[0] == 0x00 && pkScript[1] == 0x20:
		return "P2WSH"
	case len(pkScript) == 34 && pkScript[0] == 0x50 && pkScript[1] == 0x20:
		return "P2TR"
	default:
		return "Unknown"
	}
}

//func PkScriptToTaprootAddress(pkScript []byte) (string, error) {
//	// Ensure the pkScript is of correct length and version
//	if len(pkScript) != 34 || pkScript[0] != 81 || pkScript[1] != 32 {
//		return "", fmt.Errorf("invalid pkScript for Taproot")
//	}
//
//	// Extract the 32-byte public key from the pkScript
//	pubKey := pkScript[2:]
//
//	convertedPubKey, err := bech32.ConvertBits(pubKey, 8, 5, false)
//	if err != nil {
//		return "", err
//	}
//
//	// Convert to bech32m
//	hrp := "bc"
//	// Note: You might need to modify the bech32 library to include Bech32m encoding if it's not yet supported.
//	address, err := bech32.EncodeM(hrp, append([]byte{0x01}, convertedPubKey...))
//	if err != nil {
//		return "", err
//	}
//
//	return address, nil
//}

//func PkScriptToTaprootAddress(pkScript []byte) (string, error) {
//	// Compute the SHA256 hash of the PkScript.
//	hash := sha256.Sum256(pkScript)
//
//	// Convert hash to bech32 encoded address with witness version 1.
//	hrp := "bc"
//	witnessVersion := byte(1)
//	conv, err := bech32.ConvertBits(hash[:], 8, 5, true)
//	if err != nil {
//		return "", err
//	}
//	data := append([]byte{witnessVersion}, conv...)
//	address, err := bech32.EncodeM(hrp, data)
//	if err != nil {
//		return "", err
//	}
//
//	return address, nil
//}

func PkScriptToTaprootAddress(pkScript []byte) (string, error) {
	if len(pkScript) != 32 {
		return "", fmt.Errorf("invalid public key length")
	}

	// Construct the Taproot internal key.
	internalKey := make([]byte, 33)
	copy(internalKey[1:], pkScript)

	// Hash the internal key.
	hash := sha256.Sum256(internalKey)

	// Convert hash to Bech32 encoded address with witness version 1.
	hrp := "bc"
	witnessVersion := byte(1)
	conv, err := bech32.ConvertBits(hash[:], 8, 5, true)
	if err != nil {
		return "", err
	}
	data := append([]byte{witnessVersion}, conv...)
	address, err := bech32.EncodeM(hrp, data)
	if err != nil {
		return "", err
	}

	return address, nil
}
