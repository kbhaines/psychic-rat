#!/bin/sh


if [ $# -ne 3 ];then
    echo "Usage: $0 <keyfile> <aws-keyname> <instance-type>"
    exit
fi

KEYFILE=$1
AWSKEY=$2
TYPE=$3

set -x
iid=`aws ec2 run-instances --image-id ami-1b791862 --security-group-ids sg-88a76cf2 --instance-type $TYPE --key-name $AWSKEY --query "Instances[0].InstanceId"`
aws ec2 wait instance-running --instance-ids $iid
host=`aws ec2 describe-instances --instance-ids $iid --query Reservations[0].Instances[*].PublicDnsName`

SSH="ssh -i $KEYFILE ubuntu@$host"
while ! $SSH sudo curl https://get.docker.com \| sh;do
    sleep 5
done

$SSH sudo adduser ubuntu docker \; sudo apt-get install sqlite3

docker save pr:latest | gzip | ssh -i $KEYFILE ubuntu@$host gunzip \| docker load
scp -i $KEYFILE pr.dat ubuntu@$host:

container_id=`$SSH docker run -d --name pr -e COOKIE_KEYS=loadtest -v /home/ubuntu/pr.dat:/pr.dat -p 8080:8080 pr /server.linux -mockcsrf -listen :8080 -cache-templates`

echo $host with instance ID $iid is ready for load test
echo container on host is $container_id

echo Starting load test
./loadtest.sh $host:8080 5s

echo Done. Retrieving logs and db files

$SSH docker logs pr &> dockerlogs
$SSH docker exec pr sqlite3 pr.dat .dump \| gzip  | gunzip -c > dbdump

aws ec2 terminate-instances --instance-ids $iid
