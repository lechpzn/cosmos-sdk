package legacytx

import (
	"encoding/json"
	"fmt"

	"sigs.k8s.io/yaml"

	"cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
)

// LegacyMsg defines the old interface a message must fulfill,
// containing Amino signing method.
// Deprecated: Please use `Msg` instead.
type LegacyMsg interface {
	sdk.Msg

	// Get the canonical byte representation of the Msg.
	GetSignBytes() []byte
}

// StdSignDoc is replay-prevention structure.
// It includes the result of msg.GetSignBytes(),
// as well as the ChainID (prevent cross chain replay)
// and the Sequence numbers for each signature (prevent
// inchain replay and enforce tx ordering per account).
type StdSignDoc struct {
	AccountNumber uint64            `json:"account_number" yaml:"account_number"`
	Sequence      uint64            `json:"sequence" yaml:"sequence"`
	TimeoutHeight uint64            `json:"timeout_height,omitempty" yaml:"timeout_height"`
	ChainID       string            `json:"chain_id" yaml:"chain_id"`
	Memo          string            `json:"memo" yaml:"memo"`
	Fee           json.RawMessage   `json:"fee" yaml:"fee"`
	Msgs          []json.RawMessage `json:"msgs" yaml:"msgs"`
}

var RegressionTestingAminoCodec *codec.LegacyAmino

// Deprecated: please delete this code eventually.
func mustSortJSON(bz []byte) []byte {
	var c any
	err := json.Unmarshal(bz, &c)
	if err != nil {
		panic(err)
	}
	js, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return js
}

// StdSignBytes returns the bytes to sign for a transaction.
// Deprecated: Please use x/tx/signing/aminojson instead.
func StdSignBytes(chainID string, accnum, sequence, timeout uint64, fee StdFee, msgs []sdk.Msg, memo string) []byte {
	if RegressionTestingAminoCodec == nil {
		panic(fmt.Errorf("must set RegressionTestingAminoCodec before calling StdSignBytes"))
	}
	msgsBytes := make([]json.RawMessage, 0, len(msgs))
	for _, msg := range msgs {
		bz := RegressionTestingAminoCodec.MustMarshalJSON(msg)
		msgsBytes = append(msgsBytes, mustSortJSON(bz))
	}

	bz, err := legacy.Cdc.MarshalJSON(StdSignDoc{
		AccountNumber: accnum,
		ChainID:       chainID,
		Fee:           json.RawMessage(fee.Bytes()),
		Memo:          memo,
		Msgs:          msgsBytes,
		Sequence:      sequence,
		TimeoutHeight: timeout,
	})
	if err != nil {
		panic(err)
	}

	return mustSortJSON(bz)
}

// StdSignature represents a sig
//
// Deprecated: will be removed in the future.
type StdSignature struct {
	cryptotypes.PubKey `json:"pub_key" yaml:"pub_key"` // optional
	Signature          []byte                          `json:"signature" yaml:"signature"`
}

// NewStdSignature is a legacy function
//
// Deprecated: will be removed in the future.
func NewStdSignature(pk cryptotypes.PubKey, sig []byte) StdSignature {
	return StdSignature{PubKey: pk, Signature: sig}
}

// GetSignature returns the raw signature bytes.
func (ss StdSignature) GetSignature() []byte {
	return ss.Signature
}

// GetPubKey returns the public key of a signature as a cryptotypes.PubKey using the
// Amino codec.
func (ss StdSignature) GetPubKey() cryptotypes.PubKey {
	return ss.PubKey
}

// MarshalYAML returns the YAML representation of the signature.
func (ss StdSignature) MarshalYAML() (any, error) {
	pk := ""
	if ss.PubKey != nil {
		pk = ss.String()
	}

	bz, err := yaml.Marshal(struct {
		PubKey    string `json:"pub_key"`
		Signature string `json:"signature"`
	}{
		pk,
		fmt.Sprintf("%X", ss.Signature),
	})
	if err != nil {
		return nil, err
	}

	return string(bz), nil
}

func (ss StdSignature) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return codectypes.UnpackInterfaces(ss.PubKey, unpacker)
}

// StdSignatureToSignatureV2 converts a StdSignature to a SignatureV2
func StdSignatureToSignatureV2(cdc *codec.LegacyAmino, sig StdSignature) (signing.SignatureV2, error) {
	pk := sig.GetPubKey()
	data, err := pubKeySigToSigData(cdc, pk, sig.Signature)
	if err != nil {
		return signing.SignatureV2{}, err
	}

	return signing.SignatureV2{
		PubKey: pk,
		Data:   data,
	}, nil
}

func pubKeySigToSigData(cdc *codec.LegacyAmino, key cryptotypes.PubKey, sig []byte) (signing.SignatureData, error) {
	multiPK, ok := key.(multisig.PubKey)
	if !ok {
		return &signing.SingleSignatureData{
			SignMode:  signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON,
			Signature: sig,
		}, nil
	}
	var multiSig multisig.AminoMultisignature
	err := cdc.Unmarshal(sig, &multiSig)
	if err != nil {
		return nil, err
	}

	sigs := multiSig.Sigs
	sigDatas := make([]signing.SignatureData, len(sigs))
	pubKeys := multiPK.GetPubKeys()
	bitArray := multiSig.BitArray
	n := multiSig.BitArray.Count()
	signatures := multisig.NewMultisig(n)
	sigIdx := 0
	for i := range n {
		if bitArray.GetIndex(i) {
			data, err := pubKeySigToSigData(cdc, pubKeys[i], multiSig.Sigs[sigIdx])
			if err != nil {
				return nil, errors.Wrapf(err, "Unable to convert Signature to SigData %d", sigIdx)
			}

			sigDatas[sigIdx] = data
			multisig.AddSignature(signatures, data, sigIdx)
			sigIdx++
		}
	}

	return signatures, nil
}
