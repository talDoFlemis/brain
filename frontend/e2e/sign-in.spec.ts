import { test, expect } from "@playwright/test";

test.describe("Sign In page tests", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/sign-in");
  });

  test.fixme(
    "Should login and be redirected to dashboard",
    async ({ page }) => {
      const validEmail = "marcelomx30@gmail.com";
      const validPassword = "mypassword";

      await page.getByLabel(/email/i).fill(validEmail);
      await page.getByLabel(/senha/i).fill(validPassword);
      await page.getByRole("button").click();

      await expect(page).toHaveURL(/dashboard/);
    },
  );

  test.fixme(
    "Should show an error if user submitted wrong password or email",
    async ({ page }) => {
      const validEmail = "marcelomx30@gmail.com";
      const wrongPassword = "12345678";

      const emailBtn = page.getByLabel(/email/i);
      const passwordBtn = page.getByLabel(/senha/i);

      await emailBtn.fill(validEmail);
      await passwordBtn.fill(wrongPassword);
      await page.getByRole("button").click();

      await expect(emailBtn).toHaveText(validEmail);
      await expect(passwordBtn).toHaveText(wrongPassword);
      await expect(page).toHaveURL(/sign-in/);
    },
  );

  test("Should click on register and go to sign up page", async ({ page }) => {
    await page.getByRole("link", { name: /registrar/i }).click();
    await expect(page).toHaveURL(/sign-up/);
  });

  test.skip("Should click on forgot password and go to forgot password page", async ({
    page,
  }) => {
    await page.getByRole("link", { name: /esqueci/i }).click();
    await expect(page).toHaveURL(/forgot-password/);
  });
});
