#!/bin/bash
set -eu

OUTPUT_DIR=${1:-$(dirname "$0")/../public/}

BOOTSTRAP=/tmp/portfoli.go-bootstrap.zip
BOOTSTRAP_ICONS=/tmp/portfoli.go-bootstrap-icons.zip
ANIME_JS=/tmp/portfoli.go-anime.zip

wget -q https://github.com/twbs/bootstrap/releases/download/v5.0.2/bootstrap-5.0.2-dist.zip -O ${BOOTSTRAP}
wget -q https://github.com/twbs/icons/releases/download/v1.10.3/bootstrap-icons-1.10.3.zip -O ${BOOTSTRAP_ICONS}
wget -q https://github.com/juliangarnier/anime/archive/refs/tags/v3.2.1.zip -O ${ANIME_JS}

for ZIP in ${BOOTSTRAP} ${BOOTSTRAP_ICONS} ${ANIME_JS}; do
    unzip ${ZIP} -d ${ZIP%.*}
    mv ${ZIP%.*}/*/* ${ZIP%.*}
done

test -d ${OUTPUT_DIR}/css/fonts || mkdir ${OUTPUT_DIR}/css/fonts
cp ${BOOTSTRAP%.*}/css/bootstrap.min.{css,css.map} ${OUTPUT_DIR}/css/
cp ${BOOTSTRAP%.*}/js/bootstrap.min.{js,js.map} ${OUTPUT_DIR}/js/
cp ${BOOTSTRAP_ICONS%.*}/bootstrap-icons.css ${OUTPUT_DIR}/css/
cp ${BOOTSTRAP_ICONS%.*}/fonts/bootstrap-icons.woff* ${OUTPUT_DIR}/css/fonts/
cp ${ANIME_JS%.*}/lib/anime.min.js ${OUTPUT_DIR}/js/

rm -rf /tmp/portfoli.go-*