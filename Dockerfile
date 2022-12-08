FROM sordfish/build-tools:1.19.3-ubuntu as build

WORKDIR /app
COPY ./* /app/
RUN go build -o ion-sfu-gstreamer-receive

FROM sordfish/ubuntu-gstreamer:latest as runtime

WORKDIR /app
COPY --from=build /app/ion-sfu-gstreamer-receive /app/

CMD ["/app/ion-sfu-gstreamer-receive"]
