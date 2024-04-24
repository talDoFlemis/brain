"use client";

import Image from "next/image";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import LoginForm, { LoginFormSchema, UseForm } from "./login-form";
import { signIn } from "next-auth/react";
import { useRouter } from "next/navigation";
import { SIGN_IN_CALLBACK_URL } from "@/utils/constants";

function LoginSection() {
  const router = useRouter();

  const submitForm = async (values: LoginFormSchema, form: UseForm) => {
    const response = await signIn("credentials", {
      username: values.username,
      password: values.password,
      redirect: false,
    });

    if (response?.error) {
      form.setError("username", {
        type: "validate",
        message: "Usuario ou senhas incorretas",
      });
      return;
    }

    form.reset();
    router.push(SIGN_IN_CALLBACK_URL);
  };

  return (
    <section className="col-span-4 lg:col-span-2 justify-self-center flex flex-col w-full py-12 px-4 gap-10 max-w-lg">
      <Image
        className="self-center"
        src="/brain-logo-white.svg"
        alt="brain.test logo"
        width={175}
        height={50}
        priority
      />
      <div className="self-center flex flex-col items-center gap-4">
        <h2 className="text-3xl text-foreground font-semibold">
          Logue em uma conta
        </h2>
        <h3 className="text-sm text-foreground/50 font-semibold">
          Digite suas credenciais
        </h3>
      </div>
      <LoginForm submitForm={submitForm} />
      <p className="text-sm text-foreground">
        Não tem uma conta ainda?
        <Button variant="link" size="sm" asChild>
          <Link href="/sign-up">Registrar</Link>
        </Button>
      </p>
    </section>
  );
}

export default LoginSection;
