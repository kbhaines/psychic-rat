#/bin/sh
TARGET=$1
DURATION=${2:-5s}

if [ "$TARGET" == "" ];then
    echo "specify target"
    exit 1
fi

sed "s/:target:/$TARGET/" targets.tmpl > targets
$HOME/go/bin/vegeta attack -rate 100 -targets targets -duration=$DURATION | tee results.bin | $HOME/go/bin/vegeta report
