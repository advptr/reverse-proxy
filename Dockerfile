FROM iron/base

ADD reverse-proxy /
ADD config.json /
ADD medlist.tmpl /
ADD certs/ certs/
ADD schema/ schema/

ENTRYPOINT [ "./reverse-proxy" ]

EXPOSE 80


