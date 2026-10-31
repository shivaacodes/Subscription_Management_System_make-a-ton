import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Link } from 'react-router-dom';
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';
import { Alert } from '@/components/ui/alert';
import { ArrowLeft } from 'lucide-react';

const forgotSchema = z.object({
  email: z.string().trim().email('Please enter a valid email address'),
});

type ForgotFormValues = z.infer<typeof forgotSchema>;

export function ForgotPasswordPage() {
  const [submitted, setSubmitted] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<ForgotFormValues>({
    resolver: zodResolver(forgotSchema),
  });

  const onSubmit = async () => {
    // In Phase 0, we acknowledge the submission honestly per platform guidelines.
    setSubmitted(true);
  };

  return (
    <div className="w-full max-w-md">
      <Card className="border-slate-200 shadow-sm">
        <CardHeader className="text-center pb-2">
          <CardTitle className="text-xl">Reset Password</CardTitle>
          <CardDescription>
            Enter your account email to receive a secure, single-use password reset link.
          </CardDescription>
        </CardHeader>

        {submitted ? (
          <CardContent className="space-y-4 pt-4">
            <Alert variant="info" title="Password Reset Workflow">
              In this environment, email delivery provider credentials (SMTP) are not yet configured.
              The password reset token infrastructure and table schema are ready in the database, and actual email delivery will be connected in a later deployment phase.
            </Alert>
            <div className="pt-2 text-center">
              <Link to="/login" className="inline-flex items-center text-sm text-indigo-600 font-medium hover:underline">
                <ArrowLeft className="h-4 w-4 mr-1" />
                Return to Login
              </Link>
            </div>
          </CardContent>
        ) : (
          <form onSubmit={handleSubmit(onSubmit)}>
            <CardContent className="space-y-4 pt-4">
              <div>
                <Label htmlFor="email">Registered Email</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="name@company.com"
                  autoComplete="email"
                  error={errors.email?.message}
                  {...register('email')}
                />
              </div>
            </CardContent>

            <CardFooter className="flex-col gap-3">
              <Button type="submit" variant="primary" className="w-full" isLoading={isSubmitting}>
                Send Reset Instructions
              </Button>
              <div className="text-center">
                <Link to="/login" className="inline-flex items-center text-xs text-slate-500 hover:text-slate-800">
                  <ArrowLeft className="h-3 w-3 mr-1" />
                  Back to Sign In
                </Link>
              </div>
            </CardFooter>
          </form>
        )}
      </Card>
    </div>
  );
}
