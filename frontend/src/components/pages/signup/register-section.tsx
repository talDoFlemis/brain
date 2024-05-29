"use client";

import Image from "next/image";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import RegisterForm, { RegisterFormSchema, UseForm } from "./register-form";
import authService from "@/services/auth-service";
import { signIn } from "next-auth/react";
import { SIGN_IN_CALLBACK_URL } from "@/utils/constants";
import { isAxiosError } from "axios";

function RegisterSection() {
  async function submitForm(values: RegisterFormSchema, form: UseForm) {
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
      if (isAxiosError(error) && error.response?.status == 409) {
        form.setError("username", {
          type: "validate",
          message: "Nome de usuario ja existe",
        });
      }
      console.log(error);
    }
  }

  return (
    <section className="col-span-4 lg:col-span-2 justify-self-center flex flex-col w-full py-12 px-4 gap-8 max-w-lg">
      <Image
        className="self-center"
        src="/brain-logo-white.svg"
        alt="brain.test logo"
        width={175}
        height={50}
        priority
      />
      <div className="self-center flex flex-col items-center gap-4">
        <h2 className="text-3xl text-card-foreground font-semibold">
          Registrar uma conta
        </h2>
        <h3 className="text-sm text-card-foreground/50 font-semibold">
          Digite suas credenciais
        </h3>
      </div>
      <RegisterForm submitForm={submitForm} />
      <p className="text-sm text-card-foreground">
        Ja possui uma conta?
        <Button className="text-accent" variant="link" size="sm" asChild>
          <Link href="/sign-in">Logar</Link>
        </Button>
      </p>
    </section>
  );
}

export default RegisterSection;
