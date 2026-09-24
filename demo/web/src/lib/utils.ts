import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

/**
 * Merge Tailwind class names.
 * The ONLY place `clsx` may be imported directly (DESIGN.md forbids `clsx`
 * elsewhere — always route through `cn`).
 */
export function cn(...inputs: ClassValue[]): string {
  return twMerge(clsx(inputs));
}
