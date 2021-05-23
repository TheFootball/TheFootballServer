# The Football Server

JUNCTION X SEOUL 2021 AWS GAMETECH 

5HYU1HUFS TEAM

## Point that we focused on
1. AWS Architecture  
    AWS tech for heavy Traffic Game. Heavy traffic comes in Streamer VS Viewers composition ( 1 vs 1000 )
2. Simple Game Design  
    Because of limited time and resource, Game itself should not be complicated.
3. Impact and Creative  
    Simple designed game can be loose and boring, so unique concept is require.


## How we planned to those Point
1. AWS Fargate, AWS Elasticache  
    To deal with heavy traffic, we choose AWS Fargate for scale out. Scale out can cause situation that Web Socket connections don't work well. 
    So we think about Elasticache and Pub/Sub feature of it for shared memory and connectivity.
2. SUBERUNKER, but each dropping object is user too.  
    SUBERUNKER is old famous flash game. This game is simple itself. So we can use more time to put in an effort in Cloud Architecture. As simple as it was, it was easy to create an unique additional concept.
3. Impact and Creative  
    As I said above, the game is play in a streamer vs. viewer composition. So we thought of a game where 1 vs 1000 is possible through chatting system like Twitch API.
    This is an effective way for many viewers to interact with the streamer.
    And the main concept of the game is a human runningback that avoids animal football players running from above.
     
## Stack Decision
Golang with docker container, Redis DB Pub/Sub, Websocket

## Result
70% Complete. Progress was slower than first planned for the reasons below. 
- First experience for many new tech feature
- No experience for game server (All developers are Web developer)
- Too many things to consider in scale-out
- Too much time spent on the first initiative planning


## Last Words of participants
- @Neulhan  
    It's too bad I couldn't complete it, but it was a good time to try various things.
- @hongdoojung  
    I learned a lot from using skills that I didn't usually use.
- @changhoi  
    We have chosen a very difficult subject.