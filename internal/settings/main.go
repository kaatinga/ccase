package settings

import "os"

var (
	Path   = "."
	DryRun = os.Getenv("CCASE_DRYRUN") == "true"
)
