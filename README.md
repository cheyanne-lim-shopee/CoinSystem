# CoinSystem

Background: 
A e-commerce system required us to help build a coin system for their customer, which can support add coin to user/Deduct coin from user

Business requirement: 
1. Provide an api to add coin to user(coin amount, user id)
    1. system has max coin limit per user, if this adding exceed the limit, we will block it
2. Provide an api to deduct coin to user(coin amount, user id)
    1. cannot deduct coin exceed the user balance
3. Provide a query api(user ids), that can query the balance of users 
    1. also it can tell u the latest update time of coin balance

Tech requirement: 
1. Build in go language 
2. Client-server structure 
3. Mysql as database 
4. TCP & protobuf as communication 