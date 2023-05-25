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

package env

import (
	"fmt"
	"os"
)

const NITRIC_STACK = "NITRIC_STACK"

func GetNitricStack() string {
	return os.Getenv(NITRIC_STACK)
}

func GetNitricStackTag() string {
	return fmt.Sprintf("x-nitric-stack-%s", GetNitricStack())
}
