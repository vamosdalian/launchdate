import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { motion } from "framer-motion"
import { AlertCircle, Check, Circle, Clock3, Loader2 } from "lucide-react"

import { cn } from "@/lib/utils"

export type TimelineSize = "sm" | "md" | "lg"
export type TimelineStatus = "completed" | "in-progress" | "pending" | "error"
export type TimelineColor = "primary" | "secondary" | "muted" | "accent" | "destructive"

const timelineVariants = cva("relative flex flex-col", {
  variants: {
    size: {
      sm: "gap-6",
      md: "gap-8",
      lg: "gap-10",
    },
  },
  defaultVariants: {
    size: "md",
  },
})

interface TimelineProps
  extends React.OlHTMLAttributes<HTMLOListElement>,
    VariantProps<typeof timelineVariants> {
  iconsize?: TimelineSize
}

interface TimelineItemProps extends Omit<React.LiHTMLAttributes<HTMLLIElement>, "title"> {
  date?: string | Date | number
  dateFormat?: Intl.DateTimeFormatOptions
  dateLabel?: React.ReactNode
  title?: React.ReactNode
  description?: React.ReactNode
  icon?: React.ReactNode
  iconColor?: TimelineColor
  status?: TimelineStatus
  connectorColor?: TimelineColor
  showConnector?: boolean
  iconsize?: TimelineSize
  loading?: boolean
  error?: React.ReactNode
  animationDelay?: number
}

interface TimelineTimeProps extends React.TimeHTMLAttributes<HTMLTimeElement> {
  date?: string | Date | number
  format?: Intl.DateTimeFormatOptions
}

const timelineItemGridClassName =
  "grid grid-cols-[auto_minmax(0,1fr)] gap-x-4 gap-y-3 md:grid-cols-[minmax(0,12rem)_auto_minmax(0,1fr)] md:gap-x-6"

const Timeline = React.forwardRef<HTMLOListElement, TimelineProps>(
  ({ className, size, children, ...props }, ref) => {
    const items = React.Children.toArray(children).filter(Boolean)

    if (items.length === 0) {
      return <TimelineEmpty />
    }

    return (
      <ol
        ref={ref}
        aria-label="Timeline"
        className={cn(
          timelineVariants({ size }),
          "relative mx-auto flex w-full max-w-4xl flex-col",
          className
        )}
        {...props}
      >
        {children}
      </ol>
    )
  }
)
Timeline.displayName = "Timeline"

