import { expect } from "vitest";
import "@testing-library/jest-dom";
import * as extendedMatchers from "jest-extended";
import { toHaveNoViolations } from "jest-axe";

expect.extend(extendedMatchers);
expect.extend(toHaveNoViolations);
