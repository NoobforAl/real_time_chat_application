# Report logs

## Service Auth

### create two user & login

#### Register EndPoint: POST `{{auth_api}}/register`

```json
{
    "username": "test1235678",
    "password": "test1235678",
    "email": "test@test.com"
},
{
    "username": "test1235678_1",
    "password": "test1235678",
    "email": "test_1@test.com"
}
```

Response message to be like this:

```json
{
    "id": "66856d0f6be8d97d67fb058b",
    "username": "test1235678",
    "email": "test@test.com",
    "password": "$2a$10$ApHLO26EFxlLaASK8TwgJuvut.kgvNN3has58lbAcrmBoORR7qEDu",
    "notification": true,
    "create_at": "2024-07-03T15:23:59.660660625Z"
},
{
    "id": "66856d356be8d97d67fb058c",
    "username": "test1235678_1",
    "email": "test_1@test.com",
    "password": "$2a$10$yc6Or5jDlnZSMEVWSj0yZ.9z284YDlVipFFfnTVGIuqGTNVIZUqPi",
    "notification": true,
    "create_at": "2024-07-03T15:24:37.78850896Z"
}
```

#### Login EndPoint: POST `{{auth_api}}/login`

login with this data for users:

```json
{
    "username": "test1235678",
    "password": "test1235678"
},
{
    "username": "test1235678_1",
    "password": "test1235678"
}
```

Response message to be like this:

```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3QxMjM1Njc4Iiwibm90aWZpY2F0aW9ucyI6ZmFsc2UsImV4cCI6MTcyMDAyMjk1NiwianRpIjoiNjY4NTZkMGY2YmU4ZDk3ZDY3ZmIwNThiIn0.Fd-jgIZ6OwvRh6lKuWfEHvhjdQVYP0DwKePcTLim8gE",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3QxMjM1Njc4Iiwibm90aWZpY2F0aW9ucyI6ZmFsc2UsImV4cCI6MTcyMDEwNjg1Nn0.nd0ZNHcva_qyowWo6AB840VocXEs9jxBEajqKSYPkbU"
},
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3QxMjM1Njc4XzEiLCJub3RpZmljYXRpb25zIjpmYWxzZSwiZXhwIjoxNzIwMDIyOTc5LCJqdGkiOiI2Njg1NmQzNTZiZThkOTdkNjdmYjA1OGMifQ.2X1kykrvuzWd4ATVjLpM5MMlGdQq9jQfwY-9tBOe2kU",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3QxMjM1Njc4XzEiLCJub3RpZmljYXRpb25zIjpmYWxzZSwiZXhwIjoxNzIwMTA2ODc5fQ.NiEXPXxm0emXpiLup6Avw4IlqES1B_CGYVr2XF3c26E"
}
```

If pogrom get error message error like this:

```json
{
    "error": "incorrect username & password"
}
```

Service Logs:

```log
time="2024-07-03T15:23:59Z" level=info msg="Request received" fields.time="2024-07-03 15:23:59.608052969 +0000 UTC m=+405.665642954" ip=172.19.0.1 method=POST path=/register
time="2024-07-03T15:24:37Z" level=info msg="Request received" fields.time="2024-07-03 15:24:37.735191595 +0000 UTC m=+443.792781577" ip=172.19.0.1 method=POST path=/register
time="2024-07-03T15:27:36Z" level=info msg="Request received" fields.time="2024-07-03 15:27:36.513780586 +0000 UTC m=+622.571370570" ip=172.19.0.1 method=POST path=/login
time="2024-07-03T15:27:59Z" level=info msg="Request received" fields.time="2024-07-03 15:27:59.751916221 +0000 UTC m=+645.809506204" ip=172.19.0.1 method=POST path=/login
time="2024-07-03T15:28:47Z" level=info msg="Request received" fields.time="2024-07-03 15:28:47.552746872 +0000 UTC m=+693.610336855" ip=172.19.0.1 method=POST path=/login
```

