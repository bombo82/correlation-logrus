/*
 * Copyright (C) 2023 Gianni Bombelli <bombo82@giannibombelli.it>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package log_utils

import (
	"context"
	"net/http"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const CorrelationIdKey = "correlation_id"
const LoggerKey = "logger_in_ctx"

func InitLogger() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(customJSONFormatter())
	logrus.SetLevel(logrus.InfoLevel)
}

func GetContextLoggerFromHttpRequest(r *http.Request) *logrus.Entry {
	return GetContextLogger(r.Context())
}

func GetContextLogger(ctx context.Context) *logrus.Entry {
	if ctx == nil {
		InitLogger()
		logrus.Fatal("Could not get context info for logger!")
	}

	var loggerEntry any = nil
	if ctx == context.Background() {
		logger := createLogger()
		loggerEntry = logrus.NewEntry(logger)
	} else {
		loggerEntry = ctx.Value(LoggerKey)
	}

	return loggerEntry.(*logrus.Entry)
}

func createLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(customJSONFormatter())
	logger.SetLevel(logrus.DebugLevel)

	return logger
}

func GetCorrelationIdFromHttpRequest(r *http.Request) string {
	return GetCorrelationId(r.Context())
}

func GetCorrelationId(ctx context.Context) string {
	return ctx.Value(CorrelationIdKey).(string)
}

func customJSONFormatter() *logrus.JSONFormatter {
	return &logrus.JSONFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			return frame.Function, fileName
		},
	}
}

func LogWithCorrelationIdMiddleware(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		requestCtx := r.Context()

		logger := logrus.WithField(CorrelationIdKey, GetCorrelationId(requestCtx))
		ctx := context.WithValue(requestCtx, LoggerKey, logger)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return fn
}

func HttpLogInterceptorMiddleware(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		logger := GetContextLogger(r.Context())
		logger.WithFields(logrus.Fields{
			"uri":      r.RequestURI,
			"method":   r.Method,
			"duration": duration.String(),
		}).Info("request completed")
	}

	return fn
}
