# webhook-server

[webhook-server][] works together with [webhook-client][] to put your localhost on the Internet.

# Example:

```console
$ webhook-client -url=http://localhost:8080
Congrats!!!
You have access to: http://localhost:8080/<your_path>
using the next url: https://webhook-server.fly.dev/326e6580-e380-41a0-8c19-786a3e4f7fd4/<your_path>
```

Then, when you access to https://webhook-server.fly.dev/326e6580-e380-41a0-8c19-786a3e4f7fd4, you will see the content in http://localhost:8080



# How it works
[webhook-server][] is a public site (https://webhook-server.fly.dev/) using RedisDB as a database.

[webhook-client][] is executed locally in your computer.

When [webhook-client][] starts, connects to the [webhook-server][] to create a token. This token (uuid) will be used in all the communications between them.

Every time the [webhook-server][] receives a request, all the request info is stored in the DB. In the other side, [webhook-client][] will be checking if there is any request pending to process. Once the [webhook-client][] receives the request it will process using the local url (ie. http://localhost:8080) and it will send back the response to the [webhook-server][].

# Use cases
- The common use case is during the development process to test third party providers (payments, shipping platforms, ...) that use webhooks to send a callback.


[webhook-server]: https://github.com/xegea/webhook_server
[webhook-client]: https://github.com/xegea/webhook_client
