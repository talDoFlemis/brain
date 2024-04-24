import { act, render, screen, waitFor } from "@/utils/vitest/utilities";
import LoginForm from "./login-form";
import { Mock } from "vitest";
import { axe } from "jest-axe";

describe("Login Form Tests", () => {
  const renderLoginForm = (fn: Mock<any, any>) => {
    const { user } = render(<LoginForm submitForm={fn} />);
    const userNameInput = screen.getByLabelText("Username");
    const passwordInput = screen.getByLabelText("Senha");
    const submitBtn = screen.getByRole("button");

    return { user, userNameInput, passwordInput, submitBtn };
  };

  const sleep = (ms: number) =>
    new Promise((resolve) => setTimeout(resolve, ms));

  it("Should render the login form and fill the form", async () => {
    // Arrange
    const fn = vi.fn(async () => {
      await sleep(0);
    });

    const { user, userNameInput, passwordInput, submitBtn } =
      renderLoginForm(fn);

    // Act
    await user.type(userNameInput, "marcelo jr");
    await user.type(passwordInput, "melikofornite123");
    await user.click(submitBtn);

    //Assert
    expect(userNameInput).toHaveValue("marcelo jr")
    expect(passwordInput).toHaveValue("melikofornite123")
    expect(screen.getByText(/esqueci minha senha/i));
    expect(fn).toHaveBeenCalledOnce();
  });

  it("Should not submit with invalid username", async () => {
    // Arrange
    const fn = vi.fn(async () => {
      await sleep(0);
    });
    const { user, userNameInput, passwordInput, submitBtn } =
      renderLoginForm(fn);
    const badUserName = "mx30";

    // Act
    await user.type(userNameInput, badUserName);
    await user.type(passwordInput, "tubias");
    await user.click(submitBtn);

    //Assert
    expect(userNameInput).toHaveValue(badUserName);
    expect(userNameInput).toHaveAccessibleErrorMessage(
      /Nome deve ter conter 5 ou mais caracteres/i,
    );
    expect(fn).not.toHaveBeenCalledOnce();
  });

  it("Should not submit with invalid password", async () => {
    // Arrange 
    const fn = vi.fn(async () => {
      await sleep(0);
    });
    const { user, userNameInput, passwordInput, submitBtn } =
      renderLoginForm(fn);
    const badPassword = "12345"
    // Act 
    await user.type(userNameInput, "marcelo jr");
    await user.type(passwordInput, badPassword);
    await user.click(submitBtn);
    
    // Assert
    expect(userNameInput).toHaveValue("marcelo jr")
    expect(passwordInput).toHaveValue(badPassword)
    expect(passwordInput).toHaveAccessibleErrorMessage(
      /Senha deve conter 8 ou mais caracteres/i,
    );
    expect(fn).not.toHaveBeenCalledOnce();
  })

  it("Should not reset form if fails", async () => {
    // Arrange
    const fn = vi.fn(async () => {
      await sleep(0);
      return;
    });
    const { user, userNameInput, passwordInput, submitBtn } =
      renderLoginForm(fn);
    const goodUserName = "marcelo jr";
    const goodPass = "melikofornite123";

    // Act
    await user.type(userNameInput, goodUserName);
    await user.type(passwordInput, goodPass);
    await user.click(submitBtn);

    //Assert
    expect(userNameInput).toHaveValue(goodUserName);
    expect(passwordInput).toHaveValue(goodPass);
    expect(submitBtn).not.toBeDisabled();
    expect(fn).toHaveBeenCalledOnce();
  });

  it("Should disable btn until promises resolves", async () => {
    const fn = vi.fn(async () => {
      await sleep(100);
    });
    const { user, userNameInput, passwordInput, submitBtn } =
      renderLoginForm(fn);
    const goodUserName = "marcelo jr";
    const goodPass = "melikofornite123";

    // Act
    await user.type(userNameInput, goodUserName);
    await user.type(passwordInput, goodPass);
    await user.click(submitBtn);

    //Assert
    await waitFor(() => {
      expect(submitBtn).not.toBeDisabled();
    });

    expect(userNameInput).toHaveValue(goodUserName);
    expect(passwordInput).toHaveValue(goodPass);
    expect(submitBtn).not.toBeDisabled();
  });

  it("Should be accessible", async () => {
    const fn = vi.fn(async () => {
      await sleep(100);
    });

    const { container } = render(<LoginForm submitForm={fn} />);
    await act(async () => {
      const result = await axe(container);
      expect(result).toHaveNoViolations();
    });
  });
});
