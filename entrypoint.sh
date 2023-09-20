#!/bin/sh

/soft/rooms_robot \
    --proxyUrl="$PROXY_URL" \
    --notice="$NOTICE_CHANNEL" \
    --noticeUrl="$NOTICE_CHANNEL_URL" \
    --noticeKey="$NOTICE_CHANNEL_KEY" \
    --taskInterval="$TASK_INTERVAL" \
    "$@"