Note: for use any service need set access token to header response:
> Access-Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3QxMjM1Njc4Iiwibm90aWZpY2F0aW9ucyI6ZmFsc2UsImV4cCI6MTcxOTkzOTc4MywianRpIjoiNjY4NDE2OTE2NjZmMDViOGVhZjQxOTk4In0._RCzX5Znoc1Eck-3qDu_TuR2Mh-zZut4hCqPWd7dwqI

## Service Room

### create & get rooms

#### Create Room EndPoint: POST `{{room_api}}/rooms/`

```json
{
    "name": "this_my_app_name",
    "description": "description"
}
```

Response message to be like this:

```json
{
    "id": "66856f49204b0a076bcde07b",
    "name": "this_my_app_name",
    "description": "description",
    "AllowUsers": null,
    "UserId": "",
    "IsOpen": false,
    "create_at": "2024-07-03T15:33:29.013230251Z",
    "CloseAt": "0001-01-01T00:00:00Z"
}
```

#### Get Room EndPoint: GET `{{room_api}}/rooms/`

Response message to be like this:

```json
[
    {
        "id": "66856f49204b0a076bcde07b",
        "name": "this_my_app_name",
        "description": "description",
        "AllowUsers": null,
        "UserId": "",
        "IsOpen": false,
        "create_at": "2024-07-03T15:33:29.013Z",
        "CloseAt": "0001-01-01T00:00:00Z"
    }
]
```

If pogrom get error message error like this:

```json
{
    "error": "failed to create new room"
}
```

Service Logs:

```log
time="2024-07-03T15:33:13Z" level=info msg="Request received" fields.time="2024-07-03 15:33:13.24456165 +0000 UTC m=+959.241718113" ip=172.19.0.1 method=POST path=/rooms/
time="2024-07-03T15:33:29Z" level=info msg="Request received" fields.time="2024-07-03 15:33:29.012390212 +0000 UTC m=+975.009546675" ip=172.19.0.1 method=POST path=/rooms/
time="2024-07-03T15:33:47Z" level=info msg="Request received" fields.time="2024-07-03 15:33:47.995139373 +0000 UTC m=+993.992295835" ip=172.19.0.1 method=POST path=/rooms/
time="2024-07-03T15:33:47Z" level=error msg="write exception: write errors: [E11000 duplicate key error collection: real_time_chat_app.room index: name_1 dup key: { name: \"this_my_app_name\" }]"
time="2024-07-03T15:33:47Z" level=error msg="write exception: write errors: [E11000 duplicate key error collection: real_time_chat_app.room index: name_1 dup key: { name: \"this_my_app_name\" }]"
time="2024-07-03T15:36:26Z" level=info msg="Request received" fields.time="2024-07-03 15:36:26.312689213 +0000 UTC m=+1152.309845676" ip=172.19.0.1 method=GET path=/rooms/
```

## Start Chat in Room

### Create Message & get all

#### Create Message EndPoint: POST `{{message_api}}/messages/:room_id`

```json
{
    "content": "sample content"
}
```

Response message to be like this:

```json
{
    "id": "6685715cd7884a2ace650163",
    "content": "sample content",
    "sender_id": "66856d0f6be8d97d67fb058b",
    "room_id": "66856f49204b0a076bcde07b",
    "timestamp": "2024-07-03T15:42:20.324584045Z"
}
```

#### Get Messages EndPoint: GET `{{message_api}}/messages/:room_id`

Response message to be like this:

```json
[
    {
        "id": "6685715cd7884a2ace650163",
        "content": "sample content",
        "sender_id": "66856d0f6be8d97d67fb058b",
        "room_id": "66856f49204b0a076bcde07b",
        "timestamp": "2024-07-03T15:42:20.324Z"
    }
]
```

#### Logs

Note: in logs we run Scheduled service, for archive message and reports.

