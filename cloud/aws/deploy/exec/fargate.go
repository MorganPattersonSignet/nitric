// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exec

import (
	"encoding/json"
	"fmt"
	"github.com/pulumi/pulumi-awsx/sdk/go/awsx/lb"

	awsecs "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecs"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	ecsx "github.com/pulumi/pulumi-awsx/sdk/go/awsx/ecs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nitrictech/nitric/cloud/aws/deploy/config"
	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
)

type FargateExecUnitArgs struct {
	//Client  ecsiface.ECSAPI
	StackID pulumi.StringInput
	// Image needs to be built and uploaded first
	DockerImage *image.Image
	Compute     *v1.ExecutionUnit
	EnvMap      map[string]string
	Config      config.AwsFargateConfig
}

type FargateExecUnit struct {
	pulumi.ResourceState

	Name     string
	Service  *ecsx.FargateService
	Balancer *lb.ApplicationLoadBalancer
	Role     *iam.Role
}

func NewFargateExecutionUnit(ctx *pulumi.Context, name string, args *FargateExecUnitArgs, opts ...pulumi.ResourceOption) (*FargateExecUnit, error) {
	res := &FargateExecUnit{Name: name}

	err := ctx.RegisterComponentResource("nitric:exec:AWSFargate", name, res, opts...)
	if err != nil {
		return nil, err
	}

	opts = append(opts, pulumi.Parent(res))

	// FIXME: BEFORE MERGE change for fargate
	tmpJSON, err := json.Marshal(map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Sid":    "",
				"Effect": "Allow",
				"Principal": map[string]interface{}{
					"Service": "ecs-tasks.amazonaws.com",
				},
				"Action": "sts:AssumeRole",
			},
		},
	})
	if err != nil {
		return nil, err
	}

	res.Role, err = iam.NewRole(ctx, name+"FargateRole", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(tmpJSON),
		Tags:             common.Tags(ctx, args.StackID, name+"FargateRole"),
	}, opts...)
	if err != nil {
		return nil, err
	}

	//_, err = iam.NewRolePolicyAttachment(ctx, name+"FargateBasicExecution", &iam.RolePolicyAttachmentArgs{
	//	PolicyArn: iam.ManagedPolicyAWSLambdaBasicExecutionRole, // TODO: Fargate change?
	//	Role:      res.Role.ID(),
	//}, opts...)
	//if err != nil {
	//	return nil, err
	//}

	telemeteryActions := []string{
		"xray:PutTraceSegments",
		"xray:PutTelemetryRecords",
		"xray:GetSamplingRules",
		"xray:GetSamplingTargets",
		"xray:GetSamplingStatisticSummaries",
		"ssm:GetParameters",
		"logs:CreateLogStream",
		"logs:PutLogEvents",
	}

	listActions := []string{
		"sns:ListTopics",
		"sqs:ListQueues",
		"dynamodb:ListTables",
		"s3:ListAllMyBuckets",
		"tag:GetResources",
		"apigateway:GET",
	}

	// Add resource list permissions
	// Currently the membrane will use list operations
	tmpJSON, err = json.Marshal(map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Action":   append(listActions, telemeteryActions...),
				"Effect":   "Allow",
				"Resource": "*",
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// TODO: Lock this SNS topics for which this function has pub definitions
	// FIXME: Limit to known resources
	_, err = iam.NewRolePolicy(ctx, name+"ListAccess", &iam.RolePolicyArgs{
		Role:   res.Role.ID(),
		Policy: pulumi.String(tmpJSON),
	}, opts...)
	if err != nil {
		return nil, err
	}

	//envVars := pulumi.StringMap{
	//	"NITRIC_ENVIRONMENT": pulumi.String("cloud"),
	//	"NITRIC_STACK":       args.StackID,
	//	"MIN_WORKERS":        pulumi.String(fmt.Sprint(args.Compute.Workers)),
	//}
	//for k, v := range args.EnvMap {
	//	envVars[k] = pulumi.String(v)
	//}

	//environment := ecsx.TaskDefinitionKeyValuePairArray{
	//	"NITRIC_ENVIRONMENT": pulumi.String("cloud"),
	//	"NITRIC_STACK":       args.StackID,
	//	"MIN_WORKERS":        pulumi.String(fmt.Sprint(args.Compute.Workers)),
	//}
	////append(environment, {Name:})

	envVars := ecsx.TaskDefinitionKeyValuePairArray{
		&ecsx.TaskDefinitionKeyValuePairArgs{
			Name:  pulumi.String("NITRIC_ENVIRONMENT"),
			Value: pulumi.String("cloud"),
		},
		&ecsx.TaskDefinitionKeyValuePairArgs{
			Name:  pulumi.String("NITRIC_STACK"),
			Value: args.StackID,
		},
		&ecsx.TaskDefinitionKeyValuePairArgs{
			Name:  pulumi.String("MIN_WORKERS"),
			Value: pulumi.String(fmt.Sprint(args.Compute.Workers)),
		},
	}
	for k, v := range args.EnvMap {
		envVars = append(envVars, &ecsx.TaskDefinitionKeyValuePairArgs{
			Name:  pulumi.String(k),
			Value: pulumi.String(v),
		})
	}

	// Create the ECS Cluster to run the Fargate Service
	//  this may be a shared cluster in the future.
	cluster, err := awsecs.NewCluster(ctx, name+"Cluster", nil)
	if err != nil {
		return nil, err
	}

	// FIXME: Do we need a VPC or other mechanism to protect against public access direct to the functions?

	res.Balancer, err = lb.NewApplicationLoadBalancer(ctx, "ALB", nil)
	if err != nil {
		return nil, err
	}

	res.Service, err = ecsx.NewFargateService(ctx, name, &ecsx.FargateServiceArgs{
		//AssignPublicIp: pulumi.Bool(false),
		Cluster: cluster.Arn,
		TaskDefinitionArgs: &ecsx.FargateServiceTaskDefinitionArgs{
			//ExecutionRole: res.Role.Arn, //TODO: determine if we need to update the returned role or provide one
			Container: &ecsx.TaskDefinitionContainerDefinitionArgs{
				Image:     args.DockerImage.URI(),
				Cpu:       pulumi.IntPtr(args.Config.Cpu),
				Memory:    pulumi.IntPtr(args.Config.Memory),
				Essential: pulumi.Bool(true),
				PortMappings: ecsx.TaskDefinitionPortMappingArray{
					&ecsx.TaskDefinitionPortMappingArgs{
						TargetGroup: res.Balancer.DefaultTargetGroup,
					},
				},
				Environment: envVars,
			},
		},
		Tags: common.Tags(ctx, args.StackID, name),
	})
	if err != nil {
		return nil, err
	}

	// FIXME: Add scaling rules

	// FIXME BEFORE MERGE: Do we need this for Fargate?
	// ensure that the lambda was deployed successfully
	//isHealthy := res.Function.Arn.ApplyT(func(arn string) (bool, error) {
	//	payload, _ := json.Marshal(map[string]interface{}{
	//		"x-nitric-healthcheck": true,
	//	})
	//
	//	err := retry.Do(func() error {
	//		_, err := args.Client.Invoke(&lambda.InvokeInput{
	//			FunctionName: aws.String(arn),
	//			Payload:      payload,
	//		})
	//
	//		return err
	//	}, retry.Attempts(3))
	//	if err != nil {
	//		return false, err
	//	}
	//
	//	return true, nil
	//})

	return res, ctx.RegisterResourceOutputs(res, pulumi.Map{
		"name":    pulumi.String(res.Name),
		"fargate": res.Service,
		//"healthy": isHealthy,
	})
}