const TimelineItem = React.forwardRef<HTMLLIElement, TimelineItemProps>(
  (
    {
      className,
      date,
      dateFormat,
      dateLabel,
      title,
      description,
      children,
      icon,
      iconColor,
      status = "completed",
      connectorColor,
      showConnector = true,
      iconsize = "md",
      loading,
      error,
      animationDelay = 0,
      ...props
    },
    ref
  ) => {
    if (loading) {
      return (
        <li ref={ref} className={cn("relative w-full", className)} role="status" {...props}>
          <div className={timelineItemGridClassName}>
            <div className="hidden pt-1 md:flex md:justify-end md:pr-2">
              <div className="h-4 w-28 animate-pulse rounded bg-muted" />
            </div>

            <div className="row-span-2 flex self-stretch flex-col items-center">
              <div className="relative z-10 mt-1 flex h-10 w-10 items-center justify-center rounded-full bg-zinc-800 shadow-sm">
                <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
              </div>
              {showConnector ? <div className="mt-3 min-h-16 w-[2px] flex-1 rounded-full bg-white/40" /> : null}
            </div>

            <div className="pt-1 md:hidden">
              <div className="h-4 w-28 animate-pulse rounded bg-muted" />
            </div>

            <div className="col-start-2 md:col-start-3">
              <div className="space-y-3 rounded-xl border border-border/60 bg-card/60 p-5 animate-pulse">
                <div className="h-6 w-56 rounded bg-muted" />
                <div className="h-4 w-32 rounded bg-muted" />
                <div className="flex gap-2">
                  <div className="h-6 w-24 rounded-full bg-muted" />
                  <div className="h-6 w-28 rounded-full bg-muted" />
                </div>
              </div>
            </div>
          </div>
        </li>
      )
    }

    if (error) {
      return (
        <li
          ref={ref}
          className={cn("relative w-full", className)}
          role="alert"
          {...props}
        >
          <div className={timelineItemGridClassName}>
            <div className="hidden pt-1 md:flex md:justify-end md:pr-2">
              {dateLabel || date ? (
                <TimelineTime date={date} format={dateFormat} className="text-right text-destructive/80">
                  {dateLabel}
                </TimelineTime>
              ) : null}
            </div>

            <div className="row-span-2 flex self-stretch flex-col items-center">
              <div className="relative z-10 mt-1 flex h-10 w-10 items-center justify-center rounded-full bg-red-500/20 shadow-sm">
                <AlertCircle className="h-5 w-5 text-destructive" />
              </div>
              {showConnector ? <div className="mt-3 min-h-16 w-[2px] flex-1 rounded-full bg-white/40" /> : null}
            </div>

            <div className="pt-1 md:hidden">
              {dateLabel || date ? (
                <TimelineTime date={date} format={dateFormat} className="text-destructive/80">
                  {dateLabel}
                </TimelineTime>
              ) : null}
            </div>

            <div className="col-start-2 space-y-2 rounded-xl border border-destructive/40 bg-destructive/10 p-5 md:col-start-3">
              {title ? <TimelineTitle className="text-destructive">{title}</TimelineTitle> : null}
              <TimelineDescription className="text-destructive/90">{error}</TimelineDescription>
            </div>
          </div>
        </li>
      )
    }

    const hasMeta = Boolean(date || dateLabel || title || description)
    const animationProps = {
      initial: { opacity: 0, y: 24 },
      animate: { opacity: 1, y: 0 },
      transition: {
        duration: 0.45,
        delay: animationDelay,
        ease: "easeOut" as const,
      },
    }

    return (
      <li ref={ref} className={cn("relative w-full", className)} {...props}>
        <div className={timelineItemGridClassName}>
          <motion.div
            className="hidden pt-1 md:flex md:justify-end md:pr-2"
            {...animationProps}
          >
            {dateLabel || date ? (
              <TimelineTime date={date} format={dateFormat} className="text-right">
                {dateLabel}
              </TimelineTime>
            ) : null}
          </motion.div>

          <motion.div
            className="row-span-2 flex self-stretch flex-col items-center"
            {...animationProps}
          >
            <div className="relative z-10 mt-1">
              <TimelineIcon icon={icon} color={iconColor} status={status} iconSize={iconsize} />
            </div>
            {showConnector ? (
              <TimelineConnector status={status} color={connectorColor} className="mt-3" />
            ) : null}
          </motion.div>

          <motion.div className="pt-1 md:hidden" {...animationProps}>
            {dateLabel || date ? (
              <TimelineTime date={date} format={dateFormat}>
                {dateLabel}
              </TimelineTime>
            ) : null}
          </motion.div>

          <div className="col-start-2 flex min-w-0 flex-col gap-4 md:col-start-3">
            {hasMeta ? (
              <TimelineContent>
                {title ? (
                  <TimelineHeader>
                    <TimelineTitle>{title}</TimelineTitle>
                  </TimelineHeader>
                ) : null}
                {description ? <TimelineDescription>{description}</TimelineDescription> : null}
              </TimelineContent>
            ) : null}
            {children}
          </div>
        </div>
      </li>
    )
  }
)
TimelineItem.displayName = "TimelineItem"

const TimelineTime = React.forwardRef<HTMLTimeElement, TimelineTimeProps>(
  ({ className, date, format, children, ...props }, ref) => {
    const value = React.useMemo(() => {
      if (date === undefined || date === null || date === "") {
        return ""
      }

      try {
        const dateValue = new Date(date)

        if (Number.isNaN(dateValue.getTime())) {
          return ""
        }

        return new Intl.DateTimeFormat("en-US", {
          year: "numeric",
          month: "long",
          day: "numeric",
          ...format,
        }).format(dateValue)
      } catch {
        return ""
      }
    }, [date, format])

    const dateTime = React.useMemo(() => {
      if (date === undefined || date === null || date === "") {
        return undefined
      }

      const dateValue = new Date(date)
      if (Number.isNaN(dateValue.getTime())) {
        return undefined
      }

      return dateValue.toISOString()
    }, [date])

    return (
      <time
        ref={ref}
        dateTime={dateTime}
        className={cn("text-sm font-medium tracking-tight text-muted-foreground", className)}
        {...props}
      >
        {children || value}
      </time>
    )
  }
)
TimelineTime.displayName = "TimelineTime"

