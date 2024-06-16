import { z } from "zod";

export const timeoutSchema = z.union([
  z.literal("15"),
  z.literal("30"),
  z.literal("60"),
]);

export const baseQuestionSchema = z.object({
  title: z
    .string()
    .min(8, { message: "O título da questão deve ser no mínimo 8" }),
  timeout: timeoutSchema,
});
