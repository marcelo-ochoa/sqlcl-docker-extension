FROM --platform=$BUILDPLATFORM node:17.7-alpine3.14 AS client-builder
WORKDIR /app/client
# https://www.oracle.com/database/sqldeveloper/technologies/sqlcl/download/
ADD https://download.oracle.com/otn_software/java/sqldeveloper/sqlcl-24.1.0.087.0929.zip .
RUN unzip -d /opt sqlcl-24.1.0.087.0929.zip
# cache packages in layer
COPY client/package.json /app/client/package.json
COPY client/package-lock.json /app/client/package-lock.json
RUN --mount=type=cache,target=/usr/src/app/.npm \
    npm set cache /usr/src/app/.npm && \
    npm ci
# install
COPY client /app/client
RUN npm run build

FROM golang:1.17-alpine AS builder
ENV CGO_ENABLED=0
WORKDIR /backend
COPY vm/go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
COPY vm/. .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o bin/service

FROM ghcr.io/graalvm/graalvm-ce:ol8-java17-22.3.3
RUN set -eux \
    && if [ "$(arch)" == "x86_64" ]; then TTYD_PKG=ttyd.i686; fi \
    && if [ "$(arch)" == "aarch64" ]; then TTYD_PKG=ttyd.aarch64; fi \
    && curl -o /usr/bin/ttyd -L https://github.com/tsl0922/ttyd/releases/download/1.7.4/${TTYD_PKG} \
    && rpm -ivh https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm && microdnf install -y tini unzip ncurses \
    && gu install js && microdnf clean all && chmod +x /usr/bin/ttyd \
    && mkdir -p /home/sqlcl \
    && echo "HOME=/home/sqlcl;cd /home/sqlcl;/opt/sqlcl/bin/sql /nolog" > /home/sql.sh \
    && chown 1000:1000 /home/sqlcl /home/sql.sh \
    && chmod u+rwx /home/sql.sh \
    && echo "sqlcl:x:1000:1000:SQLcl:/home/sqlcl:/bin/bash" >> /etc/passwd \
    && echo "sqlcl:x:1000:sqlcl" >> /etc/group

LABEL org.opencontainers.image.title="Oracle SQLcl client tool"
LABEL org.opencontainers.image.description="Docker Extension for using an embedded version of Oracle SQLcl client tool."
LABEL org.opencontainers.image.vendor="Marcelo Ochoa"
LABEL com.docker.desktop.extension.api.version=">= 0.2.3"
LABEL com.docker.extension.categories="database,utility-tools"
LABEL com.docker.extension.screenshots="[{\"alt\":\"Sample usage using scott user\", \"url\":\"https://raw.githubusercontent.com/marcelo-ochoa/sqlcl-docker-extension/main/docs/images/screenshot2.png\"},\
    {\"alt\":\"Some formating options\", \"url\":\"https://raw.githubusercontent.com/marcelo-ochoa/sqlcl-docker-extension/main/docs/images/screenshot3.png\"},\
    {\"alt\":\"Explain Plan\", \"url\":\"https://raw.githubusercontent.com/marcelo-ochoa/sqlcl-docker-extension/main/docs/images/screenshot4.png\"}]"
LABEL com.docker.extension.publisher-url="https://github.com/marcelo-ochoa/sqlcl-docker-extension"
LABEL com.docker.extension.additional-urls="[{\"title\":\"Documentation\",\"url\":\"https://github.com/marcelo-ochoa/sqlcl-docker-extension/blob/main/README.md\"},\
    {\"title\":\"License\",\"url\":\"https://github.com/marcelo-ochoa/sqlcl-docker-extension/blob/main/LICENSE\"}]"
LABEL com.docker.extension.detailed-description="Docker Extension for using Oracle SQLcl client tool"
LABEL com.docker.extension.changelog="See full <a href=\"https://github.com/marcelo-ochoa/sqlcl-docker-extension/blob/main/CHANGELOG.md\">change log</a>"
LABEL com.docker.desktop.extension.icon="https://raw.githubusercontent.com/marcelo-ochoa/sqlcl-docker-extension/main/client/public/favicon.ico"
LABEL com.docker.extension.detailed-description="Oracle SQLcl (SQL Developer Command Line) is a Java-based command line interface for Oracle Database. \
    Using SQLcl, you can execute SQL and PL/SQL statements in interactive or batch mode. \
    SQLcl provides inline editing, statement completion, command recall, and also supports your existing SQL*Plus scripts."
COPY sqlcl.svg metadata.json docker-compose.yml /

COPY --from=client-builder /app/client/dist /ui
COPY --from=client-builder /opt/sqlcl /opt/sqlcl
COPY --from=builder /backend/bin/service /
COPY --chown=1000:1000 login.sql /home/sqlcl

ENTRYPOINT ["/usr/bin/tini", "--", "/service", "-socket", "/run/guest-services/sqlcl-docker-extension.sock"]