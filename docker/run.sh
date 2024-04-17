while true; do
    echo "[run.sh] Starting debugging..."

    go run server.go app:serve &

    PID=$!

    inotifywait -e modify -e move -e create -e delete -e attrib --exclude '(__debug_bin|\.git|\.idea)' -r .

    echo "[run.sh] Stopping process id: $PID"
    
    kill -9 $PID
    pkill -f __debug_bin
    ls -al
done