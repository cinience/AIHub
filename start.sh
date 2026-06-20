#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT_DIR"

# Example env overrides (uncomment to use)
# export REQUEST_UPSTREAM_OVERRIDE_ENABLED=true
# export REQUEST_UPSTREAM_OVERRIDE_ALLOWLIST='["*"]'
# export REQUEST_UPSTREAM_PROXY_MAP='{"api.example.com":"http://corp-proxy:8080","*.example.com":"http://corp-proxy:8080",".example.com":"http://corp-proxy:8080"}'

PORT="${PORT:-19000}"
export PORT
export GOWORK=off
DAEMON=false

while getopts ":d" opt; do
  case "${opt}" in
    d)
      DAEMON=true
      ;;
    *)
      ;;
  esac
done

kill_port_process() {
  local port="$1"
  local pid=""

  if command -v lsof >/dev/null 2>&1; then
    pid="$(lsof -ti "tcp:${port}" 2>/dev/null || true)"
  elif command -v ss >/dev/null 2>&1; then
    pid="$(ss -ltnp "sport = :${port}" 2>/dev/null | awk -F'pid=' 'NR>1 {print $2}' | awk -F',' '{print $1}' | head -n 1)"
  fi

  if [ -n "${pid}" ]; then
    echo "[start] port ${port} is in use by pid ${pid}, stopping it..."
    kill "${pid}" 2>/dev/null || true
    sleep 1
    if kill -0 "${pid}" 2>/dev/null; then
      echo "[start] pid ${pid} still running, forcing stop..."
      kill -9 "${pid}" 2>/dev/null || true
    fi
  fi
}

echo "[start] building frontend..."
if [ ! -f "${ROOT_DIR}/web/default/dist/index.html" ] || [ ! -f "${ROOT_DIR}/web/classic/dist/index.html" ]; then
  make build-all-frontends
else
  echo "[start] frontend dist already exists, skipping frontend build"
fi

echo "[start] building backend..."
mkdir -p "${ROOT_DIR}/bin"
go build -o "${ROOT_DIR}/bin/new-api" .

kill_port_process "${PORT}"

echo "[start] starting backend on port ${PORT}..."
if [ "${DAEMON}" = "true" ]; then
  mkdir -p "${ROOT_DIR}/logs"
  nohup "${ROOT_DIR}/bin/new-api" --port "${PORT}" >> "${ROOT_DIR}/logs/server.log" 2>&1 &
  echo "[start] backend started in daemon mode (pid $!)"
else
  exec "${ROOT_DIR}/bin/new-api" --port "${PORT}"
fi
