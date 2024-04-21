import { StartedPostgreSqlContainer } from "@testcontainers/postgresql";
import {
  GenericContainer,
  StartedTestContainer,
  AbstractStartedContainer,
  Wait,
} from "testcontainers";

export class SeedContainer extends GenericContainer {
  constructor(pathToBackend: string, container: StartedPostgreSqlContainer) {
    super("golang:1.22.1");
    this.withCopyDirectoriesToContainer([
      {
        source: pathToBackend,
        target: "/app",
      },
    ])
      .withEnvironment({ BRAIN_POSTGRES_HOST: container.getHost() })
      .withEnvironment({
        BRAIN_POSTGRES_PORT: container.getMappedPort(5432).toString(),
      })
      .withEnvironment({ BRAIN_POSTGRES_USER: container.getUsername() })
      .withEnvironment({ BRAIN_POSTGRES_PASSWORD: container.getPassword() })
      .withEnvironment({ BRAIN_POSTGRES_DATABASE: container.getDatabase() })
      .withCommand(["go", "run", "./cmd/seed"])
      .withWaitStrategy(Wait.forLogMessage(/seed complete/i, 1));
  }

  public override async start(): Promise<StartedSeedContainer> {
    return new StartedSeedContainer(await super.start());
  }
}

export class StartedSeedContainer extends AbstractStartedContainer {
  constructor(startedTestContainer: StartedTestContainer) {
    super(startedTestContainer);
  }
}
