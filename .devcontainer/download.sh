#!/bin/bash
set -eu

OUTPUT_DIR=${1:-$(dirname "$0")/../public/}

BOOTSTRAP_ICONS=/tmp/portfoli.go-bootstrap-icons.zip
ANIME_JS=/tmp/portfoli.go-anime.zip
INTER=/tmp/portfoli.go-inter.zip

wget -q https://github.com/twbs/icons/releases/download/v1.10.3/bootstrap-icons-1.10.3.zip -O ${BOOTSTRAP_ICONS}
wget -q https://github.com/juliangarnier/anime/archive/refs/tags/v3.2.1.zip -O ${ANIME_JS}
wget -q https://github.com/rsms/inter/releases/download/v4.1/Inter-4.1.zip -O ${INTER}

for ZIP in ${BOOTSTRAP_ICONS} ${ANIME_JS}; do
    unzip -q ${ZIP} -d ${ZIP%.*}
    mv ${ZIP%.*}/*/* ${ZIP%.*}
done
unzip -q ${INTER} -d ${INTER%.*}

test -d ${OUTPUT_DIR}/css/fonts || mkdir ${OUTPUT_DIR}/css/fonts
cp ${BOOTSTRAP_ICONS%.*}/bootstrap-icons.css ${OUTPUT_DIR}/css/
cp ${BOOTSTRAP_ICONS%.*}/fonts/bootstrap-icons.woff* ${OUTPUT_DIR}/css/fonts/
cp ${ANIME_JS%.*}/lib/anime.min.js ${OUTPUT_DIR}/js/
# The variable font ships nested under a versioned subfolder whose exact path
# has changed across releases, so locate it instead of hardcoding the path.
cp "$(find ${INTER%.*} -iname 'InterVariable.woff2' | head -n1)" ${OUTPUT_DIR}/css/fonts/

rm -rf /tmp/portfoli.go-*