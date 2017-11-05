#!/usr/bin/env bash

if [ -z ${TRAVIS_BRANCH+x} ]; then
    GIT_REV=`git rev-parse --verify --short HEAD`
    VERSION="${GIT_REV}-SNAPSHOT"
else
	VERSION=$TRAVIS_BRANCH
fi

echo "Packaging ssh-aliases ${VERSION}"

DIST="dist"
OUT_DIR="${DIST}/out"
OUT_BINARY="${OUT_DIR}/ssh-aliases"

OS_LIST="darwin linux"
ARCH_LIST="386 amd64"
ADD_FILES="LICENSE README.md"

mkdir -p ${OUT_DIR}

for FILE in ${ADD_FILES}; do
    cp ${FILE} ${OUT_DIR}
done

for OS in ${OS_LIST}; do
    for ARCH in ${ARCH_LIST}; do
        echo "Making binary for ${OS}/${ARCH}"

        env GOOS=${OS} GOARCH=${ARCH} go build -o ${OUT_BINARY} -ldflags "-X main.VERSION=${VERSION}"
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
