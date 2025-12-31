#!/bin/bash

# ============================================
# Andaria Uploads Backup Script
# ============================================
# This script backs up uploaded files
# and rotates old backups (keeps last 7 days)
# ============================================

set -e  # Exit on error

# Configuration
BACKUP_DIR="./backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="uploads_$DATE.tar.gz"
RETENTION_DAYS=7

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Andaria Uploads Backup Script${NC}"
echo -e "${GREEN}========================================${NC}"

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}Error: Docker is not running${NC}"
    exit 1
fi

# Check if uploads volume exists
if ! docker volume ls | grep -q andaria-backend-uploads; then
    echo -e "${RED}Error: Uploads volume does not exist${NC}"
    echo "Volume 'andaria-backend-uploads' not found"
    exit 1
fi

echo -e "${YELLOW}Starting uploads backup...${NC}"
echo "Backup file: $BACKUP_FILE"

# Perform backup using a temporary Alpine container
if docker run --rm \
    -v andaria-backend-uploads:/data:ro \
    -v "$(pwd)/$BACKUP_DIR:/backup" \
    alpine tar czf "/backup/$BACKUP_FILE" -C /data .; then
    echo -e "${GREEN}Backup completed successfully!${NC}"
    echo "Location: $BACKUP_DIR/$BACKUP_FILE"

    # Get file size
    SIZE=$(du -h "$BACKUP_DIR/$BACKUP_FILE" | cut -f1)
    echo "Size: $SIZE"
else
    echo -e "${RED}Backup failed!${NC}"
    exit 1
fi

# Rotate old backups
echo -e "${YELLOW}Cleaning up old backups (keeping last $RETENTION_DAYS days)...${NC}"
DELETED=$(find "$BACKUP_DIR" -name "uploads_*.tar.gz" -mtime +$RETENTION_DAYS -delete -print | wc -l)

if [ "$DELETED" -gt 0 ]; then
    echo -e "${GREEN}Deleted $DELETED old backup(s)${NC}"
else
    echo "No old backups to delete"
fi

# List current backups
echo ""
echo -e "${GREEN}Current backups:${NC}"
ls -lh "$BACKUP_DIR"/uploads_*.tar.gz 2>/dev/null || echo "No backups found"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Backup completed successfully!${NC}"
echo -e "${GREEN}========================================${NC}"

exit 0
