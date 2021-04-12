echo "Input 0-2 to kill node 0-2"
echo "Input 3 to kill all nodes"
target_list=()
read -p "Please input who u want kill :" target
if [ $target -eq "3" ];
then  
    target_list[0]="0"
    target_list[1]="1"
    target_list[2]="2"
else
    target_list[0]="${target}"
fi
echo
echo ---------------------------------------
echo target nodes : ${target_list[*]}
echo ---------------------------------------
echo
echo "Input -1 to shutdown"
echo "Input a number >= 0 to set duration before restart"
read -p "how many seconds do u want the restart to be delayed? " delay
echo 
echo "shutdown all targets"
echo ----------------------
# kill all targets
for t in "${target_list[@]}";do
    docker stop "peer-""${t}"
done
echo

if [ $delay -ne -1 ];
then
    echo "delay ${delay}s to restart..."
    echo -------------------------------
    sleep $delay
    echo "restarting..."
    # restart targets
    for t in "${target_list[@]}"; do
        docker restart "peer-""${t}"
    done
fi
    

