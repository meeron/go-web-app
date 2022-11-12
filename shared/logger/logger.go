package logger

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"time"
	"web-app/shared"
	"web-app/shared/config"
)

const (
	logDir        = "./logs"
	logFile       = "logs/app.log"
	loggerConsole = "console"
	loggerFile    = "file"
)

var writer io.Writer

type Logger struct {
	name   string
	writer io.Writer
}

func (l Logger) Info(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	timeStr := time.Now().UTC().Format(time.RFC3339)
	fmt.Fprint(l.writer, fmt.Sprintf("[%s] [INFO] %s %s\n", l.name, timeStr, msg))
}

func Init() io.Writer {
	writers := make([]io.Writer, 0)

	loggerConfig := config.GetAppLogger()

	for key, val := range loggerConfig {
		if key == loggerConsole {
			writers = append(writers, configureConsole(val))
		}

		if key == loggerFile {
			writers = append(writers, configureFile(val))
		}
	}

	// Fallback to standard output if no logger has been configured
	if len(writers) == 0 {
		writers = append(writers, configureConsole(&struct{}{}))
	}

	writer = io.MultiWriter(writers...)
	return writer
}

func Create(name string) *Logger {
	return &Logger{
		name:   name,
		writer: writer,
	}
}

func configureConsole(params interface{}) io.Writer {
	return os.Stdout
}

func configureFile(params interface{}) io.Writer {
	fileConfig, ok := params.(map[string]interface{})

	if !ok {
		return configureFileWriter()
	}

	filePathValue, exists := fileConfig["path"]
	if !exists {
		return configureFileWriter()
	}

	filePath, ok := filePathValue.(string)
	if !ok {
		return configureFileWriter()
	}

	return configureFileWriter(filePath)
}

func configureFileWriter(filePathParams ...string) io.Writer {
	filePath := logFile

	if len(filePathParams) > 0 {
		filePath = filePathParams[0]
	}

	_, dirErr := os.Stat(logDir)
	if os.IsNotExist(dirErr) {
		dirErr := os.Mkdir(logDir, fs.ModeDir)
		if dirErr != nil {
			panic(dirErr)
		}
	}

	var file *os.File

	_, statErr := os.Stat(filePath)
	if os.IsNotExist(statErr) {
		file = shared.Unwrap(os.Create(filePath))
	} else {
		// TODO: Split file on each N-bytes
		file = shared.Unwrap(os.OpenFile(filePath, os.O_APPEND, os.ModePerm))
	}

	return file
}
