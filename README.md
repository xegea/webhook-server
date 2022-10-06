# webhook-server

<a href="https://github.com/xegea/webhook_server">webhook_server</a> works together with <a href="https://github.com/xegea/webhook_client">webhook-client</a> to put your localhost in the Internet.

# Example:

$ webhook-client -url=http://localhost:8080

<code>Congratulations!!!</code>

<code>You have access to: http://localhost:8080/<your_path></code>

<code>using the next url: https://webhook-server.fly.dev/326e6580-e380-41a0-8c19-786a3e4f7fd4/<your_path></code>



Then, when you access to https://webhook-server.fly.dev/326e6580-e380-41a0-8c19-786a3e4f7fd4, you will see the content in http://localhost:8080



# How it works
<a href="https://github.com/xegea/webhook_server">webhook_server</a> is a public site (https://webhook-server.fly.dev/) using RedisDB as a database.

<a href="https://github.com/xegea/webhook_client">webhook-client</a> is executed locally in your computer.

When <a href="https://github.com/xegea/webhook_client">webhook-client</a> starts, connects to the <a href="https://github.com/xegea/webhook_server">webhook_server</a> to create a token. This token (uuid) will be used in all the communications between them.

Every time the <a href="https://github.com/xegea/webhook_server">webhook_server</a> receives a request, all the request info is stored the DB. In the other side, <a href="https://github.com/xegea/webhook_client">webhook-client</a> will be checking if there is any request pending to process. Once the <a href="https://github.com/xegea/webhook_client">webhook-client</a> receives the request it will process using the local url (ie. http://localhost:8080) and it will send back the response to the <a href="https://github.com/xegea/webhook_server">webhook_server</a>.

# Use cases
- The common use case is during the development process to test third party providers (payments, shipping platforms, ...) that use webhooks to send a callback.
