FROM ubuntu
ENV TZ Asia/Shanghai
COPY conf/conf.json /root/conf.json
COPY logagent /root/logagent
RUN chmod u+x /root/logagent
WORKDIR /root
CMD  ["./logagent","-c","conf.json"]