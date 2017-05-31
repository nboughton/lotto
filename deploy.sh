#!/bin/bash
# Rebuild code
echo "Re-compiling code"
go build -o site.app

# Stop currently running service
echo "Stopping service"
ssh server "/var/www/utils/stop lotto.nboughton.uk"

# Rsync to live
echo "Syncing new release to server"
rsync -aWvL --delete --exclude-from=exclude.rsync . server:/var/www/sites/lotto.nboughton.uk

# Start service
echo "Restarting service"
ssh server "/var/www/utils/start lotto.nboughton.uk"
