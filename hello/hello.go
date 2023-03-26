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

package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/rs/xid"

	"prova/greetings"
	"prova/log_utils"
)

const CorrelationHeaderName = "X-Correlation-Id"

func main() {
	log_utils.InitLogger()

	mux := http.NewServeMux()

	mux.Handle("/",
		correlationIdMiddleware(
			log_utils.LogWithCorrelationIdMiddleware(
				log_utils.HttpLogInterceptorMiddleware(
					homePageHandler,
				),
			),
		),
	)

	err := http.ListenAndServe(":10000", mux)
	logrus.Fatal(err)
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	logger := log_utils.GetContextLoggerFromHttpRequest(r)

	hello, err := greetings.Hello(r.Context(), "Man")
	if err != nil {
		fmt.Fprintf(w, "ERRORE")
		logger.Error(err)
	}
	fmt.Fprintf(w, hello)
	logger.Info("Endpoint Hit: homePage")
}

func correlationIdMiddleware(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var correlationId string
		if correlationId = r.Header.Get(CorrelationHeaderName); correlationId == "" {
			correlationId = xid.New().String()
			r.Header.Set(CorrelationHeaderName, correlationId)
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, log_utils.CorrelationIdKey, correlationId)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return fn
}
