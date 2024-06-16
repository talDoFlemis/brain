"use client";

import { UseFormReturn, useFieldArray } from "react-hook-form";
import { z } from "zod";
import { CreateGameFormSchema } from "./create-quiz-form";
import { quizQuestionSchema } from "./quiz-question";
import {
  defaultTrueFalseQuestion,
  trueFalseQuestionSchema,
} from "./true-false-question";
import { FaPlus } from "react-icons/fa6";
import QuestionFactory from "./question-factory";

export const questionSchema = z.discriminatedUnion("kind", [
  quizQuestionSchema,
  trueFalseQuestionSchema,
]);

export type QuestionSchema = z.infer<typeof questionSchema>;

export const questionsSchema = questionSchema
  .array()
  .min(1, { message: "Deve existir ao mínimo uma questão" });

export type QuestionsSchema = z.infer<typeof questionsSchema>;

type QuestionsProps = {
  form: UseFormReturn<CreateGameFormSchema>;
};

// TODO: make green a variable
function Questions({ form }: QuestionsProps) {
  const {
    fields: questions,
    append,
    remove,
    update,
  } = useFieldArray({
    control: form.control,
    name: "questions",
  });
  return (
    <div className="flex flex-col gap-8">
      <ul className="flex flex-col gap-8">
        {questions.map((_, index) => (
          <QuestionFactory
            key={`question-${index}`}
            form={form}
            index={index}
            remove={remove}
            update={update}
          />
        ))}
      </ul>
      <div
        role="presentation"
        className="flex items-center justify-center p-4 bg-foreground rounded-lg cursor-pointer group"
        onClick={() => append(defaultTrueFalseQuestion)}
      >
        <FaPlus className="text-white text-2xl group-hover:text-secondary" />
      </div>
    </div>
  );
}

export default Questions;
