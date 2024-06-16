import { z } from "zod";
import { baseQuestionSchema } from "./base-question";
import { UseFormReturn } from "react-hook-form";
import { CreateGameFormSchema } from "./create-quiz-form";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";

export const trueFalseQuestionSchema = z
  .object({
    kind: z.literal("true_false"),
    trueAlternative: z
      .string()
      .min(1, "A alternative correta deve ter pelo menos um caractere"),
    falseAlternative: z
      .string()
      .min(1, "A alternative falsa deve ter pelo menos um caractere"),
  })
  .merge(baseQuestionSchema);

export type TrueFalseQuestionSchema = z.infer<typeof trueFalseQuestionSchema>;

export const defaultTrueFalseQuestion: TrueFalseQuestionSchema = {
  kind: "true_false",
  title: "Verdadeiro ou Falso",
  timeout: "60",
  trueAlternative: "Verdade",
  falseAlternative: "Falso",
};

type TrueFalseQuestionProps = {
  form: UseFormReturn<CreateGameFormSchema>;
  index: number;
};

function TrueFalseQuestion({ form, index }: TrueFalseQuestionProps) {
  return (
    <>
      <FormField
        control={form.control}
        name={`questions.${index}.trueAlternative`}
        render={({ field }) => (
          <FormItem>
            <FormLabel
              htmlFor={`question-${index}-trueAlternative`}
              className="text-card-foreground"
            >
              Alternativa Verdadeira
            </FormLabel>
            <FormControl>
              <Input
                className="text-card-foreground bg-foreground border-accent"
                type="text"
                id={`question-${index}-trueAlternative`}
                aria-invalid={!!form.formState.errors.questions}
                aria-errormessage={`question-${index}-trueAlternative-error`}
                {...field}
              />
            </FormControl>
            <FormMessage id={`question-${index}-trueAlternative-error`} />
          </FormItem>
        )}
      />
      <FormField
        control={form.control}
        name={`questions.${index}.falseAlternative`}
        render={({ field }) => (
          <FormItem>
            <FormLabel
              htmlFor={`question-${index}-falseAlternative`}
              className="text-card-foreground"
            >
              Alternativa Falsa
            </FormLabel>
            <FormControl>
              <Input
                className="text-card-foreground bg-foreground border-accent"
                type="text"
                id={`question-${index}-falseAlternative`}
                aria-invalid={!!form.formState.errors.questions}
                aria-errormessage={`question-${index}-falseAlternative-error`}
                {...field}
              />
            </FormControl>
            <FormMessage id={`question-${index}-falseAlternative-error`} />
          </FormItem>
        )}
      />
    </>
  );
}

export default TrueFalseQuestion;
