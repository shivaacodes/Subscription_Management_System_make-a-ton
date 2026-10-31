import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from './AuthContext';
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';
import { Alert } from '@/components/ui/alert';

const loginSchema = z.object({
  email: z.string().trim().email('Please enter a valid email address'),
  password: z.string().min(1, 'Password is required'),
});

type LoginFormValues = z.infer<typeof loginSchema>;

export function LoginPage() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [serverError, setServerError] = useState<string | null>(null);

  const from = (location.state as { from?: { pathname: string } })?.from?.pathname;

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  });

  const onSubmit = async (values: LoginFormValues) => {
    setServerError(null);
    try {
      await login(values);
      if (from) {
        navigate(from, { replace: true });
      } else {
        navigate('/admin', { replace: true });
      }
    } catch (err: any) {
      setServerError(err.message || 'Login failed. Please verify your credentials.');
    }
  };

  return (
    <div className="w-full max-w-md">
      <Card className="border-slate-200 shadow-sm">
        <CardHeader className="text-center pb-2">
          <CardTitle className="text-xl">Sign in to Veyra</CardTitle>
          <CardDescription>
            Enter your credentials to access the subscription management platform.
          </CardDescription>
        </CardHeader>

        <form onSubmit={handleSubmit(onSubmit)}>
          <CardContent className="space-y-4 pt-4">
            {serverError && (
              <Alert variant="danger">
                {serverError}
              </Alert>
            )}

            <div>
              <Label htmlFor="email">Work Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="name@company.com"
                autoComplete="email"
                error={errors.email?.message}
                {...register('email')}
              />
            </div>

            <div>
              <div className="flex items-center justify-between mb-1">
                <Label htmlFor="password">Password</Label>
                <Link
                  to="/forgot-password"
                  className="text-xs text-indigo-600 hover:text-indigo-800 transition-colors"
                >
                  Forgot password?
                </Link>
              </div>
              <Input
                id="password"
                type="password"
                placeholder="••••••••"
                autoComplete="current-password"
                error={errors.password?.message}
                {...register('password')}
              />
            </div>
          </CardContent>

          <CardFooter className="flex-col gap-3">
            <Button type="submit" variant="primary" className="w-full" isLoading={isSubmitting}>
              Sign In
            </Button>
            <p className="text-xs text-center text-slate-500">
              Need a customer account?{' '}
              <Link to="/register" className="text-indigo-600 font-medium hover:underline">
                Create one now
              </Link>
            </p>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}