const TimelineConnector = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement> & {
    status?: TimelineStatus
    color?: TimelineColor
  }
>(({ className, status = "completed", color, ...props }, ref) => {
  const colorClass = color
    ? {
        primary: "bg-white/40",
        secondary: "bg-white/40",
        muted: "bg-white/40",
        accent: "bg-white/40",
        destructive: "bg-white/40",
      }[color]
    : {
        completed: "bg-white/40",
        "in-progress": "bg-white/40",
        pending: "bg-white/40",
        error: "bg-white/40",
      }[status]

  return <div ref={ref} className={cn("min-h-16 w-[2px] flex-1 rounded-full", colorClass, className)} {...props} />
})
TimelineConnector.displayName = "TimelineConnector"

const TimelineHeader = React.forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <div ref={ref} className={cn("flex items-start gap-3", className)} {...props} />
  )
)
TimelineHeader.displayName = "TimelineHeader"

const TimelineTitle = React.forwardRef<HTMLHeadingElement, React.HTMLAttributes<HTMLHeadingElement>>(
  ({ className, ...props }, ref) => (
    <h3 ref={ref} className={cn("text-xl font-semibold leading-tight tracking-tight", className)} {...props} />
  )
)
TimelineTitle.displayName = "TimelineTitle"

const TimelineDescription = React.forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <div ref={ref} className={cn("text-sm text-muted-foreground", className)} {...props} />
  )
)
TimelineDescription.displayName = "TimelineDescription"

const TimelineContent = React.forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <div ref={ref} className={cn("flex min-w-0 flex-col gap-3", className)} {...props} />
  )
)
TimelineContent.displayName = "TimelineContent"

const TimelineEmpty = React.forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
  ({ className, children, ...props }, ref) => (
    <div
      ref={ref}
      className={cn(
        "flex min-h-48 w-full items-center justify-center rounded-xl border border-dashed border-border bg-card/50 p-8 text-center text-muted-foreground",
        className
      )}
      {...props}
    >
      {children || "No timeline items to display"}
    </div>
  )
)
TimelineEmpty.displayName = "TimelineEmpty"

function TimelineIcon({
  icon,
  color,
  status = "completed",
  iconSize = "md",
}: {
  icon?: React.ReactNode
  color?: TimelineColor
  status?: TimelineStatus
  iconSize?: TimelineSize
}) {
  const sizeClass = {
    sm: "h-8 w-8",
    md: "h-10 w-10",
    lg: "h-12 w-12",
  }[iconSize]

  const iconClass = {
    sm: "h-4 w-4",
    md: "h-5 w-5",
    lg: "h-6 w-6",
  }[iconSize]

  const colorClass = color
    ? {
        primary: "bg-[var(--color-button-bg)] text-[var(--color-button-fg)]",
        secondary: "bg-[var(--color-card-bg)] text-[var(--color-card-fg)]",
        muted: "bg-[var(--color-card-bg)] text-[var(--color-card-fg)]",
        accent: "bg-[var(--color-button-bg)] text-[var(--color-button-fg)]",
        destructive: "bg-red-500 text-white",
      }[color]
    : {
        completed: "bg-[var(--color-button-bg)] text-[var(--color-button-fg)]",
        "in-progress": "bg-[var(--color-button-bg)] text-[var(--color-button-fg)]",
        pending: "bg-[var(--color-button-bg)] text-[var(--color-button-fg)]",
        error: "bg-red-500 text-white",
      }[status]

  const fallbackIcon = {
    completed: <Check className={iconClass} />,
    "in-progress": <Clock3 className={iconClass} />,
    pending: <Circle className={iconClass} />,
    error: <AlertCircle className={iconClass} />,
  }[status]

  return (
    <div
      className={cn(
        "flex items-center justify-center rounded-full shadow-sm",
        sizeClass,
        colorClass
      )}
    >
      {icon ? <span className="flex items-center justify-center">{icon}</span> : fallbackIcon}
    </div>
  )
}

export {
  Timeline,
  TimelineConnector,
  TimelineContent,
  TimelineDescription,
  TimelineEmpty,
  TimelineHeader,
  TimelineItem,
  TimelineTime,
  TimelineTitle,
}