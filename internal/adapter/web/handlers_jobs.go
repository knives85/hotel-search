package web

import "net/http"

// Handlers for the /jobs routes.

func (s *Server) handleJobsIndex(w http.ResponseWriter, _ *http.Request) {
	notImplemented(w, "GET /jobs")
}

func (s *Server) handleJobRow(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, "GET /jobs/"+r.PathValue("id")+"/row")
}

func (s *Server) handleJobDownload(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, "GET /jobs/"+r.PathValue("id")+"/download")
}
