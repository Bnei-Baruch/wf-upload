package common

import "os"

var (
	PORT      = os.Getenv("LISTEN_ADDRESS")
	ACC_URL   = os.Getenv("ACC_URL")
	SKIP_AUTH = os.Getenv("SKIP_AUTH") == "true"
	LOG_PATH    = os.Getenv("LOG_PATH")
)

