sudo docker stop Forum &
pid=$!
wait $pid
echo "-----"
sudo docker rm Forum &
pid=$!
wait $pid
sudo docker rmi test &
pid=$!
wait $pid
echo "-----"
sudo docker images &
pid=$!
wait $pid
echo "-----"
sudo docker ps -a 
echo "-----"
echo "tests terminés"