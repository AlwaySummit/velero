/*
Copyright The Velero Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/vmware-tanzu/velero/internal/credentials"
	"github.com/vmware-tanzu/velero/internal/credentials/mocks"
	velerov1api "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	"github.com/vmware-tanzu/velero/pkg/generated/clientset/versioned/scheme"
)

type NewUploaderProviderTestCase struct {
	Description   string
	UploaderType  string
	RequestorType string
	ExpectedError string
	needFromFile  bool
}

func TestNewUploaderProvider(t *testing.T) {
	// Mock objects or dependencies
	ctx := context.Background()
	client := fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()
	repoIdentifier := "repoIdentifier"
	bsl := &velerov1api.BackupStorageLocation{}
	backupRepo := &velerov1api.BackupRepository{}
	credGetter := &credentials.CredentialGetter{}
	repoKeySelector := &v1.SecretKeySelector{}
	log := logrus.New()

	testCases := []NewUploaderProviderTestCase{
		{
			Description:   "When requestorType is empty, it should return an error",
			UploaderType:  "kopia",
			RequestorType: "",
			ExpectedError: "requestor type is empty",
		},
		{
			Description:   "When FileStore credential is uninitialized, it should return an error",
			UploaderType:  "kopia",
			RequestorType: "requestor",
			ExpectedError: "uninitialized FileStore credentail",
		},
		{
			Description:   "When uploaderType is kopia, it should return a KopiaUploaderProvider",
			UploaderType:  "kopia",
			RequestorType: "requestor",
			needFromFile:  true,
			ExpectedError: "invalid credentials interface",
		},
		{
			Description:   "When uploaderType is not kopia, it should return a ResticUploaderProvider",
			UploaderType:  "restic",
			RequestorType: "requestor",
			needFromFile:  true,
			ExpectedError: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			if testCase.needFromFile {
				mockFileGetter := &mocks.FileStore{}
				mockFileGetter.On("Path", &v1.SecretKeySelector{}).Return("", nil)
				credGetter.FromFile = mockFileGetter

			}
			_, err := NewUploaderProvider(ctx, client, testCase.UploaderType, testCase.RequestorType, repoIdentifier, bsl, backupRepo, credGetter, repoKeySelector, log)
			if testCase.ExpectedError == "" {
				assert.Nil(t, err)
			} else {
				assert.Contains(t, err.Error(), testCase.ExpectedError)
			}
		})
	}
}
