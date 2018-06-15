curl -X POST \
	  http://127.0.0.1:8080/signup \
	    -H 'cache-control: no-cache' \
	      -H 'content-type: application/x-www-form-urlencoded' \
		  -d 'username=yzs9000&password=yzsyzsyzs'

curl -X POST \
	  http://127.0.0.1:8080/signup \
	    -H 'cache-control: no-cache' \
	      -H 'content-type: application/x-www-form-urlencoded' \
		  -d 'username=yzs8000&password=yzsyzsyzs'

curl -X POST \
	  http://127.0.0.1:8080/upload \
	    -H 'cache-control: no-cache' \
	      -H 'content-type: application/x-www-form-urlencoded' \
		  -d 'username=yzs8000&password=yzsyzsyzs&musicname=Yuzuki%20-%20you(Vocal)&difficulty=3&rank=SS&score=7000'
		  
curl -X POST \
	  http://127.0.0.1:8080/upload \
	    -H 'cache-control: no-cache' \
	      -H 'content-type: application/x-www-form-urlencoded' \
		  -d 'username=yzs9000&password=yzsyzsyzs&musicname=Yuzuki%20-%20you(Vocal)&difficulty=3&rank=SSS&score=9000'

curl -X POST \
  http://127.0.0.1:8080/rank \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/x-www-form-urlencoded' \
  -d 'musicname=Yuzuki%20-%20you(Vocal)&difficulty=3'
  
curl -X POST \
  http://127.0.0.1:8080/rank \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/x-www-form-urlencoded' \
  -d 'musicname=EastRed&difficulty=2'
curl -X POST \
  http://127.0.0.1:8080/search \
  -H 'cache-control: no-cache' \
  -d you

