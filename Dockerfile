FROM zephinzer/golang-dev:latest as development
COPY . /go/src/app
ENTRYPOINT [ "start" ]

FROM development as build
RUN build
ENTRYPOINT [ "/go/src/app/app" ]

FROM scratch as production
COPY --from=build /go/src/app/app /
COPY /html /html/
COPY /cert /cert/
EXPOSE 80
EXPOSE 443
ENTRYPOINT [ "/app" ]
