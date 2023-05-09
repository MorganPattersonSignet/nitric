package exec

import (
	"fmt"
	"sync"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"

	"github.com/nitrictech/nitric/cloud/aws/deploy/config"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
)

type mocks int

func (mocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	outputs := args.Inputs.Mappable()
	fmt.Println(args.TypeToken)
	if args.TypeToken == "aws:s3/bucket:Bucket" {
			outputs["bucket"] = args.Name
	}
	return args.Name + "_id", resource.NewPropertyMapFromMap(outputs), nil
}

func (mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	outputs := map[string]interface{}{}
	return resource.NewPropertyMapFromMap(outputs), nil
}

func TestLambda(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		res, err := NewLambdaExecutionUnit(ctx, "test-collection", &LambdaExecUnitArgs{
			EnvMap: map[string]string{
				"TEST_ENV": "true",
			},
			Compute: &v1.ExecutionUnit{
				Workers: 2,
				 
			},
			Config: config.AwsLambdaConfig{
				Memory: 1024,
				Timeout: 36,
				ProvisionedConcurreny: 0,
			},
			StackID: pulumi.String("test-stack"),
		})
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		wg.Add(2)


		pulumi.All(res.Function.URN(), res.Name).ApplyT(func(all []interface{}) error {
			urn := all[0].(pulumi.URN)
			resourceName := all[1].(string)

			assert.Equal(t, "test-exec", resourceName, "name is invalid on function %v", urn)

			wg.Done()
			return nil
		})

		pulumi.All(res.Function.URN(), res.Function.Tags).ApplyT(func(all []interface{}) error {
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