// package api

// import (
// 	"sync"
// 	"testing"

// 	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
// 	"github.com/stretchr/testify/assert"
// )

// // ... mocks as shown above

// func TestAwsApiGateway(t *testing.T) {
// 	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
// 		api, err := NewAwsApiGateway(ctx, "test-api", )
// 		assert.NoError(t, err)

// 		var wg sync.WaitGroup
// 		wg.Add(3)

// 		// TODO(check 1): Instances have a Name tag.
// 		// TODO(check 2): Instances must not use an inline userData script.
// 		// TODO(check 3): Instances must not have SSH open to the Internet.

// 		wg.Wait()

// 		pulumi.All(api.Api.Name, api.Api.Tags).ApplyT(func(all []interface{}) error {
// 			urn := all[0].(pulumi.URN)
// 			tags := all[1].(map[string]interface{})
		
// 			assert.Containsf(t, tags, "Name", "missing a Name tag on server %v", urn)
// 			wg.Done()
// 			return nil
// 		})
// 		return nil
// 	}, pulumi.WithMocks("project", "stack", mocks(0)))
// 	assert.NoError(t, err)
// }