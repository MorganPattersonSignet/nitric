package storage_plugin

import (
	"context"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/storage"
	"github.com/nitric-dev/membrane/plugins/gcp/adapters"
	"github.com/nitric-dev/membrane/plugins/gcp/ifaces"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type StoragePlugin struct {
	sdk.UnimplementedStoragePlugin
	client    ifaces.StorageClient
	projectID string
}

func (s *StoragePlugin) getBucketByName(bucket string) (ifaces.BucketHandle, error) {
	buckets := s.client.Buckets(context.Background(), s.projectID)

	for {
		b, err := buckets.Next()

		if err != nil {
			return nil, fmt.Errorf("Unable to find bucket: %s; %v", bucket, err)
		}
		// We'll label the buckets by their name in the nitric.yaml file and use this...
		if b.Labels["x-nitric-name"] == bucket {
			bucketHandle := s.client.Bucket(b.Name)
			return bucketHandle, nil
		}
	}
}

/**
 * Retrieves a previously stored object from a Google Cloud Storage Bucket
 */
func (s *StoragePlugin) Get(bucket string, key string) ([]byte, error) {
	bucketHandle, err := s.getBucketByName(bucket)
	if err != nil {
		return nil, err
	}

	reader, err := bucketHandle.Object(key).NewReader(context.Background())
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

/**
 * Stores a new Item in a Google Cloud Storage Bucket
 */
func (s *StoragePlugin) Put(bucket string, key string, object []byte) error {
	bucketHandle, err := s.getBucketByName(bucket)

	if err != nil {
		return err
	}

	writer := bucketHandle.Object(key).NewWriter(context.Background())

	if _, err := writer.Write(object); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}

/**
 * Creates a new Storage Plugin for use in GCP
 */
func New() (sdk.StoragePlugin, error) {
	ctx := context.Background()

	credentials, credentialsError := google.FindDefaultCredentials(ctx, storage.ScopeReadWrite)
	if credentialsError != nil {
		return nil, fmt.Errorf("GCP credentials error: %v", credentialsError)
	}
	// Get the
	client, err := storage.NewClient(ctx, option.WithCredentials(credentials))

	if err != nil {
		return nil, fmt.Errorf("storage client error: %v", err)
	}

	return &StoragePlugin{
		client: adapters.AdaptStorageClient(client),
	}, nil
}

func NewWithClient(client ifaces.StorageClient) (sdk.StoragePlugin, error) {
	return &StoragePlugin{
		client: client,
	}, nil
}
