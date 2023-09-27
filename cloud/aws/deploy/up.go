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
	"fmt"
	"io"
	"runtime/debug"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nitrictech/nitric/cloud/aws/deploy/config"
	commonDeploy "github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/output/interactive"
	"github.com/nitrictech/nitric/cloud/common/deploy/output/noninteractive"
	pulumiutils "github.com/nitrictech/nitric/cloud/common/deploy/pulumi"
	"github.com/nitrictech/nitric/cloud/common/env"
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/events"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optrefresh"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Up - Deploy requested infrastructure for a stack
func (d *DeployServer) Up(request *deploy.DeployUpRequest, stream deploy.DeployService_UpServer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			err = fmt.Errorf("recovered panic: %+v\n Stack: %s", r, stack)
		}
	}()

	details, err := commonDeploy.CommonStackDetailsFromAttributes(request.Attributes.AsMap())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	config, err := config.ConfigFromAttributes(request.Attributes.AsMap())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Bad stack configuration: %s", err)
	}

	// If we're interactive then we want to provide
	outputStream := &pulumiutils.UpStreamMessageWriter{
		Stream: stream,
	}

	// Default to the non-interactive writer
	var logWriter io.Writer = noninteractive.NewNonInterativeOutput(outputStream)
	pulumiUpExtraOpts := []optup.Option{}

	if env.IsInteractive() {
		pulumiEventChan := make(chan events.EngineEvent)
		deployModel := interactive.NewOutputModel(make(chan tea.Msg), pulumiEventChan)
		teaProgram := tea.NewProgram(deployModel, tea.WithOutput(&pulumiutils.UpStreamMessageWriter{
			Stream: stream,
		}))

		pulumiUpExtraOpts = append(pulumiUpExtraOpts, optup.EventStreams(pulumiEventChan))
		logWriter = deployModel

		// Run the output in a goroutine
		// TODO: Run non-interactive version as well...
		go teaProgram.Run()
		// Close the program when we're done
		defer teaProgram.Quit()
	}

	pulumiStack, err := NewUpProgram(context.TODO(), details, config, request.Spec)
	if err != nil {
		return err
	}

	_ = pulumiStack.SetConfig(context.TODO(), "aws:region", auto.ConfigValue{Value: details.Region})

	if config.Refresh {
		_, err := pulumiStack.Refresh(context.TODO(), optrefresh.ProgressStreams(logWriter))
		if err != nil {
			return err
		}
	}

	// Run the program
	// res, err := pulumiStack.Up(context.TODO(), optup.ProgressStreams(logWriter), optup.EventStreams(pulumiEventChan))
	res, err := pulumiStack.Up(context.TODO(), optup.ProgressStreams(logWriter))
	if err != nil {
		return err
	}

	// Send terminal message
	err = stream.Send(pulumiutils.PulumiOutputsToResult(res.Outputs))

	return err
}
