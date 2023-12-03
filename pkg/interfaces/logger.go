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
package interfaces

import "log/slog"

type Logger interface {
	// Write facilitates creation and writing of a LogEntry of a specified LogLevel. The client application
	// can also supply a message and a flexible list of additional arguments. These additional arguments are
	// optional. If provided, they should be treated as a key/value pair where the key is of type LogKey.
	//
	// Write flushes the LogEntry to StdOut in JSON format.
	Write(level slog.Level, message string, args ...any)
	// Error facilitates creation and writing of a LogEntry at the Error LogLevel. The client application
	// can also supply a message and a flexible list of additional arguments. These additional arguments are
	// optional. If provided, they should be treated as a key/value pair where the key is of type LogKey.
	//
	// Write flushes the LogEntry to StdErr in JSON format.
	Error(message string, args ...any)
}
