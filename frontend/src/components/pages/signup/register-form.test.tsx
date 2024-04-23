import { act, render, screen, waitFor } from "@/utils/vitest/utilities";
import { Mock } from "vitest";
import RegisterForm from "./register-form";
import { axe } from "jest-axe";
import userEvent from "@testing-library/user-event";

const getRegisterFormInputs = () => {
  const usernameInput = screen.getByLabelText(/username/i);
  const emailInput = screen.getByLabelText(/email/i);
  const passwordInput = screen.getByLabelText(/^senha/i);
  const confirmPasswordInput = screen.getByLabelText(/confirmar senha/i);

  return {
    usernameInput,
    emailInput,
    passwordInput,
    confirmPasswordInput,
  }
}

const renderRegisterForm = (fn: Mock<any, any>) => {
  const { user } = render(<RegisterForm submitForm={fn} />);
  const usernameInput = screen.getByLabelText(/username/i);
  const emailInput = screen.getByLabelText(/email/i);
  const passwordInput = screen.getByLabelText(/^senha/i);
  const confirmPasswordInput = screen.getByLabelText(/confirmar senha/i);
  const submitBtn = screen.getByRole("button");

  return {
    user,
    usernameInput,
    emailInput,
    passwordInput,
    confirmPasswordInput,
    submitBtn,
  };
};

const generateSubmitFn = (ms: number = 0, value: boolean = true) => {
  return vi.fn(async () => {
    await new Promise((resolve) => setTimeout(resolve, ms));
    
    if (!value) return;

    Object.values(getRegisterFormInputs()).forEach((input) => {
      userEvent.clear(input)
    })
  });
};

const validInputs = () => {
  const validUsername = "lgreidelas";
  const validEmail = "lgreidelas@gmail.com";
  const validPass = "12345678";

  return {
    validUsername,
    validEmail,
    validPass,
  };
};

