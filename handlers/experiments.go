package handlers

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/labstack/echo/v4"
	"net/http"
)

/*
This experiment proves that we can unwrap an unsigned PBST to ensure the inputs
are the same UTXOs as selected by the maker/taker.

We are now also extracting the public key from the PkSript (P2TR) outputs of the unsignedTx of a
psbt to hopefully cross-reference the trade to ensure nothing malicious has occurred
*/

func PSBTFromUnsignedTx(c echo.Context) error {
	//psbtHex := "70736274ff0100fd85010200000004cba2e1d1fb3d04af9ada6b5cf91f58fa1a977a644b6f398b68e251c99da03e500000000000ffffffff8a688114865023cb24f2a336f0480bf3c53d6819ecbeb3014e466cf18e162ad60000000000ffffffff90b2b9e55b3b4b8e6757d20109672fb0ef9369d8cc4de274bf88bd5dd4e84d310200000000ffffffff76ede0712b48a913f0d040c199c4e2bbc7201936eeba7ce3bd8add8df7945e260200000000ffffffff058f200000000000002251206ba9b63de4e17c5ec52351ac79d40de6a46b0ca28b3b818d76e8e20d51743dce220200000000000022512058b49d8f87981075be5f5459ecd7679f8d79caf31a87649c5532092fd5145310e8030000000000002251205cba2279dc9a998530a4fe43e78cdbf7793e5371c27877bbca3d5364f4bc1fc6270800000000000022512058b49d8f87981075be5f5459ecd7679f8d79caf31a87649c5532092fd5145310100f0000000000002251206ba9b63de4e17c5ec52351ac79d40de6a46b0ca28b3b818d76e8e20d51743dce000000000001012b8f2000000000000022512058b49d8f87981075be5f5459ecd7679f8d79caf31a87649c5532092fd51453100108420140b3666e07b2600608c04ba502951c2be75d3472bdd99c68f2cf779762629ee4d33a40797b50e8a86b3a20b3b50f83a638e623aae1faedfdc08aae46df9ff83dcd0001012b22020000000000002251206ba9b63de4e17c5ec52351ac79d40de6a46b0ca28b3b818d76e8e20d51743dce0001012bbc1600000000000022512058b49d8f87981075be5f5459ecd7679f8d79caf31a87649c5532092fd5145310010842014098d6bea9f200a047131180610035f618b1f693dce5281ed7ad31be3a5c9c680c0130d5e19902c0d1581244a917b79b65a1fdb240c4bae38e6930ae6e8f9ba2400001012ba51d0000000000002251206ba9b63de4e17c5ec52351ac79d40de6a46b0ca28b3b818d76e8e20d51743dce000000000000"

	// DEV: "_id" : ObjectId("64bf861ec8c118335ae29262") -- Confirmed Tx
	psbtHex := "70736274ff0100fd850102000000048bc746a710caaa3b1e294d3a4c6c9d1dd8debf4313d631f8f48cbe1f52c174d30100000000ffffffffa8a206c57a58d7a549f0e81213c6c4ae7e57b4e1231d25a7faed3948b3a9f1d00f00000000ffffffff7876ea3af1120f5d53e93eb8a033283dd701be3db0595e32908f3e91c933a5460400000000ffffffff8543eb84e01860321a15404d9d04a0d3ade15c5bf6c9fb09b0246ee2718f57f30500000000ffffffff052202000000000000225120f4f18b6b4e5de72fab1fc0c95d255c34b133cfae2d33a19fa8798a21e7667bfec303000000000000225120a058c4c973bfe2affd6b83b171ac8d2eda26c43e206053afe46dd4dde852d201e8030000000000002251205cba2279dc9a998530a4fe43e78cdbf7793e5371c27877bbca3d5364f4bc1fc63641000000000000225120a058c4c973bfe2affd6b83b171ac8d2eda26c43e206053afe46dd4dde852d201022f000000000000225120f4f18b6b4e5de72fab1fc0c95d255c34b133cfae2d33a19fa8798a21e7667bfe000000000001012b2202000000000000225120a058c4c973bfe2affd6b83b171ac8d2eda26c43e206053afe46dd4dde852d201010842014083afdb354e35a669f554b806c8ed2b56bcae994e3c292c687964164770632993a8c6a8b3f0eb5064ecb918d03c7757c899bb0363271c5859235fac5e6d40cc960001012bc303000000000000225120f4f18b6b4e5de72fab1fc0c95d255c34b133cfae2d33a19fa8798a21e7667bfe0001012bb34c000000000000225120a058c4c973bfe2affd6b83b171ac8d2eda26c43e206053afe46dd4dde852d20101084201402d91b32a861fa0fc9ed3c68d4a9fbaf6e01abd82cecab1b023fb11652adfa7e46e47b91b389db6335eb0699a1102c3a3151a0e209b8e4b2d275a82ea42f354e10001012b7f3a000000000000225120f4f18b6b4e5de72fab1fc0c95d255c34b133cfae2d33a19fa8798a21e7667bfe000000000000"

	//psbtHex := ""

	psbtBytes, err := hex.DecodeString(psbtHex)
	if err != nil {
		c.Logger().Errorf("Failed to decode hex: %v", err)
	}

	psbtIOReader := bytes.NewReader(psbtBytes)
	p, err := psbt.NewFromRawBytes(psbtIOReader, false)
	if err != nil {
		c.Logger().Errorf("Failed to decode PSBT: %v", err)
	}

	fmt.Println("++++++++++++++++++++++++")
	fmt.Printf("packet: %+v \n", p)
	fmt.Println("++++++++++++++++++++++++")

	fmt.Println("||||||||||||||||||||||||")
	fmt.Printf("UnsignedTx: %+v \n", *p.UnsignedTx)
	for i, in := range p.UnsignedTx.TxIn {
		fmt.Printf("in %v: %+v \n", i, *in)
	}
	for i, out := range p.UnsignedTx.TxOut {
		fmt.Println(len(out.PkScript))
		pubKey, err := extractTaprootPublicKeyHex(out.PkScript)
		if err != nil {
			c.Logger().Errorf("Failed to extract pubKey from PkScript: %v", err)
		}
		fmt.Printf("out %v: %+v \n", i, pubKey)
	}
	fmt.Println("||||||||||||||||||||||||")

	fmt.Println("========================")
	for i, in := range p.Inputs {
		fmt.Printf("in %v: %+v \n", i, *in.WitnessUtxo)
	}
	fmt.Println("========================")

	fmt.Println("========================")
	for i, out := range p.Outputs {
		fmt.Printf("out %v: %+v \n", i, out)
	}
	fmt.Println("========================")

	return c.JSON(http.StatusOK, nil)
}

