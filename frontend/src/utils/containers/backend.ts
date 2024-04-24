import { StartedPostgreSqlContainer } from "@testcontainers/postgresql";
import { GenericContainer, Wait } from "testcontainers";

export async function createNewBackendContainer(
  pathToBackend: string,
  container: StartedPostgreSqlContainer,
): Promise<GenericContainer> {
  const backend = await GenericContainer.fromDockerfile(pathToBackend).build(
    "brain-backend",
    { deleteOnExit: false },
  );

  backend
    .withCopyFilesToContainer([
      { source: pathToBackend + "/config.json", target: "/app/config.json" },
    ])
    .withEnvironment({ BRAIN_POSTGRES_HOST: container.getHost() })
    .withEnvironment({
      BRAIN_POSTGRES_PORT: container.getPort().toString(),
    })
    .withEnvironment({ BRAIN_POSTGRES_USER: container.getUsername() })
    .withEnvironment({ BRAIN_POSTGRES_PASSWORD: container.getPassword() })
    .withEnvironment({ BRAIN_POSTGRES_DATABASE: container.getDatabase() })
    .withExposedPorts(42069, 42069)
    .withWaitStrategy(Wait.forLogMessage(/starting server/i, 1));

  return backend;
}
