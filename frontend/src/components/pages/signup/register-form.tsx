"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import authService from "@/services/auth-service";

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
import { signIn } from "next-auth/react";
import { SIGN_IN_CALLBACK_URL } from "@/utils/constants";

const formSchema = z.object({
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
});

type RegisterFormSchema = z.infer<typeof formSchema>;

function RegisterForm() {
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
    if (values.password !== values.confirm) {
      form.setError("confirm", {
        type: "validate",
        message: "Senhas distintas",
      });
      return;
    }

    try {
      await authService.signUp({
        username: values.username,
        email: values.email,
        password: values.password,
      });
      
      // Signs the user in, and redirects him to the home page
      await signIn("credentials", {
        username: values.username,
        password: values.password,
        callbackUrl: SIGN_IN_CALLBACK_URL,
      });
    } catch (error) {
      // TODO: Error handling
      console.log(error);
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-3">
        <FormField
          control={form.control}
          name="username"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="text-foreground">Username</FormLabel>
              <FormControl>
                <Input
                  className="text-foreground"
                  type="text"
                  placeholder="marcelo jr"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="text-foreground">Email</FormLabel>
              <FormControl>
                <Input
                  className="text-foreground"
                  type="email"
                  placeholder="marcelo@example.com"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="text-foreground">Senha</FormLabel>
              <FormControl>
                <Input
                  className="text-foreground"
                  type="password"
                  placeholder="*****"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="confirm"
          render={({ field }) => (
            <FormItem className="pb-3">
              <FormLabel className="text-foreground">Confirmar senha</FormLabel>
              <FormControl>
                <Input
                  className="text-foreground"
                  type="password"
                  placeholder="*****"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button size="full" type="submit">
          Registrar
        </Button>
      </form>
    </Form>
  );
}

export default RegisterForm;
