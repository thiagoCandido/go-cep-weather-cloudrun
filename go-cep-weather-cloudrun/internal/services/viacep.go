package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ViaCEP struct {
	client  *http.Client
	baseURL string
}

func NewViaCEP(client *http.Client, baseURL string) *ViaCEP {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	return &ViaCEP{client: client, baseURL: baseURL}
}

type viaCEPResponse struct {
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
	Erro       bool   `json:"erro"`
}

func (v *ViaCEP) Lookup(r *http.Request, cep string) (city string, uf string, found bool, err error) {
	url := fmt.Sprintf("%s/ws/%s/json/", v.baseURL, cep)
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, url, nil)
	if err != nil {
		return "", "", false, err
	}
	resp, err := v.client.Do(req)
	if err != nil {
		return "", "", false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// treat as not found for this challenge
		return "", "", false, nil
	}

	var data viaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", "", false, err
	}
	if data.Erro {
		return "", "", false, nil
	}

	city = strings.TrimSpace(data.Localidade)
	uf = strings.TrimSpace(data.UF)
	if city == "" || uf == "" {
		return "", "", false, nil
	}
	return city, uf, true, nil
}
