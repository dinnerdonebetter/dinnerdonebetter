# Code Organization

## Common libraries

When multiple apps have similar needs in a library, it should be implemented in `packages`, in a way that allows for
some minor customization of details. For instance, every app needs tracing support, but they need to specify the name of
the service being traced. In that instance, we'll have an `apps/<app>/src/tracing` folder which calls the
`packages/tracing` code with a string argument identifying the service. In an ideal world any sufficiently useful
library that you might put in `apps/<app>/src` can be implemented in `packages` instead, unless of course it's
hyper-specific to that app.
