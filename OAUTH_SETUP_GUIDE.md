# OAuth Setup Guide for InsightIQ

This guide will walk you through setting up GitHub and Google OAuth for InsightIQ.

## Prerequisites

- InsightIQ running via Docker Compose
- Access to GitHub and Google Cloud Console

---

## 1. GitHub OAuth Setup

### Step 1: Create GitHub OAuth App

1. Go to [GitHub Developer Settings](https://github.com/settings/developers)
2. Click **"New OAuth App"**
3. Fill in the application details:
   - **Application name**: `InsightIQ Local` (or your preferred name)
   - **Homepage URL**: `http://localhost:3000`
   - **Authorization callback URL**: `http://localhost:8080/auth/callback/github`
   - **Application description**: (optional) "Self-hosted analytics platform"

4. Click **"Register application"**

### Step 2: Get Client Credentials

1. After creating the app, you'll see the **Client ID** on the app page
2. Click **"Generate a new client secret"**
3. Copy both the **Client ID** and **Client Secret**

### Step 3: Update .env File

Open your `.env` file and update:

```bash
GITHUB_CLIENT_ID=your_actual_github_client_id_here
GITHUB_CLIENT_SECRET=your_actual_github_client_secret_here
```

---

## 2. Google OAuth Setup

### Step 1: Create Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Name it something like `InsightIQ`

### Step 2: Enable Google+ API

1. In your project, go to **APIs & Services** > **Library**
2. Search for "Google+ API"
3. Click on it and click **"Enable"**

### Step 3: Configure OAuth Consent Screen

1. Go to **APIs & Services** > **OAuth consent screen**
2. Choose **External** user type (unless you have a Google Workspace)
3. Fill in the required information:
   - **App name**: `InsightIQ`
   - **User support email**: Your email
   - **Developer contact information**: Your email
4. Click **"Save and Continue"**
5. On the Scopes page, click **"Save and Continue"** (default scopes are fine)
6. On Test users page, add your email address if using "External" type
7. Click **"Save and Continue"**

### Step 4: Create OAuth Client ID

1. Go to **APIs & Services** > **Credentials**
2. Click **"Create Credentials"** > **"OAuth client ID"**
3. Choose **"Web application"**
4. Fill in:
   - **Name**: `InsightIQ Web Client`
   - **Authorized JavaScript origins**:
     - `http://localhost:3000`
   - **Authorized redirect URIs**:
     - `http://localhost:8080/auth/callback/google`

5. Click **"Create"**

### Step 5: Get Client Credentials

1. A popup will show your **Client ID** and **Client Secret**
2. Copy both values

### Step 6: Update .env File

Open your `.env` file and update:

```bash
GOOGLE_CLIENT_ID=your_actual_google_client_id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your_actual_google_client_secret
```

---

## 3. Restart Services

After updating the `.env` file with real credentials:

```bash
# Restart the backend and frontend services
docker compose restart agent frontend

# Or restart all services
docker compose restart
```

---

## 4. Test OAuth Login

1. Open your browser to `http://localhost:3000/signup` or `http://localhost:3000/login`
2. Click on **"Continue with GitHub"** or **"Continue with Google"**
3. You should be redirected to GitHub/Google for authentication
4. After authorization, you'll be redirected back to InsightIQ and logged in

---

## Troubleshooting

### "Google login failed" or "GitHub login failed"

**Cause**: OAuth credentials are not configured or incorrect

**Solutions**:
1. Verify the Client ID and Secret in `.env` are correct
2. Check that the callback URLs match exactly:
   - GitHub: `http://localhost:8080/auth/callback/github`
   - Google: `http://localhost:8080/auth/callback/google`
3. Restart the services after updating `.env`
4. Check browser console (F12) for detailed error messages
5. Check backend logs: `docker compose logs agent --tail=50`

### "Redirect URI mismatch" error

**Cause**: The callback URL configured in GitHub/Google doesn't match the one being used

**Solution**:
- Ensure callback URLs in OAuth apps match the ones in this guide
- For GitHub: Must be exactly `http://localhost:8080/auth/callback/github`
- For Google: Must be exactly `http://localhost:8080/auth/callback/google`

### OAuth works but user data not saved

**Cause**: Database or SuperTokens configuration issue

**Solution**:
1. Check SuperTokens logs: `docker compose logs supertokens --tail=50`
2. Check postgres logs: `docker compose logs postgres --tail=50`
3. Verify SuperTokens is healthy: `curl http://localhost:3567/hello`

---

## Production Deployment

For production deployment, you'll need to:

1. Update the callback URLs to use your production domain:
   - GitHub: `https://yourdomain.com/auth/callback/github`
   - Google: `https://yourdomain.com/auth/callback/google`

2. Update `.env` file:
```bash
NEXT_PUBLIC_WEBSITE_DOMAIN=https://yourdomain.com
NEXT_PUBLIC_API_DOMAIN=https://yourdomain.com
OAUTH_CALLBACK_URL=https://yourdomain.com/auth/callback
```

3. Use secure secrets for production:
```bash
SUPERTOKENS_API_KEY=generate_a_strong_random_key_here
JWT_SECRET=generate_a_strong_random_jwt_secret_here
SECRET_KEY=generate_a_strong_random_secret_here
```

4. Re-create OAuth apps with production URLs or update existing ones

---

## Alternative: Email/Password Only

If you don't want to set up OAuth, you can still use the traditional email/password authentication:

1. **Login**: Use the existing JWT-based login at `/login`
   - Default credentials: `admin@insightiq.local` / `admin123456`

2. **Signup**: Use the email/password form on `/signup` (skip the social login buttons)

The OAuth integration is optional and the app works perfectly fine without it!

---

## Need Help?

- Check the [SuperTokens Documentation](https://supertokens.com/docs/thirdparty/introduction)
- Check the [GitHub OAuth Documentation](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/creating-an-oauth-app)
- Check the [Google OAuth Documentation](https://developers.google.com/identity/protocols/oauth2)
