import { test, expect } from "@playwright/test";

test.describe("Sign Up page tests", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/sign-up");
  });

  test("Should sign up and be redirected to dashboard", async ({ page }) => {
    const validEmail = "lgreidelas@gmail.com";
    const validUsername = "lgreidelas";
    const validPassword = "mypassword";

    await page.getByLabel(/Username/).fill(validUsername);
    await page.getByLabel(/email/i).fill(validEmail);
    await page.getByLabel(/Senha/).fill(validPassword);
    await page.getByLabel(/confirmar/i).fill(validPassword);
    await page.getByRole("button").click();

    await expect(page).toHaveURL(/dashboard/);
  });

  test("Should show an error if user submitted wrong password or email", async ({
    page,
  }) => {
    const username = "marcelomx30";
    const usedEmail = "marcelomx30@gmail.com";
    const password = "marcelinhogameplays";

    const usernameInput = page.locator("#username-input");
    const emailInput = page.locator("#email-input");
    const passwordInput = page.locator("#password-input");
    const confirmPasswordInput = page.locator("#password-confirm-input");

    await usernameInput.fill(username);
    await emailInput.fill(usedEmail);
    await passwordInput.fill(password);
    await confirmPasswordInput.fill(password);
    await page.getByRole("button").click();

    await expect(page.locator("#username-error")).toHaveText(
      "Nome de usuario ja existe",
    );
    await expect(page).toHaveURL(/sign-up/);
  });

  test("Should click on sign in and go to sign in page", async ({ page }) => {
    await page.getByRole("link", { name: /logar/i }).click();
    await expect(page).toHaveURL(/sign-in/);
  });
});
