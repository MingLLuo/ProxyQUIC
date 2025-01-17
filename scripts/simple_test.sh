#!/usr/bin/env bash

###########################################################################
# simple_test.sh
###########################################################################

# ========== 1. 定义颜色变量 ==========
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ========== x. 启动 Http Proxy ==========

PROXY_ADDRESS="127.0.0.1:8082"
export HTTP_PROXY="http://${PROXY_ADDRESS}"
export HTTPS_PROXY="http://${PROXY_ADDRESS}"
echo -e "${RED}[INFO] Starting proxy (red color) ...${NC}"
go run ./cmd/http-proxy/ -mode=simple 2>&1 | \
  awk -v red="${RED}" -v reset="${NC}" '{print red $0 reset}' &
PROXY_PID=$!
sleep 1

# ========== 2. 启动 Server 并对输出加上颜色标注 ==========

echo -e "${YELLOW}[INFO] Starting server (green color) ...${NC}"
go run ./cmd/server/ -mode=simple 2>&1 | \
  awk -v green="${GREEN}" -v reset="${NC}" '{print green $0 reset}' &
SERVER_PID=$!
sleep 1

# ========== 3. 启动 Client 并对输出加上颜色标注 ==========

echo -e "${YELLOW}[INFO] Starting client (blue color)...${NC}"
go run ./cmd/client/ -mode=simple 2>&1 | \
  awk -v blue="${BLUE}" -v nc="${NC}" '{print blue $0 nc}' &
CLIENT_PID=$!

wait $CLIENT_PID
echo -e "${YELLOW}[INFO] Client done. Now press Ctrl+C to stop server, or wait...${NC}"

# ========== 4. 阻塞等待 Server 结束 ==========

wait $SERVER_PID
echo -e "${YELLOW}[INFO] Server process has ended.${NC}"

echo -e "${YELLOW}[INFO] Now press Ctrl+C to stop proxy, or wait...${NC}"

wait $PROXY_PID
echo -e "${RED}[INFO] Proxy process has ended.${NC}"
unset HTTP_PROXY
unset HTTPS_PROXY
