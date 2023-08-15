package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"ntc-services/models"
	"ntc-services/services"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/labstack/echo/v4"
	"github.com/tnakagawa/goref/bech32m"
)

/*
This experiment proves that we can unwrap an unsigned PBST to ensure the inputs are the same
UTXOs as selected by the maker/taker.

We are now also extracting the public key from the PkSript (P2TR) outputs of the unsignedTx of a
PSBT to hopefully cross-reference the trade to ensure nothing malicious has occurred. We are
successfully taking the PkScripts from the outputs and verifying the proper receiver addresses.
*/

func PSBTFromUnsignedTx(c echo.Context) error {
	//psbtHex := "70736274ff0100fd85010200000004cba2e1d1fb3d04af9ada6b5cf91f58fa1a977a644b6f398b68e251c99da03e500000000000ffffffff8a688114865023cb24f2a336f0480bf3c53d6819ecbeb3014e466cf18e162ad60000000000ffffffff90b2b9e55b3b4b8e6757d20109672fb0ef9369d8cc4de274bf88bd5dd4e84d310200000000ffffffff76ede0712b48a913f0d040c199c4e2bbc7201936eeba7ce3bd8add8df7945e260200000000ffffffff058f200000000000002251206ba9b63de4e17c5ec52351ac79d40de6a46b0ca28b3b818d76e8e20d51743dce220200000000000022512058b49d8f87981075be5f5459ecd7679f8d79caf31a87649c5532092fd5145310e8030000000000002251205cba2279dc9a998530a4fe43e78cdbf7793e5371c27877bbca3d5364f4bc1fc6270800000000000022512058b49d8f87981075be5f5459ecd7679f8d79caf31a87649c5532092fd5145310100f0000000000002251206ba9b63de4e17c5ec52351ac79d40de6a46b0ca28b3b818d76e8e20d51743dce000000000001012b8f2000000000000022512058b49d8f87981075be5f5459ecd7679f8d79caf31a87649c5532092fd51453100108420140b3666e07b2600608c04ba502951c2be75d3472bdd99c68f2cf779762629ee4d33a40797b50e8a86b3a20b3b50f83a638e623aae1faedfdc08aae46df9ff83dcd0001012b22020000000000002251206ba9b63de4e17c5ec52351ac79d40de6a46b0ca28b3b818d76e8e20d51743dce0001012bbc1600000000000022512058b49d8f87981075be5f5459ecd7679f8d79caf31a87649c5532092fd5145310010842014098d6bea9f200a047131180610035f618b1f693dce5281ed7ad31be3a5c9c680c0130d5e19902c0d1581244a917b79b65a1fdb240c4bae38e6930ae6e8f9ba2400001012ba51d0000000000002251206ba9b63de4e17c5ec52351ac79d40de6a46b0ca28b3b818d76e8e20d51743dce000000000000"

	// DEV: "_id" : ObjectId("64bf861ec8c118335ae29262") -- Confirmed Tx
	psbtHex := "70736274ff0100fd850102000000048bc746a710caaa3b1e294d3a4c6c9d1dd8debf4313d631f8f48cbe1f52c174d30100000000ffffffffa8a206c57a58d7a549f0e81213c6c4ae7e57b4e1231d25a7faed3948b3a9f1d00f00000000ffffffff7876ea3af1120f5d53e93eb8a033283dd701be3db0595e32908f3e91c933a5460400000000ffffffff8543eb84e01860321a15404d9d04a0d3ade15c5bf6c9fb09b0246ee2718f57f30500000000ffffffff052202000000000000225120f4f18b6b4e5de72fab1fc0c95d255c34b133cfae2d33a19fa8798a21e7667bfec303000000000000225120a058c4c973bfe2affd6b83b171ac8d2eda26c43e206053afe46dd4dde852d201e8030000000000002251205cba2279dc9a998530a4fe43e78cdbf7793e5371c27877bbca3d5364f4bc1fc63641000000000000225120a058c4c973bfe2affd6b83b171ac8d2eda26c43e206053afe46dd4dde852d201022f000000000000225120f4f18b6b4e5de72fab1fc0c95d255c34b133cfae2d33a19fa8798a21e7667bfe000000000001012b2202000000000000225120a058c4c973bfe2affd6b83b171ac8d2eda26c43e206053afe46dd4dde852d201010842014083afdb354e35a669f554b806c8ed2b56bcae994e3c292c687964164770632993a8c6a8b3f0eb5064ecb918d03c7757c899bb0363271c5859235fac5e6d40cc960001012bc303000000000000225120f4f18b6b4e5de72fab1fc0c95d255c34b133cfae2d33a19fa8798a21e7667bfe0001012bb34c000000000000225120a058c4c973bfe2affd6b83b171ac8d2eda26c43e206053afe46dd4dde852d20101084201402d91b32a861fa0fc9ed3c68d4a9fbaf6e01abd82cecab1b023fb11652adfa7e46e47b91b389db6335eb0699a1102c3a3151a0e209b8e4b2d275a82ea42f354e10001012b7f3a000000000000225120f4f18b6b4e5de72fab1fc0c95d255c34b133cfae2d33a19fa8798a21e7667bfe000000000000"

	// DEV: "_id" : ObjectId("64b60a8946adc479c1592e4b") -- Created Tx
	//psbtHex := "70736274ff0100fdb0010200000004d3e4a0f7a2bb5565eed1e4a6a60eb2ef7cf2235c6ab1e38d3d8ff40bf52fe9b70000000000ffffffff30c0421848cdf1fc1e0144149e84a13179f617e76dd3075b56abf19e62dbaa170300000000ffffffff5f84a67cc95f25dc20c2c429b6b1efda8abea423c8021105f12472b6cd42dfd90200000000ffffffff764d57569f688c044a7eb0d6f160e4e4f29a5728a803358e5c38550e4b6a98600000000000ffffffff062202000000000000225120edc3820dcdf065791c2690963fa16855dd18853808b58200ecd12cad3d426684a0860100000000002251203e1122de97f793bc9c250f61475bd0d76c37bd64e80e4182ae156b8b3704970ae803000000000000225120af97df56f8fdc7460ab77e68086e759f0bc220ff5a64be766e24f5cb149faa8f7d280c0000000000225120edc3820dcdf065791c2690963fa16855dd18853808b58200ecd12cad3d426684c263040000000000225120edc3820dcdf065791c2690963fa16855dd18853808b58200ecd12cad3d4266840f1b0000000000002251203e1122de97f793bc9c250f61475bd0d76c37bd64e80e4182ae156b8b3704970a000000000001012b22020000000000002251203e1122de97f793bc9c250f61475bd0d76c37bd64e80e4182ae156b8b3704970a0001012b1daf0d0000000000225120edc3820dcdf065791c2690963fa16855dd18853808b58200ecd12cad3d42668401084201408a6584b472c8f623da23dee7571d1767a01c5ef40958ecb9de31063742c857d8ff6f67b1da3b616eb1fab36a2117d31720c76bab70f65357e58a216733d274e70001012bc36f040000000000225120edc3820dcdf065791c2690963fa16855dd18853808b58200ecd12cad3d4266840108420140942e2c2efdfe68ae168bd5c03fd210164375ee1dd4910f45bb448288881b124d476d800138d418cbba2c19c258aebd18d4c699d82c2146b26ab3891260e08b690001012b10270000000000002251203e1122de97f793bc9c250f61475bd0d76c37bd64e80e4182ae156b8b3704970a00000000000000"

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

	pJSON, _ := json.MarshalIndent(p, "", "  ")

	fmt.Println(p.B64Encode())

	fmt.Println("++++++++++++++++++++++++")
	fmt.Println(string(pJSON))
	fmt.Printf("packet: %+v \n", p)
	fmt.Printf("witnessUTXO: %+v \n", *p.Inputs[0].WitnessUtxo)
	fmt.Println("++++++++++++++++++++++++")

	fmt.Println("||||||||||||||||||||||||")
	fmt.Printf("UnsignedTx: %+v \n", *p.UnsignedTx)
	for i, in := range p.UnsignedTx.TxIn {
		fmt.Printf("in %v: %+v \n", i, *in)
	}
	for i, out := range p.UnsignedTx.TxOut {
		fmt.Println(len(out.PkScript))
		pubKeyB, err := extractTaprootPublicKey(out.PkScript)
		if err != nil {
			c.Logger().Errorf("Failed to extract pubKey from PkScript: %v", err)
		}
		pubKey, err := extractTaprootPublicKeyHex(out.PkScript)
		if err != nil {
			c.Logger().Errorf("Failed to extract pubKey from PkScript: %v", err)
		}
		pkScriptType := determinePublicKeyType(out.PkScript)

		fmt.Printf("ScriptType %v: %+v \n", i, pkScriptType)
		fmt.Printf("out %v: %+v \n", i, out.PkScript)
		fmt.Printf("out %v: %+v \n", i, pubKey)
		fmt.Printf("out %v: %+v \n", i, pubKeyB)
		pkScriptHex := hex.EncodeToString(out.PkScript)
		fmt.Printf("out.PkScriptHex %v, %+v \n", i, pkScriptHex)
		tapRootPubKeyHex := pkScriptHex[2:]
		fmt.Printf("tapRootPubKeyHex %v, %+v \n", i, tapRootPubKeyHex)
	}
	fmt.Println("||||||||||||||||||||||||")

	fmt.Println("======================")
	// An attempt to get the public key from the P2TR PkScript
	for i, out := range p.UnsignedTx.TxOut {
		pkScriptHex := hex.EncodeToString(out.PkScript)

		fmt.Printf("out: %v: PkScript %+v \n", i, out.PkScript)
		fmt.Printf("out: %v: PkScriptHex %+v \n", i, pkScriptHex)

		address, err := bech32m.SegwitAddrEncode("bc", 0x01, out.PkScript[2:])
		if err != nil {
			c.Logger().Errorf("Failed convert PkScript to Segwit: %v", err)
		}
		fmt.Printf("out: %v: Address %+v \n", i, address)

		//fmt.Println(pkScriptHex)
		//fmt.Println(len(pkScriptHex))
		//
		//// Check if the PkScript is valid and extract the public key
		//if strings.HasPrefix(pkScriptHex, "51") && len(pkScriptHex[2:]) == 66 {
		//	pubKey := out.PkScript[2:]
		//	// Convert the public key into a bech32m encoded Taproot address
		//	address := bech32m.Encode("bc", pubKey, 2)
		//	fmt.Println("Taproot Address:", address)
		//} else {
		//	fmt.Println("Invalid PkScript")
		//}

		//tapRootPubKey, err := extractTaprootPublicKey(out.PkScript)
		//if err != nil {
		//	c.Logger().Errorf("Failed to extract taproot addr: %v", err)
		//}
		//
		//bech32mPubKey := convertTo5bit(tapRootPubKey)
		//segwit, err := bech32.Encode("bc", bech32mPubKey)
		//if err != nil {
		//	c.Logger().Errorf("Failed to bech32m endcode taproot pubkey: %v", err)
		//}

		//fmt.Println(len(out.PkScript))
		//segwit, err := PkScriptToTaprootAddress(out.PkScript[2:])
		//if err != nil {
		//	c.Logger().Errorf("Failed to bech32m endcode taproot pubkey: %v", err)
		//}
		//
		//fmt.Printf("out %v: %+v \n", i, segwit)
	}

	fmt.Println("======================")

	//fmt.Println("========================")
	//for i, in := range p.Inputs {
	//	fmt.Printf("in %v: %+v \n", i, *in.WitnessUtxo)
	//}
	//fmt.Println("========================")
	//
	//fmt.Println("========================")
	//for i, out := range p.Outputs {
	//	fmt.Printf("out %v: %+v \n", i, out)
	//}
	//fmt.Println("========================")

	return c.JSON(http.StatusOK, nil)
}

