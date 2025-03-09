import CardAuthWrapper from "@/components/layout/auth/CardWrapper";
import AuthLayout from "@/components/layout/AuthLayout";
import LoginForm from "@/components/LoginForm";

export default function Login() {
  return (
    <AuthLayout>
      <CardAuthWrapper
        label="Create an account"
        title="Login"
        backButtonHref="/register"
        backButtonLable="Already have an account?"
      >
        <LoginForm />
      </CardAuthWrapper>
    </AuthLayout>
  )
}
