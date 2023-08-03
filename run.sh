#!/bin/bash
function cleanup()
{
        local pids=`jobs -p`
        if [[ "$pids" != "" ]]; then
                kill $pids >/dev/null 2>/dev/null
        fi
}

trap cleanup EXIT
logagent -c conf.json >> /root/logagent.out 2>&1
