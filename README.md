# AutoHelm

AutoHelm is a lightweight Docker-based tool that monitors a remote `docker-compose.yml` file and automatically updates your local setup when changes are detected. After updating, it gracefully restarts the relevant containers using `docker compose`.

Ideal for self-hosted environments, including Raspberry Pi and x86 machines.

---

## 🚀 Features

- Polls a remote `docker-compose.yml` file once daily at 3:00 AM.
- Compares it against a local version.
- If different, replaces the local file and restarts containers using `docker compose`.
- Works with multi-arch platforms (x86_64, arm64, Raspberry Pi 3/4).
- Minimal image based on distroless + Docker CLI and Compose plugin.

---

## 🛠️ Usage

### 🐳 Docker Compose

```yaml
version: "3.8"

services:
  autohelm:
    image: your-dockerhub-username/autohelm:latest
    container_name: autohelm
    restart: unless-stopped

    environment:
      - COMPOSE_URL=https://raw.githubusercontent.com/youruser/yourrepo/main/docker-compose.yml
      # Optional: override default local file path
      - LOCAL_FILE=docker-compose.yml
      - TZ=America/Los_Angeles

    volumes:
      - docker-compose.yml:/data/docker-compose.yml
      - /var/run/docker.sock:/var/run/docker.sock
```

---

### 🌍 Environment Variables
| Variable      | Description                                                | Default              |
|---------------|------------------------------------------------------------|----------------------|
| `COMPOSE_URL` | **(Required)** URL to the remote `docker-compose.yml` file | —                    |
| `LOCAL_FILE`  | Path to the local file to compare and replace              | `docker-compose.yml` |
| `TZ`          | Timezone for scheduling daily checks (optional)            | `UTC`                |

---

## 🧰 Future Plans
* Customizable poll schedule
* Dry-run / preview mode
* External trigger to check for updates

---

## 🙏 Credits
Inspired by Watchtower but with Compose yaml file syncing.
