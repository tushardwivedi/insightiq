# Environment Variables Setup

## Security First! 🔒

All credentials and sensitive configuration have been moved to environment variables. **Never commit real credentials to git.**

## Quick Start

1. **Copy example files:**
   ```bash
   cp .env.example .env
   cp backend/.env.example backend/.env
   cp frontend/.env.example frontend/.env
   ```

2. **Update credentials in `.env`:**
   ```bash
   # Change these values!
   POSTGRES_PASSWORD=your_secure_password_here
   SUPERSET_PASSWORD=your_secure_admin_password
   SECRET_KEY=your_random_secret_key_here
   ```

3. **Start services:**
   ```bash
   docker compose up -d
   ```

## Environment Variables

### Database Configuration
- `POSTGRES_DB` - Database name (default: superset)
- `POSTGRES_USER` - Database user (default: superset)
- `POSTGRES_PASSWORD` - Database password ⚠️ **CHANGE THIS**
- `POSTGRES_URL` - Full connection string

### Superset Configuration
- `SUPERSET_USERNAME` - Admin username (default: admin)
- `SUPERSET_PASSWORD` - Admin password ⚠️ **CHANGE THIS**
- `SUPERSET_EMAIL` - Admin email
- `SUPERSET_URL` - Superset service URL

### Service URLs
- `OLLAMA_URL` - LLM service URL
- `WHISPER_URL` - Speech-to-text service URL
- `NEXT_PUBLIC_API_URL` - Frontend API endpoint

### Security
- `SECRET_KEY` - Superset secret key ⚠️ **CHANGE THIS**

## Production Security

🚨 **Before deploying to production:**

1. Generate strong passwords:
   ```bash
   # Generate secure passwords
   openssl rand -base64 32  # For passwords
   openssl rand -hex 32     # For secret keys
   ```

2. Use environment variables or secrets management
3. Never commit `.env` files to version control
4. Regularly rotate credentials

## Files Structure

```
.env.example          # Template for all services
backend/.env.example  # Backend-specific template
frontend/.env.example # Frontend-specific template
.env                 # Your actual config (GITIGNORED)
.gitignore           # Protects your secrets
```

## Troubleshooting

- If services fail to start, check your `.env` file exists
- Verify database credentials match between services
- Check Docker logs: `docker compose logs [service-name]`