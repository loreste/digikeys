import { cn } from '@/lib/utils';
import { InputHTMLAttributes, forwardRef } from 'react';
interface InputProps extends InputHTMLAttributes<HTMLInputElement> { label?: string; error?: string; }
const Input = forwardRef<HTMLInputElement, InputProps>(({ className, label, error, id, ...props }, ref) => (
  <div className="space-y-1">
    {label && <label htmlFor={id} className="block text-sm font-medium text-gray-700">{label}</label>}
    <input ref={ref} id={id} className={cn('block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm shadow-sm focus:border-emerald-500 focus:outline-none focus:ring-1 focus:ring-emerald-500 disabled:bg-gray-50', error && 'border-red-500', className)} {...props} />
    {error && <p className="text-sm text-red-600">{error}</p>}
  </div>
));
Input.displayName = 'Input';
export { Input };
