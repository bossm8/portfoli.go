#!/bin/bash

OUTPUT_DIR=$(dirname "$0")/../static/
BOOTSTRAP=/tmp/portfoli.go-bootstrap.zip
BOOTSTRAP_ICONS=/tmp/portfoli.go-bootstrap-icons.zip
BOOTSTRAP_U=/tmp/portfoli.go-bootstrap
BOOTSTRAP_ICONS_U=/tmp/portfoli.go-bootstrap-icons
ANIME_JS=/tmp/portfoli.go-anime.zip
ANIME_JS_U=/tmp/portfoli.go-anime

wget https://github.com/twbs/bootstrap/releases/download/v5.0.2/bootstrap-5.0.2-dist.zip -O ${BOOTSTRAP}
wget https://github.com/twbs/icons/releases/download/v1.10.3/bootstrap-icons-1.10.3.zip -O ${BOOTSTRAP_ICONS}
wget https://github.com/juliangarnier/anime/archive/refs/tags/v3.2.1.zip -O ${ANIME_JS}

for ZIP in ${BOOTSTRAP} ${BOOTSTRAP_ICONS} ${ANIME_JS}; do
    unzip ${ZIP} -d ${ZIP%.*}
    mv ${ZIP%.*}/*/* ${ZIP%.*}
done

test -d ${OUTPUT_DIR}/css/fonts || mkdir ${OUTPUT_DIR}/css/fonts
cp ${BOOTSTRAP_U}/css/bootstrap.min.css ${OUTPUT_DIR}/css/
cp ${BOOTSTRAP_U}/js/bootstrap.min.js ${OUTPUT_DIR}/js/
cp ${BOOTSTRAP_ICONS_U}/bootstrap-icons.css ${OUTPUT_DIR}/css
cp ${BOOTSTRAP_ICONS_U}/fonts/bootstrap-icons.woff* ${OUTPUT_DIR}/css/fonts

rm -rf /tmp/portfoli.go-*