# DecisionMaker

# This is a smart  tree for Golang.
# You can run an action based system or smart decision making

It utilizes a decision tree to build and resolve questions. 

It runs a service in  'http://localhost:1323/'
It contains a RestAPI and a Tester

The initial Tree is like this chart saved in file “treedb.json",
The schematic structure of the tree is like as following
 
 # You can run main.go and test it by an api tester app
 
 



![help](https://user-images.githubusercontent.com/11687423/127484604-1694e286-70a6-46bf-952f-fb09f492702e.jpg)

A sample API test:
![image](https://user-images.githubusercontent.com/11687423/127484926-b4631f5b-0d25-4a6c-96d2-1dd1d494927a.png)

TEST1- CURL:

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "feeling sick ?":"yes",
   "Are you in a strong pain?":"yes",
  "do you have any of the following conditions? if yes select one of them":"no, i dont have any of the above conditions",
  "How old are you?":"23"
}' \
 'http://localhost:1323/tree/resolve'



TEST2- CURL:


curl -i -X GET \
 'http://localhost:1323/'
