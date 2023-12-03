/*******************************************************************************
 * Copyright 2023 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/
package logging

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
)

type ConsoleLogger struct {
	slog.Logger
}

func NewConsoleLogger(cfg config.LoggingInfo) ConsoleLogger {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	hostnameAtt := slog.Attr{Key: hostnameAttrKey, Value: slog.StringValue(hostname)}

	level := cfg.MinLogLevel
	if !isValidLogLevel(level) {
		level = slog.LevelInfo
	}

	return ConsoleLogger{
		Logger: *slog.New(
			newCustomJsonHandler(level).
				WithAttrs([]slog.Attr{hostnameAtt}).(*slog.JSONHandler),
		),
	}
}

func (l ConsoleLogger) Write(level slog.Level, message string, args ...any) {
	if !isValidLogLevel(level) {
		level = slog.LevelInfo
	}

	lineAtt := getLineAtt()
	applicationAtt := getApplicationAtt()

	args = append(args, applicationAtt, lineAtt)

	switch level {
	case slog.LevelInfo:
		l.Info(message, args...)
	case slog.LevelDebug:
		l.Debug(message, args...)
	case slog.LevelWarn:
		l.Warn(message, args...)
	case slog.LevelError:
		l.Logger.Error(message, args...)
	}
}

func (l ConsoleLogger) Error(message string, args ...any) {
	lineAtt := getLineAtt()
	applicationAtt := getApplicationAtt()

	args = append(args, applicationAtt, lineAtt)

	l.Logger.Error(message, args...)
}

func getLineAtt() slog.Attr {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	fileName := filepath.Base(file)
	return slog.Attr{Key: lineNumberAttrKey, Value: slog.StringValue(fileName + ":" + strconv.Itoa(line))}
}

func getApplicationAtt() slog.Attr {
	return slog.Attr{Key: applicationAttrKey, Value: slog.StringValue(os.Args[0])}
}

func isValidLogLevel(level slog.Level) bool {
	switch level {
	case slog.LevelInfo, slog.LevelDebug, slog.LevelWarn, slog.LevelError:
		return true
	default:
		return false
	}
}

func newCustomJsonHandler(minLogLevel slog.Level) *slog.JSONHandler {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     minLogLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				a.Key = "timestamp"
				inputTime := a.Value.Any().(time.Time)
				outputTime := inputTime.Format("2006-01-02T15:04:05Z")
				a.Value = slog.StringValue(outputTime)
			case slog.MessageKey:
				a.Key = "message"
			case slog.LevelKey:
				a.Key = "log-level"
				level := a.Value.Any().(slog.Level)
				a.Value = slog.StringValue(level.String())
			}
			return a
		},
	})
	return handler
}
