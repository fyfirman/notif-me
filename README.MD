
## How to run

1. Make sure you have Golang installed on your machine. You can download it from the official Golang website: https://golang.org/dl/

2. Clone the repository to your local machine:
   ```
   git clone <repository-url>
   ```

4. Build the project:
   ```
   go build
   ```

5. Run the executable:
   ```
   ./notif-me
   ```

You can customize these instructions based on your specific project setup and requirements.

## Deploy

```sh
docker build --platform linux/amd64 -t vm-do.fyfirman.com:23003/fyfirman/notif-me:0.4.1 .
docker push vm-do.fyfirman.com:23003/fyfirman/notif-me:0.4.1
```
