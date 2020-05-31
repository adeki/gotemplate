package server

import (
  "net/http"
  "fmt"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
  renderJSON(w, http.StatusOK, map[string]string{"message": "Hello world"})
}
