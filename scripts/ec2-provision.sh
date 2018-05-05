#!/bin/sh

if [ $# -ne 3 ];then
    echo "Usage: $0 <keyfile> <aws-keyname> <instance-type>"
    exit
fi

KEYFILE=$1
AWSKEY=$2
TYPE=$3
IMAGE_ID=ami-f4f21593

set -x
set -e

start_instance() {
    iid=`aws ec2 run-instances --image-id $IMAGE_ID --security-groups webserver --iam-instance-profile Name=InstanceIAM --instance-type $TYPE --key-name $AWSKEY --query "Instances[0].InstanceId"`
    aws ec2 wait instance-running --instance-ids $iid
    HOST=`aws ec2 describe-instances --instance-ids $iid --query Reservations[0].Instances[*].PublicDnsName`
    SSH="ssh -i $KEYFILE ubuntu@$HOST"
    while ! $SSH sudo curl https://get.docker.com \| sh;do
        sleep 5
    done
    $SSH sudo adduser ubuntu docker \; sudo apt-get -y install sqlite3 awscli
}

start_instance

SSH="ssh -i $KEYFILE ubuntu@$HOST"

#docker save pr:latest | gzip | $SSH gunzip \| docker load

cat <<-'EOF' > container_init
#!/bin/sh

if [ $# -ne 1 ];then
    echo "Usage: $0 <domain>"
    exit 0
fi

DOMAIN=$1

rm -rf ./res

aws s3 --region eu-west-2 cp s3://psychic-images/pr.img.gz - | gunzip | docker load
docker create --name res pr 
docker cp res:/res ./res 
docker rm res

container_id=`docker run -d --name pr --env-file ./pr.env -v /home/ubuntu/cert:/cert -v /home/ubuntu/res:/res -v /home/ubuntu/pr.dat:/pr.dat -p 80:8080 -p 443:4443 pr /server.linux -listen :8080  -listenSSL :4443 -cert /cert/live/$DOMAIN/fullchain.pem -key /cert/live/$DOMAIN/privkey.pem`

docker run --rm -v /home/ubuntu/cert:/etc/letsencrypt -v $PWD/res:/res certbot/certbot:latest certonly --staging --agree-tos --email kevin@rat.me.uk -n --webroot -w /res/ -d $DOMAIN

docker restart pr
EOF

chmod +x container_init
scp -i $KEYFILE pr.dat pr.env container_init ubuntu@$HOST:

echo press enter when DNS is ready to serve from $HOST
nslookup $HOST
read

$SSH ./container_init www.rat.me.uk
