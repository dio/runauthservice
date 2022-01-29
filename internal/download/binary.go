// Copyright 2022 Dhi Aurrahman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package download

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/bazelbuild/bazelisk/httputil"
	"github.com/codeclysm/extract"
)

const (
	DefautArchivedBinaryName = "auth_server.stripped"
	DefautBinaryName         = "auth_server"
	// This is similar to: https://github.com/dio/authservice/releases/download/v0.6.0-rc0/auth_server_0.6.0-rc0_darwin_amd64.tar.gz.
	archiveURLPattern = "https://github.com/dio/authservice/releases/download/v%s/auth_server_%s_%s_amd64.tar.gz"
)

// VersionedBinary returns the downloaded binary file path.
func VersionedBinary(ctx context.Context, version, destDir, destFile string) (string, error) {
	err := os.MkdirAll(destDir, 0o750)
	if err != nil {
		return "", fmt.Errorf("could not create directory %s: %v", destDir, err)
	}

	downloadedBinaryPath := filepath.Join(destDir, DefautArchivedBinaryName)
	if _, err := os.Stat(downloadedBinaryPath); err != nil {
		downloadURL := GetArchiveURL(version)
		// TODO(dio): Streaming the bytes from remote file. We decided to use this for skipping copying
		// the retry logic that has already implemented in github.com/bazelbuild/bazelisk/httputil.
		data, _, err := httputil.ReadRemoteFile(downloadURL, "")
		if err != nil {
			return "", fmt.Errorf("failed to read remote file: %s: %w", downloadURL, err)
		}
		buffer := bytes.NewBuffer(data)
		err = extract.Gz(ctx, buffer, destDir, nil)
		if err != nil {
			return "", err
		}
		if _, err = os.Stat(downloadedBinaryPath); err != nil {
			return "", fmt.Errorf("failed to extract the remote file from: %s: %w", downloadURL, err)
		}
		if err = os.Chmod(downloadedBinaryPath, 0o755); err != nil { //nolint:gosec
			return "", fmt.Errorf("could not chmod file %s: %v", downloadedBinaryPath, err)
		}
	}

	destinationPath := filepath.Join(destDir, destFile)
	if _, err := os.Stat(destinationPath); err != nil {
		if err = os.Rename(downloadedBinaryPath, destinationPath); err != nil {
			return "", fmt.Errorf("failed to rename to: %s: %w", destinationPath, err)
		}
	}
	return destinationPath, nil
}

// GetArchiveURL renders the archive URL pattern to return the actual archive URL.
func GetArchiveURL(version string) string {
	return fmt.Sprintf(archiveURLPattern, version, version, runtime.GOOS) // We always do amd64, ignore the GOARCH for now.
}
