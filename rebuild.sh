set -e

# clean
killall -q -9 notes-backend || true
rm -f notes.db

# start
go install
notes-backend -init -database notes.db -dev -address ":8085" &
sleep 1

# test
ROOT="http://localhost:8085/api"
http() { curl -H 'Content-Type: application/json' "$@" ; }

get()    { http -X GET            "$ROOT/$1" ; }
post()   { http -X POST   -d "$2" "$ROOT/$1" ; }
put()    { http -X PUT    -d "$2" "$ROOT/$1" ; }
delete() { http -X DELETE         "$ROOT/$1" ; }

post notes '{ "title": "test", "body":  "ok" }'
post notes '{ "title": "test2", "body":  "ok2" }'
get notes
put notes/1 '{ "title": "updated title", "body":  "updated body" }'
# delete notes/1
get notes/1
get notes
