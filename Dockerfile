FROM golang:1.5.0

MAINTAINER Cai Guanhao (caiguanhao@gmail.com)

RUN python2.7 -c 'from urllib import urlopen; from json import loads; \
    print(loads(urlopen("http://ip-api.com/json").read().decode("utf-8" \
    ).strip())["timezone"])' > /etc/timezone && \
    dpkg-reconfigure -f noninteractive tzdata

ADD AliyunReminder /

CMD ["/AliyunReminder"]
