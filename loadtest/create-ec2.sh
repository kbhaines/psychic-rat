#!/bin/sh

set -x
iid=`aws ec2 run-instances --image-id ami-1b791862 --security-group-ids sg-88a76cf2 --instance-type t2.nano --key-name macpro --query "Instances[0].InstanceId"`
aws ec2 wait instance-running --instance-ids $iid
host=`aws ec2 describe-instances --instance-ids $iid --query Reservations[0].Instances[*].PublicDnsName`

SSH="ssh -i $HOME/macpro.pem ubuntu@$host"
while ! $SSH sudo curl https://get.docker.com \| sh;do
    sleep 5
done

$SSH sudo adduser ubuntu docker \; sudo apt-get install sqlite3

docker save pr:latest | gzip | ssh -i ~/macpro.pem ubuntu@$host gunzip \| docker load
scp -i ~/macpro.pem pr.dat ubuntu@$host:

container_id=`$SSH docker run -d --name pr -e COOKIE_KEYS=loadtest -v /home/ubuntu/pr.dat:/pr.dat -p 8080:8080 pr`

echo $host with instance ID $iid is ready for load test
echo container on host is $container_id

sed "s/target:/$host:/" targets.tmpl > targets
echo Starting load test
./loadtest.sh 5s

echo Done. Retrieving logs and db files

$SSH docker logs pr &> dockerlogs
$SSH docker exec pr sqlite3 pr.dat .dump \| gzip  | gunzip -c > dbdump

aws ec2 terminate-instances --instance-ids $iid
