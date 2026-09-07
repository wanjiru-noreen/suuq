FROM nginx:1.27-alpine

COPY frontend/ /usr/share/nginx/html/
COPY development/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
