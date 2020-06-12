package types

import (
	"github.com/Futuremine-chain/futuremine/futuremine/common/kit"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/crypto/ecc/secp256k1"
)

// Signature information, including the result of the
// signature and the public key.
type Signature struct {
	Bytes  []byte `json:"bytes"`
	PubKey []byte `json:"pubkey"`
}

// Sign the hash with the private key
func Sign(key *secp256k1.PrivateKey, hash arry.Hash) (*Signature, error) {
	signature, err := key.Sign(hash.Bytes())
	if err != nil {
		return nil, err
	}
	return &Signature{signature.Serialize(), key.PubKey().SerializeCompressed()}, nil
}

// Verify signature by hash and signature result
func Verify(hash arry.Hash, signScript *Signature) bool {
	if signScript == nil || signScript.PubKey == nil || signScript.Bytes == nil {
		return false
	}
	pubkey, err := secp256k1.ParsePubKey(signScript.PubKey)
	if err != nil {
		return false
	}
	signature, err := secp256k1.ParseSignature(signScript.Bytes, secp256k1.S256())
	return signature.Verify(hash.Bytes(), pubkey)
}

// Verify whether the signers are consistent through the public key
func VerifySigner(network string, signer arry.Address, pubKey []byte) bool {
	key, err := secp256k1.ParsePubKey(pubKey)
	if err != nil {
		return false
	}
	generateAddress, err := kit.GenerateAddress(network, key)
	if err != nil {
		return false
	}
	return generateAddress.IsEqual(signer)
}