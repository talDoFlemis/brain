"use client";

import React from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import GeneralInfoSection from "./general-info-section";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form } from "@/components/ui/form";
import Questions, { questionsSchema } from "./questions";
import { createGame } from "./actions";

const formSchema = z.object({
  general: z.object({
    title: z
      .string()
      .min(8, { message: "O título não pode ter menos de 8 caracteres" }),
    description: z.string().optional(),
    tags: z.string().array().optional(),
  }),
  questions: questionsSchema,
});

export type CreateGameFormSchema = z.infer<typeof formSchema>;

function CreateQuizForm() {
  const form = useForm<CreateGameFormSchema>({
    resolver: zodResolver(formSchema),
  });

  return (
    <Form {...form}>
      <form
        action={form.handleSubmit(createGame)}
        className="flex flex-col space-y-4"
      >
        <GeneralInfoSection form={form} />
        <Questions form={form} />
      </form>
    </Form>
  );
}

export default CreateQuizForm;
