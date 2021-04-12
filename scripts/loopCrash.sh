
read -p "Plaese input how many times u want to crash :" times
prefix="http://127.0.0.1:800"
suffix="/crash"
for cnt in $( seq 1 ${times} )
do
    for k in $( seq 0 2 )
    do
        curl   ""${prefix}${k}${suffix}"" 1>/dev/null 2>/dev/null
        echo "kill node "${k}""
        sleep 15
    done
    
done

