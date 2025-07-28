# l3http

This level3 package provides an implementation of the graceful HTTP server.

It is designed to be ingress-friendly and used in a container environment, k8s or otherwise.

Use Ready() as an another healthcheck for the ingress.

## Startup

As soon as the server is created (via New()), it starts accepting the connections.  
The connections hang in a waiting state until the Run() call; then it passes the requests to the provided HTTP server.

This is done so that:
- You won't lose your clients that came in while the service was starting up. This may happen because of possible ingress misconfiguration.
- The app would exit early (`New()` would return an error) if the ports are already in use. Don't you hate it when you wait for two minutes for the app to start only to find out that the ports are already in use? :)

## Shutdown

The server is stopped via the `Server.GracefulShutdown(ctx)` call. It blocks until all the requests are processed or the context is canceled.

Keep in mind that there's no delay between `Server.GracefulShutdown(ctx)` and the server being stopped.

To ensure zero downtime for your infra, mark some other healthcheck as not ready and wait some time (healthcheck period*2) before calling `Server.GracefulShutdown(ctx)`. This will ensure that the ingress will stop routing traffic to the server before it stops accepting the connections.

This logic is implemented in the level6 packages. TODO which?