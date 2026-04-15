# LaunchDate Deployment Guide

## GitHub Release Automation

The repository now includes a GitHub Actions workflow at `.github/workflows/backend-release.yml`.

When a GitHub Release is published:

- The workflow builds the backend image from `services/backend/Dockerfile`.
- The image is pushed to GitHub Container Registry as `ghcr.io/<owner>/launchdate-backend`.
- The published image always receives the release tag, for example `ghcr.io/<owner>/launchdate-backend:v1.2.0`.
- Non-prerelease versions also receive the `latest` tag.
- If the release body is empty, the workflow asks GitHub to generate the default release notes and writes them back to the release.

Repository requirements:

- GitHub Actions must be enabled.
- The workflow must be allowed to use `GITHUB_TOKEN` with read and write permissions.
- The actor publishing the release must have permission to publish packages for the repository namespace.

Optional registry secrets:

- `GHCR_TOKEN`: optional override for package publishing when the default `GITHUB_TOKEN` cannot write to GHCR.
- `GHCR_USERNAME`: optional username paired with `GHCR_TOKEN`.

If release publishing fails with `permission_denied: write_package`:

- Check repository Settings -> Actions -> General -> Workflow permissions and set it to `Read and write permissions`.
- If the package `ghcr.io/<owner>/launchdate-backend` already exists, open its package settings and grant this repository Actions access.
- If your organization or existing package policy blocks `GITHUB_TOKEN`, create `GHCR_TOKEN` with package write permission and add `GHCR_USERNAME` for the token owner.

Release publishing recommendations:

- Create a tag such as `v1.2.0` before publishing the release.
- If you want GitHub's default release notes, leave the release description empty when publishing; the workflow will fill it automatically.
- If you manually write release notes, the workflow will preserve your content and only publish the container image.

This document describes a production deployment from zero for the current monorepo, using the following architecture:

- MongoDB: Mongo Atlas
- Object storage: Cloudflare R2
- Backend: Docker on your own server
- Backend public access: Cloudflare Tunnel
- Frontends: Cloudflare Pages for both `apps/web` and `apps/admin`

## 1. Target Architecture

Recommended domain layout:

| Service | Recommended domain | Deployment target |
| --- | --- | --- |
| Web | `launch-date.com` | Cloudflare Pages |
| Web alias | `www.launch-date.com` | Cloudflare Pages |
| Admin | `admin.launch-date.com` | Cloudflare Pages |
| API | `api.launch-date.com` | Cloudflare Tunnel -> backend container |
| Image CDN | `img.launch-date.com` | Cloudflare R2 custom domain |

This layout matches the backend CORS whitelist already present in the codebase for:

- `https://launch-date.com`
- `https://www.launch-date.com`
- `https://admin.launch-date.com`

If you use different frontend domains, update the whitelist in `services/backend/internal/middleware/cors.go` before going live.

## 2. What You Need Before Starting

- One Linux server with Docker installed
- One Cloudflare account managing your domain
- One Mongo Atlas project and cluster
- One Cloudflare R2 bucket for images
- This repository pushed to GitHub or another Git provider supported by Cloudflare Pages

Recommended server baseline:

- Ubuntu 22.04 or 24.04
- 2 vCPU+
- 4 GB RAM+
- A fixed public IP if you want to restrict Mongo Atlas access by IP

## 3. Prepare Mongo Atlas

1. Create a project in Mongo Atlas.
2. Create a cluster.
3. Create a database user with username and password.
4. In Network Access, allow your server IP.
5. If your server IP is not fixed yet, you can temporarily allow `0.0.0.0/0`, but tighten this later.
6. Copy the connection string.

Example:

```text
mongodb+srv://<db_user>:<db_password>@<cluster-url>/launchdate_db?retryWrites=true&w=majority&appName=launchdate
```

Recommended database name:

```text
launchdate_db
```

You will use this value for:

- `MONGODB_URL`
- `MONGODB_DATABASE`

## 4. Prepare Cloudflare R2

1. Open Cloudflare Dashboard -> R2.
2. Create a bucket, for example `launchdate-images`.
3. Create an API token or access key pair with read and write permissions for that bucket.
4. Record your Account ID.
5. Optional but recommended: bind a custom domain such as `img.launch-date.com` to the bucket.

