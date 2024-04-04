while true; do
    echo "[run.sh] Starting debugging..."
    rm -rf .git
    dlv debug --headless --log --listen=:2345 --api-version=2 --accept-multiclient --continue -- app:serve &

    PID=$!

    inotifywait -e modify -e move -e create -e delete -e attrib --exclude '(__debug_bin|\.git|\.idea)' -r .

    echo "[run.sh] Stopping process id: $PID"
    
    kill -9 $PID
    pkill -f __debug_bin
    ls -al
done