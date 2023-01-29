#!/bin/sh

ME=$(basename "$0")

if [[ -z "${DIST_PATH}" ]] || [[ -z "${CONF_PATH}" ]]; then
    echo "[ERROR] ${ME}: requires the envs DIST_PATH and CONF_PATH to be set"
    exit 1
fi

if [[ -z "${SRV_BASE_PATH}" ]]; then
    echo "[WARNING] ${ME}: no base path set, using '/'"
    echo "[INFO] ${ME}: change this behaviour with setting the env SRV_BASE_BATH"
fi

STATIC_PATH=/var/www/portfoli.go/public

cp -rp ${STATIC_PATH} ${DIST_PATH}/static
mv ${DIST_PATH}/static/favicon.ico ${DIST_PATH}
portfoli-go -dist -dist.dir ${DIST_PATH} -config.dir ${CONF_PATH} -srv.base ${SRV_BASE_PATH:-"/"}