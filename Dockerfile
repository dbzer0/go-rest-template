#
# Этап сборки (builder)
#
FROM --platform=linux/amd64 golang:1.23 AS builder
ENV CGO_ENABLED=0

WORKDIR /go/src/github.com/dbzer0/go-rest-template
COPY . .

# Вычисляем версию; если git не выдаст тег, используем "dev"
RUN version=$(git describe --abbrev=6 --always --tag || echo "dev") && \
    echo "version=$version" && \
    cd app && \
    go build -a -tags PROJECTNAME -installsuffix PROJECTNAME -mod=vendor \
      -ldflags "-X main.version=${version} -s -w" \
      -o /go/bin/PROJECTNAME

#
# Этап базового образа (base) с сертификатами и настройкой пользователя
#
FROM --platform=linux/amd64 alpine AS base
RUN apk --no-cache add ca-certificates && \
    addgroup -S PROJECTNAME && \
    adduser -S PROJECTNAME -G PROJECTNAME

#
# Финальный образ рантайма (final)
#
FROM scratch

# Копируем бинарник
COPY --from=builder /go/bin/PROJECTNAME /bin/PROJECTNAME

# Копируем сертификаты из базового образа
COPY --from=base /etc/ssl/certs /etc/ssl/certs

# Если есть документация, её можно скопировать из базового образа
# COPY --from=base /usr/share/PROJECTNAME /usr/share/PROJECTNAME

# Копируем файлы с информацией о пользователях и группах
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group

USER PROJECTNAME
ENTRYPOINT ["/bin/PROJECTNAME", "run"]
