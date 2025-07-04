# Go Cloud SQL SSL CRUD App

A simple web app written in Go with CRUD operations for user data, connecting securely to Cloud SQL (PostgreSQL or MySQL) using SSL certificates. The app can run locally or be deployed to Cloud Run or GKE.

---

## ✨ Features

- CRUD User Data
- Cloud SQL PostgreSQL or MySQL support
- SSL connection using Cloud SQL client certificates
- Switch DB type via config (`postgres` or `mysql`)
- Ready for local development, Cloud Run, or GKE

---

## 📦 Prerequisites

- Go 1.18+
- GCP Project
- Cloud SQL instance (PostgreSQL or MySQL)
- Service account with IAM roles:
  - Cloud SQL Client
  - Storage Object Viewer

---

## ⚙️ Configuration

Edit `config.yaml`:

```yaml
database:
  driver: postgres  # or mysql
  host: database-host
  port: 5432
  name: your-database-name
  user: your-db-user
  password: your-db-password
ssl:
  enabled: true
  ca_cert: ./cert/server-ca.pem
  client_cert: ./cert/client-cert.pem
  client_key: ./cert/client-key.pem
```

---

## 🔐 Get Cloud SQL SSL Certificates

1. Go to **Cloud SQL > your instance** in GCP Console
2. Go to **Connections > Security > SSL**
3. Create client cert if needed
4. Download:
   - `server-ca.pem`
   - `client-cert.pem`
   - `client-key.pem`
5. Place them in `./cert/` folder

Run:
```bash
chmod 600 ./cert/client-key.pem
```

---

## 🚀 Run Locally

```bash
git clone https://github.com/miftahul-huda/go-cloud-ssl.git
cd go-cloud-ssl
go mod tidy
go run main.go
```

Open [http://localhost:8080](http://localhost:8080)

---

## ☁️ Deploy to Cloud Run

```bash
gcloud builds submit --tag gcr.io/YOUR_PROJECT_ID/go-cloud-ssl

gcloud run deploy go-cloud-ssl   --image gcr.io/YOUR_PROJECT_ID/go-cloud-ssl   --region asia-southeast2   --add-cloudsql-instances YOUR_PROJECT_ID:asia-southeast2:your-instance   --set-env-vars INSTANCE_CONNECTION_NAME=YOUR_PROJECT_ID:asia-southeast2:your-instance   --allow-unauthenticated
```

---

## 🐳 Deploy to GKE

```bash
docker build -t gcr.io/YOUR_PROJECT_ID/go-cloud-ssl .
docker push gcr.io/YOUR_PROJECT_ID/go-cloud-ssl

kubectl create secret generic cloudsql-client-certs   --from-file=cert/server-ca.pem   --from-file=cert/client-cert.pem   --from-file=cert/client-key.pem
```

Use `deployment.yaml` to mount the secret and deploy.

---

## 📋 Sample users Table

```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100),
  email VARCHAR(100)
);
```

---

## 🔧 Troubleshooting

| Error | Solution |
|-------|----------|
| `pq: SSL is not enabled` | Enable SSL on your Cloud SQL instance |
| `x509: cannot validate certificate` | Use DNS in host field or regenerate cert |
| `client-key.pem has world-access` | Run `chmod 600 client-key.pem` |

---