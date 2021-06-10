PROJECT=REPLACE
LOCATION=REPLACE
PREFIX=REPLACE

go run cmd/balancer/main.go \
--project $PROJECT \
--location $LOCATION \
--prefix $PREFIX \
--cluster cluster-1 \
--cluster cluster-2 \
--namespace default \
--hpa workload \
--value 110 \
--idle 30

