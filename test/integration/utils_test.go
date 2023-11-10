package integration_test

import (
	"fmt"
	"os"
	"path"
)

// testFlagStartAddrEnv environment variable name for "--flag" used during testing.
const testFlagStartAddrEnv = "EDSRV_TEST_FLAG_START_ADDR"

// staticFilePayload static file payload (content), the file get copied into any
// file that's going to be modified by the edit server.
var staticFilePayload = []byte("payload")

// getAddrFlagFromEnvOrDefault retrives the value to be sed on "--addr" flag from
// the environment variable, when empty the function returns  the default instead.
func getAddrFlagFromEnvOrDefault() string {
	addr := os.Getenv(testFlagStartAddrEnv)
	if addr == "" {
		addr = "127.0.0.1:8929"
	}
	return addr
}

// createStaticFileOnTmpDir creates a new temporary directory with a static file
// inside, the static file full path is returned and also error, when applicable.
func createStaticFileOnTmpDir() (string, error) {
	dir, err := os.MkdirTemp("/tmp", fmt.Sprintf("%s-", suiteName))
	if err != nil {
		return "", err
	}

	tmpFile := path.Join(dir, suiteName)
	if err = os.WriteFile(tmpFile, staticFilePayload, 0o600); err != nil {
		return "", err
	}
	return tmpFile, nil
}

// stringMapToSlice flatten a string based map into a slice.
func stringMapToSlice(m map[string]string) []string {
	slice := []string{}
	for k, v := range m {
		slice = append(slice, fmt.Sprintf("%s=%s", k, v))
	}
	return slice
}
