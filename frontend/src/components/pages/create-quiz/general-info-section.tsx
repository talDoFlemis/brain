"use client";

import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import React from "react";
import { UseFormReturn } from "react-hook-form";
import { CreateGameFormSchema } from "./create-quiz-form";
import { Button } from "@/components/ui/button";

function GeneralInfoSection({
  form,
}: {
  form: UseFormReturn<CreateGameFormSchema>;
}) {
  return (
    <div className="flex flex-col gap-8">
      <div className="flex flex-row justify-between items-center">
        <h1 className="text-4xl text-secondary-foreground">Criar Quiz</h1>
        <Button type="submit">Criar</Button>
      </div>
      <div className="flex flex-col gap-4 bg-foreground p-4 rounded-lg">
        <FormField
          control={form.control}
          name="general.title"
          render={({ field }) => (
            <FormItem>
              <FormLabel
                htmlFor="general-title-input"
                className="text-card-foreground"
              >
                Título
              </FormLabel>
              <FormControl>
                <Input
                  className="text-card-foreground dark:bg-[#27262c] border-none"
                  type="text"
                  id="general-title-input"
                  aria-invalid={!!form.formState.errors.general?.title}
                  aria-errormessage="general-title-error"
                  {...field}
                />
              </FormControl>
              <FormMessage id="general-title-error" />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="general.description"
          render={({ field }) => (
            <FormItem>
              <FormLabel
                htmlFor="general-description-input"
                className="text-card-foreground"
              >
                Descrição
              </FormLabel>
              <FormControl>
                <Input
                  className="text-card-foreground dark:bg-[#27262c] border-none"
                  type="text"
                  id="general-description-input"
                  aria-invalid={!!form.formState.errors.general?.description}
                  aria-errormessage="general-description-error"
                  {...field}
                />
              </FormControl>
              <FormMessage id="general-description-error" />
            </FormItem>
          )}
        />
      </div>
    </div>
  );
}

export default GeneralInfoSection;
