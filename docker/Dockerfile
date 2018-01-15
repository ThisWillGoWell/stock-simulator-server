FROM phusion/baseimage:latest
ADD stock-simulator-server /opt/server/
ADD static/ /opt/server/static/
RUN chmod +x /opt/server/stock-simulator-server
EXPOSE 8000
CMD ["./opt/server/stock-simulator-server"]