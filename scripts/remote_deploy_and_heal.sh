#!/usr/bin/env bash
set -euo pipefail

APP_DIR="${APP_DIR:-/root/qa_test_server}"
ARCHIVE_PATH="${ARCHIVE_PATH:-/root/qa_test_server_deploy.tar.gz}"
SERVICE_NAME="${SERVICE_NAME:-qa_test_server}"
HTTP_PORT="${HTTP_PORT:-8080}"
TCP_PORT="${TCP_PORT:-4001}"
PROXY_FROM_PORT="${PROXY_FROM_PORT:-4003}"
PROXY_TO_PORT="${PROXY_TO_PORT:-7777}"
DB_NAME="${DB_NAME:-go_test}"
DB_USER="${DB_USER:-qa_user}"
DB_PASS="${DB_PASS:-L1nFen9.com}"
AUTH_SECRET="${AUTH_SECRET:-qa_test_server_prod_secret_change_me}"
ADMIN_USER="${ADMIN_USER:-admin}"
ADMIN_PASS="${ADMIN_PASS:-Admin@123456}"
ENABLE_NGINX_PROXY="${ENABLE_NGINX_PROXY:-0}"
HEAL_MODE="${HEAL_MODE:-deploy}" # deploy|heal|verify
MAX_APT_RETRY="${MAX_APT_RETRY:-3}"

SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
LAST_BACKUP_DIR=""

log() {
  printf '[%s] %s\n' "$(date '+%Y-%m-%d %H:%M:%S')" "$*"
}

retry() {
  local max_attempts="$1"
  shift
  local attempt=1
  while true; do
    if "$@"; then
      return 0
    fi
    if (( attempt >= max_attempts )); then
      return 1
    fi
    log "command failed, retry ${attempt}/${max_attempts}: $*"
    sleep $(( attempt * 2 ))
    attempt=$(( attempt + 1 ))
  done
}

apt_install() {
  local pkgs=("$@")
  if [[ "${#pkgs[@]}" -eq 0 ]]; then
    return 0
  fi
  retry "${MAX_APT_RETRY}" apt-get update -y
  retry "${MAX_APT_RETRY}" apt-get install -y --no-install-recommends "${pkgs[@]}"
}

ensure_base_tools() {
  export DEBIAN_FRONTEND=noninteractive
  apt_install curl ca-certificates gnupg lsb-release tar sed
}

ensure_node() {
  if command -v node >/dev/null 2>&1 && command -v npm >/dev/null 2>&1; then
    return 0
  fi
  log "installing Node.js"
  curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
  apt_install nodejs
}

ensure_go() {
  if command -v go >/dev/null 2>&1; then
    return 0
  fi
  log "installing Go"
  apt_install golang-go
}

ensure_mariadb() {
  if ! command -v mysql >/dev/null 2>&1; then
    log "installing MariaDB"
    apt_install mariadb-server
  fi
  systemctl enable --now mariadb
}

lockfile_major_from_file() {
  local lock_file="$1"
  if [[ ! -f "$lock_file" ]]; then
    echo "8"
    return
  fi
  local ver
  ver="$(awk '/^lockfileVersion:/{print $2; exit}' "$lock_file" | tr -d "\"'")"
  case "$ver" in
    6|6.*) echo "8" ;;
    *) echo "10" ;;
  esac
}

ensure_pnpm() {
  local desired_major="$1"
  local current_major=""

  if command -v pnpm >/dev/null 2>&1; then
    current_major="$(pnpm -v | cut -d. -f1 || true)"
  fi

  if [[ -z "$current_major" || "$current_major" != "$desired_major" ]]; then
    log "installing pnpm@${desired_major}"
    npm install -g "pnpm@${desired_major}"
  fi
}

normalize_shell_scripts() {
  if [[ -d "$APP_DIR" ]]; then
    find "$APP_DIR" -maxdepth 1 -type f -name "*.sh" -print0 | while IFS= read -r -d '' f; do
      sed -i 's/\r$//' "$f"
      chmod +x "$f" || true
    done
  fi
}

