package main

//go:generate xo mysql://$MYSQL_USER:$MYSQL_PASSWORD@($MYSQL_HOST:$MYSQL_PORT)/database -o ./models --template-path ./template/
