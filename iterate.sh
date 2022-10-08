#!/bin/bash

for i in {1..10}
do
   printf "User $i: "
   curl -s -k 'GET' "http://localhost:8081/user?userId=$i" | jq
   printf "\n"

   printf "Images\n"
   curl -s -k 'GET' -H "userId:$i" 'http://localhost:8081/images' | jq
   printf "\n\n"
done


