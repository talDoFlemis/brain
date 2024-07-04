import { z } from "zod";
import { baseQuestionSchema } from "./base-question";
import { UseFormReturn, useWatch } from "react-hook-form";
import { CreateGameFormSchema } from "./create-quiz-form";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";

export const quizQuestionSchema = z
  .object({
    kind: z.literal("quiz"),
    alternatives: z
      .object({
        data: z
          .string()
          .min(1, { message: "A resposta deve ter no mínimo um caractere" }),
        is_correct: z.boolean(),
      })
      .array()
      .length(4),
  })
  .merge(baseQuestionSchema);

export type QuizQuestionSchema = z.infer<typeof quizQuestionSchema>;

export const defaultQuizQuestion: QuizQuestionSchema = {
  kind: "quiz",
  alternatives: [
    { data: "Alternativa 1", is_correct: true },
    { data: "Alternativa 2", is_correct: false },
    { data: "Alternativa 3", is_correct: false },
    { data: "Alternativa 4", is_correct: false },
  ],
  title: "Multipla escolha",
  timeout: "15",
};

type QuizQuestionProps = {
  form: UseFormReturn<CreateGameFormSchema>;
  index: number;
};

function QuizQuestion({ form, index }: QuizQuestionProps) {
  const alternatives = useWatch({
    control: form.control,
    name: `questions.${index}.alternatives`,
  });
  return (
    <ul className="flex flex-col gap-4">
      {alternatives.map((_, altIndex) => (
        <li
          key={`questions-${index}-alternatives-${altIndex}`}
          className="flex flex-row items-start gap-4"
        >
          <FormField
            control={form.control}
            name={`questions.${index}.alternatives.${altIndex}.is_correct`}
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Checkbox
                    checked={field.value}
                    onCheckedChange={field.onChange}
                  />
                </FormControl>
              </FormItem>
            )}
          />
          <div className="flex flex-col w-full space-y-3">
            <FormLabel
              htmlFor={`questions-${index}-alternatives-${altIndex}-answer`}
              className="text-card-foreground"
            >
              Alternativa {altIndex + 1}
            </FormLabel>
            <FormField
              control={form.control}
              name={`questions.${index}.alternatives.${altIndex}.data`}
              render={({ field }) => (
                <FormItem className="w-full">
                  <FormControl>
                    <Input
                      className="text-card-foreground bg-foreground border-accent"
                      type="text"
                      id={`questions-${index}-alternatives-${altIndex}-data`}
                      aria-invalid={!!form.formState.errors.questions}
                      aria-errormessage={`questions-${index}-alternatives-${altIndex}-data-error`}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage
                    id={`questions-${index}-alternatives-${altIndex}-data-error`}
                  />
                </FormItem>
              )}
            />
          </div>
        </li>
      ))}
    </ul>
  );
}

export default QuizQuestion;
