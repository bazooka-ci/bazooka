# bzk service

The `service` command of the CLI manage the installation, upgrade, and management of all necessary bazooka containers

```
Usage: bzk service  COMMAND [arg...]

Manage bazooka service (start, stop, status, upgrade...)

Commands:
  start        Start bazooka
  restart      Restart bazooka
  upgrade      Upgrade bazooka to the latest version
  stop         Stop bazooka
  status       Get bazooka status
```

## Start Bazooka

You can run `bzk service start` with multiple options

```
Usage: bzk service start [--home|--scm-key|--registry|--docker-sock]... [--tag]

Start bazooka

Options:
  --home=""          Bazooka's work directory ($BZK_HOME)
  --scm-key=""       Location of the private SSH Key Bazooka will use for SCM Fetch ($BZK_SCM_KEYFILE)
  --mongo-uri=""     URI of a MongoDB server ($BZK_MONGO_URI)
  --registry=""      ($BZK_REGISTRY)
  --docker-sock=""   Location of the Docker unix socket, usually /var/run/docker.sock ($BZK_DOCKERSOCK)
  --tag=""           The bazooka version to run
```

## Restart Bazooka

Restart Bazooka with the options previously set with `bzk run start`

```
Usage: bzk service restart

Restart bazooka
```

## Upgrade Bazooka

Upgrading Bazooka means downloading newer images of Bazooka from the Docker hub. Be aware that this action may take some time

```
Usage: bzk service upgrade

Upgrade bazooka to the latest version
```

## Stop Bazooka

`bzk service stop` will stop the Docker containers of Bazooka

```
Usage: bzk service stop

Stop bazooka
```

## Bazooka status

`bzk service status` will print the status of Bazooka containers

```
Usage: bzk service status

Get bazooka status
```
