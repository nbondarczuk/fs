package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "fs/service/config"
	"fs/service/metrics"
)

// ReadSystemHealth responds with system health feedback  for K8S
func ReadSystemHealth(w http.ResponseWriter, r *http.Request) {
	Log.Debug("Handling request " + fmt.Sprintf("[%s] %s %s %s",
		r.Method,
		r.Host,
		r.URL.Path,
		r.URL.RawQuery))

	var status int
	if metrics.IsHealthy() {
		status = http.StatusOK
	} else {
		status = http.StatusServiceUnavailable
	}

	writeResponseWithJson(w, status, nil)
}

// ReadSystemAlive responds with alivenes feedback for K8S
func ReadSystemAlive(w http.ResponseWriter, r *http.Request) {
	Log.Info("Handling request " + fmt.Sprintf("[%s] %s %s %s",
		r.Method,
		r.Host,
		r.URL.Path,
		r.URL.RawQuery))

	var status int
	if metrics.IsAlive() {
		status = http.StatusOK
	} else {
		status = http.StatusServiceUnavailable
	}

	writeResponseWithJson(w, status, nil)
}

// ReadSystemStat responds with status info about processing status info
func ReadSystemStat(w http.ResponseWriter, r *http.Request) {
	Log.Info("Handling request " + fmt.Sprintf("[%s] %s %s %s",
		r.Method,
		r.Host,
		r.URL.Path,
		r.URL.RawQuery))

	dataReplyResource := StatResource{
		Status: true,
		Data:   Stat(metrics.StatInfo()),
	}

	jstr, err := json.Marshal(dataReplyResource)
	if err != nil {
		displayAppError(w, err,
			"Error json encoding version info",
			http.StatusInternalServerError)
		return
	}

	writeResponseWithJson(w, http.StatusOK, jstr)
}

// ReadSystemVersion responds with version info
func ReadSystemVersion(w http.ResponseWriter, r *http.Request) {
	Log.Info("Handling request " + fmt.Sprintf("[%s] %s %s %s",
		r.Method,
		r.Host,
		r.URL.Path,
		r.URL.RawQuery))

	dataReplyResource := VersionResource{
		Status: true,
		Data:   Version{OneLineVersionInfo()},
	}

	jstr, err := json.Marshal(dataReplyResource)
	if err != nil {
		displayAppError(w, err,
			"Error json encoding version info",
			http.StatusInternalServerError)
		return
	}

	writeResponseWithJson(w, http.StatusOK, jstr)
}
