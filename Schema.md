### Schema

##### Ueer

```sh
User = {
    name,
    email,
    password
}
```

#### Friendship

```sh
friendship={
user1:userid,
user2:userid,
}
```

#### Message

```sh
message={
    friendshipId:friendshipid,
    last_id:number,
    text:string,
    from:userId,
    seen:bool,
    sent:bool
    createAt:""
}
```

## Operations - pseudo-sync local db and production db

### 1. First time login

- Fetch all friends, put them on db n display on screen
- Fetch all messages, save in local db
- count not-seen messages display that counter besides username
- As user click on friend name, fetch relevent message paint them screen
-
- Store in localDB and render on the screen
- Sort the friends in descending order of last message's createdAt time

### 2.Chatting started

suppose userA wants to send message to userB

- UserB would create and listen for a channel named on its objectid
- When userA send message to this userB,

  - first that message would be saved in the DB and push to the userB's channel in pub/sub, if userB is online
  - after that when userA recive back the 200 status, if would save the message in its local db

- On userB side, after recievinf the message from userA, it re-order the friends list, push the message obj in localDB

### 3. UserB come online

- UserB hit pendingchats api endpoint, to fetch the messages since the last online time

## Questions

1. How to scan a dir in
