# InsightIQ - Self-Hosted AI Analytics Platform

<div align="center">

![InsightIQ Logo](https://via.placeholder.com/150x150/6366f1/ffffff?text=InsightIQ)

**Chat with your databases using natural language, voice, and SQL**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://www.docker.com/)
[![Go](https://img.shields.io/badge/Go-1.24-00ADD8.svg)](https://golang.org/)
[![Next.js](https://img.shields.io/badge/Next.js-14-black.svg)](https://nextjs.org/)

[Features](#features) • [Quick Start](#quick-start) • [Configuration](#configuration) • [Documentation](#documentation) • [Support](#support)

</div>

---

## 🎯 Overview

InsightIQ is a **self-hosted, AI-powered analytics platform** that allows you to interact with your data using:
- 💬 **Natural Language** - Ask questions in plain English
- 🎤 **Voice Commands** - Speak your queries
- 💻 **SQL Editor** - Write traditional SQL with syntax highlighting
- 📊 **Visual Results** - Automatic chart generation

**100% Privacy-Focused**: All data stays on your infrastructure. No cloud dependencies.

---

## ✨ Features

### 🔐 Authentication
- **Email/Password Login** - Traditional authentication with JWT
- **Social Login** - OAuth integration with GitHub and Google
- **Secure Sessions** - Powered by SuperTokens
- **Auto Admin Setup** - Default admin account created on first run

### 🤖 AI-Powered Analytics
- **Natural Language Queries** - Powered by Ollama (local LLM)
- **Voice Input** - OpenAI Whisper for speech-to-text
- **SQL Generation** - AI converts natural language to SQL
- **Smart Insights** - Automatic data analysis and recommendations

### 📊 Data Visualization
- **Auto Chart Generation** - Pie charts, bar charts, line graphs
- **SQL Syntax Highlighting** - Beautiful code editor
- **Query History** - Track and replay past queries
- **Result Export** - Download data in multiple formats

### 🔌 Database Connectors
- **PostgreSQL** - Native support
- **MySQL** - Full compatibility
- **Multiple Connections** - Manage different databases
- **Connection Testing** - Verify configs before saving

### 🛠️ Developer Friendly
- **Docker Compose** - One-command deployment
- **API Access** - RESTful API for integrations
- **Extensible** - Add custom connectors and plugins
- **Open Source** - MIT License

---

## 🚀 Quick Start

### Prerequisites

- **Docker** & **Docker Compose** (v2.0+)
- **4GB RAM** minimum (8GB recommended)
- **5GB disk space** for Docker images

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/insightiq.git
cd insightiq
```

2. **Configure environment variables**
```bash
cp .env.example .env
```

Edit `.env` and update:
```bash
# Required: Change these passwords
POSTGRES_PASSWORD=your_secure_password_here
JWT_SECRET=your_jwt_secret_change_in_production
SECRET_KEY=your_secret_key_change_in_production

# Optional: Admin user (created automatically on first run)
ADMIN_EMAIL=admin@insightiq.local
ADMIN_PASSWORD=change_this_password
ADMIN_NAME=Admin User
```

3. **Start the services**
```bash
docker compose up -d
```

4. **Wait for services to be ready** (30-60 seconds)
```bash
docker compose ps
```

All services should show `healthy` status.

5. **Access InsightIQ**
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Login with**: `admin@insightiq.local` / `change_this_password`

🎉 **That's it!** You're ready to start querying your data.

---

## 📖 Documentation

### Table of Contents

1. [Configuration](#configuration)
2. [OAuth Setup (Optional)](#oauth-setup-optional)
3. [Database Connectors](#database-connectors)
4. [Using InsightIQ](#using-insightiq)
5. [API Documentation](#api-documentation)
6. [Troubleshooting](#troubleshooting)
7. [Production Deployment](#production-deployment)

---

## ⚙️ Configuration

### Environment Variables

All configuration is done via the `.env` file:

#### Database Configuration
```bash
POSTGRES_DB=insightiq
POSTGRES_USER=insightiq_user
POSTGRES_PASSWORD=your_postgres_password_here
POSTGRES_URL=postgres://insightiq_user:your_postgres_password_here@postgres:5432/insightiq?sslmode=disable
```

#### Authentication & Security
```bash
# JWT Secret for token signing
JWT_SECRET=your_jwt_secret_key_change_in_production

# General secret key
SECRET_KEY=your_secret_key_here_change_in_production

# SuperTokens API Key
SUPERTOKENS_API_KEY=your_supertokens_api_key_change_in_production

# Initial Admin User (created on first run)
ADMIN_EMAIL=admin@insightiq.local
ADMIN_PASSWORD=change_this_password
ADMIN_NAME=Admin User
```

#### AI Services Configuration
```bash
# Ollama (Local LLM)
OLLAMA_URL=http://ollama:11434

# Whisper (Speech-to-Text)
WHISPER_URL=http://whisper:9000

# Qdrant (Vector Database)
QDRANT_URL=http://qdrant:6333
```

#### OAuth Configuration (Optional)
```bash
# GitHub OAuth
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret

# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your_google_client_secret
```

#### Frontend Configuration
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WEBSITE_DOMAIN=http://localhost:3000
NEXT_PUBLIC_API_DOMAIN=http://localhost:8080
```

---

## 🔑 OAuth Setup (Optional)

Social login (GitHub/Google) is **optional**. Email/password authentication works without it.

### Why Set Up OAuth?

- **Better User Experience** - One-click login for your users
- **No Password Management** - Users don't need to remember another password
- **Secure** - Delegated authentication to trusted providers

### Important: You vs Your Users

**You (Self-Hosting Admin):**
- Create OAuth apps **once** (5-10 minutes)
- Configure credentials in `.env`
- Your users never see this configuration

**Your Users (End Users):**
- Just click "Login with GitHub/Google"
- No OAuth setup needed for them
- Works exactly like any other website

### Step-by-Step OAuth Setup

#### 1. GitHub OAuth

1. Go to [GitHub Developer Settings](https://github.com/settings/developers)
2. Click **"New OAuth App"**
3. Fill in:
   - **Application name**: `InsightIQ`
   - **Homepage URL**: `http://localhost:3000`
   - **Authorization callback URL**: `http://localhost:8080/auth/callback/github`
4. Click **"Register application"**
5. Copy **Client ID** and generate **Client Secret**
6. Update `.env`:
```bash
GITHUB_CLIENT_ID=your_actual_client_id
GITHUB_CLIENT_SECRET=your_actual_client_secret
```

#### 2. Google OAuth

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project: `InsightIQ`
3. Go to **APIs & Services** > **OAuth consent screen**
   - Choose **External** user type
   - App name: `InsightIQ`
   - Add your email for support and developer contact
4. Go to **Credentials** > **Create Credentials** > **OAuth client ID**
   - Application type: **Web application**
   - Name: `InsightIQ Web Client`
   - Authorized JavaScript origins: `http://localhost:3000`
   - Authorized redirect URIs: `http://localhost:8080/auth/callback/google`
5. Copy **Client ID** and **Client Secret**
6. Update `.env`:
```bash
GOOGLE_CLIENT_ID=your_actual_client_id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your_actual_client_secret
```

#### 3. Restart Services

```bash
docker compose restart agent frontend
```

#### 4. Test OAuth Login

1. Open http://localhost:3000/login
2. Click **"Continue with GitHub"** or **"Continue with Google"**
3. Authorize the app
4. You'll be redirected back and logged in!

**For detailed instructions, see [OAUTH_SETUP_GUIDE.md](./OAUTH_SETUP_GUIDE.md)**

---

## 🔌 Database Connectors

### Adding a Database Connection

1. **Login to InsightIQ**
2. **Click** the connector sidebar (left side)
3. **Click** "Add Connector"
4. **Fill in** connection details:
   - Name: `My PostgreSQL DB`
   - Type: `PostgreSQL`
   - Host: `your-db-host.com`
   - Port: `5432`
   - Database: `your_database`
   - Username: `your_user`
   - Password: `your_password`
5. **Click** "Test Connection"
6. **Click** "Save"

### Supported Databases

- ✅ **PostgreSQL** (Fully supported)
- ✅ **MySQL** (Fully supported)
- 🚧 **MongoDB** (Coming soon)
- 🚧 **SQL Server** (Coming soon)

### Connection Security

- ✅ Passwords encrypted at rest
- ✅ TLS/SSL support
- ✅ Connection pooling
- ✅ Query timeout protection

---

## 🎮 Using InsightIQ

### 1. Natural Language Queries

**Type your question in plain English:**
```
"Show me total sales by region for last month"
"What are the top 10 customers by revenue?"
"Find all orders above $1000"
```

**InsightIQ will:**
1. Convert your question to SQL
2. Execute the query
3. Show results in a table
4. Generate appropriate charts

### 2. Voice Commands

1. **Click** the microphone icon
2. **Speak** your query: *"Show me user signups this week"*
3. **Wait** for transcription
4. Query executes automatically

### 3. SQL Editor

For advanced users who prefer writing SQL:

```sql
SELECT
    region,
    SUM(sales_amount) as total_sales
FROM orders
WHERE order_date >= NOW() - INTERVAL '30 days'
GROUP BY region
ORDER BY total_sales DESC;
```

**Features:**
- ✨ Syntax highlighting
- 🔍 Auto-completion
- ⚡ Execute with `Ctrl+Enter`
- 📋 Copy results

### 4. Viewing Results

Results are displayed as:
- **📊 Auto Charts** - Pie, bar, line charts based on data
- **📋 Data Table** - Sortable, searchable grid
- **💾 Export** - CSV, JSON formats

---

## 🔧 API Documentation

### Authentication

All API requests require authentication via JWT token.

**Login:**
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "admin@insightiq.local",
  "password": "your_password"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "email": "admin@insightiq.local",
    "name": "Admin User"
  }
}
```

**Use Token in Requests:**
```bash
GET /api/connectors
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

### API Endpoints

#### Analytics
```bash
# Text query (natural language)
POST /api/query
{
  "query": "Show me total sales",
  "connector_id": "uuid"
}

# SQL query
POST /api/sql
{
  "sql": "SELECT * FROM users LIMIT 10",
  "connector_id": "uuid"
}

# Voice query
POST /api/voice
Content-Type: multipart/form-data
audio: <audio file>
connector_id: uuid
```

#### Connectors
```bash
# List connectors
GET /api/connectors

# Create connector
POST /api/connectors
{
  "name": "My DB",
  "type": "postgres",
  "host": "localhost",
  "port": 5432,
  "database": "mydb",
  "username": "user",
  "password": "pass"
}

# Test connector
POST /api/connectors/test-config
{
  "type": "postgres",
  "host": "localhost",
  ...
}

# Delete connector
DELETE /api/connectors/:id
```

**Full API documentation:** [API.md](./docs/API.md)

---

## 🐛 Troubleshooting

### Services Won't Start

**Check logs:**
```bash
docker compose logs agent
docker compose logs frontend
docker compose logs postgres
```

**Common fixes:**
```bash
# Rebuild containers
docker compose down
docker compose build --no-cache
docker compose up -d

# Check disk space
df -h

# Check Docker resources
docker system df
```

### Login Not Working

**Issue**: "Invalid credentials" error

**Solutions:**
1. Check admin credentials in `.env`:
```bash
grep ADMIN .env
```

2. Reset admin password:
```bash
docker compose restart agent
# Admin user recreated on restart
```

3. Check backend logs:
```bash
docker compose logs agent | grep -i auth
```

### OAuth Errors

**Issue**: "Google login failed" or "GitHub login failed"

**Cause**: OAuth credentials not configured or incorrect

**Solutions:**

1. **Check .env file:**
```bash
grep CLIENT_ID .env
grep CLIENT_SECRET .env
```

2. **Verify callback URLs match:**
   - GitHub: `http://localhost:8080/auth/callback/github`
   - Google: `http://localhost:8080/auth/callback/google`

3. **Restart services:**
```bash
docker compose restart agent frontend
```

4. **Check browser console** (F12) for detailed errors

5. **Skip OAuth - use email/password instead**

### Database Connection Failed

**Issue**: Can't connect to external database

**Solutions:**

1. **Test from container:**
```bash
docker exec -it insightiq-agent-1 /bin/sh
# Try connecting to your DB from inside container
```

2. **Check firewall rules** - Container needs network access

3. **Use host.docker.internal** for databases on host machine:
```bash
# Instead of localhost, use:
host: host.docker.internal
```

4. **Verify credentials** with database CLI first

### Frontend Can't Reach Backend

**Issue**: "Network error" or "Failed to fetch"

**Solutions:**

1. **Check backend is running:**
```bash
curl http://localhost:8080/api/health
```

2. **Check frontend environment:**
```bash
docker compose logs frontend | grep API
```

3. **Verify CORS settings** - Check `backend/internal/http/middleware.go`

### Ollama/AI Not Working

**Issue**: Queries fail or use "mock" responses

**Cause**: Ollama models not downloaded

**Solutions:**

1. **Download models:**
```bash
docker exec -it insightiq-ollama-1 ollama pull llama2
docker exec -it insightiq-ollama-1 ollama pull nomic-embed-text
```

2. **Check Ollama status:**
```bash
curl http://localhost:11434/api/tags
```

3. **Restart services:**
```bash
docker compose restart ollama agent
```

### Port Conflicts

**Issue**: "Port already in use"

**Solutions:**

1. **Find what's using the port:**
```bash
lsof -i :3000  # Frontend
lsof -i :8080  # Backend
lsof -i :5432  # PostgreSQL
```

2. **Stop conflicting service or change ports** in `docker-compose.yml`

### More Help

- Check [OAUTH_SETUP_GUIDE.md](./OAUTH_SETUP_GUIDE.md) for OAuth issues
- See [readmeclaude.md](./readmeclaude.md) for development notes
- Open an issue on GitHub

---

## 🚀 Production Deployment

### Security Checklist

Before deploying to production:

- [ ] **Change all default passwords** in `.env`
- [ ] **Use strong secrets** (32+ random characters)
- [ ] **Enable HTTPS** with valid SSL certificates
- [ ] **Update OAuth callback URLs** to production domain
- [ ] **Configure firewall rules**
- [ ] **Enable database backups**
- [ ] **Set up monitoring** (health checks, alerts)
- [ ] **Review resource limits** in `docker-compose.yml`
- [ ] **Disable debug mode**: `DEBUG=false`

### Production .env Example

```bash
# Production Database
POSTGRES_PASSWORD=$(openssl rand -base64 32)
POSTGRES_URL=postgres://insightiq_user:${POSTGRES_PASSWORD}@postgres:5432/insightiq?sslmode=require

# Production Secrets
JWT_SECRET=$(openssl rand -base64 32)
SECRET_KEY=$(openssl rand -base64 32)
SUPERTOKENS_API_KEY=$(openssl rand -base64 32)

# Production URLs
NEXT_PUBLIC_API_URL=https://yourdomain.com
NEXT_PUBLIC_WEBSITE_DOMAIN=https://yourdomain.com
NEXT_PUBLIC_API_DOMAIN=https://yourdomain.com

# Production OAuth Callbacks
OAUTH_CALLBACK_URL=https://yourdomain.com/auth/callback

# Production Admin
ADMIN_EMAIL=admin@yourdomain.com
ADMIN_PASSWORD=$(openssl rand -base64 24)
```

### HTTPS Setup

Use a reverse proxy like Nginx or Caddy:

**Caddy Example (easiest):**
```
yourdomain.com {
    reverse_proxy localhost:3000
}

api.yourdomain.com {
    reverse_proxy localhost:8080
}
```

**Nginx Example:**
```nginx
server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Backup Strategy

**Database Backups:**
```bash
# Automated backup script
docker exec insightiq-postgres-1 pg_dump -U insightiq_user insightiq > backup_$(date +%Y%m%d).sql

# Restore
docker exec -i insightiq-postgres-1 psql -U insightiq_user insightiq < backup_20250101.sql
```

### Monitoring

**Health Checks:**
```bash
# Check all services
curl https://yourdomain.com/api/health

# Monitor in cron
*/5 * * * * curl -f https://yourdomain.com/api/health || alert-script.sh
```

**Resource Monitoring:**
```bash
# Docker stats
docker stats

# Service logs
docker compose logs -f --tail=100
```

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        User Browser                          │
│                     (localhost:3000)                         │
└────────────┬───────────────────────────────────────┬─────────┘
             │                                       │
             │ HTTP/WebSocket                       │ OAuth
             │                                       │
     ┌───────▼────────┐                    ┌────────▼─────────┐
     │   Next.js      │                    │  GitHub/Google   │
     │   Frontend     │                    │  OAuth Provider  │
     │   (Port 3000)  │                    └──────────────────┘
     └───────┬────────┘
             │
             │ API Calls (JWT Auth)
             │
     ┌───────▼────────┐
     │   Go Backend   │
     │   (Port 8080)  │
     │                │
     │ ┌────────────┐ │
     │ │SuperTokens │ │◄────────┐
     │ │(Port 3567) │ │         │
     │ └────────────┘ │         │
     └───────┬────────┘         │
             │                  │
             │                  │
    ┌────────▼────────────┐     │
    │                     │     │
    │  Storage Layer      │     │
    │                     │     │
    │  ┌───────────────┐  │     │
    │  │  PostgreSQL   │  │     │
    │  │  (Port 5432)  │──┼─────┘
    │  └───────────────┘  │
    │                     │
    │  ┌───────────────┐  │
    │  │   Qdrant      │  │
    │  │  (Port 6333)  │  │
    │  └───────────────┘  │
    └─────────────────────┘
             │
             │
    ┌────────▼────────────┐
    │                     │
    │   AI Services       │
    │                     │
    │  ┌───────────────┐  │
    │  │    Ollama     │  │
    │  │  (Port 11434) │  │
    │  └───────────────┘  │
    │                     │
    │  ┌───────────────┐  │
    │  │   Whisper     │  │
    │  │  (Port 9000)  │  │
    │  └───────────────┘  │
    └─────────────────────┘
```

---

## 🤝 Contributing

We welcome contributions! Here's how:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

**Development Setup:**
```bash
# Clone repo
git clone https://github.com/yourusername/insightiq.git

# Start development environment
docker compose up -d

# Backend development
cd backend
go run cmd/agent/main.go

# Frontend development
cd frontend
npm install
npm run dev
```

---

## 📄 License

InsightIQ is open-source software licensed under the [MIT License](LICENSE).

---

## 🙏 Acknowledgments

Built with amazing open-source technologies:

- **[Next.js](https://nextjs.org/)** - React framework
- **[Go](https://golang.org/)** - Backend language
- **[SuperTokens](https://supertokens.com/)** - Authentication
- **[Ollama](https://ollama.ai/)** - Local LLM
- **[Whisper](https://github.com/openai/whisper)** - Speech-to-text
- **[PostgreSQL](https://www.postgresql.org/)** - Database
- **[Qdrant](https://qdrant.tech/)** - Vector database
- **[Docker](https://www.docker.com/)** - Containerization

---

## 📞 Support

- **Documentation**: This README + [OAUTH_SETUP_GUIDE.md](./OAUTH_SETUP_GUIDE.md)
- **Issues**: [GitHub Issues](https://github.com/yourusername/insightiq/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/insightiq/discussions)
- **Email**: support@insightiq.dev

---

<div align="center">

**⭐ Star us on GitHub if you find InsightIQ useful!**

Made with ❤️ by the InsightIQ Team

[Website](https://insightiq.dev) • [Documentation](./docs) • [Blog](https://insightiq.dev/blog)

</div>
