// Copyright 2021 Nitric Pty Ltd.
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

package firestore_service_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	firestore_service "github.com/nitric-dev/membrane/pkg/plugins/document/firestore"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
	test "github.com/nitric-dev/membrane/tests/plugins/document"
	. "github.com/onsi/ginkgo"
)

func startFirestoreProcess() *exec.Cmd {
	// Start Local DynamoDB
	os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")

	// Create Firestore Process
	args := []string{
		"beta",
		"emulators",
		"firestore",
		"start",
		"--host-port=localhost:8080",
	}
	cmd := exec.Command("gcloud", args[:]...)
	if err := cmd.Start(); err != nil {
		panic(fmt.Sprintf("Error starting Firestore Emulator %v : %v", cmd, err))
	}
	// Makes process killable
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return cmd
}

func stopFirestoreProcess(cmd *exec.Cmd) {
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
		fmt.Printf("\nFailed to kill Firestore %v : %v \n", cmd.Process.Pid, err)
	}
}

func createFirestoreClient(ctx context.Context) *firestore.Client {
	client, err := firestore.NewClient(ctx, "test")
	if err != nil {
		fmt.Printf("NewClient error: %v \n", err)
		panic(err)
	}

	return client
}

func deleteCollection(ctx context.Context, client *firestore.Client, collection string) error {
	colRef := client.Collection(collection)

	for {
		iter := colRef.Limit(100).Documents(ctx)
		numDeleted := 0

		batch := client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
	}
}

var _ = Describe("Firestore", func() {
	defer GinkgoRecover()

	// Start Firestore Emulator
	firestoreCmd := startFirestoreProcess()

	ctx := context.Background()
	db := createFirestoreClient(ctx)

	AfterEach(func() {
		deleteCollection(ctx, db, "customers")
		deleteCollection(ctx, db, "users")
		deleteCollection(ctx, db, "items")
		deleteCollection(ctx, db, "parentItems")
	})

	AfterSuite(func() {
		stopFirestoreProcess(firestoreCmd)
	})

	docPlugin, err := firestore_service.NewWithClient(db, ctx)
	if err != nil {
		panic(err)
	}

	test.GetTests(docPlugin)
	test.SetTests(docPlugin)
	test.DeleteTests(docPlugin)
	test.QueryTests(docPlugin)
})