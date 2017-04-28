FROM iron/base

ADD reverse-proxy /
ADD config.json /
ADD medlist.tmpl /
ADD test/cert/ test/cert/

ENTRYPOINT [ "./reverse-proxy" ]

EXPOSE 80


