#!/bin/bash

OUTPUT_DIR=$(dirname "$0")/../static/
BOOTSTRAP=/tmp/gofolio-bootstrap.zip
BOOTSTRAP_ICONS=/tmp/gofolio-bootstrap-icons.zip
BOOTSTRAP_U=/tmp/gofolio-bootstrap
BOOTSTRAP_ICONS_U=/tmp/gofolio-bootstrap-icons

wget https://github.com/twbs/bootstrap/releases/download/v5.0.2/bootstrap-5.0.2-dist.zip -O ${BOOTSTRAP}
wget https://github.com/twbs/icons/releases/download/v1.10.3/bootstrap-icons-1.10.3.zip -O ${BOOTSTRAP_ICONS}

unzip ${BOOTSTRAP} -d ${BOOTSTRAP_U}
mv ${BOOTSTRAP_U}/*/* ${BOOTSTRAP_U}
unzip ${BOOTSTRAP_ICONS} -d ${BOOTSTRAP_ICONS_U}
mv ${BOOTSTRAP_ICONS_U}/*/* ${BOOTSTRAP_ICONS_U}

cp ${BOOTSTRAP_U}/css/bootstrap.min.css ${OUTPUT_DIR}/css/
cp ${BOOTSTRAP_U}/js/bootstrap.min.js ${OUTPUT_DIR}/js/
cp ${BOOTSTRAP_ICONS_U}/bootstrap-icons.css ${OUTPUT_DIR}/css
test -d ${OUTPUT_DIR}/css/fonts || mkdir ${OUTPUT_DIR}/css/fonts
cp ${BOOTSTRAP_ICONS_U}/fonts/bootstrap-icons.woff* ${OUTPUT_DIR}/css/fonts

rm -rf /tmp/gofolio-bootstrap*