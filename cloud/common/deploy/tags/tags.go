package tags

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Tags(ctx *pulumi.Context, stackID string, name string) map[string]string {
	return map[string]string{
		// Locate the unique stack by the presence of the key and the resource by its name
		fmt.Sprintf("x-nitric-stack-%s", stackID): name,
	}
}
