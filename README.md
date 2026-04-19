# DIGIKEYS

**Carte Consulaire Biometrique**

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)
![Next.js](https://img.shields.io/badge/Next.js-14-000000?logo=next.js&logoColor=white)
![React Native](https://img.shields.io/badge/React_Native-0.73-61DAFB?logo=react&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)
![License](https://img.shields.io/badge/License-Proprietary-red)

---

## Overview

DIGIKEYS is a **biometric consular card platform** for the diaspora of Burkina Faso and the Democratic Republic of Congo. The system manages the complete lifecycle of consular identity cards -- from biometric enrollment at embassies and consulates worldwide, through card production and delivery, to verification and renewal.

> *"Plus qu'une carte consulaire, un veritable outil au service du developpement"*

The platform enables embassies to enroll citizens with biometric data (fingerprints, photo, signature), generate ICAO-compliant machine-readable zones (MRZ), manage card printing and delivery, and provide financial services through the Fonds de Solidarite Burkinabe (FSB) and partner bank account opening.

---

## Vision

**Une carte pour tous, tous pour le Burkina / Congo**

DIGIKEYS transforms the consular card from a simple identity document into a comprehensive tool for diaspora development:

- **Identity** -- Secure biometric identification for citizens abroad
- **Financial Inclusion** -- Bank account opening and mobile money integration
- **Solidarity** -- Automatic FSB/FSC contributions with every card issuance
- **Development** -- Data-driven insights on diaspora populations for policy making
- **Services** -- Communication channel between embassies and citizens

---

## Key Metrics

| Metric | Target |
|--------|--------|
| **Cards to Issue** | 6,000,000 across all embassies |
| **FSB Revenue** | 9 billion FCFA from solidarity contributions |
| **Savings to Government** | 72 billion FCFA in fraud prevention and administrative efficiency |
| **Countries Served** | 2 (Burkina Faso + DRC), scalable to more |

---

## Architecture

DIGIKEYS uses a **centralized intranet architecture** (no blockchain) optimized for embassy environments with potentially limited connectivity. The system supports offline-first mobile enrollment with batch synchronization.

```
                     +----------------------------+
                     |     Next.js Admin Portal    |
                     |   (Embassy Dashboard/HQ)    |
                     +-------------+--------------+
                                   |
                     +-------------+--------------+
                     |  React Native Mobile App   |
                     |  (Biometric Enrollment)     |
                     |  [Offline-First]            |
                     +-------------+--------------+
                                   |
                              HTTP / REST
                                   |
              +--------------------+--------------------+
              |              Go Backend (Chi)           |
              |                                         |
              |  +-----------------------------------+  |
              |  |       Application Services        |  |
              |  |  Auth | Citizen | Enrollment      |  |
              |  |  Card | Verify | Transfer | FSB   |  |
              |  |  Statistics | MRZ                 |  |
              |  +---+-------------+-------------+---+  |
              |      |             |             |      |
              |  +---+---+   +----+----+   +----+---+  |
              |  | Ports |   | Domain  |   | Ports  |  |
              |  | (In)  |   | Models  |   | (Out)  |  |
              |  +---+---+   +---------+   +----+---+  |
              |      |                          |      |
              +------+--------------------------+------+
                     |                          |
      +--------------+--+    +------------------+----------+
      |   PostgreSQL    |    |   External Integrations     |
      |   (All State)   |    |                             |
      +-----------------+    |  - MinIO (photos/docs)      |
                             |  - Banking APIs             |
                             |  - SMS Gateway              |
                             |  - Printing Service         |
                             |  - National ID Validation   |
                             +-----------------------------+
```

---

## Features

### Biometric Enrollment
- Four-fingerprint capture (right thumb, right index, left thumb, left index)
- Quality scoring per fingerprint
- Facial photo capture with hash integrity
- Digital signature capture
- AES-256-GCM encryption of all biometric data at rest
- Capture device and location tracking
- Mobile team support with GPS coordinates

### Citizen Registration
- Full identity data capture (name, DOB, place of birth, gender, nationality)
- Passport and national ID linkage
- Country of residence and address abroad
- Province and commune of origin tracking
- Embassy assignment
- Duplicate detection

### Card Lifecycle Management
- Complete state machine from request to delivery
- Batch printing support with print batch IDs
- Renewal with previous card linkage
- Suspension and revocation

```
Card State Machine:

  +----------+     +-----------+     +-----------+     +-----------+
  | pending  +---->| approved  +---->| printing  +---->| printed   |
  +----------+     +-----------+     +-----------+     +-----+-----+
                                                             |
  +----------+     +-----------+     +-----------+     +-----+-----+
  | expired  |     | revoked   |<----+ suspended |<----+ delivered |
  +----------+     +-----------+     +-----------+     +-----+-----+
                                                             |
                                                       +-----+-----+
                                                       |  active   |
                                                       +-----+-----+
                                                             |
                                                       +-----+-----+
                                                       |  renew    +--> (new pending)
                                                       +-----------+
```

### MRZ Generation (ICAO 9303 TD1)
- Machine-readable zone for ID card format (3 lines x 30 characters)
- Full character transliteration (accented to ASCII)
- Check digit computation (mod 10, weights 7-3-1)
- Composite check digit validation
- Country-specific nationality codes (BFA / COD)

### Financial Services
- **Fonds de Solidarite Burkinabe (FSB)** / **Fonds de Solidarite Congolais (FSC)** -- Automatic solidarity contributions
- Bank account opening through partner banks
- Mobile money transfers (Orange Money, Moov Money, Airtel Money, etc.)
- Transfer tracking: savings, FSB contributions, remittances, withdrawals
- Country-specific contribution amounts (1,500 FCFA for BF, 5,000 CDF for DRC)

### Communications
- Embassy-to-citizen communications (SMS, email)
- Targeted messaging by country of residence
- Recipient count tracking
- Communication audit trail

### Administration
- Role-based access control with 7 roles
- Embassy-scoped data isolation
- Dashboard statistics
- Audit logging

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Backend** | Go 1.22+, Chi router, pgx (PostgreSQL driver) |
| **Database** | PostgreSQL 16 |
| **Object Storage** | MinIO (S3-compatible) for photos, signatures, documents |
| **Frontend** | Next.js 14, React, TypeScript |
| **Mobile** | React Native 0.73 (offline-first enrollment) |
| **Authentication** | JWT (access + refresh tokens), bcrypt |
| **Biometric Security** | AES-256-GCM encryption |
| **MRZ** | ICAO 9303 TD1 compliant generator |
| **Banking** | Partner bank API adapters |
| **Printing** | Print service adapter for card production |
| **SMS** | SMS gateway integration |
| **Containerization** | Docker, Docker Compose |

---

## Project Structure

```
carteconsulaire/
+-- backend/
|   +-- cmd/server/               # Application entrypoint (serve, migrate, seed)
|   +-- config/                   # Configuration loading
|   +-- internal/
|   |   +-- adapters/             # Infrastructure adapters
|   |   |   +-- banking/          # Partner bank integrations
|   |   |   +-- biometric/        # Biometric template engine
|   |   |   +-- http/             # HTTP router, handlers, middleware
|   |   |   +-- mrz/              # ICAO 9303 TD1 MRZ generator
|   |   |   +-- national_id/      # National ID validation
|   |   |   +-- postgres/         # Database repos & migrations (10 migrations)
|   |   |   +-- printing/         # Card printing service adapter
|   |   |   +-- storage/          # MinIO object storage
|   |   +-- application/          # Use cases / application services
|   |   |   +-- auth_service.go
|   |   |   +-- citizen_service.go
|   |   |   +-- enrollment_service.go
|   |   |   +-- card_service.go
|   |   |   +-- verification_service.go
|   |   |   +-- transfer_service.go
|   |   |   +-- fsb_service.go
|   |   |   +-- statistics_service.go
|   |   |   +-- mrz_service.go
|   |   +-- domain/               # Core domain models
|   |   |   +-- biometric.go      # Biometric data model
|   |   |   +-- card.go           # Card model & lifecycle states
|   |   |   +-- citizen.go        # Citizen registration model
|   |   |   +-- communication.go  # Embassy communications
|   |   |   +-- country.go        # BF/DRC configurations
|   |   |   +-- embassy.go        # Embassy/consulate model
|   |   |   +-- enrollment.go     # Enrollment with offline support
|   |   |   +-- transfer.go       # Financial transfers & bank accounts
|   |   |   +-- user.go           # User model & roles
|   |   +-- ports/                # Port interfaces
|   +-- pkg/                      # Shared utilities
+-- deploy/
|   +-- env/                      # Environment templates (.env.production.bf)
+-- mobile/                       # React Native enrollment app
+-- web/                          # Next.js admin dashboard
+-- docker-compose.yml
+-- Makefile
```

---

## Getting Started

### Prerequisites

- **Go** 1.22+
- **Node.js** 18+ and npm
- **Docker** and Docker Compose
- **PostgreSQL** 16 (or use Docker)

### 1. Clone the Repository

```bash
git clone https://github.com/digikeys/carteconsulaire.git
cd carteconsulaire
```

### 2. Start the Development Stack

```bash
# Start PostgreSQL and MinIO, run migrations
make dev-up
```

### 3. Seed Initial Data

```bash
make backend-seed
```

This creates a default super admin user (`admin@carteconsulaire.bf` / `admin123456`).

### 4. Run the Backend

```bash
make backend-run
```

The API server starts on `http://localhost:8081`.

### 5. Run the Frontend

```bash
make web-dev
```

### 6. Run the Mobile App

```bash
# Android
make mobile-dev-android

# iOS
make mobile-dev-ios
```

---

## API Endpoints

Base URL: `/api/v1`

### Public Routes

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/auth/login` | Authenticate and receive JWT |
| `POST` | `/auth/register` | Register a new user |
| `POST` | `/auth/refresh` | Refresh access token |
| `GET` | `/verify/{cardNumber}` | Public card verification |

### Citizens (embassy_admin, enrollment_agent, super_admin)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/citizens` | List citizens (paginated, filterable) |
| `POST` | `/citizens` | Register a new citizen |
| `GET` | `/citizens/search` | Search citizens |
| `GET` | `/citizens/{id}` | Get citizen details |
| `PUT` | `/citizens/{id}` | Update citizen information |

### Enrollments (embassy_admin, enrollment_agent, super_admin)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/enrollments` | List enrollments |
| `POST` | `/enrollments` | Create biometric enrollment |
| `POST` | `/enrollments/sync` | Sync offline enrollments |
| `GET` | `/enrollments/{id}` | Get enrollment details |
| `POST` | `/enrollments/{id}/review` | Review enrollment (approve/reject) |

### Cards (embassy_admin, print_operator, super_admin)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/cards` | List cards (filterable by status/embassy) |
| `POST` | `/cards` | Request a new card |
| `GET` | `/cards/{id}` | Get card details |
| `POST` | `/cards/{id}/approve` | Approve card request |
| `POST` | `/cards/{id}/print` | Queue for printing |
| `POST` | `/cards/{id}/printed` | Mark as printed |
| `POST` | `/cards/{id}/delivered` | Mark as delivered |
| `POST` | `/cards/{id}/activate` | Activate card |
| `POST` | `/cards/{id}/suspend` | Suspend card |
| `POST` | `/cards/{id}/revoke` | Revoke card |
| `POST` | `/cards/{id}/renew` | Initiate renewal |

### Admin (super_admin)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/admin/statistics` | Dashboard statistics (filterable by embassy) |

---

## Card Lifecycle

Each consular card progresses through a defined state machine:

```
1. PENDING      Citizen enrolled, card requested
2. APPROVED     Embassy admin reviewed and approved
3. PRINTING     Card sent to print queue
4. PRINTED      Physical card produced
5. DELIVERED    Card handed to citizen
6. ACTIVE       Card activated and in use
7. SUSPENDED    Temporarily disabled (lost/stolen report)
8. REVOKED      Permanently invalidated
9. EXPIRED      Validity period ended
```

### State Transitions

```
pending ----[approve]----> approved
approved ---[print]------> printing
printing ---[printed]----> printed
printed ----[delivered]--> delivered
delivered --[activate]---> active
active -----[suspend]----> suspended
suspended --[activate]---> active
active -----[revoke]-----> revoked
active -----[renew]------> (creates new pending card)
```

---

## MRZ Generation

DIGIKEYS generates **ICAO 9303 TD1** machine-readable zones for ID-card format consular cards. Each MRZ consists of 3 lines of 30 characters.

### TD1 Format

```
Line 1: [Type:2][Country:3][DocNumber:9][Check:1][Optional:15]
Line 2: [DOB:6][Check:1][Sex:1][Expiry:6][Check:1][Nationality:3][Optional:11][Composite:1]
Line 3: [Surname<<GivenNames padded to 30]
```

### Example Output

```
IDBFA123456789<<<<<<<<<<<<<<<<
8501011M3012315BFA<<<<<<<<<<<2
OUEDRAOGO<<AHMED<<<<<<<<<<<<<
```

### Features
- Full transliteration of French accented characters (e, e, a, c, etc.)
- Padding with `<` filler characters
- Check digit computation using ICAO mod-10 algorithm (weights: 7, 3, 1)
- Composite check digit across document number, DOB, and expiry
- Country-specific nationality codes: `BFA` (Burkina Faso), `COD` (DRC)

---

## Biometric Security

### Encryption
- All fingerprint templates are encrypted with **AES-256-GCM** before storage
- Encryption key stored at `BIOMETRIC_ENCRYPTION_KEY_PATH` (outside database)
- Fingerprint data fields are `[]byte` -- never serialized to JSON (`json:"-"`)
- Photo integrity verified via SHA-256 hash

### Fingerprint Capture
- Four fingerprints: right thumb, right index, left thumb, left index
- Per-finger quality score (0-100)
- Capture device identification
- Capture location with GPS coordinates
- Agent identification for audit trail

### Data Protection
- Biometric data is never included in API JSON responses
- Separate encryption key per deployment
- Key rotation support via `encryption_key_id` field
- Database-level SSL (`DB_SSLMODE=require`)

---

## Separate Deployments

DIGIKEYS supports multi-country deployment with the `APP_COUNTRY` environment variable controlling country-specific behavior:

| Feature | Burkina Faso (`BF`) | DRC (`CD`) |
|---------|---------------------|------------|
| **Domain** | carte.consulaire.bf | carte.consulaire.cd |
| **Server Port** | 8081 | 8082 |
| **Database** | digikeys_bf | digikeys_cd |
| **Currency** | XOF (FCFA) | CDF (Franc Congolais) |
| **Card Prefix** | CC (Carte Consulaire) | CD (Carte Diaspora) |
| **Nationality Code** | BFA | COD |
| **Solidarity Fund** | FSB (1,500 FCFA) | FSC (5,000 CDF) |
| **Mobile Money** | Orange Money, Moov Money | Orange Money RDC, Airtel Money, Africell Money, Vodacom |
| **Partner Banks** | Coris Bank, BOA Burkina, Ecobank | Rawbank, Equity BCDC, TMB, FBN Bank |
| **SMS Sender ID** | DIGIKEYS-BF | DIGIKEYS-CD |
| **MinIO Bucket** | digikeys-bf | digikeys-cd |
| **Biometric Key** | /opt/digikeys/keys/biometric.key | /opt/digikeys/keys/biometric.key |

### Deploying for a Specific Country

```bash
# Burkina Faso
cp deploy/env/.env.production.bf /opt/digikeys/.env

# DRC
cp deploy/env/.env.production.cd /opt/digikeys/.env
```

---

## Mobile Enrollment App

The React Native mobile app is designed for **mobile enrollment teams** who travel to diaspora communities for on-site registration.

### Key Features
- **Offline-first** -- Full enrollment data captured without internet
- **Batch sync** -- Upload completed enrollments when connectivity is available
- **GPS tracking** -- Automatic location tagging of enrollment sessions
- **Biometric capture** -- Integration with fingerprint scanners and cameras
- **Team management** -- Enrollments tagged with team ID and agent ID

### Screens
- Login / Agent authentication
- Citizen registration form
- Photo capture
- Fingerprint capture (4 fingers)
- Signature capture
- Enrollment review and submission
- Sync status dashboard
- Offline queue management

### Sync Workflow

```
1. Agent enrolls citizen offline (data stored locally)
2. App queues enrollment with sync_status = "pending"
3. When online, agent triggers sync (POST /enrollments/sync)
4. Server receives and stores enrollment
5. sync_status updated to "synced", synced_at timestamp set
6. Embassy admin reviews enrollment (POST /enrollments/{id}/review)
```

---

## User Roles

| Role | Code | Permissions |
|------|------|-------------|
| **Super Admin** | `super_admin` | Full platform access, all embassies |
| **Embassy Admin** | `embassy_admin` | Manage citizens, enrollments, cards for their embassy |
| **Enrollment Agent** | `enrollment_agent` | Register citizens and capture biometrics |
| **Print Operator** | `print_operator` | Manage card printing queue |
| **Bank Agent** | `bank_agent` | Process bank account openings |
| **Verifier** | `verifier` | Verify card authenticity |
| **Read-only** | `readonly` | View-only access |

---

## Testing

### Run All Tests

```bash
make test
```

### Backend Tests with Coverage

```bash
make backend-test
```

This runs `go test ./... -cover` across all packages.

### Test Suites

The backend includes tests for:
- Authentication service (`auth_service_test.go`)
- Card service (`card_service_test.go`)
- MRZ generation (`mrz_service_test.go`, `generator_test.go`)
- Configuration loading (`config_test.go`)

---

## Deployment

### Production Setup

1. **Provision infrastructure:**
   - Ubuntu 22.04+ server (per country)
   - PostgreSQL 16
   - MinIO instance

2. **Copy and configure environment:**

```bash
cp deploy/env/.env.production.bf /opt/digikeys/.env
chmod 600 /opt/digikeys/.env
# Edit all CHANGE_ME values
```

3. **Build:**

```bash
make backend-build
cp backend/bin/server /opt/digikeys/server
make web-build
```

4. **Run migrations and seed:**

```bash
/opt/digikeys/server migrate
/opt/digikeys/server seed
```

### systemd Service

Create `/etc/systemd/system/digikeys.service`:

```ini
[Unit]
Description=DIGIKEYS Carte Consulaire API Server
After=network.target postgresql.service

[Service]
Type=simple
User=digikeys
Group=digikeys
WorkingDirectory=/opt/digikeys
EnvironmentFile=/opt/digikeys/.env
ExecStart=/opt/digikeys/server serve
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable digikeys
sudo systemctl start digikeys
```

### Nginx Reverse Proxy

```nginx
server {
    listen 443 ssl;
    server_name carte.consulaire.bf;

    ssl_certificate     /etc/letsencrypt/live/carte.consulaire.bf/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/carte.consulaire.bf/privkey.pem;

    location /api/ {
        proxy_pass http://127.0.0.1:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location / {
        root /opt/digikeys/web;
        try_files $uri $uri/ /index.html;
    }
}
```

---

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `APP_COUNTRY` | Country code (`BF` or `CD`) | -- | Yes |
| `APP_BASE_URL` | Public base URL | -- | Yes |
| `APP_ENV` | Environment (`development`/`production`) | `development` | Yes |
| `SERVER_HOST` | HTTP server bind address | `0.0.0.0` | Yes |
| `SERVER_PORT` | HTTP server port | `8081` | Yes |
| `DB_HOST` | PostgreSQL host | `localhost` | Yes |
| `DB_PORT` | PostgreSQL port | `5432` | Yes |
| `DB_USER` | PostgreSQL user | -- | Yes |
| `DB_PASSWORD` | PostgreSQL password | -- | Yes |
| `DB_NAME` | PostgreSQL database name | -- | Yes |
| `DB_SSLMODE` | PostgreSQL SSL mode | `require` | Yes |
| `JWT_SECRET` | JWT signing secret (64+ chars) | -- | Yes |
| `STORAGE_ENDPOINT` | MinIO server endpoint | `localhost:9000` | No |
| `STORAGE_ACCESS_KEY` | MinIO access key | -- | No |
| `STORAGE_SECRET_KEY` | MinIO secret key | -- | No |
| `STORAGE_BUCKET` | MinIO bucket name | `digikeys-bf` | No |
| `BIOMETRIC_ENCRYPTION_KEY_PATH` | Path to AES-256 key for biometric encryption | -- | Yes |
| `SMS_SENDER_ID` | SMS sender identification | `DIGIKEYS-BF` | No |

---

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Follow the hexagonal architecture conventions
4. Write tests for new services
5. Ensure `make test` passes
6. Submit a pull request

### Code Style
- Go: follow `gofmt` and `go vet`
- Frontend: ESLint + Prettier
- Commit messages: conventional commits (`feat:`, `fix:`, `docs:`, etc.)

### Domain Conventions
- Use French names for domain concepts matching the real-world process (Citoyen, Carte, Ambassade)
- Keep English for code identifiers (Citizen, Card, Embassy)
- Document all state transitions in domain model comments

---

## License

**Proprietary** -- Government of Burkina Faso / Government of the Democratic Republic of Congo

This software is developed for and jointly owned by the participating governments. Unauthorized distribution, modification, or use is prohibited.

Copyright (c) 2024-2026 DIGIKEYS. All rights reserved.
