
FROM scratch
MAINTAINER phyng

COPY server /
COPY 17monipdb.dat /

ENTRYPOINT ["/server"]
