# CEP -> Clima (Go) — Cloud Run Ready

# Endpoint para teste
https://go-cep-weather-cloudrun-88102390617.us-central1.run.app/weather?cep=04144000

API HTTP em Go que recebe CEP (8 dígitos), resolve cidade/UF via ViaCEP e retorna temperatura atual em:
- Celsius (temp_C)
- Fahrenheit (temp_F)
- Kelvin (temp_K)

## Endpoint
- `GET /weather?cep=01001000`
- Porta local padrão: 8080 (no Cloud Run usa `$PORT`)

## Configuração
Crie `.env` a partir do exemplo:

```bash
cp .env.example .env
# edite .env e coloque WEATHER_API_KEY
```

## Rodar em dev

Local:
```bash
go run ./cmd/server
```

Docker Compose:
```bash
docker compose up --build
```

## Testes
```bash
go test ./...
```

## Deploy no Google Cloud Run (passo a passo)

Eu não consigo publicar na sua conta GCP daqui (precisa de credenciais/projeto), então não dá pra eu te entregar uma URL ativa.
Mas o projeto está pronto para você subir no Cloud Run.

Exemplo usando gcloud:

```bash
gcloud auth login
gcloud config set project SEU_PROJETO

gcloud builds submit --tag gcr.io/SEU_PROJETO/cep-weather

gcloud run deploy cep-weather \
  --image gcr.io/SEU_PROJETO/cep-weather \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=SEU_TOKEN
```
