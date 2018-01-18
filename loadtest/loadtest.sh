#/bin/sh
DURATION=${1:-5s}
/Users/khaines/go/bin/vegeta attack -rate 100 -targets targets -duration=$DURATION | tee results.bin | /Users/khaines/go/bin/vegeta report
