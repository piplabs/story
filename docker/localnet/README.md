# Usage Instructions

## Steps to Follow

1. Place the `geth` and `story` binary files for `linux/amd64` into the `binary` folder.
2. Check the configuration in the `config` folder to see if any adjustments are necessary.
3. Start the application using the command:
   ```bash
   docker compose up -d
   ```

## Important Notes

1. The current Docker container will reset data upon each restart.
2. In the configuration file, `10.22.22.1` serves as a placeholder IP for `bootnode1`, while `10.22.33.1` to `10.22.33.4` are placeholders for `validator1` to `validator4`. These will be dynamically replaced with container IPs during startup.
3. If you are developing on macOS and seeking optimal performance, you can use the `linux/arm64` binary files. However, you must modify the `docker-compose.yaml` file to set the platform to `linux/arm64`.
