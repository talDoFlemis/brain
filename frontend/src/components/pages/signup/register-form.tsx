"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
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

const formSchema = z
  .object({
    username: z
      .string()
      .min(5, { message: "Nome deve ter conter 5 ou mais caracteres" })
      .max(50, { message: "Nome deve ter no maximo 50 caracteres" }),
    email: z.string().email({ message: "Insira um email valido" }),
    password: z
      .string()
      .min(8, { message: "Senha deve conter 8 ou mais caracteres" }),
    confirm: z
      .string()
      .min(8, { message: "Senha deve conter 8 ou mais caracteres" }),
  })
  .refine((data) => data.password === data.confirm, {
    message: "Senhas distintas",
    path: ["confirm"],
  });

export type RegisterFormSchema = z.infer<typeof formSchema>;

export type RegisterFormProps = {
  submitForm: (v: RegisterFormSchema) => Promise<boolean>;
};

function RegisterForm({ submitForm }: RegisterFormProps) {
  const form = useForm<RegisterFormSchema>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "",
      email: "",
      password: "",
      confirm: "",
    },
  });

  async function onSubmit(values: RegisterFormSchema) {
    const isOk = await submitForm(values);
    if (isOk) form.reset();
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-3">
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
                  aria-invalid={!!form.formState.errors.username}
                  aria-errormessage="username-error"
                  placeholder="marcelo jr"
                  {...field}
                />
              </FormControl>
              <FormMessage id="username-error" />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel htmlFor="email-input" className="text-foreground">
                Email
              </FormLabel>
              <FormControl>
                <Input
                  className="text-foreground"
                  type="email"
                  id="email-input"
                  aria-invalid={!!form.formState.errors.email}
                  aria-errormessage="email-error"
                  placeholder="marcelo@example.com"
                  {...field}
                />
              </FormControl>
              <FormMessage id="email-error" />
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
                  aria-invalid={!!form.formState.errors.password}
                  aria-errormessage="password-error"
                  placeholder="*****"
                  {...field}
                />
              </FormControl>
              <FormMessage id="password-error" />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="confirm"
          render={({ field }) => (
            <FormItem className="pb-3">
              <FormLabel
                htmlFor="password-confirm-input"
                className="text-foreground"
              >
                Confirmar senha
              </FormLabel>
              <FormControl>
                <Input
                  id="password-confirm-input"
                  className="text-foreground"
                  type="password"
                  placeholder="*****"
                  aria-invalid={!!form.formState.errors.confirm}
                  aria-errormessage="confirm-error"
                  {...field}
                />
              </FormControl>
              <FormMessage id="confirm-error" />
            </FormItem>
          )}
        />
        <Button
          disabled={form.formState.isSubmitting}
          size="full"
          type="submit"
        >
          Registrar
        </Button>
      </form>
    </Form>
  );
}

export default RegisterForm;
