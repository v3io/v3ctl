# v3ctl
v3ctl is a command line utility for the Iguazio Data Science Platform (the "platform"). With v3ctl, you can control various aspects of the data layer (i.e. list containers, create streams, get objects) and control layer (i.e. list events, create users - future release).

# Usage

Use v3ctl's `--help` argument to see arguments and commands. To see the root help:  
```sh
./v3ctl --help
```

To see the arguments and subcommands of a command:
```sh
./v3ctl create stream --help
```

When running from outside the platform, you will need to provide `--webapi-url` and `--access-key` (or conversely, the `V3IO_API` and `V3IO_ACCESS_KEY` environment variables respectively). When running from within the platform, this information is inferred.