//func extractTaprootPublicKey(pkScript []byte) ([]byte, error) {
//	// Check that the script is the correct length for P2TR (1-byte version + 32-byte key)
//	if len(pkScript) != 33 {
//		return nil, fmt.Errorf("invalid pkScript length for P2TR")
//	}
//
//	// Check that the version byte is 0x50 (indicating P2TR)
//	if pkScript[0] != 0x50 {
//		return nil, fmt.Errorf("invalid version byte for P2TR")
//	}
//
//	// The remaining 32 bytes are the taproot public key
//	return pkScript[1:], nil
//}

func extractTaprootPublicKey(pkScript []byte) ([]byte, error) {
	// Check that the script is the correct length for P2TR (1-byte version + 1-byte length + 32-byte key)
	if len(pkScript) != 34 {
		return nil, fmt.Errorf("invalid pkScript length for P2TR")
	}

	// Check that the version byte is 0x51 (indicating P2TR)
	if pkScript[0] != 0x51 {
		return nil, fmt.Errorf("invalid version byte for P2TR")
	}

	// Check that the length byte is 0x20 (indicating 32 bytes)
	if pkScript[1] != 0x20 {
		return nil, fmt.Errorf("invalid length byte for P2TR public key")
	}

	// The remaining 32 bytes are the taproot public key
	return pkScript[2:], nil
}

func extractTaprootPublicKeyHex(pkScript []byte) (string, error) {
	// Check that the script is the correct length for P2TR (1-byte version + 1-byte length + 32-byte key)
	if len(pkScript) != 34 {
		return "", fmt.Errorf("invalid pkScript length for P2TR")
	}

	// Check that the version byte is 0x51 (indicating P2TR)
	if pkScript[0] != 0x51 {
		return "", fmt.Errorf("invalid version byte for P2TR")
	}

	// Check that the length byte is 0x20 (indicating 32 bytes)
	if pkScript[1] != 0x20 {
		return "", fmt.Errorf("invalid length byte for P2TR public key")
	}

	// Convert the 32-byte taproot public key to hexadecimal format
	pubKeyHex := hex.EncodeToString(pkScript[2:])
	return pubKeyHex, nil
}
