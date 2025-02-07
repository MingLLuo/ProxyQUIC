#!/usr/bin/env bash

###########################################################################
# generate_cert.sh
###########################################################################

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

if [[ "$1" == "help" || "$1" == "--help" ]]; then
    echo -e "${YELLOW}[INFO] Displaying help information...${NC}"
    go run ./cmd/cert/ -h
    exit 0
fi

default_host="localhost,127.0.0.1"
default_start_date=""
default_duration="8760h"         # 365 * 24 hours
default_ca="false"
default_rsa_bits="2048"
default_ecdsa_curve="P256"
default_ed25519="false"
default_cert="cert.pem"
default_key="key.pem"

echo -e "${YELLOW}[INFO] Certificate Generation Script${NC}"

read -p "Do you want to provide custom parameters? (y/N): " custom_params

if [[ "$custom_params" =~ ^[Yy] ]]; then
    read -p "Enter host (default: ${default_host}): " input_host
    host="${input_host:-$default_host}"

    read -p "Enter start date (default: ${default_start_date}): " input_start_date
    start_date="${input_start_date:-$default_start_date}"

    read -p "Enter duration (default: ${default_duration}): " input_duration
    duration="${input_duration:-$default_duration}"

    read -p "Generate CA certificate? (true/false, default: ${default_ca}): " input_ca
    ca="${input_ca:-$default_ca}"

    read -p "Enter RSA bits (default: ${default_rsa_bits}): " input_rsa_bits
    rsa_bits="${input_rsa_bits:-$default_rsa_bits}"

    read -p "Enter ECDSA curve (default: ${default_ecdsa_curve}): " input_ecdsa_curve
    ecdsa_curve="${input_ecdsa_curve:-$default_ecdsa_curve}"

    read -p "Generate Ed25519 key? (true/false, default: ${default_ed25519}): " input_ed25519
    ed25519="${input_ed25519:-$default_ed25519}"

    read -p "Enter certificate output path (default: ${default_cert}): " input_cert
    cert="${input_cert:-$default_cert}"

    read -p "Enter key output path (default: ${default_key}): " input_key
    key="${input_key:-$default_key}"
else
    host="$default_host"
    start_date="$default_start_date"
    duration="$default_duration"
    ca="$default_ca"
    rsa_bits="$default_rsa_bits"
    ecdsa_curve="$default_ecdsa_curve"
    ed25519="$default_ed25519"
    cert="$default_cert"
    key="$default_key"
fi

echo -e "${BLUE}[INFO] Using the following parameters:${NC}"
echo -e "${BLUE}host:         ${host}${NC}"
echo -e "${BLUE}start-date:   ${start_date}${NC}"
echo -e "${BLUE}duration:     ${duration}${NC}"
echo -e "${BLUE}ca:           ${ca}${NC}"
echo -e "${BLUE}rsa-bits:     ${rsa_bits}${NC}"
echo -e "${BLUE}ecdsa-curve:  ${ecdsa_curve}${NC}"
echo -e "${BLUE}ed25519:      ${ed25519}${NC}"
echo -e "${BLUE}cert:         ${cert}${NC}"
echo -e "${BLUE}key:          ${key}${NC}"

echo -e "${YELLOW}[INFO] Starting certificate generation...${NC}"
go run ./cmd/cert/ \
  -host="$host" \
  -start-date="$start_date" \
  -duration="$duration" \
  -ca="$ca" \
  -rsa-bits="$rsa_bits" \
  -ecdsa-curve="$ecdsa_curve" \
  -ed25519="$ed25519" \
  -cert="$cert" \
  -key="$key" 2>&1 | \
  awk -v green="${GREEN}" -v reset="${NC}" '{print green $0 reset}'
