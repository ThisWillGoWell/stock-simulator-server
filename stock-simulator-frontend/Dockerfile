
# base image
FROM nginx:1.13.9-alpine

# expose port 80
EXPOSE 80

COPY . /usr/share/nginx/html/
# run nginx
CMD ["nginx", "-g", "daemon off;"]