```log
time="2024-07-03T15:17:14Z" level=info msg="Scheduled task with entry ID: 1b651e04-9bd6-4cf6-941d-5473e829c840\n"
time="2024-07-03T15:17:14Z" level=info msg="Scheduler starting"
time="2024-07-03T15:17:14Z" level=info msg="Scheduler timezone is set to UTC"
time="2024-07-03T15:17:14Z" level=info msg="Send signal TERM or INT to stop the scheduler"
time="2024-07-03T15:41:03Z" level=info msg="Request received" fields.time="2024-07-03 15:41:03.966270985 +0000 UTC m=+1429.945773183" ip=172.19.0.1 method=POST path=/messages/668431477b18eb0dc4bb97db
time="2024-07-03T15:41:30Z" level=info msg="Request received" fields.time="2024-07-03 15:41:30.569756614 +0000 UTC m=+1456.549258797" ip=172.19.0.1 method=POST path=/messages/668431477b18eb0dc4bb97db
time="2024-07-03T15:41:36Z" level=info msg="Request received" fields.time="2024-07-03 15:41:36.048746448 +0000 UTC m=+1462.028248641" ip=172.19.0.1 method=POST path=/messages/668431477b18eb0dc4bb97db
time="2024-07-03T15:41:51Z" level=info msg="Request received" fields.time="2024-07-03 15:41:51.662852 +0000 UTC m=+1477.642354192" ip=172.19.0.1 method=POST path=/messages/668431477b18eb0dc4bb97db
time="2024-07-03T15:42:20Z" level=info msg="Request received" fields.time="2024-07-03 15:42:20.32360031 +0000 UTC m=+1506.303102494" ip=172.19.0.1 method=POST path=/messages/66856f49204b0a076bcde07b
time="2024-07-03T15:43:48Z" level=info msg="Request received" fields.time="2024-07-03 15:43:48.747459456 +0000 UTC m=+1594.726961652" ip=172.19.0.1 method=GET path=/messages/66856f49204b0a076bcde07b
^Ccontext cancele
```

### Ws Message & Notification

#### WebSocket EndPoint: WS `{{message_api}}/messages/ws/:room_id`

Send message log:

```json
{"id":"","content":"test content!","sender_id":"66856d0f6be8d97d67fb058b","room_id":"66856f49204b0a076bcde07b","timestamp":"2024-07-03T15:53:31.973605507Z"}
19:23:31
{ "content": "test content!" }
19:23:31
Connected to ws://localhost:3002/messages/ws/66856f49204b0a076bcde07b
```

If notification is `on`: recived message:

```json
Search
{"id":"","content":"test content!","sender_id":"66856d0f6be8d97d67fb058b","room_id":"66856f49204b0a076bcde07b","timestamp":"2024-07-03T15:57:31.773308037Z"}
{"id":"","content":"test content!","sender_id":"66856d0f6be8d97d67fb058b","room_id":"66856f49204b0a076bcde07b","timestamp":"2024-07-03T15:57:31.333993247Z"}
{"id":"","content":"test content!","sender_id":"66856d0f6be8d97d67fb058b","room_id":"66856f49204b0a076bcde07b","timestamp":"2024-07-03T15:57:30.535358342Z"}
{"id":"","content":"test content!","sender_id":"66856d356be8d97d67fb058c","room_id":"66856f49204b0a076bcde07b","timestamp":"2024-07-03T15:56:18.029047061Z"}
{ "content": "test content!" }
```

#### Notification EndPoint: WS `{{message_api}}/notification/ws/:room_id`

Turn on Notification: (send string on or off)

```json
send: on
19:24:36
Connected to ws://localhost:3002/notification/ws/66856f49204b0a076bcde07b
```

Now send new message an receive this:

```json
{"id":"","sender_id":"66856d356be8d97d67fb058c","room_id":"66856f49204b0a076bcde07b","notification_room_id":"","Content":"received message: test content!","create_at":"2024-07-03T15:56:18.295285291Z",
send: on
Connected to ws://localhost:3002/notification/ws/66856f49204b0a076bcde07b
```

