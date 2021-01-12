# shellcheck disable=SC2046
docker volume rm $(docker volume ls -q);
s=$(docker ps -a | grep 'subd-proj' | cut -d ' ' -f 1);  docker kill "$s";
docker ps -a | grep  Exit| cut -d ' ' -f 1 | xargs  docker rm; docker ps -a;
docker build -t subd-proj ../.; docker run -p 5000:5000 --name subd-proj -t  subd-proj;