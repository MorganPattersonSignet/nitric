package bucket

import (
	"sync"
	"testing"

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

type mocks int

func (mocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	outputs := args.Inputs.Mappable()
	if args.TypeToken == "aws:s3/bucket:Bucket" {
			outputs["bucket"] = args.Name + "-s3"
	}
	return args.Name + "_id", resource.NewPropertyMapFromMap(outputs), nil
}

func (mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	outputs := map[string]interface{}{}
	return resource.NewPropertyMapFromMap(outputs), nil
}

func TestS3Bucket(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		res, err := NewS3Bucket(ctx, "test-bucket", &S3BucketArgs{
			Bucket: &v1.Bucket{},
			StackID: pulumi.String("test-stack"),
		})
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		wg.Add(2)

		pulumi.All(res.S3.URN(), res.S3.Bucket).ApplyT(func(all []interface{}) error {
			urn := all[0].(pulumi.URN)
			resourceName := all[1].(string)

			assert.Equal(t, "test-bucket-s3", resourceName, "name is invalid on bucket %v", urn)

			wg.Done()
			return nil
		})

		pulumi.All(res.S3.URN(), res.S3.Tags).ApplyT(func(all []interface{}) error {
			urn := all[0].(pulumi.URN)
			tags := all[1].(map[string]string)

			assert.Contains(t, tags, "x-nitric-name", "missing x-nitric-name on bucket %v", urn)
			assert.Equal(t, tags["x-nitric-name"], "test-bucket")

			assert.Contains(t, tags, "x-nitric-stack", "missing x-nitric-stack on bucket %v", urn)
			assert.Equal(t, tags["x-nitric-stack"], "test-stack")

			wg.Done()
			return nil
		})

		wg.Wait()

		return nil
	}, pulumi.WithMocks("project", "stack", mocks(0)))
	assert.NoError(t, err)
}