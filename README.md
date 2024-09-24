To run this project in Docker, navigate to the root directory of the project and execute the following command:

```bash
docker compose up --build
```

This command will build the Docker image and start the project inside a container. Once the process is complete, the project will be accessible at `http://localhost:5173/`.

If you plan to deploy this project to a different environment, you should update the `VITE_API_URL` in the `compose.yaml` file to point to your server's IP address.