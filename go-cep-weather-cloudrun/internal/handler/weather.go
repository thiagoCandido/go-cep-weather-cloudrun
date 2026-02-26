package handler

import (
	"encoding/json"
	"net/http"
	"regexp"
)

type ViaCEP interface {
	Lookup(r *http.Request, cep string) (city string, uf string, found bool, err error)
}

type Weather interface {
	CurrentTempC(r *http.Request, city string, uf string) (float64, error)
}

type Deps struct {
	ViaCEP  ViaCEP
	Weather Weather
}

type Handler struct {
	deps Deps
}

func New(deps Deps) *Handler { return &Handler{deps: deps} }

var cepRe = regexp.MustCompile(`^\d{8}$`)

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if !cepRe.MatchString(cep) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("invalid zipcode"))
		return
	}

	city, uf, found, err := h.deps.ViaCEP.Lookup(r, cep)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal error"))
		return
	}
	if !found || city == "" || uf == "" {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("can not find zipcode"))
		return
	}

	tempC, err := h.deps.Weather.CurrentTempC(r, city, uf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal error"))
		return
	}

	resp := map[string]float64{
		"temp_C": round1(tempC),
		"temp_F": round1(cToF(tempC)),
		"temp_K": round1(cToK(tempC)),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func cToF(c float64) float64 { return c*1.8 + 32 }
func cToK(c float64) float64 { return c + 273 } // as specified in challenge

func round1(v float64) float64 {
	if v >= 0 {
		return float64(int64(v*10+0.5)) / 10
	}
	return float64(int64(v*10-0.5)) / 10
}
