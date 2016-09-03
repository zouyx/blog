#!/bin/bash
#

# source function library
. /etc/init.d/functions

APP=blog
APP_DIR=/appdata/deploy/blog
WORK_DIR=$APP_DIR/work/
PID_FILE=$WORK_DIR/pid

RETVAL=0

#
wait_for_pid () {
	try=0
	while test $try -lt 30; do
		case "$1" in
			'created')
			if [ -f "$2" ] && [ "x$(cat $2)" != 'x' ]; then
				try=''
				break
			fi ;;
			'removed')
			if [ ! -f "$2" ] ; then
				try=''
				break
			fi ;;
		esac
		echo -n .
		try=`expr $try + 1`
		sleep 1
	done
}

#
start()
{
    PID=$(cat $PID_FILE 2>/dev/null)
    if checkpid "$PID"; then
        status -p $PID_FILE $APP
        return 0
    fi

    mkdir -p $WORK_DIR/
    run_cmd="$APP_DIR/$APP"
    $run_cmd >>$WORK_DIR/stdout.log 2>&1 &
    echo $! >"$PID_FILE"

	if [ $? -gt 0 ]; then
		failure
        return 1
	fi

	wait_for_pid created $PID_FILE
	if [ -n "$try" ] ; then
		failure
        return 1
	else
		success && echo "Starting $APP"
	fi
}

#
stop()
{
	PID=$(cat $PID_FILE 2>/dev/null)
	if checkpid "$PID"; then
		kill -KILL $PID && rm $PID_FILE 2>/dev/null
		wait_for_pid removed $PID_FILE
		if [ -n "$try" ] ; then
			failure
			return 1
		else
			success
		fi
	else
		failure
		rm $PID_FILE 2>/dev/null
	fi
	echo "Stopping $APP"
}

restart()
{
	stop
	start
}

case "$1" in
  start)
	start ;;
  stop)
	stop ;;
  restart|force-reload|reload)
	restart ;;
  status)
	status -p $PID_FILE $APP ;;
  *)
    echo $"Usage: $0 {start|stop|status|restart|reload|force-reload}"
    exit 1
esac

RETVAL=$?
exit $RETVAL