#/bin/sh
TARGET=$1
DURATION=${2:-5s}

RATE=50
USER=user041

if [ "$TARGET" == "" ];then
    echo "specify target"
    exit 1
fi

# Generate cookie
curl http://$USER:@$TARGET/callback\?p=basic -c cookie > /dev/null
curl http://$USER:@$TARGET/pledge -b cookie -c cookie > /dev/null

COOKIE=`awk '/localhost/{ print $7}' cookie`

sed "s/:target:/$TARGET/;s/:cookie:/$COOKIE/" targets.tmpl > targets
$HOME/go/bin/vegeta attack -rate $RATE -targets targets -duration=$DURATION | tee results.bin | $HOME/go/bin/vegeta report
