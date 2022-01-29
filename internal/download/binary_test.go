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

package download_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bazelbuild/bazelisk/httputil"
	"github.com/stretchr/testify/require"

	"github.com/dio/authservicebinary/internal/download"
)

var (
	seed = time.Now()
)

type fakeClock struct {
	now          time.Time
	SleepPeriods []time.Duration
}

func newFakeClock() *fakeClock {
	return &fakeClock{now: seed}
}

func (fc *fakeClock) Sleep(d time.Duration) {
	fc.now = fc.now.Add(d)
	fc.SleepPeriods = append(fc.SleepPeriods, d)
}

func (fc *fakeClock) Now() time.Time {
	return fc.now
}

func (fc *fakeClock) TimesSlept() int {
	return len(fc.SleepPeriods)
}

func setUp() (*httputil.FakeTransport, *fakeClock) {
	transport := httputil.NewFakeTransport()
	httputil.DefaultTransport = transport

	clock := newFakeClock()
	httputil.RetryClock = clock
	return transport, clock
}

// TestDownloadVersionedBinarySuccessOnFirstTry tests if we can download the archive and get the
// extracted file.
func TestDownloadVersionedBinarySuccessOnFirstTry(t *testing.T) {
	transport, _ := setUp()
	data, err := os.ReadFile(filepath.Join("testdata", "archive.tar.gz"))
	require.NoError(t, err)
	transport.AddResponse(download.GetArchiveURL("0.6.0-rc0"), 200, string(data), nil)
	downloaded, err := download.VersionedBinary(context.Background(), "0.6.0-rc0", t.TempDir(), "auth_server")
	require.NoError(t, err)
	require.FileExists(t, downloaded)
}
