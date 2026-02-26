package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/go-cep-weather-cloudrun/internal/handler"
)

type fakeVia struct {
	city  string
	uf    string
	found bool
	err   error
}

func (f fakeVia) Lookup(r *http.Request, cep string) (string, string, bool, error) {
	return f.city, f.uf, f.found, f.err
}

type fakeWeather struct {
	tempC float64
	err   error
}

func (f fakeWeather) CurrentTempC(r *http.Request, city string, uf string) (float64, error) {
	return f.tempC, f.err
}

func TestInvalidCEPReturns422(t *testing.T) {
	h := handler.New(handler.Deps{ViaCEP: fakeVia{}, Weather: fakeWeather{}})
	req := httptest.NewRequest("GET", "/weather?cep=123", nil)
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)
	if rw.Code != 422 {
		t.Fatalf("expected 422, got %d", rw.Code)
	}
	if rw.Body.String() != "invalid zipcode" {
		t.Fatalf("unexpected body: %q", rw.Body.String())
	}
}

func TestNotFoundReturns404(t *testing.T) {
	h := handler.New(handler.Deps{ViaCEP: fakeVia{found: false}, Weather: fakeWeather{}})
	req := httptest.NewRequest("GET", "/weather?cep=01001000", nil)
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)
	if rw.Code != 404 {
		t.Fatalf("expected 404, got %d", rw.Code)
	}
	if rw.Body.String() != "can not find zipcode" {
		t.Fatalf("unexpected body: %q", rw.Body.String())
	}
}

func TestSuccess200ReturnsTemps(t *testing.T) {
	h := handler.New(handler.Deps{
		ViaCEP:  fakeVia{city: "São Paulo", uf: "SP", found: true},
		Weather: fakeWeather{tempC: 28.5},
	})
	req := httptest.NewRequest("GET", "/weather?cep=01001000", nil)
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)
	if rw.Code != 200 {
		t.Fatalf("expected 200, got %d", rw.Code)
	}

	var resp map[string]float64
	if err := json.Unmarshal(rw.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json: %v", err)
	}
	if resp["temp_C"] != 28.5 {
		t.Fatalf("temp_C expected 28.5 got %v", resp["temp_C"])
	}
	if resp["temp_F"] != 83.3 {
		t.Fatalf("temp_F expected 83.3 got %v", resp["temp_F"])
	}
	if resp["temp_K"] != 301.5 {
		t.Fatalf("temp_K expected 301.5 got %v", resp["temp_K"])
	}
}
