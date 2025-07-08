# Go Cloud SQL SSL CRUD App

A simple web app written in Go with CRUD operations for user data, connecting securely to Cloud SQL (PostgreSQL) using 
Cloudsql Connector. The app can run locally or be deployed to Cloud Run or GKE.

---

## ✨ Features

- CRUD User Data
- Cloud SQL PostgreSQL support
- Ready for local development, Cloud Run, or GKE

---

## 📦 Prerequisites

- Go 1.18+
- GCP Project
- Cloud SQL instance (PostgreSQL)
- Service account (For example 'cloud-sql-user@lv-playground-appdev.iam.gserviceaccount.com') with IAM roles:
  - Cloud SQL Client 
- Cloud SQL instance database user using the service account.
- Cloud SQL instance enabled for IAM authentication
---

## Database Service Account Permission
We have to set permission in the Postgresql database to allow service account to view, add, update, and delete data
in the database.

```sql

GRANT SELECT, INSERT, UPDATE, DELETE
ON ALL TABLES IN SCHEMA public
TO "cloud-sql-user@lv-playground-appdev.iam";

DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN
        SELECT schemaname, sequencename
        FROM pg_sequences
    LOOP
        EXECUTE format(
            'GRANT USAGE, UPDATE ON SEQUENCE %I.%I TO "%s";',
            r.schemaname,
            r.sequencename,
            'cloud-sql-user@lv-playground-appdev.iam'
        );
    END LOOP;
END $$;



ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE
ON TABLES TO "cloud-sql-user@lv-playground-appdev.iam";


ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT USAGE, SELECT ON SEQUENCES
TO "cloud-sql-user@lv-playground-appdev.iam";
```
---
## ⚙️ Configuration

Edit `.env`:

```
driver=postgres 
instance_connection_name=<PROJECT_ID>:<REGION>:<INSTANCE_NAME>
name=your-database-name
user=your-service-account-name #Without the .gserviceaccount.com , for example: cloud-sql-user@lv-playground-appdev.iam
```

---

## 🚀 Run Locally


### Activate the service account in shell

---

#### 📌 Prerequisites

- Google Cloud SDK (`gcloud`) installed
- Sufficient permissions to create service accounts and assign roles

---

#### ✅ Steps


##### 1. Create and Download the Key File

```bash
gcloud iam service-accounts keys create ~/key.json \
  --iam-account=my-sa-name@YOUR_PROJECT_ID.iam.gserviceaccount.com
```

This will download the key file to your home directory as `key.json`.

---

##### 2. Authenticate with the Service Account

```bash
gcloud auth activate-service-account \
  --key-file=~/key.json
```

This command activates the service account for use with `gcloud`.

---

##### 3. Verify the Authentication

```bash
gcloud auth list
```

Look for a `*` next to the active service account.

---

##### 4. (Optional) Set the Active Project

```bash
gcloud config set project YOUR_PROJECT_ID
```
### Run the application

The application will run using active service account.
---

```bash
gcloud config 
git clone https://github.com/miftahul-huda/go-cloud-ssl.git
cd go-cloud-ssl
go mod tidy
go run main.go
```

Open [http://localhost:8080](http://localhost:8080)

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