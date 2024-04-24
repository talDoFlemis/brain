"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm, UseFormReturn } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import Link from "next/link";

const formSchema = z.object({
  username: z
    .string()
    .min(5, { message: "Nome deve ter conter 5 ou mais caracteres" }),
  password: z
    .string()
    .min(8, { message: "Senha deve conter 8 ou mais caracteres" }),
});

export type LoginFormSchema = z.infer<typeof formSchema>;
export type UseForm = UseFormReturn<LoginFormSchema>;

export type LoginFormProps = {
  submitForm: (values: LoginFormSchema, form: UseForm) => Promise<void>;
};

function LoginForm({ submitForm }: LoginFormProps) {
  const form = useForm<LoginFormSchema>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "",
      password: "",
    },
  });

  async function onSubmit(values: LoginFormSchema) {
    await submitForm(values, form);
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="username"
          render={({ field }) => (
            <FormItem>
              <FormLabel htmlFor="username-input" className="text-foreground">
                Username
              </FormLabel>
              <FormControl>
                <Input
                  className="text-foreground"
                  type="text"
                  id="username-input"
                  placeholder="marcelo jr"
                  aria-invalid={!!form.formState.errors.username}
                  aria-errormessage="username-error"
                  {...field}
                />
              </FormControl>
              <FormMessage id="username-error" />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <FormLabel htmlFor="password-input" className="text-foreground">
                Senha
              </FormLabel>
              <FormControl>
                <Input
                  className="text-foreground"
                  type="password"
                  id="password-input"
                  placeholder="*****"
                  aria-invalid={!!form.formState.errors.password}
                  aria-errormessage="password-error"
                  {...field}
                />
              </FormControl>
              <FormMessage id="password-error" />
            </FormItem>
          )}
        />

        <Button className="text-xs px-0" variant="link" size="sm" asChild>
          <Link href="#">Esqueci minha senha</Link>
        </Button>
        <Button
          disabled={form.formState.isSubmitting}
          size="full"
          type="submit"
        >
          Login
        </Button>
      </form>
    </Form>
  );
}

export default LoginForm;
