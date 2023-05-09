package collection

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
	if args.TypeToken == "aws:dynamodb/table:Table" {
			outputs["name"] = args.Name + "-dynamodb"
	}
	return args.Name + "_id", resource.NewPropertyMapFromMap(outputs), nil
}

func (mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	outputs := map[string]interface{}{}
	return resource.NewPropertyMapFromMap(outputs), nil
}

func TestDynamoDB(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		res, err := NewDynamodbCollection(ctx, "test-collection", &DynamodbCollectionArgs{
			Collection: &v1.Collection{},
			StackID: pulumi.String("test-stack"),
		})
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		wg.Add(2)


		pulumi.All(res.Table.URN(), res.Table.Name).ApplyT(func(all []interface{}) error {
			urn := all[0].(pulumi.URN)
			resourceName := all[1].(string)

			assert.Equal(t, "test-collection-dynamodb", resourceName, "name is invalid on collection %v", urn)

			wg.Done()
			return nil
		})

		pulumi.All(res.Table.URN(), res.Table.Tags).ApplyT(func(all []interface{}) error {
			urn := all[0].(pulumi.URN)
			tags := all[1].(map[string]string)

			assert.Contains(t, tags, "x-nitric-name", "missing x-nitric-name on collection %v", urn)
			assert.Equal(t, tags["x-nitric-name"], "test-collection")

			assert.Contains(t, tags, "x-nitric-stack", "missing x-nitric-stack on collection %v", urn)
			assert.Equal(t, tags["x-nitric-stack"], "test-stack")

			wg.Done()
			return nil
		})

		wg.Wait()

		return nil
	}, pulumi.WithMocks("project", "stack", mocks(0)))
	assert.NoError(t, err)
}