describe("Register Form tests", () => {
  it("Should render a register form and fill with success", async () => {
    // Arrange
    const fn = generateSubmitFn();

    const {
      user,
      usernameInput,
      emailInput,
      passwordInput,
      confirmPasswordInput,
      submitBtn,
    } = renderRegisterForm(fn);
    const { validUsername, validEmail, validPass } = validInputs();

    // Act
    await user.type(usernameInput, validUsername);
    await user.type(emailInput, validEmail);
    await user.type(passwordInput, validPass);
    await user.type(confirmPasswordInput, validPass);
    await user.click(submitBtn);

    await waitFor(() => {
      expect(submitBtn).toBeEnabled();
    });

    // Assert
    await waitFor(() => {
      expect(usernameInput).toHaveValue("");
      expect(emailInput).toHaveValue("");
      expect(passwordInput).toHaveValue("");
      expect(confirmPasswordInput).toHaveValue("");
    })
    expect(fn).toHaveBeenCalledOnce();
  });

  it("Should not submit with invalid email", async () => {
    const fn = generateSubmitFn();

    const {
      user,
      usernameInput,
      emailInput,
      passwordInput,
      confirmPasswordInput,
      submitBtn,
    } = renderRegisterForm(fn);
    const { validUsername, validPass } = validInputs();
    const badEmail = "tubias@g";

    // Act
    await user.type(usernameInput, validUsername);
    await user.type(emailInput, badEmail);
    await user.type(passwordInput, validPass);
    await user.type(confirmPasswordInput, validPass);
    await user.click(submitBtn);

    // Assert
    expect(usernameInput).toHaveValue(validUsername);
    expect(emailInput).toHaveValue(badEmail);
    expect(passwordInput).toHaveValue(validPass);
    expect(confirmPasswordInput).toHaveValue(validPass);
    expect(emailInput).toHaveAccessibleErrorMessage(/insira um email valido/i);
    expect(fn).not.toHaveBeenCalled();
  });

  it("Should not submit with invalid username", async () => {
    const fn = generateSubmitFn();

    const {
      user,
      usernameInput,
      emailInput,
      passwordInput,
      confirmPasswordInput,
      submitBtn,
    } = renderRegisterForm(fn);
    const { validEmail, validPass } = validInputs();
    const badUsername = "gab";

    // Act
    await user.type(usernameInput, badUsername);
    await user.type(emailInput, validEmail);
    await user.type(passwordInput, validPass);
    await user.type(confirmPasswordInput, validPass);
    await user.click(submitBtn);

    // Assert
    expect(usernameInput).toHaveValue(badUsername);
    expect(emailInput).toHaveValue(validEmail);
    expect(passwordInput).toHaveValue(validPass);
    expect(confirmPasswordInput).toHaveValue(validPass);
    expect(usernameInput).toHaveAccessibleErrorMessage(/^nome deve ter/i);
    expect(fn).not.toHaveBeenCalled();
  });

  it("Should not submit with invalid passwords", async () => {
    const fn = generateSubmitFn();

    const {
      user,
      usernameInput,
      emailInput,
      passwordInput,
      confirmPasswordInput,
      submitBtn,
    } = renderRegisterForm(fn);
    const { validUsername, validEmail } = validInputs();
    const badPass = "gabri";

    // Act
    await user.type(usernameInput, validUsername);
    await user.type(emailInput, validEmail);
    await user.type(passwordInput, badPass);
    await user.type(confirmPasswordInput, badPass);
    await user.click(submitBtn);

    // Assert
    expect(usernameInput).toHaveValue(validUsername);
    expect(emailInput).toHaveValue(validEmail);
    expect(passwordInput).toHaveValue(badPass);
    expect(confirmPasswordInput).toHaveValue(badPass);
    expect(passwordInput).toHaveAccessibleErrorMessage(/^senha deve conter/i);
    expect(fn).not.toHaveBeenCalled();
  });

  it("Should not submit if passwords don't match", async () => {
    const fn = generateSubmitFn();

    const {
      user,
      usernameInput,
      emailInput,
      passwordInput,
      confirmPasswordInput,
      submitBtn,
    } = renderRegisterForm(fn);
    const { validPass, validUsername, validEmail } = validInputs();
    const badConfirm = "gabri1234";

    // Act
    await user.type(usernameInput, validUsername);
    await user.type(emailInput, validEmail);
    await user.type(passwordInput, validPass);
    await user.type(confirmPasswordInput, badConfirm);
    await user.click(submitBtn);

    // Assert
    expect(usernameInput).toHaveValue(validUsername);
    expect(emailInput).toHaveValue(validEmail);
    expect(passwordInput).toHaveValue(validPass);
    expect(confirmPasswordInput).toHaveValue(badConfirm);
    expect(confirmPasswordInput).toHaveAccessibleErrorMessage(
      /^senhas distintas/i,
    );
    expect(fn).not.toHaveBeenCalled();
  });

  it("Should not reset form if promise fails", async () => {
    const fn = generateSubmitFn(0, false);

    const {
      user,
      usernameInput,
      emailInput,
      passwordInput,
      confirmPasswordInput,
      submitBtn,
    } = renderRegisterForm(fn);
    const { validPass, validUsername, validEmail } = validInputs();

    // Act
    await user.type(usernameInput, validUsername);
    await user.type(emailInput, validEmail);
    await user.type(passwordInput, validPass);
    await user.type(confirmPasswordInput, validPass);
    await user.click(submitBtn);

    // Assert
    expect(usernameInput).toHaveValue(validUsername);
    expect(emailInput).toHaveValue(validEmail);
    expect(passwordInput).toHaveValue(validPass);
    expect(confirmPasswordInput).toHaveValue(validPass);
    expect(fn).toHaveBeenCalled();
  });

  it("Should disable btn until promises resolves", async () => {
    const fn = generateSubmitFn(100);

    const {
      user,
      usernameInput,
      emailInput,
      passwordInput,
      confirmPasswordInput,
      submitBtn,
    } = renderRegisterForm(fn);
    const { validPass, validUsername, validEmail } = validInputs();

    // Act
    await user.type(usernameInput, validUsername);
    await user.type(emailInput, validEmail);
    await user.type(passwordInput, validPass);
    await user.type(confirmPasswordInput, validPass);
    await user.click(submitBtn);

    // Assert
    expect(usernameInput).toHaveValue(validUsername);
    expect(emailInput).toHaveValue(validEmail);
    expect(passwordInput).toHaveValue(validPass);
    expect(confirmPasswordInput).toHaveValue(validPass);
    expect(fn).toHaveBeenCalledOnce();
    expect(submitBtn).toBeDisabled();

    await waitFor(() => {
      expect(submitBtn).toBeEnabled();
    });
  });

  it("Should be accessible", async () => {
    const fn = generateSubmitFn();
    const { container } = render(<RegisterForm submitForm={fn} />);

    // Act
    await act(async () => {
      const result = await axe(container);
      expect(result).toHaveNoViolations();
    });
  });
});
