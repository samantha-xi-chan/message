#!/bin/sh

# Supervisor sends TERM to services when stopped.
# This wrapper has to pass the signal to it's child.
# Note that we send TERM (graceful) instead of KILL (immediate).
_term() {
    kill -TERM "$child" 2>/dev/null
    exit 1
}

trap _term SIGTERM

# Execute console.php with whatever arguments were specified to this script
"$@" &
child=$!
wait "$child"
rc=$?

# Delay to prevent supervisor from restarting too fast on failure
sleep 3

# Return with the exit code of the wrapped process
exit $rc