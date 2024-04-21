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
import Link from "next/link";

const formSchema = z.object({
  email: z.string().email({ message: "Insira um email valido" }),
  password: z.string(),
});

export type LoginFormSchema = z.infer<typeof formSchema>;

export type LoginFormProps = {
  submitForm: (v: LoginFormSchema) => Promise<boolean>;
};

function LoginForm({ submitForm }: LoginFormProps) {
  const form = useForm<LoginFormSchema>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  async function onSubmit(values: LoginFormSchema) {
    const isOk = await submitForm(values);
    if (isOk) form.reset();
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
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
                  placeholder="marcelo@example.com"
                  aria-invalid={!!form.formState.errors.email}
                  aria-errormessage="email-error"
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
                  placeholder="*****"
                  {...field}
                />
              </FormControl>
              <FormMessage />
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
