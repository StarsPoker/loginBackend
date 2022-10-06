#!/bin/sh

# create directory for Delve logs, we use it to know that Delve
# debugger is running
port_debug=40000

mkdir -p /tmp/dlv_log


runServer() {
  echo Running server

  # create directory and file to

  touch /tmp/dlv_log/output.log

  # run server with debug
  dlv --listen=:${port_debug} --headless=true --accept-multiclient --continue --api-version=2 exec \
   /server | tee -a /tmp/dlv_log/output.log &

  # wait for Delve to modify log files - means /server is running
  inotifywait -e MODIFY /tmp/dlv_log/output.log &>/dev/null

  echo Delve PID: $(pidof dlv), Server PID: $(pidof server)
  pidof dlv > /tmp/dlv.pid
  pidof server > /tmp/server.pid
}

killRunningServer() {
  if [ -f /tmp/dlv.pid ]
  then
    echo killing old Delve, PID: $(cat /tmp/dlv.pid)
    kill $(cat /tmp/dlv.pid)
    rm -f /tmp/dlv.pid
  fi

  if [ -f /tmp/server.pid ]
  then
    echo killing old server, PID: $(cat /tmp/server.pid)
    kill $(cat /tmp/server.pid)
    rm -f /tmp/server.pid
  fi
}

function buildServer() {
  echo Building server

  go build -o /server -gcflags -N -gcflags -l -buildvcs=false
  return $?
}

rerunServer () {
  killRunningServer
  buildServer
  resultBuild=$?
  if [ $resultBuild -eq "0" ]; 
  then
    echo "Deu certo"
    runServer
  else
    echo "Erro ao buildar"
  fi
 
}

lockBuild() {
  # check lock file existence
  if [ -f /tmp/server.lock ]
  then
    # waiting for the file to delete
    inotifywait -e DELETE /tmp/server.lock
  fi
  touch /tmp/server.lock
}

unlockBuild() {
  # remove lock file
  rm -f /tmp/server.lock
}

# run the server for the first time
rerunServer

inotifywait -e MODIFY -r -m /app |
  while read path action file; do
    lockBuild
      ext=${file: -3}
      if [[ "$ext" == ".go" ]]; then
        echo File changed: $file
        rerunServer
      fi
    unlockBuild
  done