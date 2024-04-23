import { act, render, screen, waitFor } from "@/utils/vitest/utilities";
import LoginForm from "./login-form";
import { Mock } from "vitest";
import { axe } from "jest-axe";

describe("Login Form Tests", () => {
  const renderLoginForm = (fn: Mock<any, any>) => {
    const { user } = render(<LoginForm submitForm={fn} />);
    const emailInput = screen.getByLabelText("Email");
    const passwordInput = screen.getByLabelText("Senha");
    const submitBtn = screen.getByRole("button");

    return { user, emailInput, passwordInput, submitBtn };
  };

  const sleep = (ms: number) =>
    new Promise((resolve) => setTimeout(resolve, ms));

  it("Should render the login form and fill the form", async () => {
    // Arrange
    const fn = vi.fn(async () => {
      await sleep(0);
      return true;
    });

    const { user, emailInput, passwordInput, submitBtn } = renderLoginForm(fn);

    // Act
    await user.type(emailInput, "marcelomx30@gmail.com");
    await user.type(passwordInput, "melikofornite123");
    await user.click(submitBtn);

    //Assert
    expect(emailInput).toHaveValue("");
    expect(passwordInput).toHaveValue("");
    expect(screen.getByText(/esqueci minha senha/i));
    expect(fn).toHaveBeenCalledOnce();
  });

  it("Should not submit with invalid email", async () => {
    // Arrange
    const fn = vi.fn(async () => {
      await sleep(0);
      return true;
    });
    const { user, emailInput, passwordInput, submitBtn } = renderLoginForm(fn);
    const badEmail = "marcelomx30@g";

    // Act
    await user.type(emailInput, badEmail);
    await user.type(passwordInput, "tubias");
    await user.click(submitBtn);

    //Assert
    expect(emailInput).toHaveValue(badEmail);
    expect(emailInput).toHaveAccessibleErrorMessage(/insira um email valido/i);
    expect(fn).not.toHaveBeenCalledOnce();
  });

  it("Should not reset form if promise fails", async () => {
    // Arrange
    const fn = vi.fn(async () => {
      await sleep(0);
      return false;
    });
    const { user, emailInput, passwordInput, submitBtn } = renderLoginForm(fn);
    const goodEmail = "marcelomx30@gmail.com";
    const goodPass = "melikofornite123";

    // Act
    await user.type(emailInput, goodEmail);
    await user.type(passwordInput, goodPass);
    await user.click(submitBtn);

    //Assert
    expect(emailInput).toHaveValue(goodEmail);
    expect(passwordInput).toHaveValue(goodPass);
    expect(submitBtn).not.toBeDisabled();
    expect(fn).toHaveBeenCalledOnce();
  });

  it("Should disable btn until promises resolves", async () => {
    const fn = vi.fn(async () => {
      await sleep(100);
      return false;
    });
    const { user, emailInput, passwordInput, submitBtn } = renderLoginForm(fn);
    const goodEmail = "marcelomx30@gmail.com";
    const goodPass = "melikofornite123";

    // Act
    await user.type(emailInput, goodEmail);
    await user.type(passwordInput, goodPass);
    await user.click(submitBtn);

    //Assert
    await waitFor(() => {
      expect(submitBtn).not.toBeDisabled();
    });

    expect(emailInput).toHaveValue(goodEmail);
    expect(passwordInput).toHaveValue(goodPass);
    expect(submitBtn).not.toBeDisabled();
  });

  it("Should be accessible", async () => {
    const fn = vi.fn(async () => {
      await sleep(100);
      return false;
    });

    const { container } = render(<LoginForm submitForm={fn} />);
    await act(async () => {
      const result = await axe(container);
      expect(result).toHaveNoViolations();
    });
  });
});
