import { cn } from '@/lib/utils';
const colors: Record<string, string> = {
  active: 'bg-emerald-100 text-emerald-800', delivered: 'bg-emerald-100 text-emerald-800',
  approved: 'bg-blue-100 text-blue-800', synced: 'bg-blue-100 text-blue-800',
  pending: 'bg-yellow-100 text-yellow-800', registered: 'bg-yellow-100 text-yellow-800',
  printing: 'bg-purple-100 text-purple-800', printed: 'bg-indigo-100 text-indigo-800',
  suspended: 'bg-orange-100 text-orange-800', needs_correction: 'bg-orange-100 text-orange-800',
  revoked: 'bg-red-100 text-red-800', rejected: 'bg-red-100 text-red-800', failed: 'bg-red-100 text-red-800',
  expired: 'bg-gray-100 text-gray-800', deceased: 'bg-gray-100 text-gray-800',
};
export function Badge({ status, className }: { status: string; className?: string }) {
  return <span className={cn('inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium', colors[status] || 'bg-gray-100 text-gray-800', className)}>{status}</span>;
}
