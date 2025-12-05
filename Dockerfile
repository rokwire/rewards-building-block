FROM golang:1.24.11-alpine AS builder

ENV CGO_ENABLED=0

RUN mkdir /rewards-app
WORKDIR /rewards-app
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
RUN make

FROM alpine:3.17

#we need timezone database
RUN apk --no-cache add tzdata

COPY --from=builder /rewards-app/bin/rewards /
COPY --from=builder /rewards-app/docs/swagger.yaml /docs/swagger.yaml

COPY --from=builder /rewards-app/driver/web/authorization_model.conf /driver/web/authorization_model.conf
COPY --from=builder /rewards-app/driver/web/authorization_policy.csv /driver/web/authorization_policy.csv

#we need timezone database
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo 

ENTRYPOINT ["/rewards"]
