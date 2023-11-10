# Use case 
A common use case that is almost everywhere -- new user sign-up/register a new account in a website/system.
E.g. Amazon/Linkedin/Google/etc...

### Use case requirements

* User fills a form and submit to the system with email
* System will send an email for verification
* User will click the link in the email to verify the account
* If not clicking, a reminder will be sent every X hours

<img width="303" alt="user case requirements" src="https://github.com/indeedeng/iwf-python-sdk/assets/4523955/356a4284-b816-42d3-9e44-b371a91834e4">

### Some old solution

With some other existing technologies, you solve it using message queue(like SQS which has timer) + Database like below:

<img width="309" alt="old solution" src="https://github.com/indeedeng/iwf-python-sdk/assets/4523955/49ef8846-9589-4a28-91bd-c575daf37dcf">

* Using visibility timeout for backoff retry
* Need to re-enqueue the message for larger backoff
* Using visibility timeout for durable timer
* Need to re-enqueue the message for once to have 24 hours timer
* Need to create one queue for every step
* Need additional storage for waiting & processing ready signal
* Only go to 3 or 4 if both conditions are met
* Also need DLQ and build tooling around

**It's complicated and hard to maintain and extend.**

### How to run this example

1. Start xdb [server](https://github.com/xdblab/xdb#option-1-use-example-docker-compose-of-xdb-with-a-database):
```shell
wget https://raw.githubusercontent.com/xdblab/xdb/main/docker-compose/docker-compose-postgres14-example.yaml && docker compose -f docker-compose-postgres14-example.yaml up -d
```

2. Start the example:
```shell
go run cmd/server/main.go
```

3. Send a request to sign up:
```shell
http://localhost:8803/signup/start?userId=test1&email=abc@c.com&firstName=Quanzheng&lastName=Long
```
If you send duplicate userId to signup, you will see error:
```shell
"StatusCode: 424 , error details: {\"detail\":\"Failed to write global attributes, please check the error message for details: pq: duplicate key value violates unique constraint \\\"sample_user_table_pkey\\\"\"}"
```
You would see some reminder "email" in the console:
```shell
sending an email to abc@c.com, title: Quanzheng, please verify your email, content: .....more content
```

4. Send a request to verify:
```shell
http://localhost:8803/signup/verify?userId=test1
```

