package env

import "os"

const ENV_INTERACTIVE = "INTERACTIVE"

func IsInteractive() bool {
	isInteractiveVal := os.Getenv(ENV_INTERACTIVE)

	// TODO: Parse as truthy value (e.g. true/false 1/0)
	return isInteractiveVal != ""
}