func extractTaprootPublicKey(pkScript []byte) ([32]byte, error) {
	// Check that the script is the correct length for P2TR (1-byte version + 1-byte length + 32-byte key)
	if len(pkScript) != 34 {
		return [32]byte{}, fmt.Errorf("invalid pkScript length for P2TR")
	}

	// Check that the version byte is 0x51 (indicating P2TR)
	if pkScript[0] != 0x51 {
		return [32]byte{}, fmt.Errorf("invalid version byte for P2TR")
	}

	// Check that the length byte is 0x20 (indicating 32 bytes)
	if pkScript[1] != 0x20 {
		return [32]byte{}, fmt.Errorf("invalid length byte for P2TR public key")
	}

	byte32TaprootPubKey, err := convertTo32ByteArray(pkScript[2:])
	if err != nil {
		return [32]byte{}, err
	}

	// The remaining 32 bytes are the taproot public key
	return byte32TaprootPubKey, nil
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

//func convertTo5bit(data [32]byte) []byte {
//	var bit5Data []byte
//	var currentByte byte
//	var nextByte byte
//	var bit5Value byte
//
//	for i := 0; i < 32; i++ {
//		currentByte = data[i]
//		if i < 31 {
//			nextByte = data[i+1]
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
//	return bit5Data[:64] // the result should have 64 5-bit values
//}

// bc17n7x3dddfethuu974vrlcre9t5y4tsxjkyx0e74c95xw5xn74p0x3gy8uave00lm0pmw05
// bc17n7x3dddfethuu974vrlcre9t5y4tsxjkyx0e74c95xw5xn74p0x3gy8uave00lm6atz2k
// bc1p5pvvfjtnhl32llttswchrtyd9mdzd3p7yps98tlydh2dm6zj6gqsfkmcnd

// bc1p7ncck66wthnjl2clcry46f2uxjcn8naw95e6r8ag0x9zremx00lyklxy8

// bc1p5pvvfjtnhl32llttswchrtyd9mdzd3p7yps98tlydh2dm6zj6gqsfkmcnd
// bc1p7ncck66wthnjl2clcry46f2uxjcn8naw95e6r8ag0x9zremx00lqvf5wve

// bc1ptjazy7wun2vc2v9ylep70rxm7aunu5m3cfu80w7284fkfa9urlrqg07a7m
// bc1ptjazy7wun2vc2v9ylep70rxm7aunu5m3cfu80w7284fkfa9urlr67w3g8
// bc1ptjazy7wun2vc2v9ylep70rxm7aunu5m3cfu80w7284fkfa9urlr0z7ad9
// bc1ptjazy7wun2vc2v9ylep70rxm7aunu5m3cfu80w7284fkfa9urlr67w3g8

// bc1ptjazy7wun2vc2v9ylep70rxm7aunu5m3cfu80w7284fkfa9urlrqg07a7m
// bc1pkd2uverm5mhsajccxezu26t303g4wgjgqsls69qp2pxddpfqug9qv9x66m
// bc1pwjj8u3np4pw0edaq286ahdwhx9kerueah22apuj7vdfjc8crjq3qn9fy2x

// bc1ptjazy7wun2vc2v9ylep70rxm7aunu5m3cfu80w7284fkfa9urlrqg07a7m
// this is the house fee addr, we need to decode out 2 to this somehow

func GeneratePSBT(c echo.Context) error {
	btcClient := services.NewBitcoinClient()

	senderTapRootAddr := "bc1p0t40pgryukh88rhwx4ffzt28cjmhxnpm56s3382vyy22zek5wpmq8rps2l"
	resp, err := services.BLOCKCHAININFO.GetUTXOsForAddr(
		senderTapRootAddr,
	)
	if err != nil {
		c.Logger().Errorf("Failed to get UTXOs from platform wallet: %v", err)
	}

	//senderSegWit := "bc1p0t40pgryukh88rhwx4ffzt28cjmhxnpm56s3382vyy22zek5wpmq8rps2l"
	//resp, err := services.BLOCKCHAININFO.GetUTXOsForAddr(
	//	senderSegWit,
	//)
	//if err != nil {
	//	c.Logger().Errorf("Failed to get UTXOs from platform wallet: %v", err)
	//}

	var bciUTXOS []*models.BciUTXO
	for k, v := range resp {
		if k == "unspent_outputs" { // TODO: move to its own func in models named ParseJSONToBciUTXOs
			for _, utxo := range v.([]interface{}) {
				bciUTXO, err := models.ParseJSONToBciUTXO(utxo.(map[string]interface{}))
				if err != nil {
					c.Logger().Errorf("Could not parse resp utxos: %v", err)
				}
				bciUTXOS = append(bciUTXOS, bciUTXO)
			}
		}
	}

	fmt.Println("----------------------")
	//fmt.Printf("psbt complete: %+v \n", pC.IsComplete())
	for i, utxo := range bciUTXOS {
		fmt.Printf("index: %v utxo: %+v \n", i, utxo)
	}
	fmt.Println("----------------------")

	receiverTapRootAddr := "bc1pxy8gsmgu5zzv0ytj7ae4pgnqkcdwaqas7xmc4szcg70zqwsj4rxss2tvhu"
	netParams := &chaincfg.MainNetParams
	//senderAddr, err := btcutil.DecodeAddress(senderTapRootAddr, netParams)
	//if err != nil {
	//	c.Logger().Errorf("Could not decode taproot address: %v", err)
	//}

	tx := wire.NewMsgTx(2)
	txHash, err := chainhash.NewHashFromStr(bciUTXOS[10].TxHashBigEndian)
	outpoint := wire.NewOutPoint(txHash, bciUTXOS[10].TxOutputN)
	tx.AddTxIn(wire.NewTxIn(outpoint, nil, nil))

	receiverAddr, err := btcutil.DecodeAddress(receiverTapRootAddr, netParams)
	if err != nil {
		c.Logger().Errorf("Could not decode taproot address: %v", err)
	}
	pkScript, err := txscript.PayToAddrScript(receiverAddr)
	outputValue := int64(546)
	tx.AddTxOut(wire.NewTxOut(outputValue, pkScript))

	packet, err := psbt.NewFromUnsignedTx(tx)
	if err != nil {
		c.Logger().Errorf("Could not create pbst: %v", err)
	}

	base64PSBT, err := packet.B64Encode()
	packet.UnsignedTx.TxHash()
	if err != nil {
		c.Logger().Errorf("Could not parse psbt to base64: %v", err)
	}

	// TODO: don't commit

	// Load private key
	//wif, err := btcutil.DecodeWIF("")
	//if err != nil {
	//	c.Logger().Errorf("Could not decode private key: %v", err)
	//}
	//privKey := wif.PrivKey
	//pubKey := privKey.PubKey()
	//if err != nil {
	//	c.Logger().Errorf("Could not parse public key: %v", err)
	//}

	privateKeyHex := ""
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		c.Logger().Errorf("Failed to decode hex: %v", err)
		return err
	}
	privKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	prevOutScript, err := hex.DecodeString(bciUTXOS[10].Script)
	if err != nil {
		c.Logger().Errorf("Could not decode prevOutScriptHex: %v", err)
	}
	//prevOutScript, err := txscript.PayToAddrScript(senderAddr)
	//if err != nil {
	//	c.Logger().Errorf("Could not decode prevOutScriptHex: %v", err)
	//}
	packet.Inputs[0].WitnessUtxo = &wire.TxOut{
		Value:    bciUTXOS[10].Value,
		PkScript: prevOutScript,
	}

	//witnessScript := append([]byte{0, 20}, senderAddr.ScriptAddress()...)
	//if err != nil {
	//	c.Logger().Errorf("Could not encode pkScript: %v", err)
	//}
	//packet.Inputs[0].WitnessScript = witnessScript

	//prevTx, err := btcClient.GetRawTransaction(&packet.UnsignedTx.TxIn[0].PreviousOutPoint.Hash)
	//prevOut := prevTx.MsgTx().TxOut[10]
	//prevOuts := []*wire.TxOut{prevOut}
	//fetcher := txscript.NewMultiPrevOutFetcher(prevOuts)

	fetcher := txscript.NewCannedPrevOutputFetcher(prevOutScript, bciUTXOS[10].Value)
	sigHashes := txscript.NewTxSigHashes(packet.UnsignedTx, fetcher)
	u, err := psbt.NewUpdater(packet)
	if err != nil {
		panic(err)
	}
	if err := u.AddInSighashType(txscript.SigHashAnyOneCanPay, 0); err != nil {
		panic(err)
	}

	//sig, err := txscript.RawTxInWitnessSignature(packet.UnsignedTx, sigHashes, 0,
	//	packet.Inputs[0].WitnessUtxo.Value, packet.Inputs[0].WitnessUtxo.PkScript,
	//	txscript.SigHashAnyOneCanPay, privKey,
	//)

	sig, err := txscript.RawTxInTaprootSignature(
		packet.UnsignedTx,
		sigHashes,
		0,
		packet.Inputs[0].WitnessUtxo.Value,
		packet.Inputs[0].WitnessUtxo.PkScript,
		nil,
		txscript.SigHashAnyOneCanPay,
		privKey,
	)

	//schnorr.Sign(privKey, )
	//schnorr.

	//signMethod, err := validateSigningMethod(&packet.Inputs[0])
	//if err != nil {
	//	panic(err)
	//}

	packetJSON, _ := json.MarshalIndent(packet, "", "  ")
	fmt.Println("++++++++++++++++++++++++++++")
	fmt.Printf("PSBT Base64: %s\n", base64PSBT)
	fmt.Printf("PSBT Hash: %s\n", packet.UnsignedTx.TxHash())
	fmt.Println(string(packetJSON))
	fmt.Printf("UnsignedTX: %+v \n", packet.UnsignedTx)
	fmt.Printf("Packet: %+v \n", packet)
	fmt.Printf("WitnessUTXO: %+v \n", packet.Inputs[0].WitnessUtxo)
	fmt.Printf("Sanity: %v \n", packet.SanityCheck())
	fmt.Printf("Ready to Sign: %+v \n", psbt.InputsReadyToSign(packet))
	//fmt.Printf("Signing Method: %+v \n", signMethod)
	fmt.Println("++++++++++++++++++++++++++++")

	//witnessScript := append([]byte{0, 20}, receiverAddr.ScriptAddress()...)
	//outputs := []*wire.TxOut{wire.NewTxOut(99500, witnessScript)}

	success, err := u.Sign(0, sig, privKey.PubKey().SerializeCompressed(), nil, nil)
	if err != nil {
		panic(err)
	}
	if success != psbt.SignSuccesful {
		panic("could not successfully sign for some reason")
	}

	packetSignedJSON, _ := json.MarshalIndent(packet, "", "  ")
	fmt.Println("++++++++++++++++++++++++++++")
	fmt.Printf("Packet: %+v\n", packet)
	fmt.Printf("PSBT Base64: %s\n", base64PSBT)
	fmt.Printf("PSBT Hash: %s\n", packet.UnsignedTx.TxHash())
	fmt.Println(string(packetSignedJSON))
	fmt.Printf("IsTaproot: %v \n", txscript.IsPayToTaproot(packet.Inputs[0].WitnessUtxo.PkScript))
	fmt.Printf("IsComplete: %v \n", packet.IsComplete())
	fmt.Println("++++++++++++++++++++++++++++")

	if err := psbt.Finalize(packet, 0); err != nil {
		c.Logger().Errorf("Could not finalize packet: %v", err)
	}

	signedTx, err := psbt.Extract(packet)
	if err != nil {
		c.Logger().Errorf("Failed to extract signedTx from packet: %v", err)
	}

	fmt.Println("||||||||||||||||||||||||||||")
	fmt.Printf("SignedTx: %+v \n", signedTx)
	fmt.Println("||||||||||||||||||||||||||||")

	hash, err := btcClient.SendRawTransaction(signedTx, true)
	if err != nil {
		c.Logger().Errorf("Failed to broadcast transaction: %v", err)
	}
	fmt.Printf("Transaction broadcasted with hash: %s", hash)

	return c.JSON(http.StatusOK, packet)
}

