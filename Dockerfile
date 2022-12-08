FROM sordfish/build-tools:1.19.3-ubuntu as build

WORKDIR /app
COPY ./* /app/
RUN go build -o ion-sfu-gstreamer-receive

FROM sordfish/ubuntu-gstreamer:latest as runtime

COPY --from=build /app/ion-sfu-gstreamer-receive /
CMD ["./ion-sfu-gstreamer-receive"]
