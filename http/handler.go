package http

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ns3777k/alertmanager-webhook-space/pkg/alertmanager"
	"github.com/ns3777k/alertmanager-webhook-space/pkg/notification"
	"github.com/ns3777k/alertmanager-webhook-space/pkg/space"
)

func unmarshalReader(r io.Reader, v interface{}) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &v)
}

type Handler struct {
	channelID   string
	spaceClient *space.Client
	Router      *mux.Router
}

func NewHandler(spaceClient *space.Client, channelID string) *Handler {
	h := &Handler{
		spaceClient: spaceClient,
		channelID:   channelID,
		Router:      mux.NewRouter(),
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

	if err := unmarshalReader(r.Body, &payload); err != nil {
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
