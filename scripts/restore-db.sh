#!/bin/bash

# ============================================
# Andaria Database Restore Script
# ============================================
# This script restores a PostgreSQL database
# from a compressed backup file
# ============================================

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Andaria Database Restore Script${NC}"
echo -e "${GREEN}========================================${NC}"

# Check if backup file is provided
if [ -z "$1" ]; then
    echo -e "${RED}Error: No backup file specified${NC}"
    echo ""
    echo "Usage: $0 <backup_file>"
    echo "Example: $0 backups/andaria_backup_20250131_020000.sql.gz"
    echo ""
    echo "Available backups:"
    ls -lh ./backups/andaria_backup_*.sql.gz 2>/dev/null || echo "No backups found"
    exit 1
fi

BACKUP_FILE="$1"

# Check if backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
    echo -e "${RED}Error: Backup file not found: $BACKUP_FILE${NC}"
    exit 1
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}Error: Docker is not running${NC}"
    exit 1
fi

# Check if postgres container is running
if ! docker ps | grep -q andaria-postgres; then
    echo -e "${RED}Error: PostgreSQL container is not running${NC}"
    echo "Please start the services with: docker-compose up -d"
    exit 1
fi

echo ""
echo -e "${YELLOW}WARNING: This will REPLACE all data in the database!${NC}"
echo "Backup file: $BACKUP_FILE"
echo ""
read -p "Are you sure you want to continue? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo -e "${YELLOW}Restore cancelled${NC}"
    exit 0
fi

echo ""
echo -e "${YELLOW}Creating safety backup of current database...${NC}"
SAFETY_BACKUP="./backups/safety_backup_$(date +%Y%m%d_%H%M%S).sql.gz"
if docker-compose exec -T postgres pg_dump -U postgres Andaria_01 | gzip > "$SAFETY_BACKUP"; then
    echo -e "${GREEN}Safety backup created: $SAFETY_BACKUP${NC}"
else
    echo -e "${RED}Failed to create safety backup!${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}Restoring database from backup...${NC}"

# Drop existing connections
echo "Terminating existing connections..."
docker-compose exec -T postgres psql -U postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'Andaria_01' AND pid <> pg_backend_pid();" > /dev/null 2>&1 || true

# Drop and recreate database
echo "Dropping existing database..."
docker-compose exec -T postgres psql -U postgres -c "DROP DATABASE IF EXISTS Andaria_01;" || true

echo "Creating new database..."
docker-compose exec -T postgres psql -U postgres -c "CREATE DATABASE Andaria_01 WITH ENCODING='UTF8';"

# Restore from backup
echo "Restoring data..."
if gunzip -c "$BACKUP_FILE" | docker-compose exec -T postgres psql -U postgres -d Andaria_01 > /dev/null 2>&1; then
    echo -e "${GREEN}Database restored successfully!${NC}"
else
    echo -e "${RED}Restore failed!${NC}"
    echo ""
    echo -e "${YELLOW}Attempting to restore from safety backup...${NC}"
    if gunzip -c "$SAFETY_BACKUP" | docker-compose exec -T postgres psql -U postgres -d Andaria_01; then
        echo -e "${GREEN}Recovered from safety backup${NC}"
    else
        echo -e "${RED}Failed to recover from safety backup!${NC}"
        echo "Manual intervention required"
        exit 1
    fi
    exit 1
fi

# Verify restore
echo ""
echo -e "${YELLOW}Verifying restore...${NC}"
TABLE_COUNT=$(docker-compose exec -T postgres psql -U postgres -d Andaria_01 -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" | tr -d ' ')
echo "Tables found: $TABLE_COUNT"

if [ "$TABLE_COUNT" -gt 0 ]; then
    echo -e "${GREEN}Verification successful!${NC}"
else
    echo -e "${RED}Warning: No tables found after restore${NC}"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Restore completed successfully!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Safety backup is available at: $SAFETY_BACKUP"
echo "You can delete it manually if the restore is successful"

exit 0
