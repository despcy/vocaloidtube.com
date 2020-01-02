项目目的：
制作一个VOCALOID 聚合网站

Parts:

1. MongoDB youtube视频增删改查的interface
2. 用户添加视频链接入库的API->request youtube api echo框架辅助
3. 前端搜索功能
4. Oauth用户登录以及session管理，用Cashbin做RBAS权限管理
5. 用户播放列表
6. 后台自动爬取机器人，miku看板娘websocket、、机器人这里不止有算法，还要设计一些红黑名单之类的

youtube:

part=topicDetails(2)  的时候返回的类型都是亚洲音乐，其中包含/m/04rlf(music) 或 /m/028sqc 就可以作为一个filter https://gist.github.com/stpe/2951130dfc8f1d0d1a2ad736bef3b703  -》这个filter交给看板娘来做

part=statistics(2) ->   
 "statistics": {
    "viewCount": "193929",
    "likeCount": "3357",
    "dislikeCount": "68",
    "favoriteCount": "0",
    "commentCount": "279"
   }

snippet(2)

quota(10000/day),每天可以爬1600首了

https://gist.github.com/dgp/1b24bf2961521bd75d6c categoryId 也可以在robot中作为筛选



