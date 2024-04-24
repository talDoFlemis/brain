import { expect, test } from "@playwright/test";

test.describe("Sign In page tests", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/sign-in");
  });

  test("Should login and be redirected to dashboard", async ({ page }) => {
    const validUsername = "marcelomx30";
    const validPassword = "mypassword";

    await page.getByLabel(/username/i).fill(validUsername);
    await page.getByLabel(/senha/i).fill(validPassword);
    await page.getByRole("button", { name: "Login" }).click();

    await expect(page).toHaveURL(/dashboard/);
  });

  test("Should show an error if user submitted wrong password or email", async ({
    page,
  }) => {
    const validUsername = "marcelomx30";
    const wrongPassword = "12345678";

    const usernameInput = page.getByLabel(/username/i);
    const passwordInput = page.getByLabel(/senha/i);

    await usernameInput.fill(validUsername);
    await passwordInput.fill(wrongPassword);
    await page.getByRole("button", { name: "Login" }).click();

    console.log(await usernameInput.innerText());
    console.log(await passwordInput.innerText());

    await expect(page.locator("#username-error")).toHaveText(
      "Usuario ou senhas incorretas",
    );
    await expect(usernameInput).toHaveValue(validUsername);
    await expect(passwordInput).toHaveValue(wrongPassword);
    await expect(page).toHaveURL(/sign-in/);
  });

  test("Should click on register and go to sign up page", async ({ page }) => {
    await page.getByRole("link", { name: /registrar/i }).click();
    await expect(page).toHaveURL(/sign-up/);
  });

  test("Should click on forgot password and go to forgot password page", async ({
    page,
  }) => {
    await page.getByRole("link", { name: /esqueci/i }).click();
    await expect(page).toHaveURL(/forgot-password/);
  });
});
