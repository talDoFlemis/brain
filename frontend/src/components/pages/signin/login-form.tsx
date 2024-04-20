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
import { signIn } from "next-auth/react";
import { SIGN_IN_CALLBACK_URL } from "@/utils/constants";
import { useRouter } from "next/navigation";

const formSchema = z.object({
  username: z
    .string()
    .min(5, { message: "Nome deve ter conter 5 ou mais caracteres" }),
  password: z
    .string()
    .min(8, { message: "Senha deve conter 8 ou mais caracteres" }),
});

type LoginFormSchema = z.infer<typeof formSchema>;

function LoginForm() {
  const form = useForm<LoginFormSchema>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "",
      password: "",
    },
  });

  const router = useRouter();

  async function onSubmit(values: LoginFormSchema) {
    const response = await signIn("credentials", {
      username: values.username,
      password: values.password,
      redirect: false,
      callbackUrl: SIGN_IN_CALLBACK_URL,
    });

    if (response?.error) {
      form.setError("username", {
        type: "validate",
        message: "Usuario ou senhas incorretas",
      });
      return;
    }

    router.push(SIGN_IN_CALLBACK_URL);
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
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

        <Button className="text-xs px-0" variant="link" size="sm" asChild>
          <Link href="#">Esqueci minha senha</Link>
        </Button>
        <Button size="full" type="submit">
          Login
        </Button>
      </form>
    </Form>
  );
}

export default LoginForm;