func UTXOs(c echo.Context) error {
	//senderTapRootAddr := "bc1p0t40pgryukh88rhwx4ffzt28cjmhxnpm56s3382vyy22zek5wpmq8rps2l"
	//resp, err := services.BLOCKCHAININFO.GetUTXOsForAddr(
	//	senderTapRootAddr,
	//)
	//if err != nil {
	//	c.Logger().Errorf("Failed to get UTXOs from platform wallet: %v", err)
	//}

	senderSegWit := "36zbLaYJ1esG5tBGmwEAQCkb4yuffeqHJy"
	resp, err := services.BLOCKCHAININFO.GetUTXOsForAddr(
		senderSegWit,
	)
	if err != nil {
		c.Logger().Errorf("Failed to get UTXOs from platform wallet: %v", err)
	}

	respJSON, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		c.Logger().Errorf("Could not marshal resp utxo: %v", err)
	}

	fmt.Println("--------------------")
	fmt.Printf("%+v \n", string(respJSON))
	fmt.Println("--------------------")

	bisInscriptions, err := services.BESTINSLOT.GetInscriptionsByWalletAddr(c, senderSegWit, 20, 0)
	if err != nil {
		c.Logger().Error(err)
	}

	bisInscriptionsJSON, err := json.MarshalIndent(bisInscriptions, "", "  ")
	if err != nil {
		c.Logger().Error(err)
	}

	fmt.Println("=====================")
	fmt.Printf("%+v \n", string(bisInscriptionsJSON))
	fmt.Println("=====================")

	return c.JSON(http.StatusOK, nil)
}

