package bucket

import (
	"sync"
	"testing"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"

	"github.com/nitrictech/nitric/cloud/aws/deploy/exec"
	awslambda "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"

	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

func TestS3Notification(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		mockBucket, err := s3.NewBucket(ctx, "test-bucket", &s3.BucketArgs{})
		assert.NoError(t, err)

		mockFunc, err := awslambda.NewFunction(ctx, "test-exec", &awslambda.FunctionArgs{
			Role: pulumi.String("test role"),
		})
		assert.NoError(t, err)

		res, err := NewS3Notification(ctx, "test-notification", &S3NotificationArgs{
			Bucket: &S3Bucket{
				Name: "test-bucket",
				S3: mockBucket,
			},
			Notification: []*deploy.BucketNotificationTarget{
				{
					Target: &deploy.BucketNotificationTarget_ExecutionUnit{
						ExecutionUnit: "test-exec",
					},
					Config: &v1.BucketNotificationConfig{
						NotificationType: v1.BucketNotificationType_Created,
						NotificationPrefixFilter: "",
					},
				},
			},
			Functions: map[string]*exec.LambdaExecUnit{
				"test-exec": {
					Name: "test-exec",
					Function: mockFunc,
				},
			},
		})
		assert.NoError(t, err)

		var wg sync.WaitGroup
		wg.Add(1)


		pulumi.All(res.Notification.URN(), res.Name).ApplyT(func(all []interface{}) error {
			urn := all[0].(pulumi.URN)
			resourceName := all[1].(string)

			assert.Equal(t, "test-notification", resourceName, "name is invalid on notification %v", urn)

			wg.Done()
			return nil
		})

		wg.Wait()

		return nil
	}, pulumi.WithMocks("project", "stack", mocks(0)))
	assert.NoError(t, err)
}