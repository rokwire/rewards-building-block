FROM golang:1.24.11-alpine AS builder

ENV CGO_ENABLED=0

# Install all required build tools
RUN apk add --no-cache make git bash

WORKDIR /rewards-app

# Copy source
COPY . .

# Build
RUN make


FROM alpine:3.21.3

# timezone database 
RUN apk --no-cache add tzdata

COPY --from=builder /rewards-app/bin/rewards /rewards
COPY --from=builder /rewards-app/docs/swagger.yaml /docs/swagger.yaml
COPY --from=builder /rewards-app/driver/web/authorization_model.conf /driver/web/authorization_model.conf
COPY --from=builder /rewards-app/driver/web/authorization_policy.csv /driver/web/authorization_policy.csv

ENTRYPOINT ["/rewards"]