// SignMethod defines the different ways a signer can sign, given a specific
// input.
type SignMethod uint8

const (
	// WitnessV0SignMethod denotes that a SegWit v0 (p2wkh, np2wkh, p2wsh)
	// input script should be signed.
	WitnessV0SignMethod SignMethod = 0

	// TaprootKeySpendBIP0086SignMethod denotes that a SegWit v1 (p2tr)
	// input should be signed by using the BIP0086 method (commit to
	// internal key only).
	TaprootKeySpendBIP0086SignMethod SignMethod = 1

	// TaprootKeySpendSignMethod denotes that a SegWit v1 (p2tr)
	// input should be signed by using a given taproot hash to commit to in
	// addition to the internal key.
	TaprootKeySpendSignMethod SignMethod = 2

	// TaprootScriptSpendSignMethod denotes that a SegWit v1 (p2tr) input
	// should be spent using the script path and that a specific leaf script
	// should be signed for.
	TaprootScriptSpendSignMethod SignMethod = 3
)

// validateSigningMethod attempts to detect the signing method that is required
// to sign for the given PSBT input and makes sure all information is available
// to do so.
func validateSigningMethod(in *psbt.PInput) (SignMethod, error) {
	script, err := txscript.ParsePkScript(in.WitnessUtxo.PkScript)
	if err != nil {
		return 0, fmt.Errorf("error detecting signing method, "+
			"couldn't parse pkScript: %v", err)
	}

	switch script.Class() {
	case txscript.WitnessV0PubKeyHashTy, txscript.ScriptHashTy,
		txscript.WitnessV0ScriptHashTy:

		return WitnessV0SignMethod, nil

	case txscript.WitnessV1TaprootTy:
		if len(in.TaprootBip32Derivation) == 0 {
			return 0, fmt.Errorf("cannot sign for taproot input " +
				"without taproot BIP0032 derivation info")
		}

		// Currently, we only support creating one signature per input.
		//
		// TODO(guggero): Should we support signing multiple paths at
		// the same time? What are the performance and security
		// implications?
		if len(in.TaprootBip32Derivation) > 1 {
			return 0, fmt.Errorf("unsupported multiple taproot " +
				"BIP0032 derivation info found, can only " +
				"sign for one at a time")
		}

		derivation := in.TaprootBip32Derivation[0]
		switch {
		// No leaf hashes means this is the internal key we're signing
		// with, so it's a key spend. And no merkle root means this is
		// a BIP0086 output we're signing for.
		case len(derivation.LeafHashes) == 0 &&
			len(in.TaprootMerkleRoot) == 0:

			return TaprootKeySpendBIP0086SignMethod, nil

		// A non-empty merkle root means we committed to a taproot hash
		// that we need to use in the tap tweak.
		case len(derivation.LeafHashes) == 0:
			// Getting here means the merkle root isn't empty, but
			// is it exactly the length we need?
			if len(in.TaprootMerkleRoot) != sha256.Size {
				return 0, fmt.Errorf("invalid taproot merkle "+
					"root length, got %d expected %d",
					len(in.TaprootMerkleRoot), sha256.Size)
			}

			return TaprootKeySpendSignMethod, nil

		// Currently, we only support signing for one leaf at a time.
		//
		// TODO(guggero): Should we support signing multiple paths at
		// the same time? What are the performance and security
		// implications?
		case len(derivation.LeafHashes) == 1:
			// If we're supposed to be signing for a leaf hash, we
			// also expect the leaf script that hashes to that hash
			// in the appropriate field.
			if len(in.TaprootLeafScript) != 1 {
				return 0, fmt.Errorf("specified leaf hash in " +
					"taproot BIP0032 derivation but " +
					"missing taproot leaf script")
			}

			leafScript := in.TaprootLeafScript[0]
			leaf := txscript.TapLeaf{
				LeafVersion: leafScript.LeafVersion,
				Script:      leafScript.Script,
			}
			leafHash := leaf.TapHash()
			if !bytes.Equal(leafHash[:], derivation.LeafHashes[0]) {
				return 0, fmt.Errorf("specified leaf hash in" +
					"taproot BIP0032 derivation but " +
					"corresponding taproot leaf script " +
					"was not found")
			}

			return TaprootScriptSpendSignMethod, nil

		default:
			return 0, fmt.Errorf("unsupported number of leaf " +
				"hashes in taproot BIP0032 derivation info, " +
				"can only sign for one at a time")
		}

	default:
		return 0, fmt.Errorf("unsupported script class for signing "+
			"PSBT: %v", script.Class())
	}
}
