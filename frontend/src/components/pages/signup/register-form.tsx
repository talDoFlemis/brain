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

const formSchema = z.object({
  username: z
    .string()
    .min(5, { message: "Nome deve ter conter 5 ou mais caracteres" })
    .max(50, { message: "Nome deve ter no maximo 50 caracteres" }),
  email: z.string().email({ message: "Insira um email valido" }),
  password: z
    .string()
    .min(5, { message: "Senha deve conter 5 ou mais caracteres" }),
  confirm: z
    .string()
    .min(5, { message: "Senha deve conter 5 ou mais caracteres" }),
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

  function onSubmit(values: RegisterFormSchema) {
    if (values.password !== values.confirm) {
      form.setError("confirm", {
        type: "validate",
        message: "Senhas distintas",
      });
      return;
    }

    console.log(values);
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
