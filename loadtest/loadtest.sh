#/bin/sh
DURATION=${1:-5s}
/Users/khaines/go/bin/vegeta attack -targets targets -duration=$DURATION | tee results.bin | /Users/khaines/go/bin/vegeta report
