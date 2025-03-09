import React from "react";
import { useForm, Controller } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Input } from "./ui/input";
import { Button } from "./ui/button";

const schema = z.object({
  email: z.string().email("Invalid email format"),
  password: z.string().min(1, "Password is required"),
});

type LoginRequest = z.infer<typeof schema>;

const LoginForm: React.FC = () => {
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginRequest>({
    resolver: zodResolver(schema),
  });

  const onSubmit = (data: LoginRequest) => {
    console.log(data);
  };

  const handleGoogleLogin = () => {
    const currentUrl = window.location.href;
    const googleLoginUrl = `${import.meta.env.VITE_API_URL}/auth/github/login?from=${currentUrl}`;
    window.location.href = googleLoginUrl;
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <div className="mb-4">
        <label htmlFor="email" className="block text-sm font-medium text-gray-700">
          Email
        </label>
        <Controller
          name="email"
          control={control}
          defaultValue=""
          render={({ field }) => (
            <Input
              {...field}
              type="email"
              id="email"
              placeholder="Enter your email"
              className={`mt-1 block w-full ${errors.email ? "border-red-500" : ""}`}
            />
          )}
        />
        {errors.email && <p className="text-red-500 text-sm mt-1">{errors.email.message}</p>}
      </div>
      <div className="mb-6">
        <label htmlFor="password" className="block text-sm font-medium text-gray-700">
          Password
        </label>
        <Controller
          name="password"
          control={control}
          defaultValue=""
          render={({ field }) => (
            <Input
              {...field}
              type="password"
              id="password"
              placeholder="Enter your password"
              className={`mt-1 block w-full ${errors.password ? "border-red-500" : ""}`}
            />
          )}
        />
        {errors.password && <p className="text-red-500 text-sm mt-1">{errors.password.message}</p>}
      </div>
      <Button type="submit" className="w-full mb-4">
        Login
      </Button>
      <Button type="button" className="w-full" onClick={handleGoogleLogin}>
        Login with Github
      </Button>
    </form>
  );
};

export default LoginForm;
