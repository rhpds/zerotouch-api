FROM registry.access.redhat.com/ubi9/go-toolset:latest AS builder
WORKDIR /zerotouch-api/

USER root
RUN chown -R ${USER_UID}:0 /zerotouch-api
USER ${USER_UID}

COPY ./ ./
RUN go build -o /zerotouch-api/zerotouch-api ./cmd

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest AS deploy
WORKDIR /zerotouch-api/
USER ${USER_UID}
COPY --from=builder /zerotouch-api/zerotouch-api ./
CMD ["./zerotouch-api"]

ENV DESCRIPTION="Zero Touch API"
LABEL name="rhpds/zerotouch-api" \
      summary="$DESCRIPTION" \
      description="$DESCRIPTION" \
      maintainer="Red Hat Demo Platform"