build_frontend() {
  local web_dir="${APP_DIR}/qa-test-web"
  if [[ ! -d "$web_dir" ]]; then
    log "frontend directory not found: $web_dir"
    return 1
  fi

  local desired_pnpm_major
  desired_pnpm_major="$(lockfile_major_from_file "${web_dir}/pnpm-lock.yaml")"
  ensure_pnpm "$desired_pnpm_major"

  cd "$web_dir"
  export npm_config_registry="${npm_config_registry:-https://registry.npmmirror.com}"
  # Always include devDependencies because the frontend build needs vue-tsc/vite.
  rm -rf node_modules
  if ! pnpm install --prod=false --frozen-lockfile; then
    log "frozen lockfile install failed, retrying without frozen lockfile"
    pnpm install --prod=false --no-frozen-lockfile
  fi
  ./node_modules/.bin/vue-tsc --noEmit
  ./node_modules/.bin/vite build

  # Sync Vite output to Gin templates dir used by backend static serving.
  mkdir -p "${APP_DIR}/templates"
  rm -rf "${APP_DIR}/templates/assets"
  cp -a "${web_dir}/dist/." "${APP_DIR}/templates/"
}

build_backend() {
  cd "$APP_DIR"
  GOPROXY="${GOPROXY:-https://goproxy.cn,direct}" go build -o "${SERVICE_NAME}" .
}

restore_templates_from_backup() {
  local source_backup=""
  if [[ -n "$LAST_BACKUP_DIR" && -d "$LAST_BACKUP_DIR/templates" ]]; then
    source_backup="$LAST_BACKUP_DIR"
  else
    source_backup="$(ls -dt "${APP_DIR}"_backup_* 2>/dev/null | head -n 1 || true)"
    if [[ -z "$source_backup" || ! -d "$source_backup/templates" ]]; then
      return 1
    fi
  fi

  mkdir -p "${APP_DIR}/templates"
  rm -rf "${APP_DIR}/templates/assets"
  cp -a "${source_backup}/templates/." "${APP_DIR}/templates/"
  log "reused templates from backup: ${source_backup}/templates"
  return 0
}

deploy_app() {
  if [[ ! -f "$ARCHIVE_PATH" ]]; then
    log "archive not found: $ARCHIVE_PATH"
    return 1
  fi

  local ts backup_dir parent_dir project_name
  ts="$(date '+%Y%m%d_%H%M%S')"
  backup_dir="${APP_DIR}_backup_${ts}"
  parent_dir="$(dirname "$APP_DIR")"
  project_name="$(basename "$APP_DIR")"

  if [[ -d "$APP_DIR" ]]; then
    mv "$APP_DIR" "$backup_dir"
    LAST_BACKUP_DIR="$backup_dir"
    log "backup created: $backup_dir"
  fi

  mkdir -p "$parent_dir"
  tar -xzf "$ARCHIVE_PATH" -C "$parent_dir"

  if [[ ! -d "$APP_DIR" ]]; then
    log "app directory missing after extract: $APP_DIR"
    return 1
  fi

  normalize_shell_scripts
  if ! build_frontend; then
    log "frontend build failed, trying backup templates fallback"
    if ! restore_templates_from_backup; then
      log "frontend fallback failed: no usable backup templates"
      return 1
    fi
  fi
  build_backend
}

