package http

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ns3777k/alertmanager-webhook-space/pkg/alertmanager"
	"github.com/ns3777k/alertmanager-webhook-space/pkg/notification"
	"github.com/ns3777k/alertmanager-webhook-space/pkg/space"
)

type Handler struct {
	channelID   string
	debug       bool
	logger      *logrus.Logger
	spaceClient *space.Client
	Router      *mux.Router
}

func NewHandler(spaceClient *space.Client, logger *logrus.Logger, channelID string, debug bool) *Handler {
	h := &Handler{
		spaceClient: spaceClient,
		channelID:   channelID,
		Router:      mux.NewRouter(),
		debug:       debug,
		logger:      logger,
	}

	h.Router.HandleFunc("/healthcheck", h.Healthcheck)
	h.Router.HandleFunc("/api/v1/webhook", h.Notify)

	return h
}

func (h *Handler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Notify(w http.ResponseWriter, r *http.Request) {
	var payload alertmanager.WebHookPayload

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.debug {
		h.logger.Info(string(b))
	}

	if err := json.Unmarshal(b, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, alert := range payload.Alerts {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

		err := h.spaceClient.SendMessage(ctx, notification.MapAlertToMessage(alert, h.channelID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			cancel()
			break
		}

		cancel()
	}

	w.WriteHeader(http.StatusOK)
}
