FROM phusion/baseimage:latest
ADD stock-simulator-server /opt/server/
RUN chmod +x /opt/server/stock-simulator-server
EXPOSE 8000
CMD ["./opt/server/stock-simulator-server"]