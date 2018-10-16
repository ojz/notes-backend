set -e

# clean
killall -q -9 notes-backend || true
rm -f notes.db

# start
go install
notes-backend -init -database notes.db -dev -address ":8085" -root "/api" &

# test
BASE="http://localhost:8085/api"
http() { curl -H 'Content-Type: application/json' "$@" ; }

get()    { http -X GET            "$BASE/$1" ; }
post()   { http -X POST   -d "$2" "$BASE/$1" ; }
put()    { http -X PUT    -d "$2" "$BASE/$1" ; }
delete() { http -X DELETE         "$BASE/$1" ; }

post notes '{ "title": "Hello World!", "body":  "First post!" }'
post notes '{ "title": "Hello Again!", "body":  "Second post!" }'
post notes '{ "title": "Notes", "body":  "This is a *real* not." }'
put notes/3 '{ "title": "Notes", "body":  "This is a *real* note." }'
delete notes/2
get notes
