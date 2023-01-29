#!/bin/sh

if [[ -z "${DIST_PATH}" ]] || [[ -z "${CONF_PATH}" ]]; then
    echo "portfoli-go-static.sh requires DIST_PATH and CONF_PATH to be set"
    exit 1
fi

STATIC_PATH=/var/www/portfoli.go/public

cp -rp ${STATIC_PATH} ${DIST_PATH}/static
portfoli-go -dist -dist.dir ${DIST_PATH} -config.dir ${CONF_PATH}