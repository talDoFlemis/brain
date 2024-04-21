import { StartedPostgreSqlContainer } from "@testcontainers/postgresql";
import {
  GenericContainer,
  StartedTestContainer,
  AbstractStartedContainer,
  Wait,
} from "testcontainers";

export class MigrationContainer extends GenericContainer {
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
      .withCommand(["go", "run", "./cmd/migrate"])
      .withWaitStrategy(Wait.forLogMessage(/migration complete/i, 1));
  }

  public override async start(): Promise<StartedMigrationContainer> {
    return new StartedMigrationContainer(await super.start());
  }
}

export class StartedMigrationContainer extends AbstractStartedContainer {
  constructor(startedTestContainer: StartedTestContainer) {
    super(startedTestContainer);
  }
}