Typical values:

```text
IMAGE_S3_ENDPOINT=https://<account-id>.r2.cloudflarestorage.com
IMAGE_S3_REGION=auto
IMAGE_S3_BUCKET=launchdate-images
IMAGE_S3_ACCESS_KEY=<r2-access-key>
IMAGE_S3_SECRET_KEY=<r2-secret-key>
IMAGE_DOMAIN=https://img.launch-date.com
```

Notes:

- R2 is S3-compatible, so the backend can use it through the existing S3 configuration.
- `IMAGE_DOMAIN` should be the public URL users will access, not the private API endpoint.
- If you do not configure a custom domain, use the R2 public bucket URL instead.

## 5. Prepare the Server

Install Docker and Compose plugin on the server.

Example for Ubuntu:

```bash
sudo apt update
sudo apt install -y ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
sudo usermod -aG docker $USER
```

Reconnect to the server after adding yourself to the `docker` group.

## 6. Upload the Project to the Server

Choose one of these approaches:

### Option A: Git clone on the server

```bash
cd /opt
sudo mkdir -p /opt/launchdate
sudo chown -R $USER:$USER /opt/launchdate
cd /opt/launchdate
git clone <your-repo-url> .
```

### Option B: Upload the code manually

Upload the repository to:

```text
/opt/launchdate
```

All commands below assume the repository root is:

```text
/opt/launchdate
```

## 7. Create Backend Production Environment File

Create a production env file on the server:

```bash
cd /opt/launchdate/services/backend
cp .env.example .env.production
```

Then edit `services/backend/.env.production` with production values:

```env
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
SERVER_ENVIRONMENT=production

MONGODB_URL=mongodb+srv://<db_user>:<db_password>@<cluster-url>/launchdate_db?retryWrites=true&w=majority&appName=launchdate
MONGODB_DATABASE=launchdate_db

LL2_URL_PREFIX=https://lldev.thespacedevs.com
LL2_REQUEST_INTERVAL=5

JWT_SECRET=<generate-with-openssl-rand-base64-32>
ACCESS_TOKEN_EXPIRE_MIN=15
REFRESH_TOKEN_EXPIRE_DAYS=7
JWT_ISSUER=launchdate-backend

IMAGE_S3_ENDPOINT=https://<account-id>.r2.cloudflarestorage.com
IMAGE_S3_REGION=auto
IMAGE_S3_BUCKET=launchdate-images
IMAGE_S3_ACCESS_KEY=<r2-access-key>
IMAGE_S3_SECRET_KEY=<r2-secret-key>
IMAGE_DOMAIN=https://img.launch-date.com
```

Generate a strong JWT secret with:

```bash
openssl rand -base64 32
```

Important:

- The backend code reads `SERVER_ENVIRONMENT`, not `ENVIRONMENT`.
- If you leave `SERVER_ENVIRONMENT` unset, production-only behavior may not work as expected.

## 8. Create Docker Compose for Backend and Tunnel

From the repo root on the server, create a deployment compose file.

Example file path:

```text
/opt/launchdate/docker-compose.prod.yml
```

Suggested content:

```yaml
services:
  backend:
    container_name: launchdate-backend
    build:
      context: ./services/backend
    env_file:
      - ./services/backend/.env.production
    restart: unless-stopped
    networks:
      - launchdate

  cloudflared:
    container_name: launchdate-cloudflared
    image: cloudflare/cloudflared:latest
    command: tunnel --no-autoupdate run --token ${CF_TUNNEL_TOKEN}
    restart: unless-stopped
    depends_on:
      - backend
    networks:
      - launchdate

networks:
  launchdate:
    driver: bridge
```

In the same directory, create a root `.env` file for Compose:

```env
CF_TUNNEL_TOKEN=<cloudflare-tunnel-token>
```

Why this setup is preferable:

- The backend is not exposed directly to the public internet.
- `cloudflared` connects outbound to Cloudflare, so you do not need to open port `80` or `443` for the API.
- The tunnel can reach the backend by Docker service name `backend:8080` inside the same Docker network.

## 9. Create the Cloudflare Tunnel

Use a remotely-managed tunnel in Cloudflare Zero Trust.

