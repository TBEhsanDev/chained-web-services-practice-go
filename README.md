# chained-web-services-practice-go</br>
‫


this is a practicing project that has 4 steps.in this project clients send post request to server and server get some information from request,save posted
data as log in file and return request body in addition request ip.
##  necessary tools:
1. nginx(used as proxy server)
2. logstash(used to read logs from file and send to db)
3. mongodb(used as db)
4. code coverage(used to check tests how much cover the code)
5. apachebench(used to check load test on project)
6. redis(used to limit count of recieved requests as rate limiter).
### step1:
 create an ip server that add client ip to response body and its related test script.then apply code coverage and apachebench .
### step2:
 add logger server between nginx and ip server to log data to jsonline file.this data concludes body request in addition request time.also should modify
 test script.then apply code coverage and apachebenc for this step too.
### step3:
 setup logstash and mongodb that records in jsonl file read by logstash and write in db.also should modify test file.
### step4:
 add rate limiter server between nginx and logger server to control number of requests.(should use redis to control).in this step also should modify test
 and code coverage is used too. and code coverage is used too.
## instruction for run and test:
> for run servers ,run bash script  **./run.sh**.for test project run bash script **./test.sh**.
## instruction for dockerize app:
>first **docker build -t ip .** to create a image that named ip.then **docker run -it --name ip ip**to create a container that named ip.this run our 
 our app.**docker ps -a** list the containers.then in another terminal tab run **docker exec -it 'container_id' /bin/bash** to create bash terminal in
 that container .then we can request in that bash.
