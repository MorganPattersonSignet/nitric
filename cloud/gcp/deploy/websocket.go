// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
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
	"fmt"

	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (a *NitricGcpPulumiProvider) Websocket(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Websocket) error {
	return fmt.Errorf("websockets aren't yet supported by Nitric on Google Cloud, remove the web sockets from your project and try again or run `nitric down` to destroy this stack")
}
