import {
  PostgreSqlContainer,
  StartedPostgreSqlContainer,
} from "@testcontainers/postgresql";
import { MigrationContainer } from "./migrate";
import { SeedContainer } from "./seed";

export function createNewPgContainer(
  username: string,
  password: string,
  database: string,
): PostgreSqlContainer {
  const container = new PostgreSqlContainer()
    .withUsername(username)
    .withPassword(password)
    .withDatabase(database);
  return container;
}

export async function migrateAndSeedPgContainer(
  pathToBackend: string,
  pgContainer: StartedPostgreSqlContainer,
): Promise<void> {
  const migrationContainer = new MigrationContainer(pathToBackend, pgContainer);
  const startedMigration = await migrationContainer.start();
  // https://github.com/testcontainers/testcontainers-node/pull/730 is not available in the current version
  // Wait for the migration to complete before seeding
  await startedMigration.stop();

  const seedContainer = new SeedContainer(pathToBackend, pgContainer);
  const startedSeed = await seedContainer.start();
  await startedSeed.stop();
}
