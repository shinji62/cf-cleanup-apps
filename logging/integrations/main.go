package main

import "github.com/shinji62/cf-cleanup-apps/logging"

func main() {
	logging.LogStd("Should OutPut this message to stdin", true)
	logging.LogError("Should OutPut this message to sdtout", "ErrorMessage")
}
