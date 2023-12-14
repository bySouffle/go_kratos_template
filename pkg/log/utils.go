package log

import "go_kratos_template/pkg/utils"

// GetLogFile get log file absolute path
func GetLogFile(filename string, suffix string) string {
	return utils.ConcatString(logDir, "/", hostname, "/", filename, suffix)
}