ensure_database() {
  ensure_mariadb
  mysql -uroot <<SQL
CREATE DATABASE IF NOT EXISTS \`${DB_NAME}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
CREATE USER IF NOT EXISTS '${DB_USER}'@'127.0.0.1' IDENTIFIED BY '${DB_PASS}';
ALTER USER '${DB_USER}'@'127.0.0.1' IDENTIFIED BY '${DB_PASS}';
GRANT ALL PRIVILEGES ON \`${DB_NAME}\`.* TO '${DB_USER}'@'127.0.0.1';
FLUSH PRIVILEGES;
SQL
}

write_service_file() {
  cat >"$SERVICE_FILE" <<SERVICE
[Unit]
Description=QA Test Server
After=network.target mariadb.service
Requires=mariadb.service

[Service]
Type=simple
WorkingDirectory=${APP_DIR}
ExecStart=${APP_DIR}/${SERVICE_NAME}
Restart=always
RestartSec=3
Environment=QA_HTTP_ADDR=:${HTTP_PORT}
Environment=QA_TCP_ADDR=:${TCP_PORT}
Environment=QA_PROXY_FROM_PORT=${PROXY_FROM_PORT}
Environment=QA_PROXY_TO_PORT=${PROXY_TO_PORT}
Environment=QA_DB_DSN=${DB_USER}:${DB_PASS}@tcp(127.0.0.1:3306)/${DB_NAME}?charset=utf8mb4&parseTime=True&loc=Local
Environment=QA_DB_AUTOMIGRATE=true
Environment=QA_AUTH_SECRET=${AUTH_SECRET}
Environment=QA_AUTH_TOKEN_TTL_HOURS=24
Environment=QA_DEFAULT_ADMIN_USERNAME=${ADMIN_USER}
Environment=QA_DEFAULT_ADMIN_PASSWORD=${ADMIN_PASS}

[Install]
WantedBy=multi-user.target
SERVICE
}

ensure_service_running() {
  if [[ ! -x "${APP_DIR}/${SERVICE_NAME}" ]]; then
    log "binary missing, rebuilding backend"
    build_backend
  fi

  write_service_file
  systemctl daemon-reload
  systemctl enable --now "${SERVICE_NAME}"
  systemctl restart "${SERVICE_NAME}"
}

configure_nginx_proxy() {
  if [[ "$ENABLE_NGINX_PROXY" != "1" ]]; then
    return 0
  fi

  if ss -lnt | awk '{print $4}' | grep -Eq '(^|:|])80$'; then
    if ! systemctl is-active --quiet nginx; then
      log "port 80 already in use by another process, skip nginx proxy"
      return 0
    fi
  fi

  apt_install nginx
  mkdir -p /etc/nginx/sites-available /etc/nginx/sites-enabled
  cat >/etc/nginx/sites-available/${SERVICE_NAME} <<NGINX
server {
    listen 80 default_server;
    server_name _;

    location / {
        proxy_pass http://127.0.0.1:${HTTP_PORT};
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
NGINX

  rm -f /etc/nginx/sites-enabled/default
  ln -sf /etc/nginx/sites-available/${SERVICE_NAME} /etc/nginx/sites-enabled/${SERVICE_NAME}
  nginx -t
  systemctl enable --now nginx
  systemctl restart nginx

  if command -v ufw >/dev/null 2>&1; then
    if ufw status | grep -q "Status: active"; then
      ufw allow 80/tcp || true
      ufw allow "${HTTP_PORT}"/tcp || true
    fi
  fi
}

verify_local() {
  systemctl is-active --quiet mariadb
  systemctl is-active --quiet "${SERVICE_NAME}"

  local i
  for i in {1..20}; do
    if ss -lnt | awk '{print $4}' | grep -Eq "(:|])${HTTP_PORT}$"; then
      break
    fi
    sleep 1
  done
  if ! ss -lnt | awk '{print $4}' | grep -Eq "(:|])${HTTP_PORT}$"; then
    log "service is not listening on ${HTTP_PORT}"
    systemctl --no-pager --full status "${SERVICE_NAME}" | sed -n '1,30p' || true
    return 1
  fi

  local root_code
  root_code="000"
  for i in {1..15}; do
    root_code="$(curl -s -o /dev/null -w '%{http_code}' --max-time 8 "http://127.0.0.1:${HTTP_PORT}/" || true)"
    if [[ "$root_code" == "200" ]]; then
      break
    fi
    sleep 1
  done
  if [[ "$root_code" != "200" ]]; then
    log "root endpoint status: ${root_code}"
    return 1
  fi

  local login_body
  login_body=""
  for i in {1..10}; do
    login_body="$(curl -sS --max-time 8 -H 'Content-Type: application/json' -X POST "http://127.0.0.1:${HTTP_PORT}/auth/login" -d "{\"username\":\"${ADMIN_USER}\",\"password\":\"${ADMIN_PASS}\"}" || true)"
    if [[ "$login_body" == *"\"success\":true"* ]]; then
      break
    fi
    sleep 1
  done
  if [[ "$login_body" != *"\"success\":true"* ]]; then
    log "login check failed: ${login_body}"
    return 1
  fi

  log "local verify passed"
  return 0
}

heal_once() {
  log "healing common failure paths"
  systemctl enable --now mariadb || true

  if [[ -d "$APP_DIR" && ! -x "${APP_DIR}/${SERVICE_NAME}" ]]; then
    build_backend || true
  fi

  if [[ -f "$SERVICE_FILE" ]]; then
    systemctl daemon-reload
    systemctl restart "${SERVICE_NAME}" || true
  fi

  if ! systemctl is-active --quiet "${SERVICE_NAME}"; then
    ensure_service_running || true
  fi

  configure_nginx_proxy || true
}

main() {
  log "mode=${HEAL_MODE}, app=${APP_DIR}, service=${SERVICE_NAME}"
  ensure_base_tools

  case "$HEAL_MODE" in
    deploy)
      ensure_node
      ensure_go
      ensure_database
      deploy_app
      ensure_service_running
      configure_nginx_proxy
      verify_local
      ;;
    heal)
      ensure_node
      ensure_go
      ensure_database
      heal_once
      verify_local
      ;;
    verify)
      verify_local
      ;;
    *)
      log "unknown mode: $HEAL_MODE"
      return 1
      ;;
  esac

  log "__REMOTE_OK__"
}

main "$@"
