#!/bin/sh

echo "running nginx"
nginx -g "daemon off;" &
NGINX_PID="$!"


echo "running backend"
/authz &
AUTHZ_PID="$!"

trap "kill $NGINX_PIDC $AUTHZ_PID" exit INT TERM

wait
