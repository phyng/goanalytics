
FROM scratch

COPY server /
COPY 17monipdb.dat /

ENTRYPOINT ["/server"]
