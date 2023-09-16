#!/bin/sh

/soft/rooms_robot \
    --notice="$NOTICE_CHANNEL" \
    --noticeUrl="$NOTICE_CHANNEL_URL" \
    --noticeKey="$NOTICE_CHANNEL_KEY" \
    "$@"