#### WS logs

```logs
time="2024-07-03T15:50:15Z" level=info msg="Request received" fields.time="2024-07-03 15:50:15.499332105 +0000 UTC m=+1981.478834284" ip=172.19.0.1 method=GET path=/messages/ws/668431307b18eb0dc4bb97d9
time="2024-07-03T15:50:15Z" level=error msg="mongo: no documents in result"
time="2024-07-03T15:50:26Z" level=info msg="Request received" fields.time="2024-07-03 15:50:26.690158398 +0000 UTC m=+1992.669660577" ip=172.19.0.1 method=GET path=/messages/ws/668431307b18eb0dc4bb97d9
time="2024-07-03T15:50:26Z" level=error msg="mongo: no documents in result"
time="2024-07-03T15:50:45Z" level=info msg="Request received" fields.time="2024-07-03 15:50:45.230853425 +0000 UTC m=+2011.210355604" ip=172.19.0.1 method=GET path=/messages/ws/668431307b18eb0dc4bb97d9
time="2024-07-03T15:50:45Z" level=error msg="mongo: no documents in result"
time="2024-07-03T15:50:53Z" level=info msg="Request received" fields.time="2024-07-03 15:50:53.90238197 +0000 UTC m=+2019.881884150" ip=172.19.0.1 method=GET path=/notification/ws/668431307b18eb0dc4bb97d9
time="2024-07-03T15:50:53Z" level=error msg="mongo: no documents in result"
time="2024-07-03T15:51:24Z" level=info msg="Request received" fields.time="2024-07-03 15:51:24.582046809 +0000 UTC m=+2050.561548988" ip=172.19.0.1 method=GET path=/messages/ws/66856f49204b0a076bcde07b
time="2024-07-03T15:51:27Z" level=info msg="Request received" fields.time="2024-07-03 15:51:27.021429784 +0000 UTC m=+2053.000931965" ip=172.19.0.1 method=GET path=/messages/ws/66856f49204b0a076bcde07b
time="2024-07-03T15:51:30Z" level=info msg="Request received" fields.time="2024-07-03 15:51:30.484165871 +0000 UTC m=+2056.463668054" ip=172.19.0.1 method=GET path=/notification/ws/66856f49204b0a076bcde07b
time="2024-07-03T15:51:32Z" level=info msg="Request received" fields.time="2024-07-03 15:51:32.206811889 +0000 UTC m=+2058.186314067" ip=172.19.0.1 method=GET path=/notification/ws/66856f49204b0a076bcde07b
time="2024-07-03T15:51:40Z" level=error msg="redis: nil"
time="2024-07-03T15:51:42Z" level=error msg="redis: nil"
time="2024-07-03T15:51:50Z" level=error msg="redis: nil"
time="2024-07-03T15:51:52Z" level=error msg="redis: nil"
time="2024-07-03T15:52:00Z" level=error msg="redis: nil"
time="2024-07-03T15:52:02Z" level=error msg="redis: nil"
time="2024-07-03T15:52:10Z" level=error msg="redis: nil"
time="2024-07-03T15:52:12Z" level=error msg="redis: nil"
```

## Grpc Service

### Logs Grpc Service

```log
time="2024-07-03T15:17:13Z" level=warning msg="not found .env file!"
time="2024-07-03T15:17:13Z" level=info msg="server listening at 172.19.0.5:50023"
time="2024-07-03T15:33:13Z" level=error msg="token is expired by 22h30m10s"
time="2024-07-03T15:41:04Z" level=error msg="token is expired by 22h18m43s"
time="2024-07-03T15:41:30Z" level=error msg="token is expired by 22h38m27s"
time="2024-07-03T15:41:36Z" level=error msg="token is expired by 22h38m33s"
```
