import { test, expect } from "@playwright/test";

test.describe("Sign Up page tests", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/sign-up");
  });

  test.fixme(
    "Should sign up and be redirected to dashboard",
    async ({ page }) => {
      const validEmail = "lgreidelas@gmail.com";
      const validUsername = "lgreidelas";
      const validPassword = "mypassword";

      await page.getByLabel(/username/i).fill(validUsername);
      await page.getByLabel(/email/i).fill(validEmail);
      await page.getByLabel(/senha/i).fill(validPassword);
      await page.getByLabel(/confirmar/i).fill(validPassword);
      await page.getByRole("button").click();

      await expect(page).toHaveURL(/dashboard/);
    },
  );

  test.fixme(
    "Should show an error if user submitted wrong password or email",
    async ({ page }) => {
      const usedEmail = "marcelomx30@gmail.com";
      const username = "marcelomx30";
      const password = "marcelinhogameplays";

      const usernameBtn = page.getByLabel(/username/i);
      const emailBtn = page.getByLabel(/email/i);
      const passwordBtn = page.getByLabel(/senha/i);
      const confirmPasswordBtn = page.getByLabel(/confirmar/i);

      await usernameBtn.fill(username);
      await emailBtn.fill(usedEmail);
      await passwordBtn.fill(password);
      await confirmPasswordBtn.fill(password);
      await page.getByRole("button").click();

      await expect(usernameBtn).toHaveText(username);
      await expect(emailBtn).toHaveText(usedEmail);
      await expect(passwordBtn).toHaveText(password);
      await expect(confirmPasswordBtn).toHaveText(password);
      await expect(page).toHaveURL(/sign-up/);
    },
  );

  test("Should click on sign in and go to sign in page", async ({ page }) => {
    await page.getByRole("link", { name: /logar/i }).click();
    await expect(page).toHaveURL(/sign-in/);
  });
});
