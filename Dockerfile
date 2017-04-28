FROM iron/base

ADD reverse-proxy /
ADD config.json /
ADD medlist.tmpl /
ADD certs/ certs/

ENTRYPOINT [ "./reverse-proxy" ]

EXPOSE 80