1. Open Cloudflare Dashboard -> Zero Trust -> Networks -> Tunnels.
2. Create a tunnel.
3. Choose Docker as the connector type.
4. Copy the generated tunnel token.
5. Put that token into `/opt/launchdate/.env` as `CF_TUNNEL_TOKEN`.
6. In the tunnel public hostname settings, add:

| Hostname | Service type | Service target |
| --- | --- | --- |
| `api.launch-date.com` | HTTP | `http://backend:8080` |

Notes:

- Because `cloudflared` runs in the same Docker network, `backend` resolves to the backend container.
- You do not need to publish backend port `8080` to the host.

## 10. Start Backend and Tunnel

From `/opt/launchdate`:

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

Check running containers:

```bash
docker compose -f docker-compose.prod.yml ps
```

Check backend logs:

```bash
docker compose -f docker-compose.prod.yml logs -f backend
```

Check tunnel logs:

```bash
docker compose -f docker-compose.prod.yml logs -f cloudflared
```

After the tunnel is healthy, verify the API:

```bash
curl https://api.launch-date.com/api/v1/health
```

Expected result should be an HTTP `200` response.

## 11. Deploy `apps/web` to Cloudflare Pages

1. Open Cloudflare Dashboard -> Workers & Pages -> Create application -> Pages.
2. Connect your Git repository.
3. Create a project for `apps/web`.

Recommended build settings:

| Setting | Value |
| --- | --- |
| Project name | `launchdate-web` |
| Production branch | your main branch |
| Framework preset | `Vite` |
| Root directory | `apps/web` |
| Build command | `npm run build` |
| Build output directory | `dist` |

Environment variables for Pages:

```env
VITE_API_BASE_URL=https://api.launch-date.com
NODE_VERSION=20
```

Then bind custom domains:

- `launch-date.com`
- `www.launch-date.com`

## 12. Deploy `apps/admin` to Cloudflare Pages

Create another Pages project for `apps/admin`.

Recommended build settings:

| Setting | Value |
| --- | --- |
| Project name | `launchdate-admin` |
| Production branch | your main branch |
| Framework preset | `Vite` |
| Root directory | `apps/admin` |
| Build command | `npm run build` |
| Build output directory | `dist` |

Environment variables for Pages:

```env
VITE_API_BASE_URL=https://api.launch-date.com
NODE_VERSION=20
```

Recommended custom domain:

- `admin.launch-date.com`

## 13. Production Checks

After all three deployments are online, verify the following:

### API

```bash
curl https://api.launch-date.com/api/v1/health
```

### Web

- Open `https://launch-date.com`
- Confirm page data loads from the API

### Admin

- Open `https://admin.launch-date.com`
- Confirm login works
- Confirm CRUD requests reach `https://api.launch-date.com`

### Images

- Confirm uploaded images are stored in R2
- Confirm public image URLs resolve through `IMAGE_DOMAIN`

## 14. Common Issues

### 1. Frontend can open, but all API requests fail with CORS

Cause:

- Frontend domain does not match the backend whitelist.

Fix:

- Update `services/backend/internal/middleware/cors.go`
- Rebuild and redeploy the backend container

### 2. Admin login fails even though the API is reachable

Check:

- `https://admin.launch-date.com` is in the CORS whitelist
- API is served over HTTPS
- Browser receives the refresh token cookie
- Frontend `VITE_API_BASE_URL` points to the production API domain

### 3. Backend cannot connect to Mongo Atlas

Check:

- Mongo Atlas network access allows your server
- Username and password are correct
- Connection string includes the correct database name

### 4. Image upload fails

Check:

- R2 endpoint is correct
- Bucket name is correct
- Access key and secret key are valid
- `IMAGE_DOMAIN` points to a public URL

### 5. Tunnel is connected, but `api.launch-date.com` still does not work

Check:

- Tunnel public hostname points to `http://backend:8080`
- `cloudflared` and `backend` containers are on the same Docker network
- The backend container is healthy and listening on `8080`

## 15. Recommended Next Improvements

- Add CI to build and redeploy the backend on push to the production branch
- Add server-side automatic restart on reboot through Docker restart policy, which the compose file above already uses
- Add database backup and restore process for Mongo Atlas
- Add Cloudflare WAF and rate limiting rules for `api.launch-date.com`
- Move secrets into a proper secret manager later if your deployment process grows more complex