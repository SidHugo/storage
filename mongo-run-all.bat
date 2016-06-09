START mongod --configsvr --replSet configReplSet --port 1488 --dbpath C:\dev\db\mongo\configdb
START mongod --port 1499 --dbpath C:\dev\db\mongo\shard1
START mongos --configdb configReplSet/127.0.0.1:1488