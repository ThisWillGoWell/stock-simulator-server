FROM phusion/baseimage:latest
ADD stock-server-linux /opt/server/
RUN chmod +x /opt/server/stock-server-linux
EXPOSE 8000
CMD ["./opt/server/stock-server-linux"]