#!/bin/bash
srcdir=$(readlink -f $(dirname 0))
utildir="/var/www/utils"
host="server"
site="lotto.nboughton.uk"

## Rebuild existing code
echo "Rebuilding all code"
cd $srcdir
go build -o site.app
cd $srcdir/public/frontend
npm run build

## Stop service
echo "Stopping Service"
ssh ${host} "$utildir/stop $site"

## Upload new code
echo "Rsyncing"
rsync -aWvL --delete --exclude-from=${srcdir}/exclude.rsync ${srcdir}/ ${host}:/var/www/sites/${site}

## Start service
echo "Restarting service"
ssh ${host} "$utildir/start $site"
