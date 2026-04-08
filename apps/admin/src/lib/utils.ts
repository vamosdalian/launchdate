import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export const buildPaginationRange = (current: number, total: number, delta = 1): (number | "ellipsis")[] => {
  const range: number[] = [];
  for (let i = 1; i <= total; i += 1) {
    if (i === 1 || i === total || (i >= current - delta && i <= current + delta)) {
      range.push(i);
    }
  }

  const result: (number | "ellipsis")[] = [];
  let previous: number | undefined;

  for (const value of range) {
    if (previous) {
      if (value - previous === 2) {
        result.push(previous + 1);
      } else if (value - previous > 2) {
        result.push("ellipsis");
      }
    }

    result.push(value);
    previous = value;
  }

  return result;
};
