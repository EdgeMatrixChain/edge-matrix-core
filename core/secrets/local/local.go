package local

import (
	"errors"
	"fmt"
	"github.com/EdgeMatrixChain/edge-matrix-core/core/crypto"
	"os"
	"path/filepath"
	"sync"

	"github.com/EdgeMatrixChain/edge-matrix-core/core/helper/common"
	"github.com/EdgeMatrixChain/edge-matrix-core/core/secrets"
	"github.com/hashicorp/go-hclog"
)

// LocalSecretsManager is a SecretsManager that
// stores secrets locally on disk
type LocalSecretsManager struct {
	// Logger object
	logger hclog.Logger

	// Path to the base working directory
	path string

	// Map of known secrets and their paths
	secretPathMap map[string]string

	// Mux for the secretPathMap
	secretPathMapLock sync.RWMutex

	secretPass string
}

// SecretsManagerFactory implements the factory method
func SecretsManagerFactory(
	_ *secrets.SecretsManagerConfig,
	params *secrets.SecretsManagerParams,
) (secrets.SecretsManager, error) {
	// Set up the base object
	localManager := &LocalSecretsManager{
		logger:        params.Logger.Named(string(secrets.Local)),
		secretPathMap: make(map[string]string),
	}

	// Grab the path to the working directory
	path, ok := params.Extra[secrets.Path]
	if !ok {
		return nil, errors.New("no path specified for local secrets manager")
	}

	localManager.path, ok = path.(string)
	if !ok {
		return nil, errors.New("invalid type assertion")
	}

	pass, ok := params.Extra[secrets.Pass]
	if ok {
		localManager.secretPass, ok = pass.(string)
		if !ok {
			return nil, errors.New("invalid type assertion")
		}
	}

	// Run the initial setup
	_ = localManager.Setup()

	return localManager, nil
}

// Setup sets up the local SecretsManager
func (l *LocalSecretsManager) Setup() error {
	// The local SecretsManager initially handles only the
	// validator and networking private keys
	l.secretPathMapLock.Lock()
	defer l.secretPathMapLock.Unlock()

	subDirectories := []string{secrets.ConsensusFolderLocal, secrets.NetworkFolderLocal}

	// Set up the local directories
	if err := common.SetupDataDir(l.path, subDirectories, 0770); err != nil {
		return err
	}

	// baseDir/consensus/validator.key
	l.secretPathMap[secrets.ValidatorKey] = filepath.Join(
		l.path,
		secrets.ConsensusFolderLocal,
		secrets.ValidatorKeyLocal,
	)

	// baseDir/libp2p/libp2p.key
	l.secretPathMap[secrets.NetworkKey] = filepath.Join(
		l.path,
		secrets.NetworkFolderLocal,
		secrets.NetworkKeyLocal,
	)

	return nil
}

// GetSecret gets the local SecretsManager's secret from disk
func (l *LocalSecretsManager) GetSecret(name string) ([]byte, error) {
	l.secretPathMapLock.RLock()
	secretPath, ok := l.secretPathMap[name]
	l.secretPathMapLock.RUnlock()

	if !ok {
		return nil, secrets.ErrSecretNotFound
	}

	// Read the secret from disk
	secret, err := os.ReadFile(secretPath)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to read secret from disk (%s), %w",
			secretPath,
			err,
		)
	}
	if l.secretPass != "" {
		decryptedSecret, err := crypto.CFBDecrypt(string(secret), l.secretPass)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error decrypting your encrypted key: %s", err))
		}
		secret = []byte(decryptedSecret)
	}

	return secret, nil
}

// SetSecret saves the local SecretsManager's secret to disk
func (l *LocalSecretsManager) SetSecret(name string, value []byte) error {
	// If the data directory is not specified, skip write
	if l.path == "" {
		return nil
	}

	l.secretPathMapLock.Lock()
	secretPath, ok := l.secretPathMap[name]
	l.secretPathMapLock.Unlock()

	if !ok {
		return secrets.ErrSecretNotFound
	}

	// Checks for existing secret
	if _, err := os.Stat(secretPath); err == nil {
		return fmt.Errorf(
			"%s already initialized",
			secretPath,
		)
	}
	// Write the secret to disk
	if err := common.SaveFileSafe(secretPath, value, 0440); err != nil {
		return fmt.Errorf(
			"unable to write secret to disk (%s), %w",
			secretPath,
			err,
		)
	}

	return nil
}

// HasSecret checks if the secret is present on disk
func (l *LocalSecretsManager) HasSecret(name string) bool {
	_, err := l.GetSecret(name)

	return err == nil
}

// RemoveSecret removes the local SecretsManager's secret from disk
func (l *LocalSecretsManager) RemoveSecret(name string) error {
	l.secretPathMapLock.Lock()
	secretPath, ok := l.secretPathMap[name]
	defer l.secretPathMapLock.Unlock()

	if !ok {
		return secrets.ErrSecretNotFound
	}

	delete(l.secretPathMap, name)

	if removeErr := os.Remove(secretPath); removeErr != nil {
		return fmt.Errorf("unable to remove secret, %w", removeErr)
	}

	return nil
}
