import { cva, type VariantProps } from "class-variance-authority";
import { forwardRef, type ButtonHTMLAttributes } from "react";

import { withPrefix } from "../constants";

const buttonVariants = cva(withPrefix("button"), {
  variants: {
    intent: {
      primary: withPrefix("button--intent-primary"),
      secondary: withPrefix("button--intent-secondary"),
      ghost: withPrefix("button--intent-ghost"),
    },
    size: {
      sm: withPrefix("button--size-sm"),
      md: withPrefix("button--size-md"),
      lg: withPrefix("button--size-lg"),
    },
    block: {
      true: withPrefix("button--block"),
      false: "",
    },
  },
  defaultVariants: {
    intent: "primary",
    size: "md",
    block: false,
  },
});

type NativeButtonProps = Omit<ButtonHTMLAttributes<HTMLButtonElement>, "color">;

export interface ButtonProps extends NativeButtonProps, VariantProps<typeof buttonVariants> {
  className?: string;
  pressed?: boolean;
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(function Button(
  { intent, size, block, className, pressed, type = "button", ...props },
  ref
) {
  return (
    <button
      {...props}
      ref={ref}
      type={type}
      aria-pressed={typeof pressed === "boolean" ? pressed : undefined}
      className={buttonVariants({ intent, size, block, className })}
    />
  );
});
