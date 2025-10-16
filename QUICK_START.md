# InsightIQ Quick Start Guide

## 🚀 Get Started in 5 Minutes

### 1. Prerequisites
- Docker & Docker Compose installed
- 4GB+ RAM available
- 5GB+ disk space

### 2. Installation
```bash
# Clone repository
git clone https://github.com/yourusername/insightiq.git
cd insightiq

# Configure environment
cp .env.example .env
# Edit .env and change passwords (required)

# Start services
docker compose up -d

# Wait 30-60 seconds for services to start
docker compose ps
```

### 3. Access InsightIQ
- **URL**: http://localhost:3000
- **Default Login**: `admin@insightiq.local` / `change_this_password`
  *(Use the password you set in `.env` file)*

### 4. First Steps

#### Option A: Try with Demo Data
1. Login to InsightIQ
2. Add a connector (PostgreSQL)
3. Use the built-in database:
   - Host: `postgres`
   - Port: `5432`
   - Database: `insightiq`
   - Username: `insightiq_user`
   - Password: *(from .env POSTGRES_PASSWORD)*

#### Option B: Connect Your Database
1. Click sidebar → "Add Connector"
2. Enter your database details
3. Click "Test Connection"
4. Click "Save"

### 5. Run Your First Query

#### Natural Language
Type: `"Show me all users"`

#### Voice
Click microphone → Say: *"Count total orders"*

#### SQL
```sql
SELECT * FROM users LIMIT 10;
```

## 🎯 Common Tasks

### Change Admin Password
Edit `.env`:
```bash
ADMIN_PASSWORD=your_new_secure_password
```
Restart: `docker compose restart agent`

### Add GitHub/Google Login
See [OAUTH_SETUP_GUIDE.md](./OAUTH_SETUP_GUIDE.md)

### Stop Services
```bash
docker compose down
```

### View Logs
```bash
docker compose logs -f agent      # Backend
docker compose logs -f frontend   # Frontend
```

### Backup Database
```bash
docker exec insightiq-postgres-1 pg_dump -U insightiq_user insightiq > backup.sql
```

## ❓ Getting Errors?

### "Login failed"
- Check `.env` file: `ADMIN_PASSWORD` matches what you're typing
- Restart services: `docker compose restart agent`

### "Port already in use"
```bash
# Find what's using the port
lsof -i :3000  # Frontend
lsof -i :8080  # Backend

# Change ports in docker-compose.yml if needed
```

### "Can't connect to database"
- Verify database is reachable from Docker containers
- Use `host.docker.internal` instead of `localhost` for host databases

### Services won't start
```bash
# Rebuild from scratch
docker compose down
docker compose build --no-cache
docker compose up -d
```

## 📚 Full Documentation

- **README.md** - Complete self-hosting guide
- **OAUTH_SETUP_GUIDE.md** - Social login setup
- **readmeclaude.md** - Development notes

## 💡 Pro Tips

1. **Change all default passwords** in `.env` before deployment
2. **Social login is optional** - email/password works without it
3. **Download Ollama models** for better AI responses:
   ```bash
   docker exec -it insightiq-ollama-1 ollama pull llama2
   ```
4. **Enable HTTPS** for production with Nginx/Caddy reverse proxy
5. **Backup regularly** - Database contains all your configs and queries

## 🆘 Need Help?

- Check [README.md](./README.md) for detailed docs
- Open an issue on GitHub
- Check logs: `docker compose logs`

---

**That's it! You're ready to start querying with InsightIQ! 🎉**
