import { render, screen } from "@/utils/vitest/utilities";
import LoginSection from "./login-section";

describe("Login Section Tests", () => {
  it("Should render the login section", async () => {
    // Arrange
    render(<LoginSection />);

    // Act
    const logHeader = screen.getByRole("heading", { level: 2 });
    const credsHeader = screen.getByRole("heading", { level: 3 });
    const noPassParagraph = screen.getByRole("paragraph");
    const registerBtn = screen.getByRole("link", { name: /registrar/i });
    screen.getByAltText(/logo/i);

    //Assert
    expect(logHeader).toHaveTextContent("Logue");
    expect(credsHeader).toHaveTextContent("Digite suas credenciais");
    expect(noPassParagraph).toHaveTextContent("ainda?");
    expect(registerBtn).not.toBeDisabled();
  });
});
