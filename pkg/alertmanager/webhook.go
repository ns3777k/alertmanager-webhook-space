package alertmanager

type WebHookPayload struct {
	Version         string   `json:"version"`
	GroupKey        string   `json:"groupKey"`
	TruncatedAlerts int      `json:"truncatedAlerts"`
	Status          string   `json:"status"`
	Receiver        string   `json:"receiver"`
	ExternalURL     string   `json:"externalURL"`
	Alerts          []*Alert `json:"alerts"`
}

type Alert struct {
	Status string `json:"status"`
	Labels struct {
		AlertName string `json:"alertname"`
		Instance  string `json:"instance"`
		Job       string `json:"job"`
		Name      string `json:"name"`
	} `json:"labels"`
	Annotations struct {
		Description string `json:"description"`
		Summary     string `json:"summary"`
	} `json:"annotations"`
	StartsAt string `json:"startsAt"`
	EndsAt   string `json:"endsAt"`
	URL      string `json:"generatorURL"`
}
