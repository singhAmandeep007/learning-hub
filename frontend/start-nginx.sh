#!/bin/sh

# Substitute environment variables in nginx config template
envsubst '${BACKEND_HOST},${BACKEND_PORT}' < /etc/nginx/conf.d/nginx.conf.template > /etc/nginx/conf.d/nginx.conf

# Remove the template file
rm /etc/nginx/conf.d/nginx.conf.template

# Start nginx in foreground
exec nginx -g "daemon off;"
