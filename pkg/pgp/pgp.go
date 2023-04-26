package pgp

import (
	"crypto"
	"io"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/pkg/errors"
)

// NewPGPEntity returns a new openGPG Entity and possibly with the identity name and email,
// with the key pair of which the paths are specified.
// It possibly returns an error.
func NewPGPEntity(authorName, authorEmail, publicKey, privateKey string) (*openpgp.Entity, error) {
	// Decode the public GPG key.
	pubKey, err := DecodePublicKeyFile(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding public GPG key")
	}

	// Decode the private GPG key.
	privKey, err := DecodePrivateKeyFile(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding private GPG key")
	}

	bits, err := privKey.BitLength()
	if err != nil {
		return nil, errors.Wrap(err, "error getting private key bit length")
	}

	// Config collects a number of parameters along with sensible defaults.
	config := packet.Config{
		DefaultHash:            crypto.SHA256,
		DefaultCipher:          packet.CipherAES256,
		DefaultCompressionAlgo: packet.CompressionZLIB,
		CompressionConfig: &packet.CompressionConfig{
			Level: 9,
		},
		RSABits: int(bits),
	}

	currentTime := config.Now()
	uid := packet.NewUserId(authorName, "", authorEmail)

	// Create an entity which represents the components of an OpenPGP key: a primary public key
	// (which must be a signing key), one or more identities claimed by that key,
	// and zero or more subkeys, which may be encryption keys.
	entity := &openpgp.Entity{
		PrimaryKey: pubKey,
		PrivateKey: privKey,
		Identities: make(map[string]*openpgp.Identity),
	}

	isPrimaryKey := false

	// Create an identity which is claimed by an entity and zero or more
	// assertions by other entities about that claim.
	entity.Identities[uid.Id] = &openpgp.Identity{
		Name:   uid.Name,
		UserId: uid,
		SelfSignature: &packet.Signature{
			CreationTime: currentTime,
			SigType:      packet.SigTypePositiveCert,
			PubKeyAlgo:   privKey.PubKeyAlgo,
			Hash:         config.Hash(),
			IsPrimaryId:  &isPrimaryKey,
			FlagsValid:   true,
			FlagSign:     true,
			FlagCertify:  true,
			IssuerKeyId:  &entity.PrimaryKey.KeyId,
		},
	}

	keyLifetimeSecs := uint32(86400 * 365)

	// Add one additional key as signing and optionally encryption key.
	entity.Subkeys = make([]openpgp.Subkey, 1)
	entity.Subkeys[0] = openpgp.Subkey{
		PublicKey:  pubKey,
		PrivateKey: privKey,
		Sig: &packet.Signature{
			CreationTime:              currentTime,
			SigType:                   packet.SigTypeSubkeyBinding,
			PubKeyAlgo:                privKey.PubKeyAlgo,
			Hash:                      config.Hash(),
			PreferredHash:             []uint8{8}, // SHA-256
			FlagsValid:                true,
			FlagEncryptStorage:        true,
			FlagEncryptCommunications: true,
			IssuerKeyId:               &entity.PrimaryKey.KeyId,
			KeyLifetimeSecs:           &keyLifetimeSecs,
		},
	}

	return entity, nil
}

func DecodePublicKeyFile(filepath string) (*packet.PublicKey, error) {
	in, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	key, err := DecodePublicKey(in)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func DecodePublicKey(keyF io.Reader) (*packet.PublicKey, error) {
	block, err := armor.Decode(keyF)
	if err != nil {
		return nil, err
	}

	if block.Type != openpgp.PublicKeyType {
		return nil, errors.New("invalid public key file")
	}

	reader := packet.NewReader(block.Body)

	pkt, err := reader.Next()
	if err != nil {
		return nil, err
	}

	key, ok := pkt.(*packet.PublicKey)
	if !ok {
		return nil, errors.New("invalid public key")
	}

	return key, nil
}

func DecodePrivateKeyFile(filepath string) (*packet.PrivateKey, error) {
	in, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	key, err := DecodePrivateKey(in)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func DecodePrivateKey(keyF io.Reader) (*packet.PrivateKey, error) {
	block, err := armor.Decode(keyF)
	if err != nil {
		return nil, err
	}

	if block.Type != openpgp.PrivateKeyType {
		return nil, errors.New("invalid private key file")
	}

	reader := packet.NewReader(block.Body)

	pkt, err := reader.Next()
	if err != nil {
		return nil, err
	}

	key, ok := pkt.(*packet.PrivateKey)
	if !ok {
		return nil, errors.New("invalid private key")
	}

	return key, nil
}
