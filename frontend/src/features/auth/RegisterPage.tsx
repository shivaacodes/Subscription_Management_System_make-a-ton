import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from './AuthContext';
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';
import { Alert } from '@/components/ui/alert';

const registerSchema = z
  .object({
    email: z.string().trim().email('Please enter a valid email address'),
    companyName: z.string().trim().min(2, 'Company name is required'),
    password: z
      .string()
      .min(9, 'Password must be strictly greater than 8 characters')
      .regex(/[A-Z]/, 'Password must include at least one uppercase letter')
      .regex(/[a-z]/, 'Password must include at least one lowercase letter')
      .regex(/[^A-Za-z0-9]/, 'Password must include at least one special character'),
    confirmPassword: z.string(),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword'],
  });

type RegisterFormValues = z.infer<typeof registerSchema>;

export function RegisterPage() {
  const { register: registerUser } = useAuth();
  const navigate = useNavigate();
  const [serverError, setServerError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      email: '',
      companyName: '',
      password: '',
      confirmPassword: '',
    },
  });

  const onSubmit = async (values: RegisterFormValues) => {
    setServerError(null);
    try {
      await registerUser({
        email: values.email,
        password: values.password,
        companyName: values.companyName,
      });
      navigate('/portal', { replace: true });
    } catch (err: any) {
      setServerError(err.message || 'Registration failed. Please check your information.');
    }
  };

  return (
    <div className="w-full max-w-md">
      <Card className="border-slate-200 shadow-sm">
        <CardHeader className="text-center pb-2">
          <CardTitle className="text-xl">Create Customer Account</CardTitle>
          <CardDescription>
            Register to manage subscriptions and invoices with Brightloop Software.
          </CardDescription>
        </CardHeader>

        <form onSubmit={handleSubmit(onSubmit)}>
          <CardContent className="space-y-4 pt-4">
            {serverError && <Alert variant="danger">{serverError}</Alert>}

            <div>
              <Label htmlFor="companyName">Company or Organization</Label>
              <Input
                id="companyName"
                placeholder="Acme Corp"
                error={errors.companyName?.message}
                {...register('companyName')}
              />
            </div>

            <div>
              <Label htmlFor="email">Work Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="you@acme.com"
                autoComplete="email"
                error={errors.email?.message}
                {...register('email')}
              />
            </div>

            <div>
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                placeholder="••••••••••"
                autoComplete="new-password"
                error={errors.password?.message}
                {...register('password')}
              />
              <p className="text-[11px] text-slate-500 mt-1">
                Must be &gt; 8 characters, include uppercase, lowercase, and special character.
              </p>
            </div>

            <div>
              <Label htmlFor="confirmPassword">Confirm Password</Label>
              <Input
                id="confirmPassword"
                type="password"
                placeholder="••••••••••"
                autoComplete="new-password"
                error={errors.confirmPassword?.message}
                {...register('confirmPassword')}
              />
            </div>
          </CardContent>

          <CardFooter className="flex-col gap-3">
            <Button type="submit" variant="primary" className="w-full" isLoading={isSubmitting}>
              Complete Registration
            </Button>
            <p className="text-xs text-center text-slate-500">
              Already have an account?{' '}
              <Link to="/login" className="text-indigo-600 font-medium hover:underline">
                Sign In
              </Link>
            </p>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}
