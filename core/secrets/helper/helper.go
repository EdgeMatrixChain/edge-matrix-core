package helper

import (
	"errors"
	"fmt"
	"github.com/EdgeMatrixChain/edge-matrix-core/core/crypto"
	"github.com/EdgeMatrixChain/edge-matrix-core/core/network"
	"github.com/EdgeMatrixChain/edge-matrix-core/core/secrets"
	"github.com/EdgeMatrixChain/edge-matrix-core/core/secrets/awsssm"
	"github.com/EdgeMatrixChain/edge-matrix-core/core/secrets/gcpssm"
	"github.com/EdgeMatrixChain/edge-matrix-core/core/secrets/hashicorpvault"
	"github.com/EdgeMatrixChain/edge-matrix-core/core/secrets/local"
	"github.com/EdgeMatrixChain/edge-matrix-core/core/types"
	"github.com/hashicorp/go-hclog"
	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

// SetupLocalSecretsManager is a helper method for boilerplate local secrets manager setup
func SetupLocalSecretsManager(dataDir string) (secrets.SecretsManager, error) {
	return local.SecretsManagerFactory(
		nil, // Local secrets manager doesn't require a config
		&secrets.SecretsManagerParams{
			Logger: hclog.NewNullLogger(),
			Extra: map[string]interface{}{
				secrets.Path: dataDir,
			},
		},
	)
}

// setupHashicorpVault is a helper method for boilerplate hashicorp vault secrets manager setup
func setupHashicorpVault(
	secretsConfig *secrets.SecretsManagerConfig,
) (secrets.SecretsManager, error) {
	return hashicorpvault.SecretsManagerFactory(
		secretsConfig,
		&secrets.SecretsManagerParams{
			Logger: hclog.NewNullLogger(),
		},
	)
}

// setupAWSSSM is a helper method for boilerplate aws ssm secrets manager setup
func setupAWSSSM(
	secretsConfig *secrets.SecretsManagerConfig,
) (secrets.SecretsManager, error) {
	return awsssm.SecretsManagerFactory(
		secretsConfig,
		&secrets.SecretsManagerParams{
			Logger: hclog.NewNullLogger(),
		},
	)
}

// setupGCPSSM is a helper method for boilerplate Google Cloud Computing secrets manager setup
func setupGCPSSM(
	secretsConfig *secrets.SecretsManagerConfig,
) (secrets.SecretsManager, error) {
	return gcpssm.SecretsManagerFactory(
		secretsConfig,
		&secrets.SecretsManagerParams{
			Logger: hclog.NewNullLogger(),
		},
	)
}

func isEncrypted(secretsManager secrets.SecretsManager, name string) (bool, error) {
	if secretsManager.HasSecret(secrets.SecureFlag + name) {
		secureFlag, err := secretsManager.GetSecret(
			secrets.SecureFlag + name,
		)
		if err != nil {
			return false, err
		}
		if secureFlag != nil {
			return true, nil
		}
	}
	return false, nil
}

func SetEncryptedKey(secretsManager secrets.SecretsManager, name string, secretsPass string, privKey string) error {
	// encrypt key
	encryptedKey, err := crypto.CFBEncrypt(privKey, secretsPass)
	if err != nil {
		return err
	}
	// Write the ICP identity private key to the secrets manager storage
	if setErr := secretsManager.SetSecret(
		name,
		[]byte(encryptedKey),
	); setErr != nil {
		return setErr
	}
	if setErr := secretsManager.SetSecret(
		secrets.SecureFlag+name,
		[]byte(secrets.SecureTrue),
	); setErr != nil {
		return setErr
	}
	return nil
}

// InitECDSAValidatorKey creates new ECDSA key and set as a validator key
func InitECDSAValidatorKey(secretsManager secrets.SecretsManager) (types.Address, error) {
	if secretsManager.HasSecret(secrets.ValidatorKey) {
		return types.ZeroAddress, fmt.Errorf(`secrets "%s" has been already initialized`, secrets.ValidatorKey)
	}

	validatorKey, validatorKeyEncoded, err := crypto.GenerateAndEncodeECDSAPrivateKey()
	if err != nil {
		return types.ZeroAddress, err
	}

	address := crypto.PubKeyToAddress(&validatorKey.PublicKey)

	// Write the validator private key to the secrets manager storage
	if setErr := secretsManager.SetSecret(
		secrets.ValidatorKey,
		validatorKeyEncoded,
	); setErr != nil {
		return types.ZeroAddress, setErr
	}

	return address, nil
}

func EncryptECDSAValidatorKey(secretsManager secrets.SecretsManager, secretsPass string) error {
	privKey := ""
	if secretsManager.HasSecret(secrets.ValidatorKey) {
		validatorKey, err := secretsManager.GetSecret(secrets.ValidatorKey)
		if err != nil {
			return err
		}
		privKey = string(validatorKey)
	} else {
		// generate validator key
		_, validatorKeyEncoded, err := crypto.GenerateAndEncodeECDSAPrivateKey()
		if err != nil {
			return err
		}
		privKey = string(validatorKeyEncoded)
	}
	err2 := SetEncryptedKey(secretsManager, secrets.ValidatorKey, secretsPass, privKey)
	if err2 != nil {
		return err2
	}
	return nil
}

func InitNetworkingPrivateKey(secretsManager secrets.SecretsManager) (libp2pCrypto.PrivKey, error) {
	if secretsManager.HasSecret(secrets.NetworkKey) {
		return nil, fmt.Errorf(`secrets "%s" has been already initialized`, secrets.NetworkKey)
	}

	// Generate the libp2p private key
	libp2pKey, libp2pKeyEncoded, keyErr := network.GenerateAndEncodeLibp2pKey()
	if keyErr != nil {
		return nil, keyErr
	}

	// Write the networking private key to the secrets manager storage
	if setErr := secretsManager.SetSecret(
		secrets.NetworkKey,
		libp2pKeyEncoded,
	); setErr != nil {
		return nil, setErr
	}

	return libp2pKey, keyErr
}

func EncryptNetworkingPrivateKey(secretsManager secrets.SecretsManager, secretsPass string) error {
	privKey := ""
	if secretsManager.HasSecret(secrets.NetworkKey) {
		p2pKey, err := secretsManager.GetSecret(secrets.NetworkKey)
		if err != nil {
			return err
		}
		privKey = string(p2pKey)
	} else {
		// generate network key for node
		_, libp2pKeyEncoded, err := network.GenerateAndEncodeLibp2pKey()
		if err != nil {
			return err
		}
		privKey = string(libp2pKeyEncoded)
	}
	err2 := SetEncryptedKey(secretsManager, secrets.NetworkKey, secretsPass, privKey)
	if err2 != nil {
		return err2
	}
	return nil
}

// LoadValidatorAddress loads ECDSA key by SecretsManager and returns validator address
func LoadValidatorAddress(secretsManager secrets.SecretsManager) (types.Address, error) {
	if !secretsManager.HasSecret(secrets.ValidatorKey) {
		return types.ZeroAddress, nil
	}

	encodedKey, err := secretsManager.GetSecret(secrets.ValidatorKey)
	if err != nil {
		return types.ZeroAddress, err
	}

	privateKey, err := crypto.BytesToECDSAPrivateKey(encodedKey)
	if err != nil {
		return types.ZeroAddress, err
	}

	return crypto.PubKeyToAddress(&privateKey.PublicKey), nil
}

// LoadNodeID loads Libp2p key by SecretsManager and returns Node ID
func LoadNodeID(secretsManager secrets.SecretsManager) (string, error) {
	if !secretsManager.HasSecret(secrets.NetworkKey) {
		return "", nil
	}

	encodedKey, err := secretsManager.GetSecret(secrets.NetworkKey)
	if err != nil {
		return "", err
	}

	parsedKey, err := network.ParseLibp2pKey(encodedKey)
	if err != nil {
		return "", err
	}

	nodeID, err := peer.IDFromPrivateKey(parsedKey)
	if err != nil {
		return "", err
	}

	return nodeID.String(), nil
}

// InitCloudSecretsManager returns the cloud secrets manager from the provided config
func InitCloudSecretsManager(secretsConfig *secrets.SecretsManagerConfig) (secrets.SecretsManager, error) {
	var secretsManager secrets.SecretsManager

	switch secretsConfig.Type {
	case secrets.HashicorpVault:
		vault, err := setupHashicorpVault(secretsConfig)
		if err != nil {
			return secretsManager, err
		}

		secretsManager = vault
	case secrets.AWSSSM:
		AWSSSM, err := setupAWSSSM(secretsConfig)
		if err != nil {
			return secretsManager, err
		}

		secretsManager = AWSSSM
	case secrets.GCPSSM:
		GCPSSM, err := setupGCPSSM(secretsConfig)
		if err != nil {
			return secretsManager, err
		}

		secretsManager = GCPSSM
	default:
		return secretsManager, errors.New("unsupported secrets manager")
	}

	return secretsManager, nil
}
