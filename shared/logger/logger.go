package logger

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
	"web-app/shared"
	"web-app/shared/config"
)

const (
	defaultFilePath = "logs/app.log"
	loggerConsole   = "console"
	loggerFile      = "file"
)

type fileWriteConfig struct {
	Path string
}

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
			writers = append(writers, configureFile(parseFileConfig(val)))
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

func configureFile(fileConfig *fileWriteConfig) io.Writer {
	logDir := filepath.Dir(fileConfig.Path)

	_, dirErr := os.Stat(logDir)
	if os.IsNotExist(dirErr) {
		dirErr := os.MkdirAll(logDir, fs.ModeDir)
		if dirErr != nil {
			panic(dirErr)
		}
	}

	var file *os.File

	_, statErr := os.Stat(fileConfig.Path)
	if os.IsNotExist(statErr) {
		file = shared.Unwrap(os.Create(fileConfig.Path))
	} else {
		// TODO: Split file on each N-bytes
		file = shared.Unwrap(os.OpenFile(fileConfig.Path, os.O_APPEND, os.ModePerm))
	}

	return file
}

func parseFileConfig(fileConfigValue interface{}) *fileWriteConfig {
	cfg := &fileWriteConfig{Path: defaultFilePath}

	fileConfig, ok := fileConfigValue.(map[string]interface{})
	if !ok {
		return cfg
	}

	filePath, ok := fileConfig["path"].(string)
	if ok {
		cfg.Path = filePath
	}

	return cfg
}
