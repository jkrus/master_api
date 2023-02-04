package middlewares

import (
	"crypto/rand"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/jkrus/master_api/internal/config"
	zaplogger "github.com/jkrus/master_api/pkg/zap-logger/v6"
	"github.com/jkrus/master_api/pkg/zap-logger/v6/fields"
	"github.com/jkrus/master_api/pkg/zap-logger/v6/reqlog"
	"github.com/jkrus/master_api/pkg/zap-logger/v6/reqlog/extractors"
)

const (
	MsgMetaDataExtracted    = "мета данные запроса извлечены"
	MsgMetaDataExtractError = "ошибка получения данных запроса для logger"
	MsgPayloadExtractError  = "ошибка получения тела запроса для logger"
)

func (mw *middlewares) GetObservabilityMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// logger
			lg := logger
			extractor := extractors.NewGroup(extractors.DefaultRequest, extractors.DefaultHeader)
			defaultFields, err := extractor.Extract(r)
			if err != nil {
				lg.Error(MsgMetaDataExtractError, zap.Error(err))
			}
			hash := fields.Hash("hash", fields.WithRandomSource(rand.Reader))
			defaultFields = append(defaultFields, hash)
			lg = lg.With(defaultFields...)
			defer zaplogger.Recover(lg)

			// Payload Extractor
			payloadExtractor := extractors.NewPayload(fields.WithConfig(config.GetConfig().PayloadConfig()))
			payloadFields, err := payloadExtractor.Extract(r)
			if err != nil {
				lg.Error(MsgPayloadExtractError, zap.Error(err))
			}
			lg = lg.With(payloadFields...)

			r = reqlog.AddToRequest(r.WithContext(r.Context()), lg)
			lg.Debug(MsgMetaDataExtracted)
			next.ServeHTTP(w, r)
		})
	}
}
