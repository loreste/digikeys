import { cn } from '@/lib/utils';
import { ButtonHTMLAttributes, forwardRef } from 'react';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
}

const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = 'primary', size = 'md', ...props }, ref) => (
    <button ref={ref} className={cn(
      'inline-flex items-center justify-center rounded-lg font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none',
      { 'bg-emerald-800 text-white hover:bg-emerald-900 focus:ring-emerald-500': variant === 'primary',
        'bg-amber-600 text-white hover:bg-amber-700 focus:ring-amber-500': variant === 'secondary',
        'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500': variant === 'danger',
        'bg-transparent hover:bg-gray-100 text-gray-700': variant === 'ghost' },
      { 'px-3 py-1.5 text-sm': size === 'sm', 'px-4 py-2 text-sm': size === 'md', 'px-6 py-3 text-base': size === 'lg' },
      className)} {...props} />
  )
);
Button.displayName = 'Button';
export { Button };
