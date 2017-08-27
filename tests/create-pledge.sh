#!/bin/sh
set -e

TEST_URL=http://localhost:8080/api/v1
CURL="curl -s "

ITEM_API=$TEST_URL/item
PLEDGE_API=$TEST_URL/pledge

itemId=`$CURL $ITEM_API?company=1 | jq -r .items[0].id`

pledgePost=`printf '{"itemId":"%s"}' $itemId`

set `$CURL -XPOST -d $pledgePost $PLEDGE_API | jq -r .pledges[length-1].id,.pledges[length-1].item.id`

[ "$itemId" = "$2" ] || (echo "expected item $itemId but got item $2 in pledge $1";exit 1)
