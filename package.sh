#!/usr/bin/env bash

if [ $# -eq 0 ]; then
    echo "No release version specified."
    exit 1
fi

VERSION=$1
echo "Packaging ssh-aliases ${VERSION}"

DIST="dist"
OUT_DIR="${DIST}/out"
OUT_BINARY="${OUT_DIR}/ssh-aliases"

OS_LIST="darwin linux"
ARCH_LIST="amd64 arm64"
ADD_FILES="LICENSE README.md"

mkdir -p ${OUT_DIR}

for FILE in ${ADD_FILES}; do
    cp ${FILE} ${OUT_DIR}
done

for OS in ${OS_LIST}; do
    for ARCH in ${ARCH_LIST}; do
        echo "Making binary for ${OS}/${ARCH}"
        env CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -a -o ${OUT_BINARY} -ldflags "-s -w -X main.Version=${VERSION}"

        if [ ${OS} == "darwin" ]; then
            zip -rj "${DIST}/ssh-aliases_${VERSION}_${OS}_${ARCH}.zip" ${OUT_DIR}
        else
            tar -C ${OUT_DIR} -czf "${DIST}/ssh-aliases_${VERSION}_${OS}_${ARCH}.tar.gz" .
        fi

        rm ${OUT_BINARY}
    done
done

for FILE in ${ADD_FILES}; do
    rm "${OUT_DIR}/${FILE}"
done
rmdir ${OUT_DIR}

ls -lh ${DIST}
