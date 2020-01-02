#https://www.techcoil.com/blog/how-to-migrate-your-mongodb-database-instance-with-mongodump-mongorestore-tar-and-scp/

mongodump --db vocaloidDB

scp -r ./dump root@vocaloidtube.com:~

mongorestore --drop --db vocaloidDB dump/vocaloidDB

