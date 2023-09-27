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

package deploy

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	commonDeploy "github.com/nitrictech/nitric/cloud/common/deploy"
	model "github.com/nitrictech/nitric/cloud/common/deploy/output/interactive"
	pulumiutils "github.com/nitrictech/nitric/cloud/common/deploy/pulumi"
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/events"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (d *DeployServer) Down(request *deploy.DeployDownRequest, stream deploy.DeployService_DownServer) error {
	details, err := commonDeploy.CommonStackDetailsFromAttributes(request.Attributes.AsMap())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	pulumiEventChan := make(chan events.EngineEvent)
	teaUpdates := make(chan tea.Msg)
	teaProgram := model.NewInteractiveOutput(teaUpdates, pulumiEventChan, &pulumiutils.DownStreamMessageWriter{
		Stream: stream,
	})
	// updateWriter := model.LogMessageSubscriptionWriter{
	// 	Sub: teaUpdates,
	// }

	// Run the output in a goroutine
	// TODO: Run non-interactive version as well...
	go teaProgram.Run()
	// Close the program when we're done
	defer teaProgram.Quit()

	// TODO: Tear down the requested stack
	// dsMessageWriter := &pulumiutils.DownStreamMessageWriter{
	// 	Stream: stream,
	// }

	s, err := auto.UpsertStackInlineSource(context.TODO(), details.FullStackName, details.Project, nil)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	// destroy the stack
	_, err = s.Destroy(context.TODO(), optdestroy.EventStreams(pulumiEventChan))
	if err != nil {
		return err
	}

	return nil
}
