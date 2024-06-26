while true; do
    echo "[run.sh] Starting debugging..."

    go run server.go app:serve &

    PID=$!

    echo "[run.sh] Server process ID: $PID"

    inotifywait -e modify -e move -e create -e delete -e attrib --exclude '(__debug_bin|\.git|\.idea)' -r .

    echo "[run.sh] Stopping process id: $PID"

    kill -9 $PID
done