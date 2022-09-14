FROM sordfish/build-tools:v1.17 as build

WORKDIR /app
COPY ./* /app/
RUN go build -o ion-sfu-gstreamer-receive


FROM sordfish/ubuntu-gstreamer:latest as runtime

COPY --from=build /app/ion-sfu-gstreamer-receive /
CMD ["./ion-sfu-gstreamer-receive"]