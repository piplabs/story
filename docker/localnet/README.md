# Usage Instructions

## Steps to Follow

1. Place the `geth` and `story` binary files for `linux/amd64` into the `binary` folder.
2. Check the configuration in the `config` folder to see if any adjustments are necessary.
3. Start the application using the command:
   ```bash
   docker compose -f ./docker-compose-1+2.yaml up -d
   ```

## Important Notes

1. In the configuration file, `10.22.22.1` serves as a placeholder IP for `bootnode1`, while `10.22.33.1` is placeholders for `validator1` to `validator4`. These will be dynamically replaced with container IPs during startup.
2. If you are developing on macOS and seeking optimal performance, you can use the `linux/arm64` binary files. However, you must modify the `docker-compose.yaml` file to set the platform to `linux/arm64`.
