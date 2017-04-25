// Package salmon implements the Salmon protocol, as defined in
// http://www.salmon-protocol.org/salmon-protocol-summary.
package salmon

import (
	"crypto"
	"encoding/xml"
	"errors"
)

// TODO: JSON schema

type MagicEnv struct {
	XMLName  xml.Name    `xml:"http://salmon-protocol.org/ns/magic-env env" json:"-"`
	Data     *MagicData  `xml:"data" json:"data"`
	Encoding string      `xml:"encoding" json:"encoding"`
	Alg      string      `xml:"alg" json:"alg"`
	Sig      []*MagicSig `xml:"sig" json:"sigs"`
}

func (env *MagicEnv) UnverifiedData() ([]byte, error) {
	switch env.Encoding {
	case "base64url":
		return decodeString(env.Data.Value)
	default:
		return nil, errors.New("salmon: unknown envelope encoding")
	}
}

func (env *MagicEnv) Verify(pk crypto.PublicKey) ([]byte, error) {
	if len(env.Sig) == 0 {
		return nil, errors.New("salmon: no signature in envelope")
	}

	// TODO: check each available signature
	if err := verify(env, pk, env.Sig[0].Value); err != nil {
		return nil, err
	}

	return env.UnverifiedData()
}

type MagicData struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

type MagicSig struct {
	KeyID string `xml:"key_id,attr" json:"key_id"`
	Value string `xml:",chardata" json:"value"`
}
