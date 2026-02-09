package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok\n")
	})
	http.HandleFunc("POST /alert", logAlert)

	slog.Info("Server starting", "port", port)
	http.ListenAndServe(":"+port, nil)
}

type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     string            `json:"startsAt"`
	EndsAt       string            `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

type AlertNotification struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	TruncatedAlerts   int               `json:"truncatedAlerts"`
	Status            string            `json:"status"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	GeneratorURL      string            `json:"generatorURL"`
	Alerts            []Alert           `json:"alerts"`
}

func logAlert(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var notification AlertNotification
	err := decoder.Decode(&notification)
	if err != nil {
		slog.Error("Error decoding webhook notification", "error", err)
		return
	}
	for _, a := range notification.Alerts {
		labels := make(map[string]string, len(a.Labels)+len(notification.CommonLabels))
		maps.Copy(labels, a.Labels)
		maps.Copy(labels, notification.CommonLabels)
		annotations := make(map[string]string, len(a.Annotations)+len(notification.CommonAnnotations))
		maps.Copy(annotations, a.Annotations)
		maps.Copy(annotations, notification.CommonAnnotations)
		slog.Info("ALERT", "generator", notification.GeneratorURL, "status", a.Status, "labels", labels, "annotations", annotations)
	}
}
